package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/tunsuy/synapse/pkg/extension"
	"github.com/tunsuy/synapse/pkg/model"
)

func TestNew(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		config  map[string]any
		wantErr bool
	}{
		{
			name: "valid config",
			config: map[string]any{
				"owner": "tunsuy",
				"repo":  "synapse",
				"token": "ghp_test123",
			},
			wantErr: false,
		},
		{
			name: "with optional fields",
			config: map[string]any{
				"owner":    "tunsuy",
				"repo":     "synapse",
				"token":    "ghp_test123",
				"branch":   "develop",
				"base_dir": "knowledge",
			},
			wantErr: false,
		},
		{
			name: "missing owner",
			config: map[string]any{
				"repo":  "synapse",
				"token": "ghp_test123",
			},
			wantErr: true,
		},
		{
			name: "missing repo",
			config: map[string]any{
				"owner": "tunsuy",
				"token": "ghp_test123",
			},
			wantErr: true,
		},
		{
			name: "missing token",
			config: map[string]any{
				"owner": "tunsuy",
				"repo":  "synapse",
			},
			wantErr: true,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			store, err := New(tc.config)
			if (err != nil) != tc.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && store == nil {
				t.Error("New() returned nil store without error")
			}
		})
	}
}

func TestGitHubStore_Name(t *testing.T) {
	t.Parallel()
	store := &GitHubStore{}
	if got := store.Name(); got != "github-store" {
		t.Errorf("Name() = %q, want %q", got, "github-store")
	}
}

func TestGitHubStore_ImplementsVersionedStore(t *testing.T) {
	t.Parallel()
	store := &GitHubStore{}

	// 验证实现了 Store 接口
	var _ extension.Store = store

	// 验证实现了 VersionedStore 接口
	var _ extension.VersionedStore = store
}

func TestGitHubStore_Read(t *testing.T) {
	t.Parallel()

	content := "---\ntype: topic\ntitle: \"Go Programming\"\ncreated: 2025-04-18T12:00:00Z\nupdated: 2025-04-18T12:00:00Z\ntags:\n  - golang\n---\n\n# Go Programming\n\nGo is awesome."
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/contents/topics/golang.md") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
		}

		resp := contentsResponse{
			Name:    "golang.md",
			Path:    "topics/golang.md",
			SHA:     "abc123",
			Size:    int64(len(content)),
			Type:    "file",
			Content: encoded,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	kf, err := store.Read(context.Background(), "topics/golang.md")
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if kf.Frontmatter.Title != "Go Programming" {
		t.Errorf("Title = %q, want %q", kf.Frontmatter.Title, "Go Programming")
	}
	if kf.Frontmatter.Type != model.PageTypeTopic {
		t.Errorf("Type = %q, want %q", kf.Frontmatter.Type, model.PageTypeTopic)
	}
	if !strings.Contains(kf.Body, "Go is awesome") {
		t.Errorf("Body should contain 'Go is awesome', got: %q", kf.Body)
	}
}

func TestGitHubStore_Read_NotFound(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message": "Not Found"}`)
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	_, err := store.Read(context.Background(), "nonexistent.md")
	if err == nil {
		t.Fatal("Read() expected error for missing file")
	}
}

func TestGitHubStore_Write_Create(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 4, 18, 12, 0, 0, 0, time.UTC)
	kf := model.KnowledgeFile{
		Path: "topics/golang.md",
		Frontmatter: model.Frontmatter{
			Type:    model.PageTypeTopic,
			Title:   "Go Programming",
			Created: now,
			Updated: now,
		},
		Body: "# Go Programming",
	}

	var putCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// getFileSHA 返回 404（新文件场景）
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message": "Not Found"}`)
		case http.MethodPut:
			putCalled = true
			var req createFileRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("decode put body: %v", err)
			}
			if req.SHA != "" {
				t.Error("SHA should be empty for new file")
			}
			if req.Branch != "main" {
				t.Errorf("Branch = %q, want 'main'", req.Branch)
			}
			if req.Content == "" {
				t.Error("Content should not be empty")
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"content": {"sha": "new123"}}`)
		default:
			t.Errorf("unexpected method: %s", r.Method)
		}
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	if err := store.Write(context.Background(), kf); err != nil {
		t.Fatalf("Write() error: %v", err)
	}
	if !putCalled {
		t.Error("PUT request was not made")
	}
}

func TestGitHubStore_Write_Update(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 4, 18, 12, 0, 0, 0, time.UTC)
	kf := model.KnowledgeFile{
		Path: "topics/golang.md",
		Frontmatter: model.Frontmatter{
			Type:    model.PageTypeTopic,
			Title:   "Go Programming Updated",
			Created: now,
			Updated: now,
		},
		Body: "# Go Programming\n\nUpdated content.",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// 返回已有文件的 SHA
			resp := contentsResponse{SHA: "existing-sha-456"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		case http.MethodPut:
			var req createFileRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("decode put body: %v", err)
			}
			if req.SHA != "existing-sha-456" {
				t.Errorf("SHA = %q, want 'existing-sha-456'", req.SHA)
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"content": {"sha": "updated-sha-789"}}`)
		default:
			t.Errorf("unexpected method: %s", r.Method)
		}
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	if err := store.Write(context.Background(), kf); err != nil {
		t.Fatalf("Write() error: %v", err)
	}
}

func TestGitHubStore_Delete(t *testing.T) {
	t.Parallel()

	var deleteCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			resp := contentsResponse{SHA: "delete-sha-123"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		case http.MethodDelete:
			deleteCalled = true
			var req deleteFileRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("decode delete body: %v", err)
			}
			if req.SHA != "delete-sha-123" {
				t.Errorf("SHA = %q, want 'delete-sha-123'", req.SHA)
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"commit": {"sha": "commit-after-delete"}}`)
		default:
			t.Errorf("unexpected method: %s", r.Method)
		}
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	if err := store.Delete(context.Background(), "topics/golang.md"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
	if !deleteCalled {
		t.Error("DELETE request was not made")
	}
}

func TestGitHubStore_List(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items := []contentsResponse{
			{Name: "golang.md", Path: "topics/golang.md", Type: "file", Size: 100},
			{Name: "rust.md", Path: "topics/rust.md", Type: "file", Size: 200},
			{Name: ".hidden.md", Path: "topics/.hidden.md", Type: "file", Size: 50},
			{Name: "readme.txt", Path: "topics/readme.txt", Type: "file", Size: 30},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	infos, err := store.List(context.Background(), "topics", model.ListOptions{})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	// 应排除 .hidden.md 和 readme.txt
	if len(infos) != 2 {
		t.Errorf("List() returned %d files, want 2", len(infos))
	}
}

func TestGitHubStore_List_WithLimit(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items := []contentsResponse{
			{Name: "golang.md", Path: "topics/golang.md", Type: "file", Size: 100},
			{Name: "rust.md", Path: "topics/rust.md", Type: "file", Size: 200},
			{Name: "python.md", Path: "topics/python.md", Type: "file", Size: 150},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	infos, err := store.List(context.Background(), "topics", model.ListOptions{Limit: 1})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(infos) != 1 {
		t.Errorf("List(limit=1) returned %d files, want 1", len(infos))
	}
}

func TestGitHubStore_List_Recursive(t *testing.T) {
	t.Parallel()

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		var items []contentsResponse
		if strings.Contains(r.URL.Path, "/contents/topics") && !strings.Contains(r.URL.Path, "/contents/topics/advanced") {
			items = []contentsResponse{
				{Name: "golang.md", Path: "topics/golang.md", Type: "file", Size: 100},
				{Name: "advanced", Path: "topics/advanced", Type: "dir"},
			}
		} else if strings.Contains(r.URL.Path, "/contents/topics/advanced") {
			items = []contentsResponse{
				{Name: "concurrency.md", Path: "topics/advanced/concurrency.md", Type: "file", Size: 300},
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	infos, err := store.List(context.Background(), "topics", model.ListOptions{Recursive: true})
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(infos) != 2 {
		t.Errorf("List(recursive) returned %d files, want 2", len(infos))
	}

	// 验证递归调用了子目录
	if callCount < 2 {
		t.Errorf("expected at least 2 API calls for recursive listing, got %d", callCount)
	}
}

func TestGitHubStore_Exists(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "existing.md") {
			resp := contentsResponse{Name: "existing.md", SHA: "abc123"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message": "Not Found"}`)
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	ctx := context.Background()

	// 存在
	exists, err := store.Exists(ctx, "existing.md")
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if !exists {
		t.Error("existing file should return true")
	}

	// 不存在
	exists, err = store.Exists(ctx, "nonexistent.md")
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if exists {
		t.Error("nonexistent file should return false")
	}
}

func TestGitHubStore_History(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/commits") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		commits := []commitInfo{
			{
				SHA: "sha-1",
				Commit: struct {
					Message string `json:"message"`
					Author  struct {
						Name string `json:"name"`
						Date string `json:"date"`
					} `json:"author"`
				}{
					Message: "update golang.md",
					Author: struct {
						Name string `json:"name"`
						Date string `json:"date"`
					}{
						Name: "tunsuy",
						Date: "2025-04-18T12:00:00Z",
					},
				},
			},
			{
				SHA: "sha-2",
				Commit: struct {
					Message string `json:"message"`
					Author  struct {
						Name string `json:"name"`
						Date string `json:"date"`
					} `json:"author"`
				}{
					Message: "create golang.md",
					Author: struct {
						Name string `json:"name"`
						Date string `json:"date"`
					}{
						Name: "tunsuy",
						Date: "2025-04-17T10:00:00Z",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(commits)
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	records, err := store.History(context.Background(), "topics/golang.md")
	if err != nil {
		t.Fatalf("History() error: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("History() returned %d records, want 2", len(records))
	}

	if records[0].Hash != "sha-1" {
		t.Errorf("records[0].Hash = %q, want 'sha-1'", records[0].Hash)
	}
	if records[0].Message != "update golang.md" {
		t.Errorf("records[0].Message = %q, want 'update golang.md'", records[0].Message)
	}
	if records[0].Author != "tunsuy" {
		t.Errorf("records[0].Author = %q, want 'tunsuy'", records[0].Author)
	}
}

func TestGitHubStore_Commit(t *testing.T) {
	t.Parallel()

	step := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		step++
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/git/refs/heads/main"):
			// Step 1: 获取分支引用
			fmt.Fprint(w, `{"ref": "refs/heads/main", "object": {"sha": "parent-sha-1"}}`)

		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/git/commits/parent-sha-1"):
			// Step 2: 获取父 commit
			fmt.Fprint(w, `{"sha": "parent-sha-1", "tree": {"sha": "tree-sha-1"}}`)

		case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/git/commits"):
			// Step 3: 创建新 commit
			var req createCommitRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("decode commit body: %v", err)
			}
			if req.Message != "synapse: batch update" {
				t.Errorf("Message = %q, want 'synapse: batch update'", req.Message)
			}
			if req.Tree != "tree-sha-1" {
				t.Errorf("Tree = %q, want 'tree-sha-1'", req.Tree)
			}
			if len(req.Parents) != 1 || req.Parents[0] != "parent-sha-1" {
				t.Errorf("Parents = %v, want ['parent-sha-1']", req.Parents)
			}
			fmt.Fprint(w, `{"sha": "new-commit-sha"}`)

		case r.Method == http.MethodPatch && strings.Contains(r.URL.Path, "/git/refs/heads/main"):
			// Step 4: 更新分支引用
			var req updateRefRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("decode ref body: %v", err)
			}
			if req.SHA != "new-commit-sha" {
				t.Errorf("SHA = %q, want 'new-commit-sha'", req.SHA)
			}
			fmt.Fprint(w, `{"ref": "refs/heads/main", "object": {"sha": "new-commit-sha"}}`)

		default:
			t.Errorf("unexpected request: %s %s (step %d)", r.Method, r.URL.Path, step)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	if err := store.Commit(context.Background(), "synapse: batch update"); err != nil {
		t.Fatalf("Commit() error: %v", err)
	}
}

func TestGitHubStore_WithBaseDir(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证路径包含 base_dir 前缀
		if !strings.Contains(r.URL.Path, "/contents/knowledge/topics/golang.md") {
			t.Errorf("path should include base_dir, got: %s", r.URL.Path)
		}

		content := "# Go"
		encoded := base64.StdEncoding.EncodeToString([]byte(content))
		resp := contentsResponse{
			Name:    "golang.md",
			Path:    "knowledge/topics/golang.md",
			SHA:     "abc123",
			Type:    "file",
			Content: encoded,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	store := &GitHubStore{
		owner:   "tunsuy",
		repo:    "synapse",
		branch:  "main",
		baseDir: "knowledge",
		token:   "test-token",
		apiBase: server.URL,
		client:  server.Client(),
	}

	_, err := store.Read(context.Background(), "topics/golang.md")
	if err != nil {
		t.Fatalf("Read() with base_dir error: %v", err)
	}
}

func TestFullPath(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		baseDir string
		path    string
		want    string
	}{
		{
			name:    "no base dir",
			baseDir: "",
			path:    "topics/golang.md",
			want:    "topics/golang.md",
		},
		{
			name:    "with base dir",
			baseDir: "knowledge",
			path:    "topics/golang.md",
			want:    "knowledge/topics/golang.md",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			store := &GitHubStore{baseDir: tc.baseDir}
			got := store.fullPath(tc.path)
			if got != tc.want {
				t.Errorf("fullPath(%q) = %q, want %q", tc.path, got, tc.want)
			}
		})
	}
}

func TestParseKnowledgeFile(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		path      string
		data      string
		wantTitle string
		wantType  model.PageType
		wantBody  string
	}{
		{
			name:      "with frontmatter",
			path:      "topics/golang.md",
			data:      "---\ntype: topic\ntitle: \"Go Programming\"\ncreated: 2025-04-18T12:00:00Z\nupdated: 2025-04-18T12:00:00Z\n---\n\n# Go Programming",
			wantTitle: "Go Programming",
			wantType:  model.PageTypeTopic,
			wantBody:  "# Go Programming",
		},
		{
			name:      "without frontmatter",
			path:      "topics/test.md",
			data:      "# Just some content",
			wantTitle: "test",
			wantType:  "",
			wantBody:  "# Just some content",
		},
		{
			name:      "with unquoted title",
			path:      "topics/unquoted.md",
			data:      "---\ntype: entity\ntitle: CodeBuddy\n---\n\nEntity content",
			wantTitle: "CodeBuddy",
			wantType:  model.PageTypeEntity,
			wantBody:  "Entity content",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			kf := parseKnowledgeFile(tc.path, []byte(tc.data))
			if kf.Frontmatter.Title != tc.wantTitle {
				t.Errorf("Title = %q, want %q", kf.Frontmatter.Title, tc.wantTitle)
			}
			if kf.Frontmatter.Type != tc.wantType {
				t.Errorf("Type = %q, want %q", kf.Frontmatter.Type, tc.wantType)
			}
			if kf.Body != tc.wantBody {
				t.Errorf("Body = %q, want %q", kf.Body, tc.wantBody)
			}
		})
	}
}

func TestIsNotFoundErr(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "404 error",
			err:  &apiError{StatusCode: 404, Message: "not found"},
			want: true,
		},
		{
			name: "500 error",
			err:  &apiError{StatusCode: 500, Message: "server error"},
			want: false,
		},
		{
			name: "non-api error",
			err:  fmt.Errorf("some error"),
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := isNotFoundErr(tc.err)
			if got != tc.want {
				t.Errorf("isNotFoundErr() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestApiError_Error(t *testing.T) {
	t.Parallel()

	err := &apiError{StatusCode: 404, Message: "not found"}
	expected := "github api error (status 404): not found"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}
