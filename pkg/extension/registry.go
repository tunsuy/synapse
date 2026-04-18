package extension

import (
	"fmt"
	"sort"
	"sync"
)

// Factory 是扩展点实现的工厂函数类型
// config 来自 .synapse/config.yaml 中该扩展点的 config 字段
type SourceFactory func(config map[string]any) (Source, error)
type ProcessorFactory func(config map[string]any) (Processor, error)
type StoreFactory func(config map[string]any) (Store, error)
type IndexerFactory func(config map[string]any) (Indexer, error)
type ConsumerFactory func(config map[string]any) (Consumer, error)
type AuditorFactory func(config map[string]any) (Auditor, error)

// Registry 是扩展点的全局注册表
// 通过 init() 自注册模式，各扩展点实现在导入时自动注册
type Registry struct {
	mu         sync.RWMutex
	sources    map[string]SourceFactory
	processors map[string]ProcessorFactory
	stores     map[string]StoreFactory
	indexers   map[string]IndexerFactory
	consumers  map[string]ConsumerFactory
	auditors   map[string]AuditorFactory
}

// globalRegistry 全局单例注册表
var globalRegistry = &Registry{
	sources:    make(map[string]SourceFactory),
	processors: make(map[string]ProcessorFactory),
	stores:     make(map[string]StoreFactory),
	indexers:   make(map[string]IndexerFactory),
	consumers:  make(map[string]ConsumerFactory),
	auditors:   make(map[string]AuditorFactory),
}

// --- Source ---

// RegisterSource 注册一个 Source 实现
func RegisterSource(name string, factory SourceFactory) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	if _, exists := globalRegistry.sources[name]; exists {
		panic(fmt.Sprintf("source %q already registered", name))
	}
	globalRegistry.sources[name] = factory
}

// GetSource 根据名称创建 Source 实例
func GetSource(name string, config map[string]any) (Source, error) {
	globalRegistry.mu.RLock()
	factory, ok := globalRegistry.sources[name]
	globalRegistry.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown source: %s", name)
	}
	return factory(config)
}

// ListSources 列出所有已注册的 Source 名称
func ListSources() []string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	return sortedKeys(globalRegistry.sources)
}

// --- Processor ---

// RegisterProcessor 注册一个 Processor 实现
func RegisterProcessor(name string, factory ProcessorFactory) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	if _, exists := globalRegistry.processors[name]; exists {
		panic(fmt.Sprintf("processor %q already registered", name))
	}
	globalRegistry.processors[name] = factory
}

// GetProcessor 根据名称创建 Processor 实例
func GetProcessor(name string, config map[string]any) (Processor, error) {
	globalRegistry.mu.RLock()
	factory, ok := globalRegistry.processors[name]
	globalRegistry.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown processor: %s", name)
	}
	return factory(config)
}

// ListProcessors 列出所有已注册的 Processor 名称
func ListProcessors() []string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	return sortedKeys(globalRegistry.processors)
}

// --- Store ---

// RegisterStore 注册一个 Store 实现
func RegisterStore(name string, factory StoreFactory) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	if _, exists := globalRegistry.stores[name]; exists {
		panic(fmt.Sprintf("store %q already registered", name))
	}
	globalRegistry.stores[name] = factory
}

// GetStore 根据名称创建 Store 实例
func GetStore(name string, config map[string]any) (Store, error) {
	globalRegistry.mu.RLock()
	factory, ok := globalRegistry.stores[name]
	globalRegistry.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown store: %s", name)
	}
	return factory(config)
}

// ListStores 列出所有已注册的 Store 名称
func ListStores() []string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	return sortedKeys(globalRegistry.stores)
}

// --- Indexer ---

// RegisterIndexer 注册一个 Indexer 实现
func RegisterIndexer(name string, factory IndexerFactory) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	if _, exists := globalRegistry.indexers[name]; exists {
		panic(fmt.Sprintf("indexer %q already registered", name))
	}
	globalRegistry.indexers[name] = factory
}

// GetIndexer 根据名称创建 Indexer 实例
func GetIndexer(name string, config map[string]any) (Indexer, error) {
	globalRegistry.mu.RLock()
	factory, ok := globalRegistry.indexers[name]
	globalRegistry.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown indexer: %s", name)
	}
	return factory(config)
}

// ListIndexers 列出所有已注册的 Indexer 名称
func ListIndexers() []string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	return sortedKeys(globalRegistry.indexers)
}

// --- Consumer ---

// RegisterConsumer 注册一个 Consumer 实现
func RegisterConsumer(name string, factory ConsumerFactory) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	if _, exists := globalRegistry.consumers[name]; exists {
		panic(fmt.Sprintf("consumer %q already registered", name))
	}
	globalRegistry.consumers[name] = factory
}

// GetConsumer 根据名称创建 Consumer 实例
func GetConsumer(name string, config map[string]any) (Consumer, error) {
	globalRegistry.mu.RLock()
	factory, ok := globalRegistry.consumers[name]
	globalRegistry.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown consumer: %s", name)
	}
	return factory(config)
}

// ListConsumers 列出所有已注册的 Consumer 名称
func ListConsumers() []string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	return sortedKeys(globalRegistry.consumers)
}

// --- Auditor ---

// RegisterAuditor 注册一个 Auditor 实现
func RegisterAuditor(name string, factory AuditorFactory) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	if _, exists := globalRegistry.auditors[name]; exists {
		panic(fmt.Sprintf("auditor %q already registered", name))
	}
	globalRegistry.auditors[name] = factory
}

// GetAuditor 根据名称创建 Auditor 实例
func GetAuditor(name string, config map[string]any) (Auditor, error) {
	globalRegistry.mu.RLock()
	factory, ok := globalRegistry.auditors[name]
	globalRegistry.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown auditor: %s", name)
	}
	return factory(config)
}

// ListAuditors 列出所有已注册的 Auditor 名称
func ListAuditors() []string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	return sortedKeys(globalRegistry.auditors)
}

// --- 辅助函数 ---

// sortedKeys 从 map 中提取键并排序
func sortedKeys[V any](m map[string]V) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
