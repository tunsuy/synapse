// Package github 实现基于 GitHub 仓库的 VersionedStore
// 通过 GitHub REST API 实现知识文件的 CRUD 和版本控制操作
package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/tunsuy/synapse/internal/store/tmpl"
	"github.com/tunsuy/synapse/pkg/extension"
	"github.com/tunsuy/synapse/pkg/model"
)

func init() {
	extension.RegisterStore("github-store", New)
}

const (
	defaultAPIBase = "https://api.github.com"
	defaultBranch  = "main"
)

// GitHubStore 基于 GitHub 仓库的 VersionedStore 实现
// 通过 GitHub REST API 进行知识文件的读写和版本控制
type GitHubStore struct {
	owner   string
	repo    string
	branch  string
	baseDir string
	token   string
	apiBase string
	client  *http.Client
}

// New 创建一个新的 GitHubStore 实例
func New(config map[string]any) (extension.Store, error) {
	owner, _ := config["owner"].(string)
	if owner == "" {
		return nil, fmt.Errorf("github-store requires 'owner' config")
	}

	repo, _ := config["repo"].(string)
	if repo == "" {
		return nil, fmt.Errorf("github-store requires 'repo' config")
	}

	token, _ := config["token"].(string)
	if token == "" {
		return nil, fmt.Errorf("github-store requires 'token' config")
	}

	branch, _ := config["branch"].(string)
	if branch == "" {
		branch = defaultBranch
	}

	baseDir, _ := config["base_dir"].(string)

	apiBase, _ := config["api_base"].(string)
	if apiBase == "" {
		apiBase = defaultAPIBase
	}

	return &GitHubStore{
		owner:   owner,
		repo:    repo,
		branch:  branch,
		baseDir: baseDir,
		token:   token,
		apiBase: apiBase,
		client:  &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// Name 返回存储后端名称
func (s *GitHubStore) Name() string {
	return "github-store"
}

// Init 通过 GitHub API 初始化远端知识库结构
func (s *GitHubStore) Init(ctx context.Context, opts extension.InitOptions) error {
	name := opts.Name
	if name == "" {
		name = "Synapse User"
	}
	now := time.Now().Format(time.RFC3339)

	// 定义需要创建的文件列表
	type fileEntry struct {
		path    string
		content string
	}

	files := []fileEntry{
		{
			path:    "profile/me.md",
			content: fmt.Sprintf("---\ntype: profile\ntitle: \"%s\"\ncreated: %s\nupdated: %s\ntags:\n  - profile\n---\n\n# 👤 %s\n\n## 简介\n\n<!-- 在这里描述自己 -->\n\n## 技术栈\n\n<!-- 列出你的主要技术栈 -->\n\n## 兴趣领域\n\n<!-- 列出你感兴趣的领域 -->\n\n## 当前关注\n\n<!-- 列出你最近在关注/学习的内容 -->\n", name, now, now, name),
		},
		{path: "topics/.gitkeep", content: ""},
		{path: "entities/.gitkeep", content: ""},
		{path: "concepts/.gitkeep", content: ""},
		{path: "inbox/.gitkeep", content: ""},
		{path: "journal/.gitkeep", content: ""},
		{
			path:    "graph/relations.json",
			content: "{\n  \"version\": \"1.0\",\n  \"nodes\": [],\n  \"edges\": [],\n  \"metadata\": {\n    \"generated\": \"auto\",\n    \"description\": \"Knowledge graph relations — auto-generated from [[wiki-links]]\"\n  }\n}\n",
		},
		{
			path:    "README.md",
			content: tmpl.GenerateReadme(name),
		},
	}

	// 如果有 schema 数据，先写入
	if len(opts.SchemaData) > 0 {
		header := "# Synapse Knowledge Schema — 知识库行为契约\n# 所有扩展点共同遵守此规范\n\n"
		files = append([]fileEntry{{
			path:    ".synapse/schema.yaml",
			content: header + string(opts.SchemaData),
		}}, files...)
	}

	// 逐个写入文件
	for _, f := range files {
		kf := model.KnowledgeFile{Path: f.path, Body: f.content}
		data := []byte(f.content)

		fullPath := s.fullPath(f.path)
		url := fmt.Sprintf("%s/repos/%s/%s/contents/%s",
			s.apiBase, s.owner, s.repo, fullPath)

		reqBody := createFileRequest{
			Message: fmt.Sprintf("synapse: init %s", f.path),
			Content: base64.StdEncoding.EncodeToString(data),
			Branch:  s.branch,
		}

		_ = kf // suppress unused

		payload, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("marshal request for %s: %w", f.path, err)
		}

		if _, err := s.doRequest(ctx, http.MethodPut, url, payload); err != nil {
			return fmt.Errorf("init %s: %w", f.path, err)
		}
	}

	return nil
}

// Initialized 检查远端知识库是否已初始化（通过检查 .synapse/schema.yaml 是否存在）
func (s *GitHubStore) Initialized(ctx context.Context) (bool, error) {
	return s.Exists(ctx, ".synapse/schema.yaml")
}

// Read 通过 GitHub Contents API 读取知识文件
func (s *GitHubStore) Read(ctx context.Context, path string) (model.KnowledgeFile, error) {
	fullPath := s.fullPath(path)

	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s",
		s.apiBase, s.owner, s.repo, fullPath, s.branch)

	body, err := s.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.KnowledgeFile{}, fmt.Errorf("read %s: %w", path, err)
	}

	var resp contentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return model.KnowledgeFile{}, fmt.Errorf("parse response for %s: %w", path, err)
	}

	data, err := base64.StdEncoding.DecodeString(resp.Content)
	if err != nil {
		return model.KnowledgeFile{}, fmt.Errorf("decode content for %s: %w", path, err)
	}

	kf, err := parseKnowledgeFile(path, data)
	if err != nil {
		return model.KnowledgeFile{}, fmt.Errorf("parse %s: %w", path, err)
	}

	return kf, nil
}

// Write 通过 GitHub Contents API 写入知识文件
func (s *GitHubStore) Write(ctx context.Context, file model.KnowledgeFile) error {
	fullPath := s.fullPath(file.Path)

	data, err := file.Marshal()
	if err != nil {
		return fmt.Errorf("marshal %s: %w", file.Path, err)
	}

	// 尝试获取现有文件的 SHA（用于更新场景）
	sha, _ := s.getFileSHA(ctx, fullPath)

	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s",
		s.apiBase, s.owner, s.repo, fullPath)

	reqBody := createFileRequest{
		Message: fmt.Sprintf("synapse: update %s", file.Path),
		Content: base64.StdEncoding.EncodeToString(data),
		Branch:  s.branch,
	}
	if sha != "" {
		reqBody.SHA = sha
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request for %s: %w", file.Path, err)
	}

	if _, err := s.doRequest(ctx, http.MethodPut, url, payload); err != nil {
		return fmt.Errorf("write %s: %w", file.Path, err)
	}

	return nil
}

// Delete 通过 GitHub Contents API 删除知识文件
func (s *GitHubStore) Delete(ctx context.Context, path string) error {
	fullPath := s.fullPath(path)

	sha, err := s.getFileSHA(ctx, fullPath)
	if err != nil {
		return fmt.Errorf("get sha for %s: %w", path, err)
	}

	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s",
		s.apiBase, s.owner, s.repo, fullPath)

	reqBody := deleteFileRequest{
		Message: fmt.Sprintf("synapse: delete %s", path),
		SHA:     sha,
		Branch:  s.branch,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal delete request for %s: %w", path, err)
	}

	if _, err := s.doRequest(ctx, http.MethodDelete, url, payload); err != nil {
		return fmt.Errorf("delete %s: %w", path, err)
	}

	return nil
}

// List 通过 GitHub Contents API 列出指定目录下的知识文件
func (s *GitHubStore) List(ctx context.Context, dir string, opts model.ListOptions) ([]model.FileInfo, error) {
	var files []model.FileInfo

	if err := s.listDir(ctx, dir, opts, &files); err != nil {
		return nil, fmt.Errorf("list %s: %w", dir, err)
	}

	return files, nil
}

// Exists 检查文件是否存在
func (s *GitHubStore) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := s.fullPath(path)

	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s",
		s.apiBase, s.owner, s.repo, fullPath, s.branch)

	_, err := s.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		if isNotFoundErr(err) {
			return false, nil
		}
		return false, fmt.Errorf("check exists %s: %w", path, err)
	}

	return true, nil
}

// --- VersionedStore 接口实现 ---

// Commit 通过创建 Git commit 提交当前变更
// GitHub Contents API 的每次 PUT/DELETE 操作已自动创建 commit，
// 此方法用于在已有 commit 基础上添加一条空的标记 commit（通过更新 branch ref）
func (s *GitHubStore) Commit(ctx context.Context, message string) error {
	// 获取当前分支的最新 commit SHA
	refURL := fmt.Sprintf("%s/repos/%s/%s/git/refs/heads/%s",
		s.apiBase, s.owner, s.repo, s.branch)

	body, err := s.doRequest(ctx, http.MethodGet, refURL, nil)
	if err != nil {
		return fmt.Errorf("get branch ref: %w", err)
	}

	var ref gitRef
	if err := json.Unmarshal(body, &ref); err != nil {
		return fmt.Errorf("parse branch ref: %w", err)
	}

	parentSHA := ref.Object.SHA

	// 获取父 commit 的 tree SHA
	commitURL := fmt.Sprintf("%s/repos/%s/%s/git/commits/%s",
		s.apiBase, s.owner, s.repo, parentSHA)

	body, err = s.doRequest(ctx, http.MethodGet, commitURL, nil)
	if err != nil {
		return fmt.Errorf("get parent commit: %w", err)
	}

	var parentCommit gitCommit
	if err := json.Unmarshal(body, &parentCommit); err != nil {
		return fmt.Errorf("parse parent commit: %w", err)
	}

	// 创建新 commit（复用 parent tree，仅更新 message）
	createCommitURL := fmt.Sprintf("%s/repos/%s/%s/git/commits",
		s.apiBase, s.owner, s.repo)

	createReq := createCommitRequest{
		Message: message,
		Tree:    parentCommit.Tree.SHA,
		Parents: []string{parentSHA},
	}

	payload, err := json.Marshal(createReq)
	if err != nil {
		return fmt.Errorf("marshal commit request: %w", err)
	}

	body, err = s.doRequest(ctx, http.MethodPost, createCommitURL, payload)
	if err != nil {
		return fmt.Errorf("create commit: %w", err)
	}

	var newCommit gitCommit
	if err := json.Unmarshal(body, &newCommit); err != nil {
		return fmt.Errorf("parse new commit: %w", err)
	}

	// 更新分支引用指向新 commit
	updateRefReq := updateRefRequest{
		SHA: newCommit.SHA,
	}

	payload, err = json.Marshal(updateRefReq)
	if err != nil {
		return fmt.Errorf("marshal update ref request: %w", err)
	}

	if _, err := s.doRequest(ctx, http.MethodPatch, refURL, payload); err != nil {
		return fmt.Errorf("update branch ref: %w", err)
	}

	return nil
}

// History 获取文件的变更历史
func (s *GitHubStore) History(ctx context.Context, path string) ([]extension.ChangeRecord, error) {
	fullPath := s.fullPath(path)

	url := fmt.Sprintf("%s/repos/%s/%s/commits?path=%s&sha=%s",
		s.apiBase, s.owner, s.repo, fullPath, s.branch)

	body, err := s.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("get history for %s: %w", path, err)
	}

	var commits []commitInfo
	if err := json.Unmarshal(body, &commits); err != nil {
		return nil, fmt.Errorf("parse commits for %s: %w", path, err)
	}

	records := make([]extension.ChangeRecord, 0, len(commits))
	for _, c := range commits {
		records = append(records, extension.ChangeRecord{
			Hash:      c.SHA,
			Message:   c.Commit.Message,
			Author:    c.Commit.Author.Name,
			Timestamp: c.Commit.Author.Date,
		})
	}

	return records, nil
}

// --- 内部辅助方法 ---

// fullPath 计算文件在仓库中的完整路径
func (s *GitHubStore) fullPath(path string) string {
	if s.baseDir == "" {
		return path
	}
	return s.baseDir + "/" + path
}

// getFileSHA 获取文件的当前 SHA 值
func (s *GitHubStore) getFileSHA(ctx context.Context, fullPath string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s",
		s.apiBase, s.owner, s.repo, fullPath, s.branch)

	body, err := s.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	var resp contentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parse sha response: %w", err)
	}

	return resp.SHA, nil
}

// listDir 递归或非递归地列出目录内容
func (s *GitHubStore) listDir(
	ctx context.Context,
	dir string,
	opts model.ListOptions,
	files *[]model.FileInfo,
) error {
	fullDir := s.fullPath(dir)

	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s",
		s.apiBase, s.owner, s.repo, fullDir, s.branch)

	body, err := s.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	var items []contentsResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return fmt.Errorf("parse directory listing: %w", err)
	}

	for _, item := range items {
		if opts.Limit > 0 && len(*files) >= opts.Limit {
			return nil
		}

		if item.Type == "dir" {
			if opts.Recursive {
				subDir := item.Path
				if s.baseDir != "" {
					subDir = strings.TrimPrefix(subDir, s.baseDir+"/")
				}
				if err := s.listDir(ctx, subDir, opts, files); err != nil {
					return err
				}
			}
			continue
		}

		// 只处理 .md 和 .json 文件
		ext := filepath.Ext(item.Name)
		if ext != ".md" && ext != ".json" {
			continue
		}

		// 跳过隐藏文件
		if strings.HasPrefix(item.Name, ".") {
			continue
		}

		relPath := item.Path
		if s.baseDir != "" {
			relPath = strings.TrimPrefix(item.Path, s.baseDir+"/")
		}

		fi := model.FileInfo{
			Path:  relPath,
			Title: strings.TrimSuffix(item.Name, ext),
			Size:  item.Size,
		}

		*files = append(*files, fi)
	}

	return nil
}

// doRequest 执行带认证的 HTTP 请求
func (s *GitHubStore) doRequest(ctx context.Context, method, url string, body []byte) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = strings.NewReader(string(body))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, &apiError{StatusCode: resp.StatusCode, Message: "not found"}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &apiError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
		}
	}

	return respBody, nil
}

// --- 辅助解析函数 ---

// parseKnowledgeFile 解析 Markdown 文件为 KnowledgeFile
// 复用 local store 的解析逻辑
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
		if title, ok := strings.CutPrefix(line, "title:"); ok {
			title = strings.TrimSpace(title)
			title = strings.Trim(title, "\"'")
			frontmatter.Title = title
		}
		if t, ok := strings.CutPrefix(line, "type:"); ok {
			t = strings.TrimSpace(t)
			frontmatter.Type = model.PageType(t)
		}
	}
}

// isNotFoundErr 判断错误是否为 404
func isNotFoundErr(err error) bool {
	if e, ok := err.(*apiError); ok {
		return e.StatusCode == http.StatusNotFound
	}
	return false
}

// --- API 请求/响应结构体 ---

// apiError GitHub API 错误
type apiError struct {
	StatusCode int
	Message    string
}

// Error 实现 error 接口
func (e *apiError) Error() string {
	return fmt.Sprintf("github api error (status %d): %s", e.StatusCode, e.Message)
}

// contentsResponse GitHub Contents API 响应
type contentsResponse struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	SHA     string `json:"sha"`
	Size    int64  `json:"size"`
	Type    string `json:"type"` // "file" or "dir"
	Content string `json:"content"`
}

// createFileRequest 创建/更新文件请求
type createFileRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
	Branch  string `json:"branch"`
	SHA     string `json:"sha,omitempty"`
}

// deleteFileRequest 删除文件请求
type deleteFileRequest struct {
	Message string `json:"message"`
	SHA     string `json:"sha"`
	Branch  string `json:"branch"`
}

// commitInfo GitHub Commits API 响应
type commitInfo struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Name string `json:"name"`
			Date string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

// gitRef Git 引用
type gitRef struct {
	Ref    string `json:"ref"`
	Object struct {
		SHA string `json:"sha"`
	} `json:"object"`
}

// gitCommit Git commit 对象
type gitCommit struct {
	SHA  string `json:"sha"`
	Tree struct {
		SHA string `json:"sha"`
	} `json:"tree"`
}

// createCommitRequest 创建 commit 请求
type createCommitRequest struct {
	Message string   `json:"message"`
	Tree    string   `json:"tree"`
	Parents []string `json:"parents"`
}

// updateRefRequest 更新分支引用请求
type updateRefRequest struct {
	SHA string `json:"sha"`
}
