// Package extension 定义 Synapse 六大扩展点接口和全局注册机制
//
// 扩展点模型是 Synapse 的核心架构，以 Store 为底座，
// 六个独立扩展点（Source/Processor/Store/Indexer/Consumer/Auditor）按需组合。
// 每个扩展点都是一个独立的 Go 接口，社区可以针对任意一个扩展点贡献新实现。
package extension

import (
	"context"

	"github.com/tunsuy/synapse/pkg/model"
)

// Source 从外部数据源获取原始内容
// 社区共建场景：想接入 RSS 订阅？实现 Source 接口即可
type Source interface {
	// Name 返回数据源名称
	Name() string

	// Fetch 获取原始内容
	Fetch(ctx context.Context, opts model.FetchOptions) ([]model.RawContent, error)
}

// Processor 将原始内容处理为结构化知识
// 社区共建场景：想用本地 LLM 处理？实现 Processor 接口即可
type Processor interface {
	// Name 返回处理引擎名称
	Name() string

	// Process 将原始内容转换为知识文件
	Process(ctx context.Context, raw []model.RawContent) ([]model.KnowledgeFile, error)
}

// Store 提供知识文件的持久化存储
// Store 是整个架构的底座，所有模块都依赖它
type Store interface {
	// Name 返回存储后端名称
	Name() string

	// Read 读取知识文件
	Read(ctx context.Context, path string) (model.KnowledgeFile, error)

	// Write 写入知识文件
	Write(ctx context.Context, file model.KnowledgeFile) error

	// Delete 删除知识文件
	Delete(ctx context.Context, path string) error

	// List 列出指定目录下的知识文件
	List(ctx context.Context, dir string, opts model.ListOptions) ([]model.FileInfo, error)

	// Exists 检查文件是否存在
	Exists(ctx context.Context, path string) (bool, error)
}

// VersionedStore 支持版本控制的存储后端（可选能力，通过类型断言检查）
type VersionedStore interface {
	Store

	// Commit 提交当前变更
	Commit(ctx context.Context, message string) error

	// History 获取文件变更历史
	History(ctx context.Context, path string) ([]ChangeRecord, error)
}

// ChangeRecord 变更记录
type ChangeRecord struct {
	// Hash 变更哈希
	Hash string `json:"hash"`

	// Message 变更说明
	Message string `json:"message"`

	// Author 作者
	Author string `json:"author"`

	// Timestamp 变更时间
	Timestamp string `json:"timestamp"`
}

// Indexer 提供知识库检索能力
// 社区共建场景：想加向量检索？实现 Indexer 接口即可
type Indexer interface {
	// Name 返回检索引擎名称
	Name() string

	// Index 索引知识文件
	Index(ctx context.Context, file model.KnowledgeFile) error

	// Build 构建/重建全量索引
	Build(ctx context.Context, store Store) error

	// Search 搜索知识库
	Search(ctx context.Context, query string, opts model.SearchOptions) ([]model.SearchResult, error)
}

// Consumer 将知识输出为特定消费形式
// 社区共建场景：想导出 Anki 闪卡？实现 Consumer 接口即可
type Consumer interface {
	// Name 返回消费端名称
	Name() string

	// Consume 消费知识库，生成输出产物
	Consume(ctx context.Context, store Store, opts model.ConsumeOptions) error
}

// Auditor 检查知识库健康状态
// 社区共建场景：想自定义审计规则？实现 Auditor 接口即可
type Auditor interface {
	// Name 返回审计器名称
	Name() string

	// Audit 执行知识库审计
	Audit(ctx context.Context, store Store) (model.AuditReport, error)

	// Fix 自动修复问题（可选，实现时不支持可返回错误）
	Fix(ctx context.Context, store Store, issues []model.AuditIssue) (model.FixResult, error)
}
