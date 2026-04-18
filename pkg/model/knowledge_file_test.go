package model

import (
	"strings"
	"testing"
	"time"
)

func TestValidPageTypes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		pageType PageType
		want     bool
	}{
		{"profile is valid", PageTypeProfile, true},
		{"topic is valid", PageTypeTopic, true},
		{"entity is valid", PageTypeEntity, true},
		{"concept is valid", PageTypeConcept, true},
		{"inbox is valid", PageTypeInbox, true},
		{"journal is valid", PageTypeJournal, true},
		{"graph is valid", PageTypeGraph, true},
		{"unknown is invalid", PageType("unknown"), false},
		{"empty is invalid", PageType(""), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ValidPageTypes[tc.pageType]
			if got != tc.want {
				t.Errorf("ValidPageTypes[%q] = %v, want %v", tc.pageType, got, tc.want)
			}
		})
	}
}

func TestKnowledgeFile_Marshal(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 4, 18, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		kf       KnowledgeFile
		contains []string
	}{
		{
			name: "basic frontmatter",
			kf: KnowledgeFile{
				Path: "topics/golang.md",
				Frontmatter: Frontmatter{
					Type:    PageTypeTopic,
					Title:   "Go Programming",
					Created: now,
					Updated: now,
				},
				Body: "# Go Programming\n\nGo is awesome.",
			},
			contains: []string{
				"---\n",
				"type: topic\n",
				`title: "Go Programming"`,
				"created: 2025-04-18T12:00:00Z\n",
				"updated: 2025-04-18T12:00:00Z\n",
				"---\n\n# Go Programming\n\nGo is awesome.",
			},
		},
		{
			name: "with tags and links",
			kf: KnowledgeFile{
				Path: "topics/concurrency.md",
				Frontmatter: Frontmatter{
					Type:    PageTypeTopic,
					Title:   "Concurrency",
					Created: now,
					Updated: now,
					Tags:    []string{"golang", "concurrency"},
					Links:   []string{"goroutine", "channel"},
				},
				Body: "Concurrency content.",
			},
			contains: []string{
				"tags:\n  - golang\n  - concurrency\n",
				"links:\n  - goroutine\n  - channel\n",
			},
		},
		{
			name: "with optional fields",
			kf: KnowledgeFile{
				Path: "entities/codebuddy.md",
				Frontmatter: Frontmatter{
					Type:       PageTypeEntity,
					Title:      "CodeBuddy",
					Created:    now,
					Updated:    now,
					Source:     "codebuddy",
					Confidence: 0.95,
					Aliases:    []string{"CB", "Tencent CodeBuddy"},
					Category:   "tool",
					Status:     "active",
				},
				Body: "CodeBuddy is an AI coding assistant.",
			},
			contains: []string{
				"source: codebuddy\n",
				"confidence: 0.95\n",
				"aliases:\n  - CB\n  - Tencent CodeBuddy\n",
				"category: tool\n",
				"status: active\n",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.kf.Marshal()
			if err != nil {
				t.Fatalf("Marshal() error: %v", err)
			}
			result := string(data)
			for _, want := range tc.contains {
				if !strings.Contains(result, want) {
					t.Errorf("Marshal() result missing %q\nGot:\n%s", want, result)
				}
			}
		})
	}
}

func TestKnowledgeFile_Marshal_StartsAndEndsCorrectly(t *testing.T) {
	t.Parallel()

	kf := KnowledgeFile{
		Path: "test.md",
		Frontmatter: Frontmatter{
			Type:    PageTypeTopic,
			Title:   "Test",
			Created: time.Now(),
			Updated: time.Now(),
		},
		Body: "Test body",
	}

	data, err := kf.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}

	result := string(data)
	if !strings.HasPrefix(result, "---\n") {
		t.Error("Marshal() should start with ---")
	}
	if !strings.Contains(result, "\n---\n\n") {
		t.Error("Marshal() should contain closing --- followed by blank line")
	}
	if !strings.HasSuffix(result, "Test body") {
		t.Error("Marshal() should end with body content")
	}
}
