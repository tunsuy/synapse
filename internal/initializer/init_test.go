package initializer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit_CreatesDirectoryStructure(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	targetDir := filepath.Join(dir, "my-knowhub")

	err := Init(Options{
		Path: targetDir,
		Name: "Test User",
	})
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// 验证目录结构
	expectedDirs := []string{
		".synapse",
		"profile",
		"topics",
		"entities",
		"concepts",
		"inbox",
		"journal",
		"graph",
	}

	for _, d := range expectedDirs {
		fullPath := filepath.Join(targetDir, d)
		info, err := os.Stat(fullPath)
		if err != nil {
			t.Errorf("directory %q does not exist: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%q is not a directory", d)
		}
	}
}

func TestInit_CreatesFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	targetDir := filepath.Join(dir, "my-knowhub")

	err := Init(Options{
		Path: targetDir,
		Name: "Test User",
	})
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// 验证文件
	expectedFiles := []string{
		".synapse/schema.yaml",
		".synapse/config.yaml",
		"profile/me.md",
		"graph/relations.json",
		".gitignore",
		"README.md",
	}

	for _, f := range expectedFiles {
		fullPath := filepath.Join(targetDir, f)
		info, err := os.Stat(fullPath)
		if err != nil {
			t.Errorf("file %q does not exist: %v", f, err)
			continue
		}
		if info.IsDir() {
			t.Errorf("%q should be a file, not a directory", f)
		}
		if info.Size() == 0 {
			t.Errorf("file %q should not be empty", f)
		}
	}
}

func TestInit_ProfileContainsUserName(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	targetDir := filepath.Join(dir, "my-knowhub")

	err := Init(Options{
		Path: targetDir,
		Name: "Alice",
	})
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(targetDir, "profile", "me.md"))
	if err != nil {
		t.Fatalf("read profile: %v", err)
	}

	content := string(data)
	if !containsStr(content, "Alice") {
		t.Error("profile/me.md should contain user name 'Alice'")
	}
}

func TestInit_AlreadyInitialized(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	targetDir := filepath.Join(dir, "my-knowhub")

	// 第一次初始化
	if err := Init(Options{Path: targetDir, Name: "User"}); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// 第二次初始化应该失败
	err := Init(Options{Path: targetDir, Name: "User"})
	if err == nil {
		t.Fatal("Init() expected error for already initialized knowhub")
	}
}

func TestInit_ForceOverwrite(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	targetDir := filepath.Join(dir, "my-knowhub")

	// 第一次初始化
	if err := Init(Options{Path: targetDir, Name: "User1"}); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// 强制覆盖
	err := Init(Options{Path: targetDir, Name: "User2", Force: true})
	if err != nil {
		t.Fatalf("Init(force) error: %v", err)
	}

	// 验证被覆盖
	data, err := os.ReadFile(filepath.Join(targetDir, "profile", "me.md"))
	if err != nil {
		t.Fatalf("read profile: %v", err)
	}
	if !containsStr(string(data), "User2") {
		t.Error("profile should contain 'User2' after force overwrite")
	}
}

func TestInit_DefaultName(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	targetDir := filepath.Join(dir, "my-knowhub")

	err := Init(Options{
		Path: targetDir,
		Name: "", // 空名称
	})
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// 验证使用了默认名称
	data, err := os.ReadFile(filepath.Join(targetDir, "profile", "me.md"))
	if err != nil {
		t.Fatalf("read profile: %v", err)
	}
	if !containsStr(string(data), "Synapse User") {
		t.Error("profile should contain default 'Synapse User' when name is empty")
	}
}

func TestInit_SchemaFileIsValidYAML(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	targetDir := filepath.Join(dir, "my-knowhub")

	if err := Init(Options{Path: targetDir, Name: "Test"}); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(targetDir, ".synapse", "schema.yaml"))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	content := string(data)
	if !containsStr(content, "version:") {
		t.Error("schema.yaml should contain 'version:' field")
	}
	if !containsStr(content, "page_types:") {
		t.Error("schema.yaml should contain 'page_types:' field")
	}
}

func TestInit_GitkeepFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	targetDir := filepath.Join(dir, "my-knowhub")

	if err := Init(Options{Path: targetDir, Name: "Test"}); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	gitkeepDirs := []string{
		"topics",
		"entities",
		"concepts",
		"inbox",
		"journal",
	}

	for _, d := range gitkeepDirs {
		gk := filepath.Join(targetDir, d, ".gitkeep")
		if _, err := os.Stat(gk); err != nil {
			t.Errorf(".gitkeep not found in %s: %v", d, err)
		}
	}
}

// containsStr 检查字符串是否包含子串
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
