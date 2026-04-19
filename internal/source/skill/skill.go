// Package skill 实现基于 AI Skill 的 Source
// 这是 M2 阶段的参考实现，支持通过 CLI 参数或管道输入采集原始内容
package skill

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tunsuy/synapse/pkg/extension"
	"github.com/tunsuy/synapse/pkg/model"
)

func init() {
	extension.RegisterSource("skill-source", New)
}

// SkillSource 基于 AI Skill 的 Source 实现
// 它接收来自 AI 对话的原始内容，转换为标准 RawContent 格式
type SkillSource struct {
	defaultSource string
}

// New 创建一个新的 SkillSource 实例
func New(config map[string]any) (extension.Source, error) {
	source := "codebuddy"
	if s, ok := config["source"].(string); ok && s != "" {
		source = s
	}
	return &SkillSource{defaultSource: source}, nil
}

// Name 返回数据源名称
func (s *SkillSource) Name() string {
	return "skill-source"
}

// Fetch 获取原始内容
// 在 Skill 模式下，原始内容通过 FetchOptions.Config 传入
// 支持的 config 参数：
//   - "content": 原始对话内容（必须）
//   - "title": 内容标题（可选）
//   - "session_id": 会话标识（可选）
//   - "source": 数据来源覆盖（可选）
//   - "suggested_topics": 建议主题列表，逗号分隔（可选）
//   - "suggested_entities": 建议实体列表，逗号分隔（可选）
//   - "suggested_concepts": 建议概念列表，逗号分隔（可选）
//   - "key_points": 关键知识点列表，逗号分隔（可选）
func (s *SkillSource) Fetch(ctx context.Context, opts model.FetchOptions) ([]model.RawContent, error) {
	content, _ := opts.Config["content"].(string)
	if content == "" {
		return nil, fmt.Errorf("skill-source requires 'content' in fetch options config")
	}

	rc := model.RawContent{
		Content:   content,
		Source:    s.defaultSource,
		Timestamp: time.Now(),
	}

	// 标题
	if title, ok := opts.Config["title"].(string); ok && title != "" {
		rc.Title = title
	} else {
		rc.Title = extractTitle(content)
	}

	// 会话ID
	if sessionID, ok := opts.Config["session_id"].(string); ok {
		rc.SessionID = sessionID
	}

	// 数据来源覆盖
	if source, ok := opts.Config["source"].(string); ok && source != "" {
		rc.Source = source
	}

	// 建议的主题
	if topics, ok := opts.Config["suggested_topics"].(string); ok && topics != "" {
		rc.SuggestedTopics = splitAndTrim(topics)
	}

	// 建议的实体
	if entities, ok := opts.Config["suggested_entities"].(string); ok && entities != "" {
		rc.SuggestedEntities = splitAndTrim(entities)
	}

	// 建议的概念
	if concepts, ok := opts.Config["suggested_concepts"].(string); ok && concepts != "" {
		rc.SuggestedConcepts = splitAndTrim(concepts)
	}

	// 关键知识点
	if keyPoints, ok := opts.Config["key_points"].(string); ok && keyPoints != "" {
		rc.KeyPoints = splitAndTrim(keyPoints)
	}

	// 扩展元数据
	rc.Metadata = make(map[string]any)
	for k, v := range opts.Config {
		switch k {
		case "content", "title", "session_id", "source",
			"suggested_topics", "suggested_entities", "suggested_concepts", "key_points":
			continue
		default:
			rc.Metadata[k] = v
		}
	}

	return []model.RawContent{rc}, nil
}

// extractTitle 从内容中提取标题
// 策略：取第一行非空文本，最多 50 个字符
func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 去掉 Markdown 标题标记
		line = strings.TrimLeft(line, "# ")
		if len(line) > 50 {
			return line[:50] + "..."
		}
		return line
	}
	return "Untitled"
}

// splitAndTrim 按逗号分割字符串，并去除空格
func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
