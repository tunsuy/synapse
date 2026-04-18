// Package engine 是 Synapse 的核心编排引擎
// 它读取配置，实例化各扩展点，协调它们的执行
package engine

import (
	"context"
	"fmt"
	"log"

	"github.com/tunsuy/synapse/internal/config"
	"github.com/tunsuy/synapse/pkg/extension"
	"github.com/tunsuy/synapse/pkg/model"
)

// Engine 核心编排引擎
// 负责读取配置 → 实例化扩展点 → 协调执行
type Engine struct {
	cfg       *config.Config
	sources   []extension.Source
	processor extension.Processor
	store     extension.Store
	indexer   extension.Indexer   // 可选
	consumers []extension.Consumer // 可选
	auditor   extension.Auditor   // 可选
}

// New 从配置文件创建引擎实例
func New(configPath string) (*Engine, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return NewFromConfig(cfg)
}

// NewFromConfig 从配置对象创建引擎实例
func NewFromConfig(cfg *config.Config) (*Engine, error) {
	e := &Engine{cfg: cfg}

	var err error

	// 1. 实例化 Store（底座，最先初始化）
	e.store, err = extension.GetStore(
		cfg.Synapse.Store.Name,
		cfg.Synapse.Store.Config,
	)
	if err != nil {
		return nil, fmt.Errorf("init store %q: %w", cfg.Synapse.Store.Name, err)
	}

	// 2. 实例化 Sources（可多个）
	for _, srcCfg := range cfg.Synapse.Sources {
		if !srcCfg.IsEnabled() {
			continue
		}
		src, err := extension.GetSource(srcCfg.Name, srcCfg.Config)
		if err != nil {
			return nil, fmt.Errorf("init source %q: %w", srcCfg.Name, err)
		}
		e.sources = append(e.sources, src)
	}

	// 3. 实例化 Processor（可选，M2 阶段才有实际实现）
	if cfg.Synapse.Processor != nil {
		e.processor, err = extension.GetProcessor(
			cfg.Synapse.Processor.Name,
			cfg.Synapse.Processor.Config,
		)
		if err != nil {
			// Processor 可选，找不到实现不阻断
			log.Printf("WARN: init processor %q: %v", cfg.Synapse.Processor.Name, err)
		}
	}

	// 4. 实例化 Indexer（可选）
	if cfg.Synapse.Indexer != nil {
		e.indexer, err = extension.GetIndexer(
			cfg.Synapse.Indexer.Name,
			cfg.Synapse.Indexer.Config,
		)
		if err != nil {
			log.Printf("WARN: init indexer %q: %v", cfg.Synapse.Indexer.Name, err)
		}
	}

	// 5. 实例化 Consumers（可多个）
	for _, conCfg := range cfg.Synapse.Consumers {
		if !conCfg.IsEnabled() {
			continue
		}
		c, err := extension.GetConsumer(conCfg.Name, conCfg.Config)
		if err != nil {
			log.Printf("WARN: init consumer %q: %v", conCfg.Name, err)
			continue
		}
		e.consumers = append(e.consumers, c)
	}

	// 6. 实例化 Auditor（可选）
	if cfg.Synapse.Auditor != nil {
		e.auditor, err = extension.GetAuditor(
			cfg.Synapse.Auditor.Name,
			cfg.Synapse.Auditor.Config,
		)
		if err != nil {
			log.Printf("WARN: init auditor %q: %v", cfg.Synapse.Auditor.Name, err)
		}
	}

	return e, nil
}

// Store 返回引擎使用的 Store 实例
func (e *Engine) Store() extension.Store {
	return e.store
}

// Collect 执行采集流程：Source.Fetch → Processor.Process → Store.Write
func (e *Engine) Collect(ctx context.Context) error {
	if len(e.sources) == 0 {
		return fmt.Errorf("no sources configured")
	}

	// 1. 从所有启用的 Source 采集原始内容
	var allRaw []model.RawContent
	for _, src := range e.sources {
		raw, err := src.Fetch(ctx, model.FetchOptions{})
		if err != nil {
			// 单个 Source 失败不阻断整体流程
			log.Printf("WARN: source %s fetch failed: %v", src.Name(), err)
			continue
		}
		allRaw = append(allRaw, raw...)
	}

	if len(allRaw) == 0 {
		log.Println("no raw content collected")
		return nil
	}

	// 2. 用 Processor 处理原始内容为结构化知识
	if e.processor == nil {
		return fmt.Errorf("no processor configured")
	}

	files, err := e.processor.Process(ctx, allRaw)
	if err != nil {
		return fmt.Errorf("process raw content: %w", err)
	}

	// 3. 写入 Store
	for _, f := range files {
		if err := e.store.Write(ctx, f); err != nil {
			log.Printf("WARN: write %s failed: %v", f.Path, err)
		}
	}

	// 4. 如果有 Indexer，更新索引
	if e.indexer != nil {
		for _, f := range files {
			if err := e.indexer.Index(ctx, f); err != nil {
				log.Printf("WARN: index %s failed: %v", f.Path, err)
			}
		}
	}

	log.Printf("collected %d raw items, produced %d knowledge files", len(allRaw), len(files))
	return nil
}

// Search 执行检索
func (e *Engine) Search(ctx context.Context, query string, opts model.SearchOptions) ([]model.SearchResult, error) {
	if e.indexer == nil {
		return nil, fmt.Errorf("no indexer configured")
	}
	return e.indexer.Search(ctx, query, opts)
}

// Audit 执行审计
func (e *Engine) Audit(ctx context.Context) (model.AuditReport, error) {
	if e.auditor == nil {
		return model.AuditReport{}, fmt.Errorf("no auditor configured")
	}
	return e.auditor.Audit(ctx, e.store)
}

// Publish 触发所有消费端
func (e *Engine) Publish(ctx context.Context) error {
	if len(e.consumers) == 0 {
		return fmt.Errorf("no consumers configured")
	}
	for _, c := range e.consumers {
		if err := c.Consume(ctx, e.store, model.ConsumeOptions{}); err != nil {
			log.Printf("WARN: consumer %s failed: %v", c.Name(), err)
		}
	}
	return nil
}
