// Package tmpl 提供 Store 初始化时使用的共享模板内容
// 所有 Store 实现（local、github 等）共用同一份 README、gitignore 等模板
package tmpl

import (
	"fmt"
	"strings"

	"github.com/tunsuy/synapse/internal/schema"
)

// GenerateReadme dynamically generates the knowledge hub README.md content based on schema.
// Directory structure, page type tables, etc. are all derived from schema.PageTypes
// to ensure consistency with the schema definition.
func GenerateReadme(name string, s *schema.Schema) string {
	if name == "" {
		name = "My"
	}
	if s == nil {
		s = schema.Default()
	}

	var b strings.Builder

	// ── Header ──
	b.WriteString(fmt.Sprintf("# 🧠 %s's Knowledge Hub\n\n", name))
	b.WriteString("> Powered by [Synapse](https://github.com/tunsuy/synapse) — Your personal knowledge hub\n\n")
	b.WriteString("This is a personal knowledge base automatically curated by AI assistants. During daily conversations,\n")
	b.WriteString("AI assistants capture valuable knowledge snippets, organize and categorize them, and build a knowledge\n")
	b.WriteString("graph through bidirectional links — making your knowledge truly accumulate, searchable, and reusable.\n\n")
	b.WriteString("---\n\n")

	// ── Directory structure (dynamically generated from schema.PageTypes) ──
	b.WriteString("## 📁 Directory Structure\n\n")
	b.WriteString("```\n")
	b.WriteString(".\n")
	b.WriteString("├── .synapse/\n")
	b.WriteString("│   └── schema.yaml       # Knowledge hub behavior contract (page types, field specs, etc.)\n")

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

		if pt.Name == "profile" {
			b.WriteString(fmt.Sprintf("%s %s/\n", prefix, dir))
			b.WriteString(fmt.Sprintf("│   └── me.md             %s\n", comment))
		} else if pt.Name == "graph" {
			b.WriteString(fmt.Sprintf("%s %s/\n", prefix, dir))
			b.WriteString(fmt.Sprintf("│   └── relations.json    %s\n", comment))
		} else {
			b.WriteString(fmt.Sprintf("%s %s/", prefix, dir))
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
	b.WriteString("└── README.md             # This file\n")
	b.WriteString("```\n\n")

	// ── Directory description table (dynamically generated from schema.PageTypes) ──
	b.WriteString("### Directory Descriptions\n\n")
	b.WriteString("| Directory | Purpose | Example |\n")
	b.WriteString("|-----------|---------|--------|\n")
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

	// ── Frontmatter field descriptions (dynamically generated from schema.Frontmatter) ──
	b.WriteString("## 📄 Knowledge Page Format\n\n")
	b.WriteString("Each knowledge page uses Markdown + YAML Frontmatter:\n\n")

	b.WriteString("**Required fields**: ")
	for i, f := range s.Frontmatter.Required {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("`%s`", f))
	}
	b.WriteString("\n\n")

	if len(s.Frontmatter.Optional) > 0 {
		b.WriteString("**Optional fields**: ")
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
	b.WriteString("title: \"Go Concurrency Model\"\n")
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
	b.WriteString("# Go Concurrency Model\n\n")
	b.WriteString("Content goes here...\n")
	b.WriteString("```\n\n")
	b.WriteString("---\n\n")

	// ── Bidirectional links ──
	b.WriteString("## 🔗 Bidirectional Links\n\n")
	b.WriteString(fmt.Sprintf("This knowledge hub uses `%s` syntax to connect related knowledge, for example:\n\n", s.LinkFormat))
	b.WriteString("```markdown\n")
	b.WriteString("## Go Concurrency Model\n\n")
	b.WriteString("Go uses [[goroutine]] and [[channel]] for concurrency...\n")
	b.WriteString("See [[csp-model]] for theoretical foundations.\n")
	b.WriteString("```\n\n")

	b.WriteString("### 🔮 Obsidian Compatible\n\n")
	b.WriteString("This knowledge hub is fully compatible with [Obsidian](https://obsidian.md/). Just open this directory in Obsidian:\n")
	b.WriteString("- Auto-detect bidirectional links\n")
	b.WriteString("- Visual knowledge graph\n")
	b.WriteString("- Full-text search\n\n")
	b.WriteString("---\n\n")

	// ── Quick start ──
	b.WriteString("## 🚀 Quick Start\n\n")
	b.WriteString("### 1. Install AI Skill\n\n")
	b.WriteString("```bash\n")
	b.WriteString("# Install Synapse Skill for your AI assistant\n")
	b.WriteString("synapse install codebuddy    # CodeBuddy\n")
	b.WriteString("synapse install claude       # Claude Code\n")
	b.WriteString("synapse install cursor       # Cursor\n")
	b.WriteString("```\n\n")

	b.WriteString("### 2. Edit Your Profile\n\n")
	b.WriteString("Edit `profile/me.md` with your tech stack and interests. AI assistants will reference it to better understand you.\n\n")

	b.WriteString("### 3. Start Accumulating Knowledge\n\n")
	b.WriteString("After installing the Skill, AI assistants will automatically help you during daily conversations:\n")
	b.WriteString("- **Capture** — Identify valuable knowledge snippets from conversations\n")
	b.WriteString("- **Compile** — Categorize into appropriate directories with tags and links\n")
	b.WriteString("- **Retrieve** — Find relevant content from your knowledge base when needed\n")
	b.WriteString("- **Audit** — Periodically check for outdated or low-quality knowledge\n\n")
	b.WriteString("---\n\n")

	// ── Footer ──
	b.WriteString("## 🤖 Powered by Synapse\n\n")
	b.WriteString("[Synapse](https://github.com/tunsuy/synapse) is an open-source personal knowledge management tool\n")
	b.WriteString("that turns your AI assistant into a knowledge management partner.\n\n")
	b.WriteString("- 📖 Docs: [github.com/tunsuy/synapse](https://github.com/tunsuy/synapse)\n")
	b.WriteString("- 🐛 Feedback: [Issues](https://github.com/tunsuy/synapse/issues)\n")

	return b.String()
}
