// Package tmpl 提供 Store 初始化时使用的共享模板内容
// 所有 Store 实现（local、github 等）共用同一份 README、gitignore 等模板
package tmpl

import (
	"fmt"
	"strings"

	"github.com/tunsuy/synapse/internal/schema"
)

// GenerateReadme 根据 schema 动态生成知识库的 README.md 内容
// 目录结构、页面类型表格等均从 schema.PageTypes 派生，保证与 schema 定义一致
func GenerateReadme(name string, s *schema.Schema) string {
	if name == "" {
		name = "My"
	}
	if s == nil {
		s = schema.Default()
	}

	var b strings.Builder

	// ── 头部 ──
	b.WriteString(fmt.Sprintf("# 🧠 %s's Knowledge Hub\n\n", name))
	b.WriteString("> Powered by [Synapse](https://github.com/tunsuy/synapse) — 你的个人知识中枢\n\n")
	b.WriteString("这是一个由 AI 助手自动策展的个人知识库。在日常对话中，AI 助手会帮你捕获有价值的知识片段，\n")
	b.WriteString("自动整理归类，并通过双向链接构建知识图谱，让你的知识真正沉淀、可检索、可复用。\n\n")
	b.WriteString("---\n\n")

	// ── 目录结构（从 schema.PageTypes 动态生成）──
	b.WriteString("## 📁 目录结构\n\n")
	b.WriteString("```\n")
	b.WriteString(".\n")
	b.WriteString("├── .synapse/\n")
	b.WriteString("│   └── schema.yaml       # 知识库行为契约（定义页面类型、字段规范等）\n")

	for i, pt := range s.PageTypes {
		dir := strings.TrimSuffix(pt.Directory, "/")
		emoji := pt.Emoji
		if emoji != "" {
			emoji = " " + emoji
		}
		comment := fmt.Sprintf("#%s %s", emoji, pt.Description)

		isLast := i == len(s.PageTypes)-1
		prefix := "├──"
		if isLast {
			prefix = "├──"
		}

		// 某些目录有子文件需要特殊展示
		if pt.Name == "profile" {
			b.WriteString(fmt.Sprintf("%s %s/\n", prefix, dir))
			b.WriteString(fmt.Sprintf("│   └── me.md             %s\n", comment))
		} else if pt.Name == "graph" {
			b.WriteString(fmt.Sprintf("%s %s/\n", prefix, dir))
			b.WriteString(fmt.Sprintf("│   └── relations.json    %s\n", comment))
		} else {
			b.WriteString(fmt.Sprintf("%s %s/", prefix, dir))
			// 用空格对齐注释
			padding := 22 - len(dir)
			if padding < 1 {
				padding = 1
			}
			b.WriteString(strings.Repeat(" ", padding))
			b.WriteString(comment)
			b.WriteString("\n")
		}
	}

	b.WriteString("├── .gitignore\n")
	b.WriteString("└── README.md             # 本文件\n")
	b.WriteString("```\n\n")

	// ── 各目录说明表格（从 schema.PageTypes 动态生成）──
	b.WriteString("### 各目录说明\n\n")
	b.WriteString("| 目录 | 用途 | 示例 |\n")
	b.WriteString("|------|------|------|\n")
	for _, pt := range s.PageTypes {
		dir := strings.TrimSuffix(pt.Directory, "/") + "/"
		example := ""
		if pt.Example != "" {
			example = fmt.Sprintf("`%s`", pt.Example)
		} else {
			example = "—"
		}
		b.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", dir, pt.Description, example))
	}
	b.WriteString("\n")
	b.WriteString("---\n\n")

	// ── Frontmatter 字段说明（从 schema.Frontmatter 动态生成）──
	b.WriteString("## 📄 知识页面格式\n\n")
	b.WriteString("每个知识页面使用 Markdown + YAML Frontmatter：\n\n")

	b.WriteString("**必填字段**: ")
	for i, f := range s.Frontmatter.Required {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("`%s`", f))
	}
	b.WriteString("\n\n")

	if len(s.Frontmatter.Optional) > 0 {
		b.WriteString("**可选字段**: ")
		for i, f := range s.Frontmatter.Optional {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("`%s`", f))
		}
		b.WriteString("\n\n")
	}

	b.WriteString("```markdown\n")
	b.WriteString("---\n")
	b.WriteString("type: topic\n")
	b.WriteString("title: \"Go 并发模型\"\n")
	b.WriteString("created: 2025-01-15T10:30:00+08:00\n")
	b.WriteString("updated: 2025-01-15T10:30:00+08:00\n")
	b.WriteString("tags:\n")
	b.WriteString("  - golang\n")
	b.WriteString("  - concurrency\n")
	b.WriteString("links:\n")
	b.WriteString("  - \"[[goroutine]]\"\n")
	b.WriteString("  - \"[[channel]]\"\n")
	b.WriteString("confidence: 0.85\n")
	b.WriteString("---\n\n")
	b.WriteString("# Go 并发模型\n\n")
	b.WriteString("正文内容...\n")
	b.WriteString("```\n\n")
	b.WriteString("---\n\n")

	// ── 双向链接 ──
	b.WriteString("## 🔗 双向链接\n\n")
	b.WriteString(fmt.Sprintf("本知识库使用 `%s` 语法连接相关知识，例如：\n\n", s.LinkFormat))
	b.WriteString("```markdown\n")
	b.WriteString("## Go 并发模型\n\n")
	b.WriteString("Go 使用 [[goroutine]] 和 [[channel]] 实现并发...\n")
	b.WriteString("参见 [[csp-model]] 了解理论基础。\n")
	b.WriteString("```\n\n")

	b.WriteString("### 🔮 兼容 Obsidian\n\n")
	b.WriteString("本知识库完全兼容 [Obsidian](https://obsidian.md/)，直接用 Obsidian 打开本目录即可：\n")
	b.WriteString("- 自动识别双向链接\n")
	b.WriteString("- 可视化知识图谱\n")
	b.WriteString("- 全文搜索\n\n")
	b.WriteString("---\n\n")

	// ── 快速开始 ──
	b.WriteString("## 🚀 快速开始\n\n")
	b.WriteString("### 1. 安装 AI Skill\n\n")
	b.WriteString("```bash\n")
	b.WriteString("# 为你使用的 AI 助手安装 Synapse Skill\n")
	b.WriteString("synapse install codebuddy    # CodeBuddy\n")
	b.WriteString("synapse install claude       # Claude Code\n")
	b.WriteString("synapse install cursor       # Cursor\n")
	b.WriteString("```\n\n")

	b.WriteString("### 2. 编辑你的画像\n\n")
	b.WriteString("编辑 `profile/me.md`，填写你的技术栈和兴趣领域，AI 助手会参考它来更好地理解你。\n\n")

	b.WriteString("### 3. 开始积累知识\n\n")
	b.WriteString("安装 Skill 后，AI 助手会在日常对话中自动帮你：\n")
	b.WriteString("- **捕获** — 识别对话中有价值的知识片段\n")
	b.WriteString("- **整理** — 归类到合适的目录，添加标签和链接\n")
	b.WriteString("- **检索** — 在需要时从知识库中找到相关内容\n")
	b.WriteString("- **审计** — 定期检查过时或低质量的知识\n\n")
	b.WriteString("---\n\n")

	// ── 尾部 ──
	b.WriteString("## 🤖 Powered by Synapse\n\n")
	b.WriteString("[Synapse](https://github.com/tunsuy/synapse) 是一个开源的个人知识管理工具，\n")
	b.WriteString("让 AI 助手成为你的知识管理伙伴。\n\n")
	b.WriteString("- 📖 文档: [github.com/tunsuy/synapse](https://github.com/tunsuy/synapse)\n")
	b.WriteString("- 🐛 反馈: [Issues](https://github.com/tunsuy/synapse/issues)\n")

	return b.String()
}
