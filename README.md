<p align="center">
  <img src="assets/logo.png" alt="Synapse Logo" width="200" />
</p>

<h1 align="center">Synapse</h1>

<p align="center">
  <strong>个人知识中枢（Personal Knowledge Hub）</strong><br/>
  从各种 AI 助手对话中自动沉淀、整理、反哺知识，让你的每一次 AI 对话都成为知识复利。
</p>

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**语言: 简体中文 | [English](README_en.md) | [日本語](README_ja.md) | [한국어](README_ko.md) | [Français](README_fr.md) | [Español](README_es.md)**

---

## 🎯 为什么需要 Synapse？

我们在日常工作和学习中使用各种 AI 助手（ChatGPT、Claude、CodeBuddy、Gemini 等），每一次对话本质上都是知识积累。但现实是：

- **知识碎片化** — 散落在各个 AI 助手中，难以回顾
- **AI 认知割裂** — AI 对你的了解是碎片化的，每次对话都像第一次
- **知识是"暗资产"** — 大量有价值的对话产物，用完就遗忘了

**Synapse 的目标**：让你与 AI 的每次对话都变成可沉淀、可检索、可反哺的知识资产。

> "Wiki 是持久的、复利增长的知识产物。" — Andrej Karpathy

---

## ✨ 核心特性

- 🔌 **扩展点模型** — 六大独立扩展点（Source / Processor / Store / Indexer / Consumer / Auditor），按需组合、独立替换
- 📥 **多源采集** — 支持从任意 AI 助手、RSS、Notion、播客等数据源零摩擦获取内容
- 🧠 **智能处理** — AI 驱动的知识提取、分类、关联，自动将原始对话编译为结构化知识
- 💾 **存储自主** — 数据存在你选择的任何后端（本地 / GitHub / S3 / WebDAV），完全自主可控
- 🔍 **灵活检索** — 可插拔的检索引擎（BM25 / 向量检索 / 图谱遍历）
- 📊 **多形态消费** — 知识可输出为静态网站、Obsidian Vault、Anki 闪卡、邮件周报等
- 🔗 **双向链接** — `[[wiki-link]]` 格式，兼容 Obsidian，构建个人知识图谱
- 📋 **Schema 驱动** — 通过 Schema 文件定义 AI 行为契约，修改 Schema 即修改所有 AI 助手的行为
- 🧩 **插件生态** — 完整的插件管理 CLI，支持多来源安装，社区可贡献任意扩展点实现

---

## 🏗️ 架构概览

Synapse 采用 **扩展点模型（Extension Point Model）**，以 Store 为底座，六个独立扩展点按需组合的星型架构：

```
                    ┌─────────────┐
                    │   Source     │  数据源（AI 对话 / RSS / Notion / ...）
                    └──────┬──────┘
                           │ RawContent
                           ▼
                    ┌─────────────┐
                    │  Processor  │  处理引擎（Skill / MCP / LocalLLM / ...）
                    └──────┬──────┘
                           │ KnowledgeFile
                           ▼
┌──────────────────────────────────────────────────────┐
│                  Store（存储底座）                      │
│        Local FS / GitHub / S3 / WebDAV / ...         │
└────────┬──────────────────┬──────────────────┬───────┘
         │                  │                  │
         ▼                  ▼                  ▼
  ┌─────────────┐   ┌─────────────┐   ┌───────────────┐
  │   Indexer    │   │   Auditor   │   │   Consumer    │
  │   检索引擎   │   │   质量审计   │   │   消费端      │
  └─────────────┘   └─────────────┘   └───────────────┘
```

> 详细架构说明请参阅 [ARCHITECTURE.md](ARCHITECTURE.md)。

---

## 🚀 快速开始

### 环境要求

- Go >= 1.21

### 安装

```bash
go install github.com/tunsuy/synapse@latest
```

### 初始化知识库

```bash
# 初始化一个新的知识库
synapse init ~/knowhub

# 查看知识库结构
tree ~/knowhub
```

### 知识库目录结构

```
knowhub/
├── .synapse/
│   ├── schema.yaml       # 知识规范（行为契约）
│   └── config.yaml       # 扩展点配置
├── profile/
│   └── me.md             # 用户画像
├── topics/               # 主题知识
│   ├── golang/
│   ├── architecture/
│   └── ...
├── entities/             # 实体页（人物、工具、项目）
├── concepts/             # 概念页（技术概念、方法论）
├── inbox/                # 待整理内容
├── journal/              # 时间线日志
└── graph/
    └── relations.json    # 知识关联图谱
```

---

## 🔌 扩展点

| 扩展点 | 职责 | 默认实现 | 社区可贡献 |
|--------|------|---------|-----------|
| **Source** | 从外部获取原始内容 | CodeBuddy Skill | RSS / Notion / Twitter / 播客 / 微信... |
| **Processor** | 原始内容 → 结构化知识 | Skill Processor | 本地 LLM / 规则引擎 / 混合处理... |
| **Store** | 知识文件的 CRUD + 版本控制 | Local Store | GitHub / S3 / WebDAV / SQLite / IPFS... |
| **Indexer** | 知识库检索 | BM25 Indexer | 向量检索 / 图谱遍历 / Elasticsearch... |
| **Consumer** | 知识输出为各种消费形式 | Hugo 网站 | VitePress / Anki / 邮件 / TUI... |
| **Auditor** | 知识库质量检查与修复 | Default Auditor | 自定义审计规则... |

---

## 🧩 插件管理

```bash
# 查看已安装插件
synapse plugin list

# 从 Go module 安装插件
synapse plugin install github.com/example/synapse-rss-source

# 从 Git 仓库安装
synapse plugin install --git https://github.com/example/synapse-vector-indexer.git

# 从本地目录安装
synapse plugin install --local ./my-custom-processor

# 启用 / 禁用插件
synapse plugin enable rss-source
synapse plugin disable rss-source

# 检查插件健康状态
synapse plugin doctor
```

---

## 📅 路线图

| 里程碑 | 内容 | 状态 |
|--------|------|------|
| **M1 基座搭建** | Schema 规范 + 扩展点接口 + CLI init | 🟡 待开始 |
| **M2 Skill 集成** | 第一个 Source + Processor + Store，跑通闭环 | 🟡 待开始 |
| **M3 MCP + 插件管理** | MCP Server + GitHub Store + BM25 Indexer + 插件 CLI | 🔵 规划中 |
| **M4 多平台适配** | Claude Code / Cursor / ChatGPT Source | 🔵 规划中 |
| **M5 Consumer 实现** | Hugo 网站 + Obsidian 兼容 + 知识图谱 | 🔵 规划中 |
| **M6+ 社区生态** | 插件市场 + 全扩展点开放 + 社区共建 | 🔵 远期规划 |

> 详细路线图请参阅 [docs/roadmap.md](docs/roadmap.md)。

---

## 🤝 参与贡献

我们欢迎任何形式的贡献！无论是提交 Bug 报告、提出新功能建议，还是直接贡献代码。

- 📖 阅读 [贡献指南](CONTRIBUTING.md) 了解如何参与
- 🏛️ 阅读 [架构说明](ARCHITECTURE.md) 了解技术设计
- 📋 阅读 [行为准则](CODE_OF_CONDUCT.md) 了解社区规范
- 🗺️ 阅读 [路线图](docs/roadmap.md) 了解项目规划

### 贡献方向

每个扩展点都欢迎社区贡献新的实现：

- 🔌 **Source 插件**：接入更多数据源（RSS、Notion、微信、播客...）
- ⚙️ **Processor 插件**：支持更多处理引擎（本地 LLM、规则引擎...）
- 💾 **Store 插件**：支持更多存储后端（S3、WebDAV、IPFS...）
- 🔍 **Indexer 插件**：支持更多检索引擎（向量检索、图谱遍历...）
- 📊 **Consumer 插件**：支持更多消费形式（VitePress、Anki、TUI...）

---

## 📄 许可证

本项目基于 [Apache License 2.0](LICENSE) 开源。

---

## 💬 联系我们

- **Issues**：[GitHub Issues](https://github.com/tunsuy/synapse/issues)
- **Discussions**：[GitHub Discussions](https://github.com/tunsuy/synapse/discussions)

---

> *Synapse — 让每一次 AI 对话都成为知识复利。*
