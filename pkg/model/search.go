package model

// SearchResult 检索结果
type SearchResult struct {
	// Path 匹配的知识文件路径
	Path string `json:"path" yaml:"path"`

	// Title 文件标题
	Title string `json:"title" yaml:"title"`

	// Score 相关性评分（0.0-1.0）
	Score float64 `json:"score" yaml:"score"`

	// Snippet 匹配的文本片段
	Snippet string `json:"snippet,omitempty" yaml:"snippet,omitempty"`

	// Highlights 高亮匹配的关键词位置
	Highlights []string `json:"highlights,omitempty" yaml:"highlights,omitempty"`
}

// SearchOptions 检索选项
type SearchOptions struct {
	// Limit 最大返回数量
	Limit int `json:"limit,omitempty" yaml:"limit,omitempty"`

	// MinScore 最小相关性评分阈值
	MinScore float64 `json:"min_score,omitempty" yaml:"min_score,omitempty"`

	// Types 限定搜索的页面类型
	Types []PageType `json:"types,omitempty" yaml:"types,omitempty"`

	// Tags 限定搜索的标签
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// AuditReport 审计报告
type AuditReport struct {
	// Score 知识库健康评分（0-100）
	Score int `json:"score" yaml:"score"`

	// Issues 发现的问题列表
	Issues []AuditIssue `json:"issues" yaml:"issues"`

	// Stats 知识库统计信息
	Stats AuditStats `json:"stats" yaml:"stats"`
}

// AuditIssueSeverity 审计问题严重性
type AuditIssueSeverity string

const (
	// SeverityError 错误（必须修复）
	SeverityError AuditIssueSeverity = "error"
	// SeverityWarning 警告（建议修复）
	SeverityWarning AuditIssueSeverity = "warning"
	// SeverityInfo 信息（可忽略）
	SeverityInfo AuditIssueSeverity = "info"
)

// AuditIssue 审计发现的问题
type AuditIssue struct {
	// Type 问题类型（如 "broken_link", "orphan_page", "stale_content", "duplicate"）
	Type string `json:"type" yaml:"type"`

	// Severity 严重性
	Severity AuditIssueSeverity `json:"severity" yaml:"severity"`

	// Path 涉及的文件路径
	Path string `json:"path" yaml:"path"`

	// Message 问题描述
	Message string `json:"message" yaml:"message"`

	// Suggestion 修复建议
	Suggestion string `json:"suggestion,omitempty" yaml:"suggestion,omitempty"`
}

// AuditStats 知识库统计
type AuditStats struct {
	// TotalFiles 总文件数
	TotalFiles int `json:"total_files" yaml:"total_files"`

	// FilesByType 各类型文件数量
	FilesByType map[PageType]int `json:"files_by_type" yaml:"files_by_type"`

	// TotalLinks 总链接数
	TotalLinks int `json:"total_links" yaml:"total_links"`

	// BrokenLinks 断链数
	BrokenLinks int `json:"broken_links" yaml:"broken_links"`

	// OrphanPages 孤儿页面数
	OrphanPages int `json:"orphan_pages" yaml:"orphan_pages"`

	// TotalTags 总标签数
	TotalTags int `json:"total_tags" yaml:"total_tags"`
}

// FixResult 自动修复结果
type FixResult struct {
	// Fixed 已修复的问题数
	Fixed int `json:"fixed" yaml:"fixed"`

	// Skipped 跳过的问题数
	Skipped int `json:"skipped" yaml:"skipped"`

	// Details 修复详情
	Details []string `json:"details,omitempty" yaml:"details,omitempty"`
}
