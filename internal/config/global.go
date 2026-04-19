// Package config 负责加载和管理 Synapse 配置
// global.go 处理全局配置文件（~/.synapse/config.yaml）的创建和加载
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// GlobalConfigDir 全局配置目录名
	GlobalConfigDir = ".synapse"

	// GlobalConfigFile 全局配置文件名
	GlobalConfigFile = "config.yaml"
)

// GlobalConfigPath 返回全局配置文件的完整路径
func GlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, GlobalConfigDir, GlobalConfigFile), nil
}

// GlobalConfigDirPath 返回全局配置目录的完整路径
func GlobalConfigDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, GlobalConfigDir), nil
}

// EnsureGlobalConfig 确保全局配置文件存在
// 如果 ~/.synapse/config.yaml 不存在，则创建配置模板
// 如果已存在，不做任何修改
// 返回 (配置文件路径, 是否新创建, error)
func EnsureGlobalConfig() (string, bool, error) {
	cfgPath, err := GlobalConfigPath()
	if err != nil {
		return "", false, err
	}

	// 检查是否已存在
	if _, err := os.Stat(cfgPath); err == nil {
		return cfgPath, false, nil
	}

	// 创建目录
	cfgDir := filepath.Dir(cfgPath)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return "", false, fmt.Errorf("create config dir %s: %w", cfgDir, err)
	}

	// 写入配置模板
	template := DefaultTemplate()
	if err := os.WriteFile(cfgPath, []byte(template), 0o644); err != nil {
		return "", false, fmt.Errorf("write config template %s: %w", cfgPath, err)
	}

	return cfgPath, true, nil
}

// LoadGlobal 加载全局配置文件
func LoadGlobal() (*Config, error) {
	cfgPath, err := GlobalConfigPath()
	if err != nil {
		return nil, err
	}
	return Load(cfgPath)
}

// DefaultTemplate 返回默认的配置文件模板内容（带注释说明）
func DefaultTemplate() string {
	return `# Synapse Configuration — 扩展点注册中心
# 全局配置文件，声明各扩展点使用哪个实现
#
# 文档: https://github.com/tunsuy/synapse/blob/main/docs/roadmap.md
#
# 使用前请根据你的需求修改以下配置：
# 1. 选择合适的 store 后端（local-store 或 github-store）
# 2. 配置对应的参数
# 3. 运行 synapse check 验证配置

synapse:
  # 配置版本（请勿修改）
  version: "1.0"

  # ──────────────────────────────────
  # 数据源（Source）— 原始内容从哪来
  # ──────────────────────────────────
  sources:
    - name: "skill-source"
      enabled: true
      # config: {}

  # ──────────────────────────────────
  # 处理引擎（Processor）— 原始内容 → 结构化知识
  # ──────────────────────────────────
  processor:
    name: "skill-processor"
    # config:
    #   default_confidence: 0.7

  # ──────────────────────────────────
  # 存储底座（Store）— 知识文件的持久化
  # 请选择一个存储后端并填写对应配置
  # ──────────────────────────────────

  ## 方案 A：本地文件系统存储（适合个人本地使用）
  # store:
  #   name: "local-store"
  #   config:
  #     path: "~/knowhub"        # 知识库本地路径

  ## 方案 B：GitHub 仓库存储（适合云端同步）
  store:
    name: "github-store"
    config:
      owner: "${GITHUB_OWNER}"   # GitHub 用户名（必填，支持环境变量）
      repo: "${GITHUB_REPO}"     # GitHub 仓库名（必填，支持环境变量）
      token: "${GITHUB_TOKEN}"   # GitHub Personal Access Token（必填，支持环境变量）
      branch: "main"             # 分支名（默认 main）

  # ──────────────────────────────────
  # 检索引擎（Indexer）— 知识库检索（可选）
  # ──────────────────────────────────
  # indexer:
  #   name: "bm25-indexer"

  # ──────────────────────────────────
  # 消费端（Consumer）— 知识输出为各种形式（可选）
  # ──────────────────────────────────
  # consumers:
  #   - name: "hugo-consumer"
  #     config:
  #       output_dir: "./public"

  # ──────────────────────────────────
  # 审计器（Auditor）— 知识库健康检查（可选）
  # ──────────────────────────────────
  # auditor:
  #   name: "default-auditor"
`
}
