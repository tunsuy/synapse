package model

import (
	"fmt"
	"strings"
	"time"
)

// PageType 知识页面类型
type PageType string

const (
	// PageTypeProfile 用户画像
	PageTypeProfile PageType = "profile"
	// PageTypeTopic 主题知识
	PageTypeTopic PageType = "topic"
	// PageTypeEntity 实体页（人物/工具/项目/组织）
	PageTypeEntity PageType = "entity"
	// PageTypeConcept 概念页（技术概念/方法论/理论）
	PageTypeConcept PageType = "concept"
	// PageTypeInbox 增量缓冲区
	PageTypeInbox PageType = "inbox"
	// PageTypeJournal 时间线
	PageTypeJournal PageType = "journal"
	// PageTypeGraph 知识关联图谱
	PageTypeGraph PageType = "graph"
)

// ValidPageTypes 所有合法的页面类型集合
var ValidPageTypes = map[PageType]bool{
	PageTypeProfile: true,
	PageTypeTopic:   true,
	PageTypeEntity:  true,
	PageTypeConcept: true,
	PageTypeInbox:   true,
	PageTypeJournal: true,
	PageTypeGraph:   true,
}

// Frontmatter 知识文件的元数据头
type Frontmatter struct {
	// Type 页面类型（必须）
	Type PageType `json:"type" yaml:"type"`

	// Title 标题（必须）
	Title string `json:"title" yaml:"title"`

	// Created 创建时间（必须）
	Created time.Time `json:"created" yaml:"created"`

	// Updated 最后更新时间（必须）
	Updated time.Time `json:"updated" yaml:"updated"`

	// Tags 标签列表
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`

	// Links 双向链接目标列表（如 ["golang", "concurrency", "goroutine"]）
	Links []string `json:"links,omitempty" yaml:"links,omitempty"`

	// Source 数据来源标识
	Source string `json:"source,omitempty" yaml:"source,omitempty"`

	// Confidence 知识置信度（0.0-1.0）
	Confidence float64 `json:"confidence,omitempty" yaml:"confidence,omitempty"`

	// Aliases 页面别名（用于双向链接匹配）
	Aliases []string `json:"aliases,omitempty" yaml:"aliases,omitempty"`

	// Category 分类（实体页面适用，如 "person", "tool", "project", "organization"）
	Category string `json:"category,omitempty" yaml:"category,omitempty"`

	// Status 页面状态（如 "draft", "active", "archived"）
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
}

// KnowledgeFile 是 Processor 输出的结构化知识文件
// 它符合 Knowledge Schema 规范，由 Frontmatter 元数据和 Markdown 正文组成
type KnowledgeFile struct {
	// Path 文件在知识库中的相对路径（如 "topics/golang/concurrency.md"）
	Path string `json:"path" yaml:"path"`

	// Frontmatter 元数据头
	Frontmatter Frontmatter `json:"frontmatter" yaml:"frontmatter"`

	// Body Markdown 正文内容（包含 [[双向链接]]）
	Body string `json:"body" yaml:"body"`
}

// FileInfo 知识文件的简要信息（用于列表展示）
type FileInfo struct {
	// Path 文件相对路径
	Path string `json:"path" yaml:"path"`

	// Title 标题
	Title string `json:"title" yaml:"title"`

	// Type 页面类型
	Type PageType `json:"type" yaml:"type"`

	// Updated 最后更新时间
	Updated time.Time `json:"updated" yaml:"updated"`

	// Size 文件大小（字节）
	Size int64 `json:"size" yaml:"size"`
}

// Marshal 将 KnowledgeFile 序列化为 Markdown 格式（含 YAML Frontmatter）
func (kf *KnowledgeFile) Marshal() ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("type: %s\n", kf.Frontmatter.Type))
	sb.WriteString(fmt.Sprintf("title: %q\n", kf.Frontmatter.Title))
	sb.WriteString(fmt.Sprintf("created: %s\n", kf.Frontmatter.Created.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("updated: %s\n", kf.Frontmatter.Updated.Format(time.RFC3339)))

	if len(kf.Frontmatter.Tags) > 0 {
		sb.WriteString("tags:\n")
		for _, tag := range kf.Frontmatter.Tags {
			sb.WriteString(fmt.Sprintf("  - %s\n", tag))
		}
	}

	if len(kf.Frontmatter.Links) > 0 {
		sb.WriteString("links:\n")
		for _, link := range kf.Frontmatter.Links {
			sb.WriteString(fmt.Sprintf("  - %s\n", link))
		}
	}

	if kf.Frontmatter.Source != "" {
		sb.WriteString(fmt.Sprintf("source: %s\n", kf.Frontmatter.Source))
	}

	if kf.Frontmatter.Confidence > 0 {
		sb.WriteString(fmt.Sprintf("confidence: %.2f\n", kf.Frontmatter.Confidence))
	}

	if len(kf.Frontmatter.Aliases) > 0 {
		sb.WriteString("aliases:\n")
		for _, alias := range kf.Frontmatter.Aliases {
			sb.WriteString(fmt.Sprintf("  - %s\n", alias))
		}
	}

	if kf.Frontmatter.Category != "" {
		sb.WriteString(fmt.Sprintf("category: %s\n", kf.Frontmatter.Category))
	}

	if kf.Frontmatter.Status != "" {
		sb.WriteString(fmt.Sprintf("status: %s\n", kf.Frontmatter.Status))
	}

	sb.WriteString("---\n\n")
	sb.WriteString(kf.Body)

	return []byte(sb.String()), nil
}
