// Package model 定义 Synapse 核心领域模型
package model

import "time"

// RawContent 是 Source 输出的标准原始内容格式
// 它是连接数据源与处理引擎的统一中间表示，任何数据源都必须将原始数据转换为此格式
type RawContent struct {
	// Title 原始内容的标题
	Title string `json:"title" yaml:"title"`

	// Content 原始内容的正文
	Content string `json:"content" yaml:"content"`

	// Source 数据来源标识（如 "codebuddy", "chatgpt", "rss"）
	Source string `json:"source" yaml:"source"`

	// SourceVersion 数据源版本
	SourceVersion string `json:"source_version,omitempty" yaml:"source_version,omitempty"`

	// SessionID 会话标识（对话类数据源适用）
	SessionID string `json:"session_id,omitempty" yaml:"session_id,omitempty"`

	// Timestamp 原始内容的产生时间
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`

	// SuggestedTopics AI 数据源建议的主题分类
	SuggestedTopics []string `json:"suggested_topics,omitempty" yaml:"suggested_topics,omitempty"`

	// SuggestedEntities AI 数据源建议的实体（人物/工具/项目/组织）
	SuggestedEntities []string `json:"suggested_entities,omitempty" yaml:"suggested_entities,omitempty"`

	// SuggestedConcepts AI 数据源建议的概念（技术概念/方法论/理论）
	SuggestedConcepts []string `json:"suggested_concepts,omitempty" yaml:"suggested_concepts,omitempty"`

	// KeyPoints 关键知识点列表
	KeyPoints []string `json:"key_points,omitempty" yaml:"key_points,omitempty"`

	// ProfileUpdates 用户画像更新建议
	ProfileUpdates []ProfileUpdate `json:"profile_updates,omitempty" yaml:"profile_updates,omitempty"`

	// Metadata 扩展元数据（数据源特有的附加信息）
	Metadata map[string]any `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// ProfileUpdate 用户画像更新建议
type ProfileUpdate struct {
	// Type 更新类型（如 "skill", "interest", "preference"）
	Type string `json:"type" yaml:"type"`

	// Content 更新内容描述
	Content string `json:"content" yaml:"content"`
}
