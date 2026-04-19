package skill

import (
	"context"
	"strings"
	"testing"

	"github.com/tunsuy/synapse/pkg/model"
)

func TestNew(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		config         map[string]any
		wantConfidence float64
	}{
		{
			name:           "default config",
			config:         nil,
			wantConfidence: 0.7,
		},
		{
			name:           "custom confidence",
			config:         map[string]any{"default_confidence": 0.9},
			wantConfidence: 0.9,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			proc, err := New(tc.config)
			if err != nil {
				t.Fatalf("New() error: %v", err)
			}
			sp := proc.(*SkillProcessor)
			if sp.defaultConfidence != tc.wantConfidence {
				t.Errorf("defaultConfidence = %v, want %v", sp.defaultConfidence, tc.wantConfidence)
			}
		})
	}
}

func TestSkillProcessor_Name(t *testing.T) {
	t.Parallel()

	proc := &SkillProcessor{}
	if got := proc.Name(); got != "skill-processor" {
		t.Errorf("Name() = %q, want %q", got, "skill-processor")
	}
}

func TestSkillProcessor_Process_WithTopics(t *testing.T) {
	t.Parallel()

	proc, _ := New(nil)
	ctx := context.Background()

	raw := []model.RawContent{
		{
			Title:           "Go Concurrency",
			Content:         "Go's concurrency model uses goroutines and channels.",
			Source:          "codebuddy",
			SuggestedTopics: []string{"golang", "concurrency"},
		},
	}

	files, err := proc.Process(ctx, raw)
	if err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("Process() returned %d files, want 2", len(files))
	}

	// 验证第一个文件（golang 主题）
	if files[0].Path != "topics/golang.md" {
		t.Errorf("files[0].Path = %q, want %q", files[0].Path, "topics/golang.md")
	}
	if files[0].Frontmatter.Type != model.PageTypeTopic {
		t.Errorf("files[0].Type = %q, want %q", files[0].Frontmatter.Type, model.PageTypeTopic)
	}
	if files[0].Frontmatter.Title != "golang" {
		t.Errorf("files[0].Title = %q, want %q", files[0].Frontmatter.Title, "golang")
	}
	if files[0].Frontmatter.Source != "codebuddy" {
		t.Errorf("files[0].Source = %q, want %q", files[0].Frontmatter.Source, "codebuddy")
	}
	if files[0].Frontmatter.Status != "active" {
		t.Errorf("files[0].Status = %q, want %q", files[0].Frontmatter.Status, "active")
	}

	// 验证 Body 包含关键内容
	if !strings.Contains(files[0].Body, "# golang") {
		t.Error("files[0].Body should contain heading")
	}
	if !strings.Contains(files[0].Body, "Go's concurrency") {
		t.Error("files[0].Body should contain original content")
	}
}

func TestSkillProcessor_Process_WithEntities(t *testing.T) {
	t.Parallel()

	proc, _ := New(nil)
	ctx := context.Background()

	raw := []model.RawContent{
		{
			Title:             "CodeBuddy Overview",
			Content:           "CodeBuddy is an AI coding tool developed by Tencent.",
			Source:            "codebuddy",
			SuggestedEntities: []string{"CodeBuddy"},
		},
	}

	files, err := proc.Process(ctx, raw)
	if err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Process() returned %d files, want 1", len(files))
	}

	if files[0].Path != "entities/codebuddy.md" {
		t.Errorf("files[0].Path = %q, want %q", files[0].Path, "entities/codebuddy.md")
	}
	if files[0].Frontmatter.Type != model.PageTypeEntity {
		t.Errorf("files[0].Type = %q, want %q", files[0].Frontmatter.Type, model.PageTypeEntity)
	}
	// CodeBuddy 应该被识别为 tool
	if files[0].Frontmatter.Category != "tool" {
		t.Errorf("files[0].Category = %q, want %q", files[0].Frontmatter.Category, "tool")
	}
}

func TestSkillProcessor_Process_WithConcepts(t *testing.T) {
	t.Parallel()

	proc, _ := New(nil)
	ctx := context.Background()

	raw := []model.RawContent{
		{
			Title:             "CSP Model",
			Content:           "CSP (Communicating Sequential Processes) is a formal language.",
			Source:            "codebuddy",
			SuggestedConcepts: []string{"CSP"},
		},
	}

	files, err := proc.Process(ctx, raw)
	if err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Process() returned %d files, want 1", len(files))
	}

	if files[0].Path != "concepts/csp.md" {
		t.Errorf("files[0].Path = %q, want %q", files[0].Path, "concepts/csp.md")
	}
	if files[0].Frontmatter.Type != model.PageTypeConcept {
		t.Errorf("files[0].Type = %q, want %q", files[0].Frontmatter.Type, model.PageTypeConcept)
	}
}

func TestSkillProcessor_Process_Inbox(t *testing.T) {
	t.Parallel()

	proc, _ := New(nil)
	ctx := context.Background()

	raw := []model.RawContent{
		{
			Title:   "Random Note",
			Content: "Some unclassified knowledge content.",
			Source:  "codebuddy",
			// 没有任何建议分类
		},
	}

	files, err := proc.Process(ctx, raw)
	if err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Process() returned %d files, want 1", len(files))
	}

	if !strings.HasPrefix(files[0].Path, "inbox/") {
		t.Errorf("files[0].Path = %q, should start with 'inbox/'", files[0].Path)
	}
	if files[0].Frontmatter.Type != model.PageTypeInbox {
		t.Errorf("files[0].Type = %q, want %q", files[0].Frontmatter.Type, model.PageTypeInbox)
	}
	if files[0].Frontmatter.Status != "draft" {
		t.Errorf("files[0].Status = %q, want %q", files[0].Frontmatter.Status, "draft")
	}
}

func TestSkillProcessor_Process_EmptyContent(t *testing.T) {
	t.Parallel()

	proc, _ := New(nil)
	ctx := context.Background()

	raw := []model.RawContent{
		{Content: ""},
	}

	files, err := proc.Process(ctx, raw)
	if err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Process() returned %d files for empty content, want 0", len(files))
	}
}

func TestSkillProcessor_Process_Multiple(t *testing.T) {
	t.Parallel()

	proc, _ := New(nil)
	ctx := context.Background()

	raw := []model.RawContent{
		{
			Title:           "Go Basics",
			Content:         "Go is a statically typed language.",
			Source:          "codebuddy",
			SuggestedTopics: []string{"golang"},
		},
		{
			Title:           "Rust Basics",
			Content:         "Rust is a systems programming language.",
			Source:          "codebuddy",
			SuggestedTopics: []string{"rust"},
		},
	}

	files, err := proc.Process(ctx, raw)
	if err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("Process() returned %d files, want 2", len(files))
	}
}

func TestSkillProcessor_Process_WithKeyPoints(t *testing.T) {
	t.Parallel()

	proc, _ := New(nil)
	ctx := context.Background()

	raw := []model.RawContent{
		{
			Title:   "Note with points",
			Content: "Content here.",
			Source:  "codebuddy",
			KeyPoints: []string{
				"Point 1",
				"Point 2",
			},
		},
	}

	files, err := proc.Process(ctx, raw)
	if err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Process() returned %d files, want 1", len(files))
	}

	if !strings.Contains(files[0].Body, "Point 1") {
		t.Error("Body should contain key point 'Point 1'")
	}
	if !strings.Contains(files[0].Body, "Point 2") {
		t.Error("Body should contain key point 'Point 2'")
	}
}

func TestSkillProcessor_Process_CrossLinks(t *testing.T) {
	t.Parallel()

	proc, _ := New(nil)
	ctx := context.Background()

	raw := []model.RawContent{
		{
			Title:             "Complex Content",
			Content:           "Discussion about Go and goroutines as a concept.",
			Source:            "codebuddy",
			SuggestedTopics:   []string{"golang"},
			SuggestedEntities: []string{"Go"},
			SuggestedConcepts: []string{"goroutine"},
		},
	}

	files, err := proc.Process(ctx, raw)
	if err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	// 应该生成 3 个文件：1 topic + 1 entity + 1 concept
	if len(files) != 3 {
		t.Fatalf("Process() returned %d files, want 3", len(files))
	}

	// 验证 golang topic 的 links 包含其他项
	golangFile := files[0]
	if golangFile.Path != "topics/golang.md" {
		t.Errorf("first file should be topics/golang.md, got %q", golangFile.Path)
	}

	// golang topic 应该链接到 Go 和 goroutine
	if len(golangFile.Frontmatter.Links) < 1 {
		t.Errorf("golang topic should have cross links, got %d", len(golangFile.Frontmatter.Links))
	}
}

func TestToSlug(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "golang", "golang"},
		{"with spaces", "Go Programming", "go-programming"},
		{"with special chars", "C++ & Rust", "c-and-rust"},
		{"with slashes", "path/to/file", "path-to-file"},
		{"uppercase", "GOLANG", "golang"},
		{"with parens", "React (library)", "react-library"},
		{"leading/trailing spaces", "  hello  ", "hello"},
		{"empty", "", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := toSlug(tc.input)
			if got != tc.want {
				t.Errorf("toSlug(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestGuessCategory(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		entity  string
		content string
		want    string
	}{
		{"tool", "CodeBuddy", "AI coding tool", "tool"},
		{"person", "Rob Pike", "person who created Go", "person"},
		{"organization", "Google", "company that develops Go", "organization"},
		{"project", "Kubernetes", "container orchestration project", "project"},
		{"unknown", "Something", "just some text", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := guessCategory(tc.entity, tc.content)
			if got != tc.want {
				t.Errorf("guessCategory(%q, ...) = %q, want %q", tc.entity, got, tc.want)
			}
		})
	}
}

func TestBuildTags(t *testing.T) {
	t.Parallel()

	rc := model.RawContent{Source: "codebuddy"}
	tags := buildTags(rc, "golang")

	if len(tags) < 1 {
		t.Fatal("buildTags() should return at least 1 tag")
	}

	found := make(map[string]bool)
	for _, tag := range tags {
		found[tag] = true
	}

	if !found["golang"] {
		t.Error("tags should contain 'golang'")
	}
	if !found["codebuddy"] {
		t.Error("tags should contain 'codebuddy'")
	}
}

func TestCollectLinks(t *testing.T) {
	t.Parallel()

	rc := model.RawContent{
		SuggestedTopics:   []string{"golang", "concurrency"},
		SuggestedEntities: []string{"Go"},
		SuggestedConcepts: []string{"CSP"},
	}

	links := collectLinks(rc, "golang")

	// golang 应该被排除
	for _, link := range links {
		if strings.ToLower(link) == "golang" {
			t.Error("collectLinks() should exclude the primary item")
		}
	}

	if len(links) < 2 {
		t.Errorf("collectLinks() returned %d links, want >= 2", len(links))
	}
}
