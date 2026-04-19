package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGlobalConfigPath(t *testing.T) {
	t.Parallel()

	path, err := GlobalConfigPath()
	if err != nil {
		t.Fatalf("GlobalConfigPath() error: %v", err)
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".synapse", "config.yaml")
	if path != expected {
		t.Errorf("GlobalConfigPath() = %q, want %q", path, expected)
	}
}

func TestGlobalConfigDirPath(t *testing.T) {
	t.Parallel()

	path, err := GlobalConfigDirPath()
	if err != nil {
		t.Fatalf("GlobalConfigDirPath() error: %v", err)
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".synapse")
	if path != expected {
		t.Errorf("GlobalConfigDirPath() = %q, want %q", path, expected)
	}
}

func TestDefaultTemplate(t *testing.T) {
	t.Parallel()

	tmpl := DefaultTemplate()

	// 验证模板包含关键配置项
	checks := []struct {
		name    string
		content string
	}{
		{"version field", "version:"},
		{"store section", "store:"},
		{"sources section", "sources:"},
		{"processor section", "processor:"},
		{"local-store option", "local-store"},
		{"github-store option", "github-store"},
		{"GITHUB_OWNER placeholder", "${GITHUB_OWNER}"},
		{"GITHUB_TOKEN placeholder", "${GITHUB_TOKEN}"},
		{"synapse check hint", "synapse check"},
	}

	for _, c := range checks {
		if !strings.Contains(tmpl, c.content) {
			t.Errorf("DefaultTemplate() should contain %s (%q)", c.name, c.content)
		}
	}

	// 验证模板非空且有足够长度
	if len(tmpl) < 200 {
		t.Errorf("DefaultTemplate() length = %d, expected > 200", len(tmpl))
	}
}

// TestEnsureGlobalConfig_CreatesTemplate 使用临时目录测试配置文件创建
// 注意：不使用 t.Parallel()，因为需要通过 t.Setenv 修改 HOME 环境变量
func TestEnsureGlobalConfig_CreatesTemplate(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	cfgPath, created, err := EnsureGlobalConfig()
	if err != nil {
		t.Fatalf("EnsureGlobalConfig() error: %v", err)
	}

	if !created {
		t.Error("EnsureGlobalConfig() created = false, want true on first call")
	}

	// 验证文件存在
	if _, err := os.Stat(cfgPath); err != nil {
		t.Errorf("config file should exist at %s: %v", cfgPath, err)
	}

	// 验证文件内容包含模板内容
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "version:") {
		t.Error("config template should contain 'version:'")
	}
	if !strings.Contains(content, "store:") {
		t.Error("config template should contain 'store:'")
	}
}

// TestEnsureGlobalConfig_DoesNotOverwrite 验证不会覆盖已有配置文件
func TestEnsureGlobalConfig_DoesNotOverwrite(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// 第一次创建
	cfgPath, created, err := EnsureGlobalConfig()
	if err != nil {
		t.Fatalf("first EnsureGlobalConfig() error: %v", err)
	}
	if !created {
		t.Fatal("first call should create the config")
	}

	// 写入自定义内容
	customContent := "# custom config\nsynapse:\n  version: \"2.0\"\n  store:\n    name: local-store\n"
	if err := os.WriteFile(cfgPath, []byte(customContent), 0o644); err != nil {
		t.Fatalf("write custom config: %v", err)
	}

	// 第二次调用不应覆盖
	_, created, err = EnsureGlobalConfig()
	if err != nil {
		t.Fatalf("second EnsureGlobalConfig() error: %v", err)
	}
	if created {
		t.Error("second call should not recreate the config")
	}

	// 验证内容未被覆盖
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(data), "version: \"2.0\"") {
		t.Error("config should still contain custom content")
	}
}

// TestLoadGlobal 验证加载全局配置
func TestLoadGlobal(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// 创建有效的配置文件
	cfgDir := filepath.Join(tmpHome, ".synapse")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("create config dir: %v", err)
	}

	validConfig := `synapse:
  version: "1.0"
  store:
    name: "local-store"
    config:
      path: "/tmp/test-knowhub"
`
	cfgPath := filepath.Join(cfgDir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(validConfig), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadGlobal()
	if err != nil {
		t.Fatalf("LoadGlobal() error: %v", err)
	}

	if cfg.Synapse.Version != "1.0" {
		t.Errorf("Version = %q, want %q", cfg.Synapse.Version, "1.0")
	}
	if cfg.Synapse.Store.Name != "local-store" {
		t.Errorf("Store.Name = %q, want %q", cfg.Synapse.Store.Name, "local-store")
	}
}

// TestLoadGlobal_MissingFile 验证缺失配置文件时返回错误
func TestLoadGlobal_MissingFile(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	_, err := LoadGlobal()
	if err == nil {
		t.Fatal("LoadGlobal() expected error for missing config file")
	}
}

// TestEnsureGlobalConfig_PathStructure 验证配置文件路径结构
func TestEnsureGlobalConfig_PathStructure(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	cfgPath, _, err := EnsureGlobalConfig()
	if err != nil {
		t.Fatalf("EnsureGlobalConfig() error: %v", err)
	}

	expectedDir := filepath.Join(tmpHome, ".synapse")
	expectedPath := filepath.Join(expectedDir, "config.yaml")
	if cfgPath != expectedPath {
		t.Errorf("config path = %q, want %q", cfgPath, expectedPath)
	}

	// 验证目录创建
	info, err := os.Stat(expectedDir)
	if err != nil {
		t.Fatalf("config dir should exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("config path should be a directory")
	}
}
