// Package config 负责加载和管理 Synapse 配置
// 配置文件 .synapse/config.yaml 是扩展点的注册中心，声明各扩展点使用哪个实现
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config Synapse 配置文件（.synapse/config.yaml）
type Config struct {
	// Synapse 配置根节点
	Synapse SynapseConfig `yaml:"synapse"`
}

// SynapseConfig 核心配置
type SynapseConfig struct {
	// Version 配置版本
	Version string `yaml:"version"`

	// Sources 数据源配置列表（可同时启用多个）
	Sources []SourceConfig `yaml:"sources,omitempty"`

	// Processor 处理引擎配置（选一个）
	Processor *ExtensionConfig `yaml:"processor,omitempty"`

	// Store 存储底座配置（选一个）
	Store ExtensionConfig `yaml:"store"`

	// Indexer 检索引擎配置（可选）
	Indexer *ExtensionConfig `yaml:"indexer,omitempty"`

	// Consumers 消费端配置列表（可同时启用多个）
	Consumers []SourceConfig `yaml:"consumers,omitempty"`

	// Auditor 审计器配置（可选）
	Auditor *ExtensionConfig `yaml:"auditor,omitempty"`
}

// SourceConfig 带启用/禁用开关的扩展点配置
type SourceConfig struct {
	// Name 扩展点实现名称
	Name string `yaml:"name"`

	// Enabled 是否启用（默认 true）
	Enabled *bool `yaml:"enabled,omitempty"`

	// Plugin 外部插件可执行文件路径（如有）
	Plugin string `yaml:"plugin,omitempty"`

	// Config 该扩展点的自定义配置
	Config map[string]any `yaml:"config,omitempty"`
}

// IsEnabled 返回该配置项是否启用
func (sc SourceConfig) IsEnabled() bool {
	if sc.Enabled == nil {
		return true // 默认启用
	}
	return *sc.Enabled
}

// ExtensionConfig 扩展点配置
type ExtensionConfig struct {
	// Name 扩展点实现名称
	Name string `yaml:"name"`

	// Plugin 外部插件可执行文件路径（如有）
	Plugin string `yaml:"plugin,omitempty"`

	// Config 该扩展点的自定义配置
	Config map[string]any `yaml:"config,omitempty"`
}

// Load 从文件加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file %s: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

// validate 校验配置的合法性
func (c *Config) validate() error {
	if c.Synapse.Version == "" {
		return fmt.Errorf("synapse.version is required")
	}
	if c.Synapse.Store.Name == "" {
		return fmt.Errorf("synapse.store.name is required")
	}
	return nil
}

// Default 返回默认配置
func Default(knowhubPath string) *Config {
	return &Config{
		Synapse: SynapseConfig{
			Version: "1.0",
			Sources: []SourceConfig{
				{Name: "skill-source"},
			},
			Processor: &ExtensionConfig{
				Name: "skill-processor",
			},
			Store: ExtensionConfig{
				Name: "local-store",
				Config: map[string]any{
					"path": knowhubPath,
				},
			},
			Indexer: nil,
			Consumers: nil,
			Auditor:   nil,
		},
	}
}

// MarshalYAML 序列化配置为 YAML 格式
func (c *Config) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(c)
}
