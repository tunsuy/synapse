// Package schema 负责加载和校验 Knowledge Schema 规范
// Schema 是所有扩展点遵循的规范，定义知识库的"行为契约"
package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Schema 知识库行为契约（.synapse/schema.yaml）
// 它定义页面模板、工作流规则、质量标准，所有扩展点共同遵守
type Schema struct {
	// Version Schema 版本
	Version string `yaml:"version"`

	// PageTypes 知识页面类型定义
	PageTypes []PageTypeDefinition `yaml:"page_types"`

	// Frontmatter 标准字段定义
	Frontmatter FrontmatterSpec `yaml:"frontmatter"`

	// LinkFormat 双向链接格式
	LinkFormat string `yaml:"link_format"`

	// Operations 支持的操作列表
	Operations []string `yaml:"operations"`

	// Quality 质量标准
	Quality QualitySpec `yaml:"quality,omitempty"`
}

// PageTypeDefinition 页面类型定义
type PageTypeDefinition struct {
	// Name 类型名称
	Name string `yaml:"name"`

	// Directory 对应的目录
	Directory string `yaml:"directory"`

	// Template 模板文件路径（可选）
	Template string `yaml:"template,omitempty"`

	// Description 类型描述
	Description string `yaml:"description"`

	// Emoji 目录展示用的 emoji 图标（可选，用于 README 生成等场景）
	Emoji string `yaml:"emoji,omitempty"`

	// Example 示例文件名（可选，用于 README 表格展示）
	Example string `yaml:"example,omitempty"`
}

// FrontmatterSpec Frontmatter 字段规范
type FrontmatterSpec struct {
	// Required 必须的字段
	Required []string `yaml:"required"`

	// Optional 可选的字段
	Optional []string `yaml:"optional"`
}

// QualitySpec 质量标准
type QualitySpec struct {
	// MinConfidence 最小置信度
	MinConfidence float64 `yaml:"min_confidence,omitempty"`

	// MaxStaledays 最大过时天数
	MaxStaleDays int `yaml:"max_stale_days,omitempty"`

	// RequireTags 是否要求所有页面都有标签
	RequireTags bool `yaml:"require_tags,omitempty"`
}

// Load 从文件加载 Schema
func Load(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read schema file %s: %w", path, err)
	}

	var s Schema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse schema file %s: %w", path, err)
	}

	if err := s.validate(); err != nil {
		return nil, fmt.Errorf("validate schema: %w", err)
	}

	return &s, nil
}

// validate 校验 Schema 的合法性
func (s *Schema) validate() error {
	if s.Version == "" {
		return fmt.Errorf("schema version is required")
	}
	if len(s.PageTypes) == 0 {
		return fmt.Errorf("at least one page type is required")
	}
	for _, pt := range s.PageTypes {
		if pt.Name == "" {
			return fmt.Errorf("page type name is required")
		}
		if pt.Directory == "" {
			return fmt.Errorf("page type %q directory is required", pt.Name)
		}
	}
	return nil
}

// Default 返回默认的 Schema 定义
func Default() *Schema {
	return &Schema{
		Version: "1.0",
		PageTypes: []PageTypeDefinition{
			{Name: "profile", Directory: "profile/", Description: "User profile", Emoji: "👤", Example: "me.md"},
			{Name: "topic", Directory: "topics/", Description: "Topic knowledge (deep content organized by topic)", Emoji: "📚", Example: "go-concurrency.md"},
			{Name: "entity", Directory: "entities/", Description: "Entity pages (people, tools, projects, organizations)", Emoji: "🏷️", Example: "docker.md"},
			{Name: "concept", Directory: "concepts/", Description: "Concept pages (technical concepts, methodologies, theories)", Emoji: "💡", Example: "cap-theorem.md"},
			{Name: "inbox", Directory: "inbox/", Description: "Incremental buffer (new knowledge pending organization)", Emoji: "📥", Example: ""},
			{Name: "journal", Directory: "journal/", Description: "Timeline (knowledge activities recorded chronologically)", Emoji: "📅", Example: "2025-01-15.md"},
			{Name: "graph", Directory: "graph/", Description: "Knowledge graph (auto-generated from [[wiki-links]])", Emoji: "🔗", Example: "relations.json"},
		},
		Frontmatter: FrontmatterSpec{
			Required: []string{"type", "title", "created", "updated"},
			Optional: []string{"tags", "links", "source", "confidence", "aliases", "category", "status"},
		},
		LinkFormat: "[[page-id]]",
		Operations: []string{"capture", "compile", "retrieve", "audit"},
		Quality: QualitySpec{
			MinConfidence: 0.5,
			MaxStaleDays:  90,
			RequireTags:   false,
		},
	}
}

// MarshalYAML 序列化 Schema 为 YAML 格式
func (s *Schema) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(s)
}
