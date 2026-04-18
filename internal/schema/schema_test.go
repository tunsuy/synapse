package schema

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	s := Default()

	if s.Version != "1.0" {
		t.Errorf("Version = %q, want %q", s.Version, "1.0")
	}

	if len(s.PageTypes) != 7 {
		t.Errorf("PageTypes count = %d, want 7", len(s.PageTypes))
	}

	// 验证所有 7 种页面类型
	expectedTypes := map[string]string{
		"profile": "profile/",
		"topic":   "topics/",
		"entity":  "entities/",
		"concept": "concepts/",
		"inbox":   "inbox/",
		"journal": "journal/",
		"graph":   "graph/",
	}

	for _, pt := range s.PageTypes {
		dir, ok := expectedTypes[pt.Name]
		if !ok {
			t.Errorf("unexpected page type: %q", pt.Name)
			continue
		}
		if pt.Directory != dir {
			t.Errorf("PageType %q directory = %q, want %q", pt.Name, pt.Directory, dir)
		}
		if pt.Description == "" {
			t.Errorf("PageType %q description should not be empty", pt.Name)
		}
	}

	// 验证 Frontmatter 规范
	if len(s.Frontmatter.Required) != 4 {
		t.Errorf("Required frontmatter fields = %d, want 4", len(s.Frontmatter.Required))
	}
	if s.LinkFormat != "[[page-id]]" {
		t.Errorf("LinkFormat = %q, want %q", s.LinkFormat, "[[page-id]]")
	}
	if len(s.Operations) != 4 {
		t.Errorf("Operations count = %d, want 4", len(s.Operations))
	}
}

func TestLoad_ValidFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "schema.yaml")

	content := `version: "1.0"
page_types:
  - name: topic
    directory: topics/
    description: Topic knowledge
frontmatter:
  required:
    - type
    - title
  optional:
    - tags
link_format: "[[page-id]]"
operations:
  - capture
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if s.Version != "1.0" {
		t.Errorf("Version = %q, want %q", s.Version, "1.0")
	}
	if len(s.PageTypes) != 1 {
		t.Errorf("PageTypes count = %d, want 1", len(s.PageTypes))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	t.Parallel()

	_, err := Load("/nonexistent/path/schema.yaml")
	if err == nil {
		t.Fatal("Load() expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "schema.yaml")

	if err := os.WriteFile(path, []byte("{{invalid yaml"), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() expected error for invalid YAML")
	}
}

func TestLoad_MissingVersion(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "schema.yaml")

	content := `page_types:
  - name: topic
    directory: topics/
    description: Topic
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() expected error for missing version")
	}
}

func TestLoad_MissingPageTypes(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "schema.yaml")

	content := `version: "1.0"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() expected error for missing page types")
	}
}

func TestLoad_MissingPageTypeDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "schema.yaml")

	content := `version: "1.0"
page_types:
  - name: topic
    description: Topic
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() expected error for missing directory")
	}
}

func TestMarshalYAML(t *testing.T) {
	t.Parallel()

	s := Default()
	data, err := s.MarshalYAML()
	if err != nil {
		t.Fatalf("MarshalYAML() error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("MarshalYAML() returned empty data")
	}
}
