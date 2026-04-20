package skill

import (
	"context"
	"testing"

	"github.com/tunsuy/synapse/pkg/model"
)

func TestNew(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		config     map[string]any
		wantSource string
	}{
		{
			name:       "default config",
			config:     nil,
			wantSource: "codebuddy",
		},
		{
			name:       "empty config",
			config:     map[string]any{},
			wantSource: "codebuddy",
		},
		{
			name:       "custom source",
			config:     map[string]any{"source": "claude"},
			wantSource: "claude",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			src, err := New(tc.config)
			if err != nil {
				t.Fatalf("New() error: %v", err)
			}
			ss := src.(*SkillSource)
			if ss.defaultSource != tc.wantSource {
				t.Errorf("defaultSource = %q, want %q", ss.defaultSource, tc.wantSource)
			}
		})
	}
}

func TestSkillSource_Name(t *testing.T) {
	t.Parallel()

	src := &SkillSource{}
	if got := src.Name(); got != "skill-source" {
		t.Errorf("Name() = %q, want %q", got, "skill-source")
	}
}

func TestSkillSource_Fetch(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		config     map[string]any
		opts       model.FetchOptions
		wantErr    bool
		wantLen    int
		wantTitle  string
		wantSource string
	}{
		{
			name:   "basic content",
			config: nil,
			opts: model.FetchOptions{
				Config: map[string]any{
					"content": "Go's concurrency model is based on goroutines and channels.",
				},
			},
			wantErr:    false,
			wantLen:    1,
			wantTitle:  "Go's concurrency model is based on goroutines and...",
			wantSource: "codebuddy",
		},
		{
			name:   "with explicit title",
			config: nil,
			opts: model.FetchOptions{
				Config: map[string]any{
					"content": "Some content here.",
					"title":   "Go Concurrency",
				},
			},
			wantErr:    false,
			wantLen:    1,
			wantTitle:  "Go Concurrency",
			wantSource: "codebuddy",
		},
		{
			name:   "with custom source",
			config: map[string]any{"source": "cursor"},
			opts: model.FetchOptions{
				Config: map[string]any{
					"content": "Some content.",
				},
			},
			wantErr:    false,
			wantLen:    1,
			wantSource: "cursor",
		},
		{
			name:   "missing content",
			config: nil,
			opts: model.FetchOptions{
				Config: map[string]any{},
			},
			wantErr: true,
		},
		{
			name:    "nil config in opts",
			config:  nil,
			opts:    model.FetchOptions{},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			src, err := New(tc.config)
			if err != nil {
				t.Fatalf("New() error: %v", err)
			}

			ctx := context.Background()
			result, err := src.Fetch(ctx, tc.opts)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Fetch() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}

			if len(result) != tc.wantLen {
				t.Fatalf("Fetch() returned %d items, want %d", len(result), tc.wantLen)
			}

			if tc.wantTitle != "" && result[0].Title != tc.wantTitle {
				t.Errorf("Title = %q, want %q", result[0].Title, tc.wantTitle)
			}
			if result[0].Source != tc.wantSource {
				t.Errorf("Source = %q, want %q", result[0].Source, tc.wantSource)
			}
		})
	}
}

func TestSkillSource_Fetch_WithSuggestions(t *testing.T) {
	t.Parallel()

	src, _ := New(nil)
	ctx := context.Background()

	result, err := src.Fetch(ctx, model.FetchOptions{
		Config: map[string]any{
			"content":            "Go concurrency with goroutines and channels.",
			"suggested_topics":   "golang, concurrency",
			"suggested_entities": "Go, goroutine",
			"suggested_concepts": "CSP, actor-model",
			"key_points":         "goroutines are lightweight, channels for communication",
			"session_id":         "sess-123",
		},
	})
	if err != nil {
		t.Fatalf("Fetch() error: %v", err)
	}

	rc := result[0]

	if rc.SessionID != "sess-123" {
		t.Errorf("SessionID = %q, want %q", rc.SessionID, "sess-123")
	}

	if len(rc.SuggestedTopics) != 2 {
		t.Fatalf("SuggestedTopics count = %d, want 2", len(rc.SuggestedTopics))
	}
	if rc.SuggestedTopics[0] != "golang" {
		t.Errorf("SuggestedTopics[0] = %q, want %q", rc.SuggestedTopics[0], "golang")
	}

	if len(rc.SuggestedEntities) != 2 {
		t.Fatalf("SuggestedEntities count = %d, want 2", len(rc.SuggestedEntities))
	}

	if len(rc.SuggestedConcepts) != 2 {
		t.Fatalf("SuggestedConcepts count = %d, want 2", len(rc.SuggestedConcepts))
	}

	if len(rc.KeyPoints) != 2 {
		t.Fatalf("KeyPoints count = %d, want 2", len(rc.KeyPoints))
	}
}

func TestSkillSource_Fetch_ExtraMetadata(t *testing.T) {
	t.Parallel()

	src, _ := New(nil)
	ctx := context.Background()

	result, err := src.Fetch(ctx, model.FetchOptions{
		Config: map[string]any{
			"content":    "Some content.",
			"custom_key": "custom_value",
			"priority":   "high",
		},
	})
	if err != nil {
		t.Fatalf("Fetch() error: %v", err)
	}

	rc := result[0]

	if rc.Metadata["custom_key"] != "custom_value" {
		t.Errorf("Metadata[custom_key] = %v, want %q", rc.Metadata["custom_key"], "custom_value")
	}
	if rc.Metadata["priority"] != "high" {
		t.Errorf("Metadata[priority] = %v, want %q", rc.Metadata["priority"], "high")
	}
	// 已知字段不应出现在 Metadata 中
	if _, ok := rc.Metadata["content"]; ok {
		t.Error("Metadata should not contain 'content'")
	}
}

func TestExtractTitle(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "short content",
			content: "Hello world",
			want:    "Hello world",
		},
		{
			name:    "long content",
			content: "This is a very long title that exceeds the fifty character limit for title extraction purposes",
			want:    "This is a very long title that exceeds the fifty " + "...",
		},
		{
			name:    "markdown heading",
			content: "# Go Programming\n\nGo is great.",
			want:    "Go Programming",
		},
		{
			name:    "empty content",
			content: "",
			want:    "Untitled",
		},
		{
			name:    "blank lines only",
			content: "\n\n\n",
			want:    "Untitled",
		},
		{
			name:    "content with leading blank lines",
			content: "\n\n  First real line",
			want:    "First real line",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := extractTitle(tc.content)
			if got != tc.want {
				t.Errorf("extractTitle() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSplitAndTrim(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "simple",
			input: "a, b, c",
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "no spaces",
			input: "a,b,c",
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "empty parts",
			input: "a,,b,,c",
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "single item",
			input: "golang",
			want:  []string{"golang"},
		},
		{
			name:  "empty string",
			input: "",
			want:  []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := splitAndTrim(tc.input)
			if len(got) != len(tc.want) {
				t.Fatalf("splitAndTrim(%q) returned %d items, want %d", tc.input, len(got), len(tc.want))
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("splitAndTrim(%q)[%d] = %q, want %q", tc.input, i, got[i], tc.want[i])
				}
			}
		})
	}
}
