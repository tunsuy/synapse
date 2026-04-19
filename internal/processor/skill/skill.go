// Package skill 实现基于规则的 Skill Processor
// 这是 M2 阶段的参考实现，将 RawContent 处理为结构化知识文件
// 不依赖外部 AI API，使用基于规则的方式进行知识分类和结构化
package skill

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/tunsuy/synapse/pkg/extension"
	"github.com/tunsuy/synapse/pkg/model"
)

func init() {
	extension.RegisterProcessor("skill-processor", New)
}

// SkillProcessor 基于规则的 Skill Processor 实现
// 将 RawContent 处理为结构化知识文件
type SkillProcessor struct {
	defaultConfidence float64
}

// New 创建一个新的 SkillProcessor 实例
func New(config map[string]any) (extension.Processor, error) {
	confidence := 0.7
	if c, ok := config["default_confidence"].(float64); ok && c > 0 {
		confidence = c
	}
	return &SkillProcessor{defaultConfidence: confidence}, nil
}

// Name 返回处理引擎名称
func (p *SkillProcessor) Name() string {
	return "skill-processor"
}

// Process 将原始内容转换为知识文件
// 处理策略：
// 1. 如果 RawContent 有 SuggestedTopics，为每个主题生成 topic 类型的知识文件
// 2. 如果有 SuggestedEntities，为每个实体生成 entity 类型的知识文件
// 3. 如果有 SuggestedConcepts，为每个概念生成 concept 类型的知识文件
// 4. 如果都没有，则将内容放入 inbox 待整理
func (p *SkillProcessor) Process(ctx context.Context, raw []model.RawContent) ([]model.KnowledgeFile, error) {
	var files []model.KnowledgeFile

	for _, rc := range raw {
		if rc.Content == "" {
			continue
		}

		generated := p.processOne(rc)
		files = append(files, generated...)
	}

	return files, nil
}

// processOne 处理单条原始内容
func (p *SkillProcessor) processOne(rc model.RawContent) []model.KnowledgeFile {
	var files []model.KnowledgeFile
	now := time.Now()

	hasClassification := false

	// 生成主题知识文件
	for _, topic := range rc.SuggestedTopics {
		kf := p.buildTopicFile(rc, topic, now)
		files = append(files, kf)
		hasClassification = true
	}

	// 生成实体知识文件
	for _, entity := range rc.SuggestedEntities {
		kf := p.buildEntityFile(rc, entity, now)
		files = append(files, kf)
		hasClassification = true
	}

	// 生成概念知识文件
	for _, concept := range rc.SuggestedConcepts {
		kf := p.buildConceptFile(rc, concept, now)
		files = append(files, kf)
		hasClassification = true
	}

	// 如果没有分类建议，放入 inbox
	if !hasClassification {
		kf := p.buildInboxFile(rc, now)
		files = append(files, kf)
	}

	return files
}

// buildTopicFile 构建主题知识文件
func (p *SkillProcessor) buildTopicFile(rc model.RawContent, topic string, now time.Time) model.KnowledgeFile {
	slug := toSlug(topic)
	path := filepath.Join("topics", slug+".md")

	// 收集与此主题相关的链接
	links := collectLinks(rc, topic)

	body := buildBody(rc, topic, "topic")

	return model.KnowledgeFile{
		Path: path,
		Frontmatter: model.Frontmatter{
			Type:       model.PageTypeTopic,
			Title:      topic,
			Created:    now,
			Updated:    now,
			Tags:       buildTags(rc, topic),
			Links:      links,
			Source:     rc.Source,
			Confidence: p.defaultConfidence,
			Status:     "active",
		},
		Body: body,
	}
}

// buildEntityFile 构建实体知识文件
func (p *SkillProcessor) buildEntityFile(rc model.RawContent, entity string, now time.Time) model.KnowledgeFile {
	slug := toSlug(entity)
	path := filepath.Join("entities", slug+".md")

	links := collectLinks(rc, entity)

	body := buildBody(rc, entity, "entity")

	return model.KnowledgeFile{
		Path: path,
		Frontmatter: model.Frontmatter{
			Type:       model.PageTypeEntity,
			Title:      entity,
			Created:    now,
			Updated:    now,
			Tags:       buildTags(rc, entity),
			Links:      links,
			Source:     rc.Source,
			Confidence: p.defaultConfidence,
			Category:   guessCategory(entity, rc.Content),
			Status:     "active",
		},
		Body: body,
	}
}

// buildConceptFile 构建概念知识文件
func (p *SkillProcessor) buildConceptFile(rc model.RawContent, concept string, now time.Time) model.KnowledgeFile {
	slug := toSlug(concept)
	path := filepath.Join("concepts", slug+".md")

	links := collectLinks(rc, concept)

	body := buildBody(rc, concept, "concept")

	return model.KnowledgeFile{
		Path: path,
		Frontmatter: model.Frontmatter{
			Type:       model.PageTypeConcept,
			Title:      concept,
			Created:    now,
			Updated:    now,
			Tags:       buildTags(rc, concept),
			Links:      links,
			Source:     rc.Source,
			Confidence: p.defaultConfidence,
			Status:     "active",
		},
		Body: body,
	}
}

// buildInboxFile 构建待整理的 inbox 文件
func (p *SkillProcessor) buildInboxFile(rc model.RawContent, now time.Time) model.KnowledgeFile {
	slug := toSlug(rc.Title)
	if slug == "" {
		slug = fmt.Sprintf("inbox-%s", now.Format("20060102-150405"))
	}
	path := filepath.Join("inbox", slug+".md")

	var body strings.Builder
	body.WriteString(fmt.Sprintf("# %s\n\n", rc.Title))

	if len(rc.KeyPoints) > 0 {
		body.WriteString("## 关键知识点\n\n")
		for _, kp := range rc.KeyPoints {
			body.WriteString(fmt.Sprintf("- %s\n", kp))
		}
		body.WriteString("\n")
	}

	body.WriteString("## 原始内容\n\n")
	body.WriteString(rc.Content)

	return model.KnowledgeFile{
		Path: path,
		Frontmatter: model.Frontmatter{
			Type:       model.PageTypeInbox,
			Title:      rc.Title,
			Created:    now,
			Updated:    now,
			Source:     rc.Source,
			Confidence: p.defaultConfidence,
			Status:     "draft",
		},
		Body: body.String(),
	}
}

// toSlug 将标题转换为 URL 友好的文件名
func toSlug(title string) string {
	s := strings.ToLower(title)
	s = strings.TrimSpace(s)

	// 替换空格和特殊字符为连字符
	replacer := strings.NewReplacer(
		" ", "-",
		"/", "-",
		"\\", "-",
		".", "-",
		"_", "-",
		":", "-",
		"(", "",
		")", "",
		"[", "",
		"]", "",
		"{", "",
		"}", "",
		"'", "",
		"\"", "",
		"?", "",
		"!", "",
		"&", "and",
		"#", "",
	)
	s = replacer.Replace(s)

	// 合并连续的连字符
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}

	// 去掉首尾连字符
	s = strings.Trim(s, "-")

	return s
}

// buildTags 构建标签列表
func buildTags(rc model.RawContent, primaryTag string) []string {
	tagSet := make(map[string]bool)
	tagSet[strings.ToLower(primaryTag)] = true

	// 添加来源作为标签
	if rc.Source != "" {
		tagSet[rc.Source] = true
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags
}

// collectLinks 收集相关的双向链接目标
func collectLinks(rc model.RawContent, exclude string) []string {
	linkSet := make(map[string]bool)
	excludeLower := strings.ToLower(exclude)

	// 将其他建议的主题/实体/概念都作为链接目标
	for _, t := range rc.SuggestedTopics {
		if strings.ToLower(t) != excludeLower {
			linkSet[t] = true
		}
	}
	for _, e := range rc.SuggestedEntities {
		if strings.ToLower(e) != excludeLower {
			linkSet[e] = true
		}
	}
	for _, c := range rc.SuggestedConcepts {
		if strings.ToLower(c) != excludeLower {
			linkSet[c] = true
		}
	}

	links := make([]string, 0, len(linkSet))
	for link := range linkSet {
		links = append(links, link)
	}
	return links
}

// buildBody 构建知识文件的 Markdown 正文
func buildBody(rc model.RawContent, title, pageType string) string {
	var body strings.Builder

	body.WriteString(fmt.Sprintf("# %s\n\n", title))

	// 关键知识点
	if len(rc.KeyPoints) > 0 {
		body.WriteString("## 关键知识点\n\n")
		for _, kp := range rc.KeyPoints {
			body.WriteString(fmt.Sprintf("- %s\n", kp))
		}
		body.WriteString("\n")
	}

	// 相关链接
	allLinks := collectLinks(rc, title)
	if len(allLinks) > 0 {
		body.WriteString("## 相关链接\n\n")
		for _, link := range allLinks {
			body.WriteString(fmt.Sprintf("- [[%s]]\n", link))
		}
		body.WriteString("\n")
	}

	// 原始内容
	body.WriteString("## 内容\n\n")
	body.WriteString(rc.Content)

	return body.String()
}

// guessCategory 猜测实体的分类
func guessCategory(entity, content string) string {
	lower := strings.ToLower(entity + " " + content)

	toolKeywords := []string{"tool", "framework", "library", "sdk", "ide", "editor", "platform"}
	for _, kw := range toolKeywords {
		if strings.Contains(lower, kw) {
			return "tool"
		}
	}

	personKeywords := []string{"person", "author", "creator", "founder", "developer"}
	for _, kw := range personKeywords {
		if strings.Contains(lower, kw) {
			return "person"
		}
	}

	orgKeywords := []string{"company", "organization", "team", "group", "community"}
	for _, kw := range orgKeywords {
		if strings.Contains(lower, kw) {
			return "organization"
		}
	}

	projectKeywords := []string{"project", "repo", "repository", "app", "application", "service"}
	for _, kw := range projectKeywords {
		if strings.Contains(lower, kw) {
			return "project"
		}
	}

	return ""
}
