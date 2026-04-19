package extension

import (
	"context"
	"testing"

	"github.com/tunsuy/synapse/pkg/model"
)

// mockSource 用于测试的 Source mock
type mockSource struct {
	name string
}

func (m *mockSource) Name() string { return m.name }
func (m *mockSource) Fetch(_ context.Context, _ model.FetchOptions) ([]model.RawContent, error) {
	return nil, nil
}

// mockProcessor 用于测试的 Processor mock
type mockProcessor struct {
	name string
}

func (m *mockProcessor) Name() string { return m.name }
func (m *mockProcessor) Process(_ context.Context, _ []model.RawContent) ([]model.KnowledgeFile, error) {
	return nil, nil
}

// mockStore 用于测试的 Store mock
type mockStore struct {
	name string
}

func (m *mockStore) Name() string { return m.name }
func (m *mockStore) Init(_ context.Context, _ InitOptions) error {
	return nil
}
func (m *mockStore) Initialized(_ context.Context) (bool, error) { return false, nil }
func (m *mockStore) Read(_ context.Context, _ string) (model.KnowledgeFile, error) {
	return model.KnowledgeFile{}, nil
}
func (m *mockStore) Write(_ context.Context, _ model.KnowledgeFile) error   { return nil }
func (m *mockStore) Delete(_ context.Context, _ string) error                { return nil }
func (m *mockStore) List(_ context.Context, _ string, _ model.ListOptions) ([]model.FileInfo, error) {
	return nil, nil
}
func (m *mockStore) Exists(_ context.Context, _ string) (bool, error) { return false, nil }

// mockIndexer 用于测试的 Indexer mock
type mockIndexer struct {
	name string
}

func (m *mockIndexer) Name() string                                       { return m.name }
func (m *mockIndexer) Index(_ context.Context, _ model.KnowledgeFile) error { return nil }
func (m *mockIndexer) Build(_ context.Context, _ Store) error              { return nil }
func (m *mockIndexer) Search(_ context.Context, _ string, _ model.SearchOptions) ([]model.SearchResult, error) {
	return nil, nil
}

// mockConsumer 用于测试的 Consumer mock
type mockConsumer struct {
	name string
}

func (m *mockConsumer) Name() string { return m.name }
func (m *mockConsumer) Consume(_ context.Context, _ Store, _ model.ConsumeOptions) error {
	return nil
}

// mockAuditor 用于测试的 Auditor mock
type mockAuditor struct {
	name string
}

func (m *mockAuditor) Name() string { return m.name }
func (m *mockAuditor) Audit(_ context.Context, _ Store) (model.AuditReport, error) {
	return model.AuditReport{}, nil
}
func (m *mockAuditor) Fix(_ context.Context, _ Store, _ []model.AuditIssue) (model.FixResult, error) {
	return model.FixResult{}, nil
}

// resetRegistry 重置全局注册表（测试辅助函数）
func resetRegistry() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.sources = make(map[string]SourceFactory)
	globalRegistry.processors = make(map[string]ProcessorFactory)
	globalRegistry.stores = make(map[string]StoreFactory)
	globalRegistry.indexers = make(map[string]IndexerFactory)
	globalRegistry.consumers = make(map[string]ConsumerFactory)
	globalRegistry.auditors = make(map[string]AuditorFactory)
}

func TestRegisterAndGetSource(t *testing.T) {
	resetRegistry()
	defer resetRegistry()

	RegisterSource("test-source", func(config map[string]any) (Source, error) {
		return &mockSource{name: "test-source"}, nil
	})

	src, err := GetSource("test-source", nil)
	if err != nil {
		t.Fatalf("GetSource() error: %v", err)
	}
	if src.Name() != "test-source" {
		t.Errorf("Name() = %q, want %q", src.Name(), "test-source")
	}
}

func TestGetSource_Unknown(t *testing.T) {
	resetRegistry()
	defer resetRegistry()

	_, err := GetSource("nonexistent", nil)
	if err == nil {
		t.Fatal("GetSource() expected error for unknown source")
	}
}

func TestRegisterSource_DuplicatePanics(t *testing.T) {
	resetRegistry()
	defer resetRegistry()

	factory := func(config map[string]any) (Source, error) {
		return &mockSource{name: "dup"}, nil
	}
	RegisterSource("dup-source", factory)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on duplicate registration")
		}
	}()
	RegisterSource("dup-source", factory)
}

func TestListSources(t *testing.T) {
	resetRegistry()
	defer resetRegistry()

	factory := func(config map[string]any) (Source, error) {
		return &mockSource{}, nil
	}
	RegisterSource("beta-source", factory)
	RegisterSource("alpha-source", factory)

	names := ListSources()
	if len(names) != 2 {
		t.Fatalf("ListSources() returned %d, want 2", len(names))
	}
	if names[0] != "alpha-source" || names[1] != "beta-source" {
		t.Errorf("ListSources() = %v, want sorted [alpha-source, beta-source]", names)
	}
}

func TestRegisterAndGetStore(t *testing.T) {
	resetRegistry()
	defer resetRegistry()

	RegisterStore("test-store", func(config map[string]any) (Store, error) {
		return &mockStore{name: "test-store"}, nil
	})

	store, err := GetStore("test-store", nil)
	if err != nil {
		t.Fatalf("GetStore() error: %v", err)
	}
	if store.Name() != "test-store" {
		t.Errorf("Name() = %q, want %q", store.Name(), "test-store")
	}
}

func TestRegisterAndGetProcessor(t *testing.T) {
	resetRegistry()
	defer resetRegistry()

	RegisterProcessor("test-processor", func(config map[string]any) (Processor, error) {
		return &mockProcessor{name: "test-processor"}, nil
	})

	proc, err := GetProcessor("test-processor", nil)
	if err != nil {
		t.Fatalf("GetProcessor() error: %v", err)
	}
	if proc.Name() != "test-processor" {
		t.Errorf("Name() = %q, want %q", proc.Name(), "test-processor")
	}
}

func TestRegisterAndGetIndexer(t *testing.T) {
	resetRegistry()
	defer resetRegistry()

	RegisterIndexer("test-indexer", func(config map[string]any) (Indexer, error) {
		return &mockIndexer{name: "test-indexer"}, nil
	})

	idx, err := GetIndexer("test-indexer", nil)
	if err != nil {
		t.Fatalf("GetIndexer() error: %v", err)
	}
	if idx.Name() != "test-indexer" {
		t.Errorf("Name() = %q, want %q", idx.Name(), "test-indexer")
	}
}

func TestRegisterAndGetConsumer(t *testing.T) {
	resetRegistry()
	defer resetRegistry()

	RegisterConsumer("test-consumer", func(config map[string]any) (Consumer, error) {
		return &mockConsumer{name: "test-consumer"}, nil
	})

	con, err := GetConsumer("test-consumer", nil)
	if err != nil {
		t.Fatalf("GetConsumer() error: %v", err)
	}
	if con.Name() != "test-consumer" {
		t.Errorf("Name() = %q, want %q", con.Name(), "test-consumer")
	}
}

func TestRegisterAndGetAuditor(t *testing.T) {
	resetRegistry()
	defer resetRegistry()

	RegisterAuditor("test-auditor", func(config map[string]any) (Auditor, error) {
		return &mockAuditor{name: "test-auditor"}, nil
	})

	aud, err := GetAuditor("test-auditor", nil)
	if err != nil {
		t.Fatalf("GetAuditor() error: %v", err)
	}
	if aud.Name() != "test-auditor" {
		t.Errorf("Name() = %q, want %q", aud.Name(), "test-auditor")
	}
}

func TestListStores_Empty(t *testing.T) {
	resetRegistry()
	defer resetRegistry()

	names := ListStores()
	if len(names) != 0 {
		t.Errorf("ListStores() returned %d items, want 0", len(names))
	}
}
