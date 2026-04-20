// Package local 实现基于本地文件系统的 Store
// 这是 M1 阶段的参考实现，提供知识文件的本地 CRUD 操作
package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/tunsuy/synapse/internal/store/tmpl"
	"github.com/tunsuy/synapse/pkg/extension"
	"github.com/tunsuy/synapse/pkg/model"
)

func init() {
	extension.RegisterStore("local-store", New)
}

// LocalStore 基于本地文件系统的 Store 实现
type LocalStore struct {
	basePath string
}

// New 创建一个新的 LocalStore 实例
func New(config map[string]any) (extension.Store, error) {
	path, _ := config["path"].(string)
	if path == "" {
		return nil, fmt.Errorf("local-store requires 'path' config")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path %q: %w", path, err)
	}
	return &LocalStore{basePath: absPath}, nil
}

// Name 返回存储后端名称
func (s *LocalStore) Name() string {
	return "local-store"
}

// Init 初始化本地知识库目录结构
func (s *LocalStore) Init(ctx context.Context, opts extension.InitOptions) error {
	// 创建知识库目录结构
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
		fullPath := filepath.Join(s.basePath, dir)
		if err := os.MkdirAll(fullPath, 0o755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}

	// 写入 schema.yaml
	if len(opts.SchemaData) > 0 {
		schemaPath := filepath.Join(s.basePath, ".synapse", "schema.yaml")
		header := []byte("# Synapse Knowledge Schema — 知识库行为契约\n# 所有扩展点共同遵守此规范，修改 Schema 即修改所有 AI 助手的行为\n#\n# 文档: https://github.com/tunsuy/synapse/blob/main/docs/roadmap.md\n\n")
		content := append(header, opts.SchemaData...)
		if err := os.WriteFile(schemaPath, content, 0o644); err != nil {
			return fmt.Errorf("write schema.yaml: %w", err)
		}
	}

	// 写入模板文件
	if err := s.writeTemplates(opts.Name); err != nil {
		return fmt.Errorf("write templates: %w", err)
	}

	// 写入 .gitignore
	if err := s.writeGitignore(); err != nil {
		return fmt.Errorf("write .gitignore: %w", err)
	}

	// 写入 README.md
	if err := s.writeReadme(opts.Name); err != nil {
		return fmt.Errorf("write README.md: %w", err)
	}

	return nil
}

// Initialized 检查本地知识库是否已初始化
func (s *LocalStore) Initialized(ctx context.Context) (bool, error) {
	schemaPath := filepath.Join(s.basePath, ".synapse", "schema.yaml")
	_, err := os.Stat(schemaPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// writeTemplates 写入知识库模板文件
func (s *LocalStore) writeTemplates(name string) error {
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

	if err := os.WriteFile(filepath.Join(s.basePath, "profile", "me.md"), []byte(profileContent), 0o644); err != nil {
		return err
	}

	// 各目录的 .gitkeep
	keepDirs := []string{"topics", "entities", "concepts", "inbox", "journal"}
	for _, dir := range keepDirs {
		if err := os.WriteFile(filepath.Join(s.basePath, dir, ".gitkeep"), []byte(""), 0o644); err != nil {
			return err
		}
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
	if err := os.WriteFile(filepath.Join(s.basePath, "graph", "relations.json"), []byte(relationsContent), 0o644); err != nil {
		return err
	}

	return nil
}

// writeGitignore 写入 .gitignore
func (s *LocalStore) writeGitignore() error {
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
	return os.WriteFile(filepath.Join(s.basePath, ".gitignore"), []byte(content), 0o644)
}

// writeReadme 写入 README.md
func (s *LocalStore) writeReadme(name string) error {
	content := tmpl.GenerateReadme(name)
	return os.WriteFile(filepath.Join(s.basePath, "README.md"), []byte(content), 0o644)
}

// Read 读取知识文件
func (s *LocalStore) Read(ctx context.Context, path string) (model.KnowledgeFile, error) {
	fullPath := filepath.Join(s.basePath, path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return model.KnowledgeFile{}, fmt.Errorf("read %s: %w", path, err)
	}

	kf, err := parseKnowledgeFile(path, data)
	if err != nil {
		return model.KnowledgeFile{}, fmt.Errorf("parse %s: %w", path, err)
	}

	return kf, nil
}

// Write 写入知识文件
func (s *LocalStore) Write(ctx context.Context, file model.KnowledgeFile) error {
	fullPath := filepath.Join(s.basePath, file.Path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("create directory for %s: %w", file.Path, err)
	}

	data, err := file.Marshal()
	if err != nil {
		return fmt.Errorf("marshal %s: %w", file.Path, err)
	}

	return os.WriteFile(fullPath, data, 0o644)
}

// Delete 删除知识文件
func (s *LocalStore) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}

// List 列出指定目录下的知识文件
func (s *LocalStore) List(ctx context.Context, dir string, opts model.ListOptions) ([]model.FileInfo, error) {
	fullDir := filepath.Join(s.basePath, dir)

	var files []model.FileInfo

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if !opts.Recursive && path != fullDir {
				return filepath.SkipDir
			}
			return nil
		}

		// 只处理 .md 和 .json 文件
		ext := filepath.Ext(path)
		if ext != ".md" && ext != ".json" {
			return nil
		}

		// 跳过隐藏文件和 .gitkeep
		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") {
			return nil
		}

		relPath, err := filepath.Rel(s.basePath, path)
		if err != nil {
			return err
		}

		fi := model.FileInfo{
			Path:    relPath,
			Title:   strings.TrimSuffix(base, ext),
			Updated: info.ModTime(),
			Size:    info.Size(),
		}

		files = append(files, fi)

		if opts.Limit > 0 && len(files) >= opts.Limit {
			return filepath.SkipAll
		}

		return nil
	}

	if err := filepath.Walk(fullDir, walkFn); err != nil {
		return nil, fmt.Errorf("walk %s: %w", dir, err)
	}

	return files, nil
}

// Exists 检查文件是否存在
func (s *LocalStore) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(s.basePath, path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// parseKnowledgeFile 解析 Markdown 文件为 KnowledgeFile
// 使用 YAML 库完整解析 Frontmatter + Markdown Body
func parseKnowledgeFile(path string, data []byte) (model.KnowledgeFile, error) {
	content := string(data)

	kf := model.KnowledgeFile{
		Path: path,
	}

	// 解析 frontmatter
	if strings.HasPrefix(content, "---\n") {
		end := strings.Index(content[4:], "\n---\n")
		if end != -1 {
			fmData := content[4 : 4+end]

			if err := yaml.Unmarshal([]byte(fmData), &kf.Frontmatter); err != nil {
				// YAML 解析失败时回退到简单解析
				parseFrontmatterSimple(fmData, &kf.Frontmatter)
			}

			kf.Body = strings.TrimSpace(content[4+end+5:])
		} else {
			kf.Body = content
		}
	} else {
		kf.Body = content
	}

	if kf.Frontmatter.Title == "" {
		kf.Frontmatter.Title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	now := time.Now()
	if kf.Frontmatter.Created.IsZero() {
		kf.Frontmatter.Created = now
	}
	if kf.Frontmatter.Updated.IsZero() {
		kf.Frontmatter.Updated = now
	}

	return kf, nil
}

// parseFrontmatterSimple 简单回退解析（当 YAML 解析失败时使用）
func parseFrontmatterSimple(fm string, frontmatter *model.Frontmatter) {
	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "title:") {
			title := strings.TrimPrefix(line, "title:")
			title = strings.TrimSpace(title)
			title = strings.Trim(title, "\"'")
			frontmatter.Title = title
		}
		if strings.HasPrefix(line, "type:") {
			t := strings.TrimPrefix(line, "type:")
			t = strings.TrimSpace(t)
			frontmatter.Type = model.PageType(t)
		}
	}
}
