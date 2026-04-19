package local

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tunsuy/synapse/pkg/extension"
	"github.com/tunsuy/synapse/pkg/model"
)

func setupTestStore(t *testing.T) (*LocalStore, string) {
	t.Helper()
	dir := t.TempDir()
	store := &LocalStore{basePath: dir}
	return store, dir
}

func TestNew(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		config  map[string]any
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  map[string]any{"path": "/tmp/test"},
			wantErr: false,
		},
		{
			name:    "missing path",
			config:  map[string]any{},
			wantErr: true,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := New(tc.config)
			if (err != nil) != tc.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestLocalStore_Name(t *testing.T) {
	t.Parallel()
	store := &LocalStore{}
	if got := store.Name(); got != "local-store" {
		t.Errorf("Name() = %q, want %q", got, "local-store")
	}
}

func TestLocalStore_WriteAndRead(t *testing.T) {
	t.Parallel()
	store, _ := setupTestStore(t)
	ctx := context.Background()

	now := time.Date(2025, 4, 18, 12, 0, 0, 0, time.UTC)
	kf := model.KnowledgeFile{
		Path: "topics/golang.md",
		Frontmatter: model.Frontmatter{
			Type:    model.PageTypeTopic,
			Title:   "Go Programming",
			Created: now,
			Updated: now,
			Tags:    []string{"golang"},
		},
		Body: "# Go Programming\n\nGo is awesome.",
	}

	// Write
	if err := store.Write(ctx, kf); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	// Read
	got, err := store.Read(ctx, "topics/golang.md")
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if got.Frontmatter.Title != "Go Programming" {
		t.Errorf("Title = %q, want %q", got.Frontmatter.Title, "Go Programming")
	}
	if got.Frontmatter.Type != model.PageTypeTopic {
		t.Errorf("Type = %q, want %q", got.Frontmatter.Type, model.PageTypeTopic)
	}
}

func TestLocalStore_Delete(t *testing.T) {
	t.Parallel()
	store, _ := setupTestStore(t)
	ctx := context.Background()

	kf := model.KnowledgeFile{
		Path: "topics/delete-me.md",
		Frontmatter: model.Frontmatter{
			Type:    model.PageTypeTopic,
			Title:   "Delete Me",
			Created: time.Now(),
			Updated: time.Now(),
		},
		Body: "Will be deleted.",
	}

	if err := store.Write(ctx, kf); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	// 验证文件存在
	exists, err := store.Exists(ctx, "topics/delete-me.md")
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if !exists {
		t.Fatal("file should exist after write")
	}

	// Delete
	if err := store.Delete(ctx, "topics/delete-me.md"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	// 验证文件已删除
	exists, err = store.Exists(ctx, "topics/delete-me.md")
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if exists {
		t.Fatal("file should not exist after delete")
	}
}

func TestLocalStore_Exists(t *testing.T) {
	t.Parallel()
	store, dir := setupTestStore(t)
	ctx := context.Background()

	// 不存在
	exists, err := store.Exists(ctx, "nonexistent.md")
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if exists {
		t.Fatal("nonexistent file should return false")
	}

	// 创建文件
	if err := os.MkdirAll(filepath.Join(dir, "topics"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "topics/test.md"), []byte("test"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	// 存在
	exists, err = store.Exists(ctx, "topics/test.md")
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if !exists {
		t.Fatal("existing file should return true")
	}
}

func TestLocalStore_List(t *testing.T) {
	t.Parallel()
	store, _ := setupTestStore(t)
	ctx := context.Background()

	now := time.Now()
	// 写入多个文件
	files := []model.KnowledgeFile{
		{
			Path:        "topics/golang.md",
			Frontmatter: model.Frontmatter{Type: model.PageTypeTopic, Title: "Go", Created: now, Updated: now},
			Body:        "Go content",
		},
		{
			Path:        "topics/rust.md",
			Frontmatter: model.Frontmatter{Type: model.PageTypeTopic, Title: "Rust", Created: now, Updated: now},
			Body:        "Rust content",
		},
		{
			Path:        "entities/codebuddy.md",
			Frontmatter: model.Frontmatter{Type: model.PageTypeEntity, Title: "CodeBuddy", Created: now, Updated: now},
			Body:        "CodeBuddy content",
		},
	}

	for _, f := range files {
		if err := store.Write(ctx, f); err != nil {
			t.Fatalf("Write(%s) error: %v", f.Path, err)
		}
	}

	// List topics
	infos, err := store.List(ctx, "topics", model.ListOptions{})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(infos) != 2 {
		t.Errorf("List(topics) returned %d files, want 2", len(infos))
	}

	// List with limit
	infos, err = store.List(ctx, "topics", model.ListOptions{Limit: 1})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(infos) != 1 {
		t.Errorf("List(topics, limit=1) returned %d files, want 1", len(infos))
	}
}

func TestLocalStore_List_Recursive(t *testing.T) {
	t.Parallel()
	store, _ := setupTestStore(t)
	ctx := context.Background()

	now := time.Now()
	files := []model.KnowledgeFile{
		{
			Path:        "topics/golang.md",
			Frontmatter: model.Frontmatter{Type: model.PageTypeTopic, Title: "Go", Created: now, Updated: now},
			Body:        "Go",
		},
		{
			Path:        "entities/codebuddy.md",
			Frontmatter: model.Frontmatter{Type: model.PageTypeEntity, Title: "CB", Created: now, Updated: now},
			Body:        "CB",
		},
	}
	for _, f := range files {
		if err := store.Write(ctx, f); err != nil {
			t.Fatalf("Write() error: %v", err)
		}
	}

	// Recursive list from root
	infos, err := store.List(ctx, "", model.ListOptions{Recursive: true})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(infos) < 2 {
		t.Errorf("List(recursive) returned %d files, want >= 2", len(infos))
	}
}

func TestLocalStore_Read_NotFound(t *testing.T) {
	t.Parallel()
	store, _ := setupTestStore(t)
	ctx := context.Background()

	_, err := store.Read(ctx, "nonexistent.md")
	if err == nil {
		t.Fatal("Read() expected error for missing file")
	}
}

func TestLocalStore_Init(t *testing.T) {
	t.Parallel()
	store, dir := setupTestStore(t)
	ctx := context.Background()

	schemaData := []byte("version: \"1.0\"\npage_types:\n  - name: topic\n    directory: topics/\n    description: 主题知识\n")

	err := store.Init(ctx, extension.InitOptions{
		Name:       "Test User",
		SchemaData: schemaData,
	})
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// 验证目录结构
	expectedDirs := []string{
		".synapse", "profile", "topics", "entities",
		"concepts", "inbox", "journal", "graph",
	}
	for _, d := range expectedDirs {
		fullPath := filepath.Join(dir, d)
		info, err := os.Stat(fullPath)
		if err != nil {
			t.Errorf("directory %q does not exist: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%q is not a directory", d)
		}
	}

	// 验证 schema.yaml
	schemaPath := filepath.Join(dir, ".synapse", "schema.yaml")
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	if !strings.Contains(string(data), "version:") {
		t.Error("schema.yaml should contain 'version:'")
	}

	// 验证 profile/me.md
	profilePath := filepath.Join(dir, "profile", "me.md")
	data, err = os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("read profile: %v", err)
	}
	if !strings.Contains(string(data), "Test User") {
		t.Error("profile/me.md should contain user name")
	}

	// 验证 .gitignore
	gitignorePath := filepath.Join(dir, ".gitignore")
	if _, err := os.Stat(gitignorePath); err != nil {
		t.Error(".gitignore should exist")
	}

	// 验证 README.md
	readmePath := filepath.Join(dir, "README.md")
	data, err = os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}
	if !strings.Contains(string(data), "Synapse") {
		t.Error("README.md should mention Synapse")
	}

	// 验证 graph/relations.json
	relationsPath := filepath.Join(dir, "graph", "relations.json")
	data, err = os.ReadFile(relationsPath)
	if err != nil {
		t.Fatalf("read relations: %v", err)
	}
	if !strings.Contains(string(data), "\"nodes\"") {
		t.Error("relations.json should contain 'nodes' field")
	}

	// 验证 .gitkeep 文件
	gitkeepDirs := []string{"topics", "entities", "concepts", "inbox", "journal"}
	for _, d := range gitkeepDirs {
		gk := filepath.Join(dir, d, ".gitkeep")
		if _, err := os.Stat(gk); err != nil {
			t.Errorf(".gitkeep not found in %s: %v", d, err)
		}
	}
}

func TestLocalStore_Init_DefaultName(t *testing.T) {
	t.Parallel()
	store, dir := setupTestStore(t)
	ctx := context.Background()

	err := store.Init(ctx, extension.InitOptions{
		Name:       "",
		SchemaData: []byte("version: \"1.0\"\npage_types:\n  - name: topic\n    directory: topics/\n    description: test\n"),
	})
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// 空名称应使用默认名称 "Synapse User"
	data, err := os.ReadFile(filepath.Join(dir, "profile", "me.md"))
	if err != nil {
		t.Fatalf("read profile: %v", err)
	}
	if !strings.Contains(string(data), "Synapse User") {
		t.Error("profile should contain default 'Synapse User' when name is empty")
	}
}

func TestLocalStore_Init_NoSchema(t *testing.T) {
	t.Parallel()
	store, dir := setupTestStore(t)
	ctx := context.Background()

	// 不传 SchemaData
	err := store.Init(ctx, extension.InitOptions{
		Name: "User",
	})
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// schema.yaml 不应被创建
	schemaPath := filepath.Join(dir, ".synapse", "schema.yaml")
	if _, err := os.Stat(schemaPath); !os.IsNotExist(err) {
		t.Error("schema.yaml should not exist when SchemaData is empty")
	}
}

func TestLocalStore_Initialized(t *testing.T) {
	t.Parallel()
	store, dir := setupTestStore(t)
	ctx := context.Background()

	// 未初始化时应返回 false
	initialized, err := store.Initialized(ctx)
	if err != nil {
		t.Fatalf("Initialized() error: %v", err)
	}
	if initialized {
		t.Error("Initialized() = true, want false before init")
	}

	// 初始化后应返回 true
	err = store.Init(ctx, extension.InitOptions{
		Name:       "User",
		SchemaData: []byte("version: \"1.0\"\npage_types:\n  - name: topic\n    directory: topics/\n    description: test\n"),
	})
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	initialized, err = store.Initialized(ctx)
	if err != nil {
		t.Fatalf("Initialized() error after init: %v", err)
	}
	if !initialized {
		t.Error("Initialized() = false, want true after init")
	}

	// 手动删除 schema.yaml 后应返回 false
	schemaPath := filepath.Join(dir, ".synapse", "schema.yaml")
	if err := os.Remove(schemaPath); err != nil {
		t.Fatalf("remove schema: %v", err)
	}

	initialized, err = store.Initialized(ctx)
	if err != nil {
		t.Fatalf("Initialized() error after remove: %v", err)
	}
	if initialized {
		t.Error("Initialized() = true, want false after schema removal")
	}
}

func TestParseKnowledgeFile(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		path      string
		data      string
		wantTitle string
		wantType  model.PageType
		wantBody  string
	}{
		{
			name: "with frontmatter",
			path: "topics/golang.md",
			data: `---
type: topic
title: "Go Programming"
created: 2025-04-18T12:00:00Z
updated: 2025-04-18T12:00:00Z
---

# Go Programming`,
			wantTitle: "Go Programming",
			wantType:  model.PageTypeTopic,
			wantBody:  "# Go Programming",
		},
		{
			name:      "without frontmatter",
			path:      "topics/test.md",
			data:      "# Just some content",
			wantTitle: "test",
			wantType:  "",
			wantBody:  "# Just some content",
		},
		{
			name: "with unquoted title",
			path: "topics/unquoted.md",
			data: `---
type: entity
title: CodeBuddy
---

Entity content`,
			wantTitle: "CodeBuddy",
			wantType:  model.PageTypeEntity,
			wantBody:  "Entity content",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			kf, err := parseKnowledgeFile(tc.path, []byte(tc.data))
			if err != nil {
				t.Fatalf("parseKnowledgeFile() error: %v", err)
			}
			if kf.Frontmatter.Title != tc.wantTitle {
				t.Errorf("Title = %q, want %q", kf.Frontmatter.Title, tc.wantTitle)
			}
			if kf.Frontmatter.Type != tc.wantType {
				t.Errorf("Type = %q, want %q", kf.Frontmatter.Type, tc.wantType)
			}
			if kf.Body != tc.wantBody {
				t.Errorf("Body = %q, want %q", kf.Body, tc.wantBody)
			}
		})
	}
}
