package model

// FetchOptions Source.Fetch 的选项
type FetchOptions struct {
	// Since 只获取此时间之后的内容
	Since string `json:"since,omitempty" yaml:"since,omitempty"`

	// Limit 最大获取条数
	Limit int `json:"limit,omitempty" yaml:"limit,omitempty"`

	// Config 数据源特有的配置参数
	Config map[string]any `json:"config,omitempty" yaml:"config,omitempty"`
}

// ProcessResult Processor.Process 的结果
type ProcessResult struct {
	// Files 生成的知识文件列表
	Files []KnowledgeFile `json:"files" yaml:"files"`

	// LinksCreated 新建的链接数
	LinksCreated int `json:"links_created" yaml:"links_created"`

	// ProfileUpdated 用户画像是否被更新
	ProfileUpdated bool `json:"profile_updated" yaml:"profile_updated"`
}

// ListOptions Store.List 的选项
type ListOptions struct {
	// Recursive 是否递归列出子目录
	Recursive bool `json:"recursive,omitempty" yaml:"recursive,omitempty"`

	// Types 按页面类型过滤
	Types []PageType `json:"types,omitempty" yaml:"types,omitempty"`

	// Limit 最大返回数量
	Limit int `json:"limit,omitempty" yaml:"limit,omitempty"`

	// Offset 偏移量
	Offset int `json:"offset,omitempty" yaml:"offset,omitempty"`
}

// ConsumeOptions Consumer.Consume 的选项
type ConsumeOptions struct {
	// OutputDir 输出目录
	OutputDir string `json:"output_dir,omitempty" yaml:"output_dir,omitempty"`

	// Config 消费端特有的配置参数
	Config map[string]any `json:"config,omitempty" yaml:"config,omitempty"`
}
