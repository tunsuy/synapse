// Package initializer 实现 synapse init 命令的核心逻辑
// 负责初始化 knowhub 仓库结构，生成 schema.yaml、config.yaml 和各目录模板
package initializer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/tunsuy/synapse/internal/config"
	"github.com/tunsuy/synapse/internal/schema"
	"github.com/tunsuy/synapse/internal/store/tmpl"
)

// Options 初始化选项
type Options struct {
	// Path 知识库根目录路径
	Path string

	// Name 知识库名称（用于 profile 等）
	Name string

	// Force 是否覆盖已有文件
	Force bool
}

// Init 初始化 knowhub 仓库结构
func Init(opts Options) error {
	absPath, err := filepath.Abs(opts.Path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	// 检查是否已初始化
	synapseDir := filepath.Join(absPath, ".synapse")
	if !opts.Force {
		if _, err := os.Stat(synapseDir); err == nil {
			return fmt.Errorf("knowhub already initialized at %s (use --force to overwrite)", absPath)
		}
	}

	fmt.Printf("🧠 Initializing Synapse knowhub at %s\n\n", absPath)

	// 1. 创建目录结构
	if err := createDirectories(absPath); err != nil {
		return fmt.Errorf("create directories: %w", err)
	}

	// 2. 生成 .synapse/schema.yaml
	if err := writeSchema(absPath); err != nil {
		return fmt.Errorf("write schema: %w", err)
	}

	// 3. 生成 .synapse/config.yaml
	if err := writeConfig(absPath); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	// 4. 生成模板文件
	if err := writeTemplates(absPath, opts.Name); err != nil {
		return fmt.Errorf("write templates: %w", err)
	}

	// 5. 生成 .gitignore
	if err := writeGitignore(absPath); err != nil {
		return fmt.Errorf("write gitignore: %w", err)
	}

	// 6. 生成 README.md
	if err := writeReadme(absPath, opts.Name); err != nil {
		return fmt.Errorf("write readme: %w", err)
	}

	fmt.Println("\n✅ Knowhub initialized successfully!")
	fmt.Println()
	fmt.Println("📁 Directory structure:")
	fmt.Printf("   %s/\n", filepath.Base(absPath))
	fmt.Println("   ├── .synapse/          # Synapse 配置")
	fmt.Println("   │   ├── schema.yaml    # 知识规范")
	fmt.Println("   │   └── config.yaml    # 扩展点配置")
	fmt.Println("   ├── profile/           # 用户画像")
	fmt.Println("   │   └── me.md")
	fmt.Println("   ├── topics/            # 主题知识")
	fmt.Println("   ├── entities/          # 实体页")
	fmt.Println("   ├── concepts/          # 概念页")
	fmt.Println("   ├── inbox/             # 增量缓冲区")
	fmt.Println("   ├── journal/           # 时间线")
	fmt.Println("   └── graph/             # 知识图谱")
	fmt.Println("       └── relations.json")
	fmt.Println()
	fmt.Println("🚀 Next steps:")
	fmt.Println("   1. Edit profile/me.md to describe yourself")
	fmt.Println("   2. Configure .synapse/config.yaml for your setup")
	fmt.Println("   3. Start accumulating knowledge with AI assistants!")

	return nil
}

// createDirectories 创建 knowhub 目录结构
func createDirectories(basePath string) error {
	dirs := []string{
		".synapse",
		"profile",
		"topics",
		"entities",
		"concepts",
		"inbox",
		"journal",
		"graph",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(basePath, dir)
		if err := os.MkdirAll(fullPath, 0o755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
		fmt.Printf("  📂 Created %s/\n", dir)
	}

	return nil
}

// writeSchema 生成 .synapse/schema.yaml
func writeSchema(basePath string) error {
	s := schema.Default()
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal schema: %w", err)
	}

	header := []byte("# Synapse Knowledge Schema — 知识库行为契约\n# 所有扩展点共同遵守此规范，修改 Schema 即修改所有 AI 助手的行为\n#\n# 文档: https://github.com/tunsuy/synapse/blob/main/docs/roadmap.md\n\n")
	content := make([]byte, 0, len(header)+len(data))
	content = append(content, header...)
	content = append(content, data...)

	path := filepath.Join(basePath, ".synapse", "schema.yaml")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	fmt.Println("  📄 Created .synapse/schema.yaml")
	return nil
}

// writeConfig 生成 .synapse/config.yaml
func writeConfig(basePath string) error {
	cfg := config.Default()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	header := []byte("# Synapse Configuration — 扩展点注册中心\n# 声明各扩展点使用哪个实现\n#\n# 文档: https://github.com/tunsuy/synapse/blob/main/docs/roadmap.md\n\n")
	content := make([]byte, 0, len(header)+len(data))
	content = append(content, header...)
	content = append(content, data...)

	path := filepath.Join(basePath, ".synapse", "config.yaml")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	fmt.Println("  📄 Created .synapse/config.yaml")
	return nil
}

// writeTemplates 生成模板文件
func writeTemplates(basePath, name string) error {
	now := time.Now().Format(time.RFC3339)

	if name == "" {
		name = "Synapse User"
	}

	// profile/me.md
	profileContent := fmt.Sprintf(`---
type: profile
title: "%s"
created: %s
updated: %s
tags:
  - profile
---

# 👤 %s

## 简介

<!-- 在这里描述自己，AI 助手会参考这个画像来更好地理解你 -->

## 技术栈

<!-- 列出你的主要技术栈 -->

## 兴趣领域

<!-- 列出你感兴趣的领域 -->

## 当前关注

<!-- 列出你最近在关注/学习的内容 -->
`, name, now, now, name)

	if err := writeFile(filepath.Join(basePath, "profile", "me.md"), profileContent); err != nil {
		return err
	}
	fmt.Println("  📄 Created profile/me.md")

	// topics/.gitkeep
	if err := writeFile(filepath.Join(basePath, "topics", ".gitkeep"), ""); err != nil {
		return err
	}

	// entities/.gitkeep
	if err := writeFile(filepath.Join(basePath, "entities", ".gitkeep"), ""); err != nil {
		return err
	}

	// concepts/.gitkeep
	if err := writeFile(filepath.Join(basePath, "concepts", ".gitkeep"), ""); err != nil {
		return err
	}

	// inbox/.gitkeep
	if err := writeFile(filepath.Join(basePath, "inbox", ".gitkeep"), ""); err != nil {
		return err
	}

	// journal/.gitkeep
	if err := writeFile(filepath.Join(basePath, "journal", ".gitkeep"), ""); err != nil {
		return err
	}

	// graph/relations.json
	relationsContent := `{
  "version": "1.0",
  "nodes": [],
  "edges": [],
  "metadata": {
    "generated": "auto",
    "description": "Knowledge graph relations — auto-generated from [[wiki-links]]"
  }
}
`
	if err := writeFile(filepath.Join(basePath, "graph", "relations.json"), relationsContent); err != nil {
		return err
	}
	fmt.Println("  📄 Created graph/relations.json")

	return nil
}

// writeGitignore 生成 .gitignore
func writeGitignore(basePath string) error {
	content := `# Synapse 生成的临时文件
.synapse/cache/
.synapse/index/

# 操作系统文件
.DS_Store
Thumbs.db

# 编辑器临时文件
*.swp
*.swo
*~
.vscode/
.idea/
`
	if err := writeFile(filepath.Join(basePath, ".gitignore"), content); err != nil {
		return err
	}
	fmt.Println("  📄 Created .gitignore")
	return nil
}

// writeReadme 生成 README.md
func writeReadme(basePath, name string) error {
	content := tmpl.GenerateReadme(name, schema.Default())

	if err := writeFile(filepath.Join(basePath, "README.md"), content); err != nil {
		return err
	}
	fmt.Println("  📄 Created README.md")
	return nil
}

// writeFile 写入文件的辅助函数
func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
