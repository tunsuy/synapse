package link

import (
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		content string
		want    []Link
	}{
		{
			name:    "simple link",
			content: "This is about [[golang]].",
			want: []Link{
				{Target: "golang", Display: "", Position: 14},
			},
		},
		{
			name:    "link with display text",
			content: "Check [[golang|Go Language]] for details.",
			want: []Link{
				{Target: "golang", Display: "Go Language", Position: 6},
			},
		},
		{
			name:    "multiple links",
			content: "Learn [[golang]] and [[concurrency]] together.",
			want: []Link{
				{Target: "golang", Position: 6},
				{Target: "concurrency", Position: 21},
			},
		},
		{
			name:    "no links",
			content: "This has no wiki links.",
			want:    []Link{},
		},
		{
			name:    "empty content",
			content: "",
			want:    []Link{},
		},
		{
			name:    "link with spaces in target",
			content: "See [[ goroutine ]] usage.",
			want: []Link{
				{Target: "goroutine", Position: 4},
			},
		},
		{
			name:    "link with spaces in display",
			content: "See [[goroutine| Go Routines ]].",
			want: []Link{
				{Target: "goroutine", Display: "Go Routines", Position: 4},
			},
		},
		{
			name:    "adjacent links",
			content: "[[a]][[b]]",
			want: []Link{
				{Target: "a", Position: 0},
				{Target: "b", Position: 5},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Parse(tc.content)

			if len(got) == 0 && len(tc.want) == 0 {
				return
			}

			if len(got) != len(tc.want) {
				t.Fatalf("Parse(%q) returned %d links, want %d", tc.content, len(got), len(tc.want))
			}

			for i := range got {
				if got[i].Target != tc.want[i].Target {
					t.Errorf("link[%d].Target = %q, want %q", i, got[i].Target, tc.want[i].Target)
				}
				if got[i].Display != tc.want[i].Display {
					t.Errorf("link[%d].Display = %q, want %q", i, got[i].Display, tc.want[i].Display)
				}
				if got[i].Position != tc.want[i].Position {
					t.Errorf("link[%d].Position = %d, want %d", i, got[i].Position, tc.want[i].Position)
				}
			}
		})
	}
}

func TestLink_DisplayText(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		link Link
		want string
	}{
		{
			name: "with display text",
			link: Link{Target: "golang", Display: "Go Language"},
			want: "Go Language",
		},
		{
			name: "without display text",
			link: Link{Target: "golang"},
			want: "golang",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.link.DisplayText(); got != tc.want {
				t.Errorf("DisplayText() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestExtractTargets(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "unique targets",
			content: "Learn [[golang]], [[rust]], and [[python]].",
			want:    []string{"golang", "rust", "python"},
		},
		{
			name:    "duplicate targets",
			content: "[[golang]] is great. [[golang]] rocks.",
			want:    []string{"golang"},
		},
		{
			name:    "no links",
			content: "No links here.",
			want:    []string{},
		},
		{
			name:    "mixed with display text",
			content: "[[go|Go Language]] and [[go|Go标准库]] are the same.",
			want:    []string{"go"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ExtractTargets(tc.content)

			if len(got) == 0 && len(tc.want) == 0 {
				return
			}

			if len(got) != len(tc.want) {
				t.Fatalf("ExtractTargets() returned %d targets, want %d", len(got), len(tc.want))
			}

			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("targets[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestReplaceLinks(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		content string
		format  func(target, display string) string
		want    string
	}{
		{
			name:    "replace with markdown links",
			content: "See [[golang]] for details.",
			format: func(target, display string) string {
				if display == "" {
					display = target
				}
				return "[" + display + "](/topics/" + target + ")"
			},
			want: "See [golang](/topics/golang) for details.",
		},
		{
			name:    "replace with display text",
			content: "Use [[golang|Go Language]].",
			format: func(target, display string) string {
				if display == "" {
					display = target
				}
				return "[" + display + "](/topics/" + target + ")"
			},
			want: "Use [Go Language](/topics/golang).",
		},
		{
			name:    "no links to replace",
			content: "No links here.",
			format:  func(target, display string) string { return target },
			want:    "No links here.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ReplaceLinks(tc.content, tc.format)
			if got != tc.want {
				t.Errorf("ReplaceLinks() = %q, want %q", got, tc.want)
			}
		})
	}
}
