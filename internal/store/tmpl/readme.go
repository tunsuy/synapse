// Package tmpl 提供 Store 初始化时使用的共享模板内容
// 所有 Store 实现（local、github 等）共用同一份 README、gitignore 等模板
package tmpl

import "fmt"

// GenerateReadme 生成知识库的 README.md 内容
func GenerateReadme(name string) string {
	if name == "" {
		name = "My"
	}
	return fmt.Sprintf(`# 🧠 %s's Knowledge Hub

> Powered by [Synapse](https://github.com/tunsuy/synapse) — 你的个人知识中枢

这是一个由 AI 助手自动策展的个人知识库。在日常对话中，AI 助手会帮你捕获有价值的知识片段，
自动整理归类，并通过双向链接构建知识图谱，让你的知识真正沉淀、可检索、可复用。

---

## 📁 目录结构

`+"```"+`
.
├── .synapse/
│   └── schema.yaml       # 知识库行为契约（定义页面类型、字段规范等）
├── profile/
│   └── me.md             # 用户画像（AI 助手参考此文件了解你）
├── topics/               # 📚 主题知识（按主题组织的深度内容）
├── entities/             # 🏷️  实体页（人物、工具、项目、组织）
├── concepts/             # 💡 概念页（技术概念、方法论、理论）
├── inbox/                # 📥 增量缓冲区（待整理的新知识）
├── journal/              # 📅 时间线（按时间记录的知识活动）
├── graph/
│   └── relations.json    # 🔗 知识图谱（从 [[wiki-links]] 自动生成）
├── .gitignore
└── README.md             # 本文件
`+"```"+`

### 各目录说明

| 目录 | 用途 | 示例 |
|------|------|------|
| `+"`profile/`"+` | 用户画像和偏好设置 | `+"`me.md`"+` — 你的技术栈、兴趣领域 |
| `+"`topics/`"+` | 按主题组织的知识 | `+"`go-concurrency.md`"+`、`+"`system-design.md`"+` |
| `+"`entities/`"+` | 具体的实体 | `+"`docker.md`"+`、`+"`kubernetes.md`"+`、`+"`team-xxx.md`"+` |
| `+"`concepts/`"+` | 抽象的概念 | `+"`cap-theorem.md`"+`、`+"`clean-architecture.md`"+` |
| `+"`inbox/`"+` | 新捕获的待整理知识 | AI 对话中提取的知识片段 |
| `+"`journal/`"+` | 按日期记录的知识活动 | `+"`2025-01-15.md`"+` |
| `+"`graph/`"+` | 知识关联图谱 | 自动从 `+"`[[wiki-links]]`"+` 生成 |

---

## 🔗 双向链接

本知识库使用 `+"`[[wiki-links]]`"+` 语法连接相关知识，例如：

`+"```markdown"+`
## Go 并发模型

Go 使用 [[goroutine]] 和 [[channel]] 实现并发...
参见 [[csp-model]] 了解理论基础。
`+"```"+`

### 🔮 兼容 Obsidian

本知识库完全兼容 [Obsidian](https://obsidian.md/)，直接用 Obsidian 打开本目录即可：
- 自动识别双向链接
- 可视化知识图谱
- 全文搜索

---

## 🚀 快速开始

### 1. 安装 AI Skill

`+"```bash"+`
# 为你使用的 AI 助手安装 Synapse Skill
synapse install codebuddy    # CodeBuddy
synapse install claude       # Claude Code
synapse install cursor       # Cursor
`+"```"+`

### 2. 编辑你的画像

编辑 `+"`profile/me.md`"+`，填写你的技术栈和兴趣领域，AI 助手会参考它来更好地理解你。

### 3. 开始积累知识

安装 Skill 后，AI 助手会在日常对话中自动帮你：
- **捕获** — 识别对话中有价值的知识片段
- **整理** — 归类到合适的目录，添加标签和链接
- **检索** — 在需要时从知识库中找到相关内容
- **审计** — 定期检查过时或低质量的知识

---

## 📄 知识页面格式

每个知识页面使用 Markdown + YAML Frontmatter：

`+"```markdown"+`
---
type: topic
title: "Go 并发模型"
created: 2025-01-15T10:30:00+08:00
updated: 2025-01-15T10:30:00+08:00
tags:
  - golang
  - concurrency
links:
  - "[[goroutine]]"
  - "[[channel]]"
confidence: 0.85
---

# Go 并发模型

正文内容...
`+"```"+`

---

## 🤖 Powered by Synapse

[Synapse](https://github.com/tunsuy/synapse) 是一个开源的个人知识管理工具，
让 AI 助手成为你的知识管理伙伴。

- 📖 文档: [github.com/tunsuy/synapse](https://github.com/tunsuy/synapse)
- 🐛 反馈: [Issues](https://github.com/tunsuy/synapse/issues)
`, name)
}
