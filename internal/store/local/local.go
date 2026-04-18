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
// 简单解析 YAML Frontmatter + Markdown Body
func parseKnowledgeFile(path string, data []byte) (model.KnowledgeFile, error) {
	content := string(data)

	kf := model.KnowledgeFile{
		Path: path,
	}

	// 解析 frontmatter
	if strings.HasPrefix(content, "---\n") {
		end := strings.Index(content[4:], "\n---\n")
		if end != -1 {
			// 简单提取 title
			fm := content[4 : 4+end]
			for _, line := range strings.Split(fm, "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "title:") {
					title := strings.TrimPrefix(line, "title:")
					title = strings.TrimSpace(title)
					title = strings.Trim(title, "\"'")
					kf.Frontmatter.Title = title
				}
				if strings.HasPrefix(line, "type:") {
					t := strings.TrimPrefix(line, "type:")
					t = strings.TrimSpace(t)
					kf.Frontmatter.Type = model.PageType(t)
				}
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
