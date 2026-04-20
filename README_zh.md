<p align="center">
  <img src="assets/logo.png" alt="Synapse Logo" width="200" />
</p>

<h1 align="center">Synapse</h1>

<p align="center">
  <strong>个人知识中枢（Personal Knowledge Hub）</strong><br/>
  从各种 AI 助手对话中自动沉淀、整理、反哺知识，让你的每一次 AI 对话都成为知识复利。
</p>

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.24-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**语言: 简体中文 | [English](README.md) | [日本語](README_ja.md) | [한국어](README_ko.md) | [Français](README_fr.md) | [Español](README_es.md)**

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

- Go >= 1.24

### 安装

```bash
go install github.com/tunsuy/synapse@latest
```

首次运行任意 synapse 命令时，会自动在 `~/.synapse/config.yaml` 创建全局配置模板：

```bash
# 触发自动创建配置模板
synapse --version

# 输出：
# 📝 Created global config template: /Users/you/.synapse/config.yaml
#    Please edit this file to configure your store and extensions.
#    Then run 'synapse check' to verify your configuration.
```

### Step 1：配置扩展点

编辑全局配置文件 `~/.synapse/config.yaml`，选择你的存储后端和其他扩展点。

#### 方案 A：本地文件系统存储（推荐新手使用）

```yaml
synapse:
  version: "1.0"

  sources:
    - name: "skill-source"
      enabled: true

  processor:
    name: "skill-processor"

  # 本地存储
  store:
    name: "local-store"
    config:
      path: "~/knowhub"        # 知识库本地路径
```

#### 方案 B：GitHub 仓库存储（适合云端同步）

```yaml
synapse:
  version: "1.0"

  sources:
    - name: "skill-source"
      enabled: true

  processor:
    name: "skill-processor"

  # GitHub 存储
  store:
    name: "github-store"
    config:
      owner: "${GITHUB_OWNER}"   # 你的 GitHub 用户名
      repo: "${GITHUB_REPO}"     # 知识库仓库名
      token: "${GITHUB_TOKEN}"   # GitHub Personal Access Token
      branch: "main"
```

> 💡 **提示**：使用 `${ENV_VAR}` 格式引用环境变量，避免在配置文件中硬编码敏感信息。

### Step 2：验证配置

```bash
synapse check
```

输出示例：

```
🔍 Checking Synapse configuration...
   Config: /Users/you/.synapse/config.yaml

   ✅ Config file exists
   ✅ Config file is valid YAML
   ✅ Version: 1.0
   ✅ Store: local-store
   ✅ Store "local-store" is registered
   ✅ Source: skill-source (registered)
   ✅ Processor: skill-processor (registered)

✅ Configuration is valid! You can now run 'synapse init' to initialize your knowledge base.
```

`check` 命令会检查以下内容：

| 检查项 | 说明 |
|--------|------|
| 配置文件存在性 | `~/.synapse/config.yaml` 是否存在 |
| YAML 合法性 | 文件是否为有效 YAML |
| 必填字段 | `synapse.version`、`synapse.store.name` 是否填写 |
| 扩展点注册 | 配置的 Store/Source/Processor 是否已注册到 Registry |
| 环境变量 | `${ENV_VAR}` 占位符是否已设置对应环境变量 |

### Step 3：初始化知识库

```bash
# 使用全局配置初始化
synapse init

# 指定知识库拥有者名称
synapse init --name "你的名字"

# 使用指定配置文件
synapse init --config /path/to/config.yaml

# 强制重新初始化（不会删除已有数据）
synapse init --force
```

`init` 命令会根据配置中指定的 Store 后端自动执行初始化：

| Store | 初始化行为 |
|-------|-----------|
| `local-store` | 在本地创建知识库目录结构和模板文件 |
| `github-store` | 通过 GitHub API 在仓库中创建知识库骨架文件 |

初始化完成后的知识库目录结构：

```
knowhub/
├── .synapse/
│   └── schema.yaml       # 知识规范（行为契约）
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

> ⚠️ **幂等性**：如果知识库已经初始化过，`init` 命令会提示并跳过，使用 `--force` 可以强制重新初始化。

### Step 4：安装 Skill 到 AI 助手

Skill 是一段预置的 Prompt 指令，安装后 AI 助手会在对话中自动帮你采集、整理、反哺知识。

```bash
# 安装到 CodeBuddy（推荐）
synapse install codebuddy

# 安装到 Claude Code
synapse install claude --target /path/to/project

# 安装到 Cursor
synapse install cursor

# 查看所有支持的 AI 助手
synapse install --list
```

安装后，你在 AI 助手中对话时可以使用以下触发词：

| 你说的话 | AI 会做什么 |
|---------|-----------|
| "记一下" / "保存到知识库" | 立即采集当前对话中的知识 |
| "检查知识库" / "审计" | 执行知识库健康检查 |
| "我知道什么关于 X" | 从知识库中检索相关内容 |
| "整理 inbox" | 帮你整理待处理内容 |

### Step 5：日常使用

#### 手动采集知识

除了通过 Skill 自动采集，你也可以用 CLI 手动采集：

```bash
# 直接传入内容
synapse collect --content "Go接口是隐式实现的" --title "Go Interfaces" \
  --topics "Go" --concepts "Duck Typing"

# 通过管道输入
echo "学习笔记内容..." | synapse collect --topics "分布式系统" --entities "Raft"
```

#### 搜索知识库

```bash
# 关键词搜索
synapse search goroutine

# 按类型过滤
synapse search --type topic "并发模型"

# 限制返回数量
synapse search --limit 5 golang
```

#### 审计知识库

```bash
synapse audit
```

审计报告包括：

| 检查项 | 说明 |
|--------|------|
| 健康评分 | 综合评分（满分 100） |
| Frontmatter 完整性 | 标题、类型等必填字段是否缺失 |
| 断链检测 | `[[双向链接]]` 指向的页面是否存在 |
| 孤儿页面 | 没有被任何其他页面链接到的页面 |
| 知识库统计 | 文件数量、链接数量、按类型分布 |

#### 管理扩展点插件

```bash
# 查看所有已注册的扩展点
synapse plugin list
```

---

## 📖 命令参考

| 命令 | 说明 | 示例 |
|------|------|------|
| `synapse init` | 初始化知识库 | `synapse init --name "张三"` |
| `synapse check` | 检查配置有效性 | `synapse check` |
| `synapse collect` | 采集知识 | `synapse collect --content "..." --topics "Go"` |
| `synapse search` | 搜索知识库 | `synapse search goroutine` |
| `synapse audit` | 审计知识库健康状态 | `synapse audit` |
| `synapse install` | 安装 Skill 到 AI 助手 | `synapse install codebuddy` |
| `synapse plugin list` | 查看已注册插件 | `synapse plugin list` |

### 全局参数

| 参数 | 说明 |
|------|------|
| `--config`, `-c` | 指定配置文件路径（默认 `~/.synapse/config.yaml`） |
| `--version`, `-v` | 显示版本号 |
| `--help`, `-h` | 显示帮助信息 |

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

Synapse 通过插件系统扩展功能。M3 阶段将全面支持插件管理：

```bash
# 查看已注册的扩展点插件（已支持）
synapse plugin list

# 以下命令将在 M3 阶段实现：
# synapse plugin install github.com/example/synapse-rss-source  # Go module 安装
# synapse plugin install --git https://github.com/example/xxx.git  # Git 仓库安装
# synapse plugin install --local ./my-custom-processor  # 本地目录安装
# synapse plugin enable rss-source   # 启用插件
# synapse plugin disable rss-source  # 禁用插件
# synapse plugin doctor              # 检查插件健康状态
```

---

## 📅 路线图

| 里程碑 | 内容 | 状态 |
|--------|------|------|
| **M1 基座搭建** | Schema 规范 + 扩展点接口 + CLI init/check | ✅ 已完成 |
| **M2 Skill 集成** | 第一个 Source + Processor + Store，跑通闭环 | ✅ 已完成 |
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
- 🔐 阅读 [安全策略](SECURITY.md) 了解漏洞报告流程

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
