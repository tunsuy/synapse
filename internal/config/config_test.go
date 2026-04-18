package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	cfg := Default("/tmp/test-knowhub")

	if cfg.Synapse.Version != "1.0" {
		t.Errorf("Version = %q, want %q", cfg.Synapse.Version, "1.0")
	}
	if cfg.Synapse.Store.Name != "local-store" {
		t.Errorf("Store.Name = %q, want %q", cfg.Synapse.Store.Name, "local-store")
	}

	path, ok := cfg.Synapse.Store.Config["path"].(string)
	if !ok || path != "/tmp/test-knowhub" {
		t.Errorf("Store.Config[\"path\"] = %q, want %q", path, "/tmp/test-knowhub")
	}

	if len(cfg.Synapse.Sources) != 1 {
		t.Errorf("Sources count = %d, want 1", len(cfg.Synapse.Sources))
	}
	if cfg.Synapse.Sources[0].Name != "skill-source" {
		t.Errorf("Sources[0].Name = %q, want %q", cfg.Synapse.Sources[0].Name, "skill-source")
	}

	if cfg.Synapse.Processor == nil {
		t.Fatal("Processor should not be nil")
	}
	if cfg.Synapse.Processor.Name != "skill-processor" {
		t.Errorf("Processor.Name = %q, want %q", cfg.Synapse.Processor.Name, "skill-processor")
	}
}

func TestSourceConfig_IsEnabled(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		enabled *bool
		want    bool
	}{
		{"nil defaults to true", nil, true},
		{"explicitly true", boolPtr(true), true},
		{"explicitly false", boolPtr(false), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sc := SourceConfig{Enabled: tc.enabled}
			if got := sc.IsEnabled(); got != tc.want {
				t.Errorf("IsEnabled() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLoad_ValidFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `synapse:
  version: "1.0"
  store:
    name: local-store
    config:
      path: /tmp/test
  sources:
    - name: test-source
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Synapse.Version != "1.0" {
		t.Errorf("Version = %q, want %q", cfg.Synapse.Version, "1.0")
	}
	if cfg.Synapse.Store.Name != "local-store" {
		t.Errorf("Store.Name = %q, want %q", cfg.Synapse.Store.Name, "local-store")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	t.Parallel()

	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("Load() expected error for missing file")
	}
}

func TestLoad_MissingVersion(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `synapse:
  store:
    name: local-store
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() expected error for missing version")
	}
}

func TestLoad_MissingStoreName(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `synapse:
  version: "1.0"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() expected error for missing store name")
	}
}

func TestMarshalYAML(t *testing.T) {
	t.Parallel()

	cfg := Default("/tmp/test")
	data, err := cfg.MarshalYAML()
	if err != nil {
		t.Fatalf("MarshalYAML() error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("MarshalYAML() returned empty data")
	}
}

func boolPtr(b bool) *bool {
	return &b
}
