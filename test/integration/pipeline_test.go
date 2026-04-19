// Package integration 集成测试
// 验证真实的 Source → Processor → Store 完整管道
package integration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tunsuy/synapse/internal/config"
	"github.com/tunsuy/synapse/internal/engine"
	"github.com/tunsuy/synapse/pkg/model"

	// 激活内置扩展点实现
	_ "github.com/tunsuy/synapse/internal/processor/skill"
	_ "github.com/tunsuy/synapse/internal/source/skill"
	_ "github.com/tunsuy/synapse/internal/store/local"
)

// TestFullPipeline_TopicCollection 测试完整管道：收集一条带主题分类的内容
func TestFullPipeline_TopicCollection(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	createDirs(t, dir)

	eng := newTestEngine(t, dir)
	ctx := context.Background()

	opts := engine.CollectOptions{
		FetchOpts: model.FetchOptions{
			Config: map[string]any{
				"content": "Go's concurrency model is based on goroutines and channels, inspired by CSP.",
				"title":   "Go Concurrency",
				"suggested_topics": "Go Concurrency",
				"key_points":       "goroutines are lightweight threads,channels enable communication",
			},
		},
	}

	if err := eng.Collect(ctx, opts); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// 验证文件已写入
	topicFile := filepath.Join(dir, "topics", "go-concurrency.md")
	data, err := os.ReadFile(topicFile)
	if err != nil {
		t.Fatalf("read topic file: %v", err)
	}

	content := string(data)

	// 验证 frontmatter
	if !strings.Contains(content, "type: topic") {
		t.Error("topic file should have type: topic")
	}
	if !strings.Contains(content, `title: "Go Concurrency"`) {
		t.Error("topic file should have correct title")
	}
	if !strings.Contains(content, "status: active") {
		t.Error("topic file should have status: active")
	}

	// 验证 body 中包含关键知识点
	if !strings.Contains(content, "goroutines are lightweight threads") {
		t.Error("topic file should contain key points")
	}
}

// TestFullPipeline_EntityCollection 测试完整管道：收集实体
func TestFullPipeline_EntityCollection(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	createDirs(t, dir)

	eng := newTestEngine(t, dir)
	ctx := context.Background()

	opts := engine.CollectOptions{
		FetchOpts: model.FetchOptions{
			Config: map[string]any{
				"content":              "Docker is a containerization platform that simplifies deployment.",
				"title":               "Docker Overview",
				"suggested_entities":   "Docker",
				"suggested_concepts":   "Containerization",
			},
		},
	}

	if err := eng.Collect(ctx, opts); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// 验证实体文件
	entityFile := filepath.Join(dir, "entities", "docker.md")
	data, err := os.ReadFile(entityFile)
	if err != nil {
		t.Fatalf("read entity file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "type: entity") {
		t.Error("entity file should have type: entity")
	}
	if !strings.Contains(content, "category: tool") {
		t.Error("Docker entity should be categorized as tool (contains 'platform')")
	}

	// 验证概念文件
	conceptFile := filepath.Join(dir, "concepts", "containerization.md")
	data, err = os.ReadFile(conceptFile)
	if err != nil {
		t.Fatalf("read concept file: %v", err)
	}

	content = string(data)
	if !strings.Contains(content, "type: concept") {
		t.Error("concept file should have type: concept")
	}
}

// TestFullPipeline_InboxFallback 测试完整管道：无分类建议时放入 inbox
func TestFullPipeline_InboxFallback(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	createDirs(t, dir)

	eng := newTestEngine(t, dir)
	ctx := context.Background()

	opts := engine.CollectOptions{
		FetchOpts: model.FetchOptions{
			Config: map[string]any{
				"content": "Some random notes from today's conversation.",
				"title":   "Daily Notes",
			},
		},
	}

	if err := eng.Collect(ctx, opts); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// 验证 inbox 文件
	inboxFile := filepath.Join(dir, "inbox", "daily-notes.md")
	data, err := os.ReadFile(inboxFile)
	if err != nil {
		t.Fatalf("read inbox file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "type: inbox") {
		t.Error("inbox file should have type: inbox")
	}
	if !strings.Contains(content, "status: draft") {
		t.Error("inbox file should have status: draft")
	}
}

// TestFullPipeline_CrossLinks 测试完整管道：交叉链接
func TestFullPipeline_CrossLinks(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	createDirs(t, dir)

	eng := newTestEngine(t, dir)
	ctx := context.Background()

	opts := engine.CollectOptions{
		FetchOpts: model.FetchOptions{
			Config: map[string]any{
				"content":            "Kubernetes uses Docker containers for orchestration.",
				"title":              "K8s and Docker",
				"suggested_topics":   "Container Orchestration",
				"suggested_entities": "Kubernetes,Docker",
			},
		},
	}

	if err := eng.Collect(ctx, opts); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// 验证 Kubernetes 实体文件有指向 Docker 和 Container Orchestration 的链接
	k8sFile := filepath.Join(dir, "entities", "kubernetes.md")
	data, err := os.ReadFile(k8sFile)
	if err != nil {
		t.Fatalf("read k8s file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "[[Docker]]") {
		t.Error("Kubernetes file should contain [[Docker]] link")
	}
	if !strings.Contains(content, "[[Container Orchestration]]") {
		t.Error("Kubernetes file should contain [[Container Orchestration]] link")
	}

	// 验证 Docker 实体文件有指向 Kubernetes 和 Container Orchestration 的链接
	dockerFile := filepath.Join(dir, "entities", "docker.md")
	data, err = os.ReadFile(dockerFile)
	if err != nil {
		t.Fatalf("read docker file: %v", err)
	}

	content = string(data)
	if !strings.Contains(content, "[[Kubernetes]]") {
		t.Error("Docker file should contain [[Kubernetes]] link")
	}
}

// TestFullPipeline_StdinLike 测试完整管道：模拟从 stdin 传入内容
func TestFullPipeline_StdinLike(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	createDirs(t, dir)

	eng := newTestEngine(t, dir)
	ctx := context.Background()

	// 模拟一段较长的 stdin 内容（如一篇 markdown 笔记）
	longContent := `# Understanding Go Interfaces

Go interfaces are implicit - you don't need to declare that a type implements an interface.
This makes Go code more flexible and composable.

## Key Points

1. Interfaces are satisfied implicitly
2. Small interfaces are preferred (io.Reader, io.Writer)
3. Empty interface (interface{}) matches any type
4. Type assertions can check runtime types

## Example

` + "```go\ntype Reader interface {\n    Read(p []byte) (n int, err error)\n}\n```"

	opts := engine.CollectOptions{
		FetchOpts: model.FetchOptions{
			Config: map[string]any{
				"content":            longContent,
				"title":              "Understanding Go Interfaces",
				"suggested_topics":   "Go Interfaces",
				"suggested_concepts": "Interface,Duck Typing",
				"key_points":         "interfaces are implicit,small interfaces preferred,empty interface matches any",
			},
		},
	}

	if err := eng.Collect(ctx, opts); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// 验证主题文件
	topicFile := filepath.Join(dir, "topics", "go-interfaces.md")
	if _, err := os.Stat(topicFile); os.IsNotExist(err) {
		t.Fatal("topic file 'go-interfaces.md' should exist")
	}

	// 验证概念文件
	interfaceFile := filepath.Join(dir, "concepts", "interface.md")
	if _, err := os.Stat(interfaceFile); os.IsNotExist(err) {
		t.Fatal("concept file 'interface.md' should exist")
	}

	duckTypingFile := filepath.Join(dir, "concepts", "duck-typing.md")
	if _, err := os.Stat(duckTypingFile); os.IsNotExist(err) {
		t.Fatal("concept file 'duck-typing.md' should exist")
	}
}

// TestFullPipeline_ReadBackWritten 测试完整管道：写入后能正确读回
func TestFullPipeline_ReadBackWritten(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	createDirs(t, dir)

	eng := newTestEngine(t, dir)
	ctx := context.Background()

	opts := engine.CollectOptions{
		FetchOpts: model.FetchOptions{
			Config: map[string]any{
				"content":          "gRPC is a high-performance RPC framework by Google.",
				"title":            "gRPC Overview",
				"suggested_topics": "gRPC",
				"source":           "test-source",
			},
		},
	}

	if err := eng.Collect(ctx, opts); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// 通过 Store 读回文件
	store := eng.Store()
	kf, err := store.Read(ctx, "topics/grpc.md")
	if err != nil {
		t.Fatalf("Store.Read() error: %v", err)
	}

	if kf.Frontmatter.Title != "gRPC" {
		t.Errorf("Read back title = %q, want %q", kf.Frontmatter.Title, "gRPC")
	}
	if kf.Frontmatter.Type != model.PageTypeTopic {
		t.Errorf("Read back type = %q, want %q", kf.Frontmatter.Type, model.PageTypeTopic)
	}
	if kf.Frontmatter.Source != "test-source" {
		t.Errorf("Read back source = %q, want %q", kf.Frontmatter.Source, "test-source")
	}
}

// --- 辅助函数 ---

// newTestEngine 创建一个用于集成测试的 Engine 实例
func newTestEngine(t *testing.T, basePath string) *engine.Engine {
	t.Helper()

	cfg := &config.Config{
		Synapse: config.SynapseConfig{
			Version: "1.0",
			Sources: []config.SourceConfig{
				{Name: "skill-source"},
			},
			Processor: &config.ExtensionConfig{
				Name: "skill-processor",
			},
			Store: config.ExtensionConfig{
				Name: "local-store",
				Config: map[string]any{
					"path": basePath,
				},
			},
		},
	}

	eng, err := engine.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("create test engine: %v", err)
	}
	return eng
}

// createDirs 创建知识库目录结构
func createDirs(t *testing.T, basePath string) {
	t.Helper()

	dirs := []string{"profile", "topics", "entities", "concepts", "inbox", "journal", "graph"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(basePath, dir), 0o755); err != nil {
			t.Fatalf("create dir %s: %v", dir, err)
		}
	}
}
