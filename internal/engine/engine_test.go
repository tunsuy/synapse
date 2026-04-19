package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tunsuy/synapse/internal/config"
	"github.com/tunsuy/synapse/internal/schema"
	"github.com/tunsuy/synapse/pkg/extension"
	"github.com/tunsuy/synapse/pkg/model"
)

// --- 测试用 Mock 实现 ---

type mockSource struct {
	name    string
	content []model.RawContent
	err     error
}

func (m *mockSource) Name() string { return m.name }
func (m *mockSource) Fetch(_ context.Context, _ model.FetchOptions) ([]model.RawContent, error) {
	return m.content, m.err
}

type mockProcessor struct {
	name  string
	files []model.KnowledgeFile
	err   error
}

func (m *mockProcessor) Name() string { return m.name }
func (m *mockProcessor) Process(_ context.Context, _ []model.RawContent) ([]model.KnowledgeFile, error) {
	return m.files, m.err
}

type mockStore struct {
	name    string
	written []model.KnowledgeFile
}

func (m *mockStore) Name() string { return m.name }
func (m *mockStore) Init(_ context.Context, _ extension.InitOptions) error {
	return nil
}
func (m *mockStore) Initialized(_ context.Context) (bool, error) { return false, nil }
func (m *mockStore) Read(_ context.Context, _ string) (model.KnowledgeFile, error) {
	return model.KnowledgeFile{}, nil
}
func (m *mockStore) Write(_ context.Context, f model.KnowledgeFile) error {
	m.written = append(m.written, f)
	return nil
}
func (m *mockStore) Delete(_ context.Context, _ string) error { return nil }
func (m *mockStore) List(_ context.Context, _ string, _ model.ListOptions) ([]model.FileInfo, error) {
	return nil, nil
}
func (m *mockStore) Exists(_ context.Context, _ string) (bool, error) { return false, nil }

// --- 测试用例 ---

func TestNewFromConfig_BasicSetup(t *testing.T) {
	t.Parallel()

	// 注册测试用的扩展点
	testSourceName := "test-source-engine"
	testProcessorName := "test-processor-engine"
	testStoreName := "test-store-engine"

	extension.RegisterSource(testSourceName, func(_ map[string]any) (extension.Source, error) {
		return &mockSource{name: testSourceName}, nil
	})
	extension.RegisterProcessor(testProcessorName, func(_ map[string]any) (extension.Processor, error) {
		return &mockProcessor{name: testProcessorName}, nil
	})
	extension.RegisterStore(testStoreName, func(_ map[string]any) (extension.Store, error) {
		return &mockStore{name: testStoreName}, nil
	})

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Sources: []config.SourceConfig{
				{Name: testSourceName},
			},
			Processor: &config.ExtensionConfig{Name: testProcessorName},
			Store:     config.ExtensionConfig{Name: testStoreName},
		},
	}

	eng, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig() error: %v", err)
	}

	if eng.Store() == nil {
		t.Fatal("Engine.Store() should not be nil")
	}
	if eng.Schema() == nil {
		t.Fatal("Engine.Schema() should not be nil")
	}
	if eng.Config() != cfg {
		t.Fatal("Engine.Config() should return the config passed in")
	}
}

func TestNewFromConfig_WithSchema(t *testing.T) {
	t.Parallel()

	testStoreName := "test-store-schema"
	extension.RegisterStore(testStoreName, func(_ map[string]any) (extension.Store, error) {
		return &mockStore{name: testStoreName}, nil
	})

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Store:   config.ExtensionConfig{Name: testStoreName},
		},
	}

	customSchema := &schema.Schema{
		Version: "custom-1.0",
		PageTypes: []schema.PageTypeDefinition{
			{Name: "custom", Directory: "custom/", Description: "Custom page"},
		},
	}

	eng, err := NewFromConfig(cfg, WithSchema(customSchema))
	if err != nil {
		t.Fatalf("NewFromConfig() error: %v", err)
	}

	if eng.Schema().Version != "custom-1.0" {
		t.Errorf("Schema.Version = %q, want %q", eng.Schema().Version, "custom-1.0")
	}
}

func TestNewFromConfig_DefaultSchema(t *testing.T) {
	t.Parallel()

	testStoreName := "test-store-default-schema"
	extension.RegisterStore(testStoreName, func(_ map[string]any) (extension.Store, error) {
		return &mockStore{name: testStoreName}, nil
	})

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Store:   config.ExtensionConfig{Name: testStoreName},
		},
	}

	eng, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig() error: %v", err)
	}

	if eng.Schema().Version != "1.0" {
		t.Errorf("Default Schema.Version = %q, want %q", eng.Schema().Version, "1.0")
	}
}

func TestCollect_NoSources(t *testing.T) {
	t.Parallel()

	testStoreName := "test-store-no-sources"
	extension.RegisterStore(testStoreName, func(_ map[string]any) (extension.Store, error) {
		return &mockStore{name: testStoreName}, nil
	})

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Store:   config.ExtensionConfig{Name: testStoreName},
		},
	}

	eng, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig() error: %v", err)
	}

	err = eng.Collect(context.Background())
	if err == nil {
		t.Fatal("Collect() expected error for no sources")
	}
	if err.Error() != "no sources configured" {
		t.Errorf("Collect() error = %q, want %q", err.Error(), "no sources configured")
	}
}

func TestCollect_NoProcessor(t *testing.T) {
	t.Parallel()

	testSourceName := "test-source-no-proc"
	testStoreName := "test-store-no-proc"

	extension.RegisterSource(testSourceName, func(_ map[string]any) (extension.Source, error) {
		return &mockSource{
			name: testSourceName,
			content: []model.RawContent{
				{Title: "Test", Content: "content"},
			},
		}, nil
	})
	extension.RegisterStore(testStoreName, func(_ map[string]any) (extension.Store, error) {
		return &mockStore{name: testStoreName}, nil
	})

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Sources: []config.SourceConfig{
				{Name: testSourceName},
			},
			Store: config.ExtensionConfig{Name: testStoreName},
		},
	}

	eng, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig() error: %v", err)
	}

	err = eng.Collect(context.Background())
	if err == nil {
		t.Fatal("Collect() expected error for no processor")
	}
}

func TestCollect_FullPipeline(t *testing.T) {
	t.Parallel()

	testSourceName := "test-source-full"
	testProcessorName := "test-processor-full"
	testStoreName := "test-store-full"

	rawContent := []model.RawContent{
		{Title: "Go Concurrency", Content: "goroutines are lightweight"},
	}
	knowledgeFiles := []model.KnowledgeFile{
		{
			Path: "topics/go-concurrency.md",
			Frontmatter: model.Frontmatter{
				Type:  model.PageTypeTopic,
				Title: "Go Concurrency",
			},
			Body: "# Go Concurrency",
		},
	}

	extension.RegisterSource(testSourceName, func(_ map[string]any) (extension.Source, error) {
		return &mockSource{name: testSourceName, content: rawContent}, nil
	})
	extension.RegisterProcessor(testProcessorName, func(_ map[string]any) (extension.Processor, error) {
		return &mockProcessor{name: testProcessorName, files: knowledgeFiles}, nil
	})

	ms := &mockStore{name: testStoreName}
	extension.RegisterStore(testStoreName, func(_ map[string]any) (extension.Store, error) {
		return ms, nil
	})

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Sources: []config.SourceConfig{
				{Name: testSourceName},
			},
			Processor: &config.ExtensionConfig{Name: testProcessorName},
			Store:     config.ExtensionConfig{Name: testStoreName},
		},
	}

	eng, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig() error: %v", err)
	}

	if err := eng.Collect(context.Background()); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if len(ms.written) != 1 {
		t.Fatalf("Store written %d files, want 1", len(ms.written))
	}
	if ms.written[0].Path != "topics/go-concurrency.md" {
		t.Errorf("Written file path = %q, want %q", ms.written[0].Path, "topics/go-concurrency.md")
	}
}

func TestCollect_WithFetchOptions(t *testing.T) {
	t.Parallel()

	testSourceName := "test-source-opts"
	testProcessorName := "test-processor-opts"
	testStoreName := "test-store-opts"

	var capturedOpts model.FetchOptions

	extension.RegisterSource(testSourceName, func(_ map[string]any) (extension.Source, error) {
		return &configCapturingSource{
			name:       testSourceName,
			capturePtr: &capturedOpts,
		}, nil
	})
	extension.RegisterProcessor(testProcessorName, func(_ map[string]any) (extension.Processor, error) {
		return &mockProcessor{
			name: testProcessorName,
			files: []model.KnowledgeFile{
				{Path: "inbox/test.md", Body: "test"},
			},
		}, nil
	})
	extension.RegisterStore(testStoreName, func(_ map[string]any) (extension.Store, error) {
		return &mockStore{name: testStoreName}, nil
	})

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Sources: []config.SourceConfig{
				{Name: testSourceName},
			},
			Processor: &config.ExtensionConfig{Name: testProcessorName},
			Store:     config.ExtensionConfig{Name: testStoreName},
		},
	}

	eng, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig() error: %v", err)
	}

	opts := CollectOptions{
		FetchOpts: model.FetchOptions{
			Config: map[string]any{
				"content": "test content",
				"title":   "test title",
			},
		},
	}

	if err := eng.Collect(context.Background(), opts); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if capturedOpts.Config["content"] != "test content" {
		t.Errorf("FetchOptions.Config[content] = %v, want %q", capturedOpts.Config["content"], "test content")
	}
}

func TestNew_WithConfigFile(t *testing.T) {
	t.Parallel()

	// 创建临时目录和配置文件
	dir := t.TempDir()
	configDir := filepath.Join(dir, ".synapse")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("create config dir: %v", err)
	}

	// 注册测试用的扩展点
	testStoreName := "test-store-config-file"
	extension.RegisterStore(testStoreName, func(cfg map[string]any) (extension.Store, error) {
		return &mockStore{name: testStoreName}, nil
	})

	configContent := `synapse:
  version: "1.0"
  store:
    name: test-store-config-file
    config:
      path: /tmp/test
`
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	eng, err := New(configPath)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	if eng.Store() == nil {
		t.Fatal("Engine.Store() should not be nil")
	}
	if eng.Schema() == nil {
		t.Fatal("Engine.Schema() should not be nil")
	}
}

func TestSearch_NoIndexer(t *testing.T) {
	t.Parallel()

	testStoreName := "test-store-search"
	extension.RegisterStore(testStoreName, func(_ map[string]any) (extension.Store, error) {
		return &mockStore{name: testStoreName}, nil
	})

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Store:   config.ExtensionConfig{Name: testStoreName},
		},
	}

	eng, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig() error: %v", err)
	}

	_, err = eng.Search(context.Background(), "test", model.SearchOptions{})
	if err == nil {
		t.Fatal("Search() expected error for no indexer")
	}
}

func TestAudit_NoAuditor(t *testing.T) {
	t.Parallel()

	testStoreName := "test-store-audit"
	extension.RegisterStore(testStoreName, func(_ map[string]any) (extension.Store, error) {
		return &mockStore{name: testStoreName}, nil
	})

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Store:   config.ExtensionConfig{Name: testStoreName},
		},
	}

	eng, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig() error: %v", err)
	}

	_, err = eng.Audit(context.Background())
	if err == nil {
		t.Fatal("Audit() expected error for no auditor")
	}
}

func TestPublish_NoConsumers(t *testing.T) {
	t.Parallel()

	testStoreName := "test-store-publish"
	extension.RegisterStore(testStoreName, func(_ map[string]any) (extension.Store, error) {
		return &mockStore{name: testStoreName}, nil
	})

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Store:   config.ExtensionConfig{Name: testStoreName},
		},
	}

	eng, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("NewFromConfig() error: %v", err)
	}

	err = eng.Publish(context.Background())
	if err == nil {
		t.Fatal("Publish() expected error for no consumers")
	}
}

// --- 测试辅助 ---

// configCapturingSource 捕获 FetchOptions 的 Source
type configCapturingSource struct {
	name       string
	capturePtr *model.FetchOptions
}

func (s *configCapturingSource) Name() string { return s.name }
func (s *configCapturingSource) Fetch(_ context.Context, opts model.FetchOptions) ([]model.RawContent, error) {
	*s.capturePtr = opts
	return []model.RawContent{
		{Title: "Captured", Content: "captured content"},
	}, nil
}
