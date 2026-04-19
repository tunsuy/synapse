# Synapse Roadmap

> **Synapse — 个人知识中枢（Personal Knowledge Hub）**
>
> 从各种 AI 助手对话中自动沉淀、整理、反哺知识，让你的每一次 AI 对话都成为知识复利。

---

## 🎯 愿景

打造一个**规范驱动、生态开放**的个人知识中枢，融合 Karpathy LLM Wiki 的**"知识编译"**思想：

> "Wiki 是持久的、复利增长的知识产物。" — Andrej Karpathy
>
> "Synapse 不是一个产品，而是一套规范 + 参考实现 + 多个独立扩展点。"

- **规范先行**：以 Knowledge Schema 为中心，定义统一的知识结构，所有扩展点共同遵守
- **扩展点模型**：六个独立扩展点（Source / Processor / Store / Indexer / Consumer / Auditor），每个可独立替换、社区共建
- **数据源可插拔**：从**任何数据源**（AI 对话、RSS、Notion、播客...）零摩擦地获取原始内容
- **处理引擎可插拔**：由**任意处理引擎**（AI Skill / MCP / 本地 LLM / 规则引擎...）将原始内容编译为结构化知识
- **存储为底座**：知识存在**用户选择的任何后端**（本地 / GitHub / S3 / WebDAV...），完全自主可控
- **检索引擎可插拔**：通过**任意检索引擎**（BM25 / 向量检索 / 图谱遍历...）找到相关知识
- **消费端可插拔**：知识可以被**任何消费端**输出（静态网站 / Obsidian / Anki / 邮件周报 / AI 反哺...）
- **质量审计**：借鉴 Karpathy 的 Lint 操作，可插拔的审计策略维护知识库健康
- **Schema 驱动**：通过 Schema 文件定义 AI 行为契约，修改 Schema 即修改所有 AI 助手的行为
- **插件生态**：完整的插件管理 CLI（install / uninstall / enable / disable / update / doctor），支持多来源安装（Go module / 本地目录 / Git 仓库），远期演进为插件市场

---

## 📅 里程碑总览

```
                2026 Q2                          2026 Q3                    2026 Q4+
├──────────────────────────────┼──────────────────────────────┼──────────────────────┤
│  M1: 基座搭建 ✅  M2: Skill  │  M3: MCP Server   M4: 多平台 │  M5: 消费端    M6+   │
│  (1周)          集成 ✅(2周) │  + 插件管理(3周)  适配(2周)   │  实现(3周)    社区   │
│                              │                              │              生态    │
│  ◆ Schema 规范              │  ◆ MCP Processor             │  ◆ Hugo 网站        │
│  ◆ 扩展点接口定义            │  ◆ GitHub Store              │  ◆ Obsidian 兼容    │
│  ◆ CLI init/check           │  ◆ BM25 Indexer              │  ◆ 知识图谱         │
│  ◆ 全局配置模板              │  ◆ Claude/Cursor Source      │  ◆ 插件市场         │
│  ◆ 第一个 Source            │  ◆ 默认 Auditor              │  ◆ Marketplace      │
│  ◆ 第一个 Processor         │  ◆ 插件管理 CLI              │  ◆ 社区共建         │
│  ◆ GitHub Store             │  ◆ 多来源安装                │                     │
└──────────────────────────────┴──────────────────────────────┴──────────────────────┘
```

---

## 🧬 生态化架构：扩展点模型（Extension Point Model）

> **核心思想**：Synapse = 知识规范（Schema）+ 存储底座（Store）+ 五个独立扩展点。不是固定分层的协议栈，而是以 Store 为中心、多个扩展点按需组合的星型架构。

### 架构全景

```
                        ┌─────────────┐
                        │   Source     │  数据源（原始内容从哪来）
                        │             │
                        │  AI 对话     │  可替换：AI 对话 / RSS / Notion / 播客 / ...
                        │  RSS        │
                        │  Notion     │
                        └──────┬──────┘
                               │ RawContent
                               ▼
                        ┌─────────────┐
                        │  Processor  │  处理引擎（原始内容 → 结构化知识）
                        │             │
                        │  Skill      │  可替换：Skill / MCP / LocalLLM / Rules / ...
                        │  MCP        │
                        │  Ollama     │
                        └──────┬──────┘
                               │ KnowledgeFile
                               ▼
    ┌──────────────────────────────────────────────────────────┐
    │                                                          │
    │                Store（存储底座）                            │
    │                                                          │
    │  可替换：Local FS / GitHub / S3 / WebDAV / IPFS / ...    │
    │                                                          │
    │  职责：知识文件的 CRUD + 可选版本控制                        │
    │                                                          │
    └────────┬──────────────────┬──────────────────┬───────────┘
             │                  │                  │
             ▼                  ▼                  ▼
    ┌─────────────┐    ┌─────────────┐    ┌───────────────┐
    │   Indexer    │    │   Auditor   │    │   Consumer    │
    │   检索引擎   │    │   质量审计   │    │   消费端       │
    │             │    │             │    │               │
    │  BM25       │    │  断链检测    │    │  Hugo 网站     │
    │  向量检索    │    │  孤儿页面    │    │  Obsidian     │
    │  图谱遍历    │    │  过时检测    │    │  Anki 闪卡     │
    │             │    │  去重        │    │  邮件周报       │
    └─────────────┘    └─────────────┘    │  AI 反哺(MCP)  │
                                          │  TUI 浏览      │
                                          └───────────────┘

    贯穿所有模块：Schema（知识规范）— 统一的"语言"
```

### 与"四层协议栈"的区别

上一版把架构设计成线性流水线（采集→编译→存储→展示），新版做了根本性调整：

| 维度 | 四层协议栈（旧） | 扩展点模型（新） |
|------|----------------|-----------------|
| **架构形态** | 线性流水线 | **星型**：Store 为底座，扩展点围绕 |
| **存储定位** | 数据流上的一个节点 | **底座**（所有模块都依赖它） |
| **采集+编译** | 两个固定层 | **两个扩展点**（Source + Processor），不强制分层 |
| **检索** | 塞在 Compiler 接口里 | **独立扩展点**（Indexer） |
| **审计** | 塞在 Compiler 接口里 | **独立扩展点**（Auditor） |
| **展示** | 只覆盖"人看" | **Consumer 扩展点**，覆盖所有消费形式 |
| **模块关系** | 层与层严格上下级 | **按需组合**，扩展点之间不强制依赖 |

### 五个扩展点详解

| 扩展点 | Go 接口 | 职责 | 官方参考实现 | 社区可共建 |
|--------|---------|------|-------------|-----------|
| **Source** | `Source` | 从外部获取原始内容 | AI Skill、CLI import | RSS / Notion / Twitter / 播客 / 微信... |
| **Processor** | `Processor` | 将原始内容处理为结构化知识 | Skill Processor、MCP Processor | 本地 LLM / 规则引擎 / 混合处理... |
| **Store** | `Store` | 知识文件的 CRUD + 可选版本控制 | Local Store、GitHub Store | S3 / WebDAV / Gitea / SQLite / IPFS... |
| **Indexer** | `Indexer` | 知识库检索 | BM25 Indexer | 向量检索 / 图谱遍历 / Elasticsearch... |
| **Consumer** | `Consumer` | 将知识输出为各种消费形式 | Hugo 网站、Obsidian 兼容 | VitePress / TUI / PDF / Anki / 邮件 / AI 反哺... |
| *Auditor*（可选） | `Auditor` | 知识库质量检查与修复 | Default Auditor | 自定义审计规则... |

### 为什么是"扩展点"而不是"层"？

1. **不强制分层**：Source 和 Processor 可以合并为一个实现（如 Skill 同时做采集和处理），也可以分开
2. **不强制线性**：Consumer 直接读 Store，不必等前面所有步骤完成
3. **按需组合**：最小集合只需要 Source + Processor + Store，Indexer / Auditor / Consumer 都是可选增强
4. **正交替换**：换 Indexer 不影响 Processor，换 Store 不影响 Source

### 插件通信方式

社区扩展点插件通过**子进程 + JSON-RPC**通信（类似 MCP stdio 模式），支持跨语言开发：

```
synapse CLI ←→ [JSON-RPC over stdin/stdout] ←→ 插件子进程（Go/Python/Node/Rust...）
```

### 配置示例

```yaml
# .synapse/config.yaml
synapse:
  version: "1.0"

  # 数据源（可同时启用多个）
  sources:
    - name: codebuddy-skill
      enabled: true
    - name: rss-reader
      enabled: true
      config:
        feeds: ["https://..."]

  # 处理引擎（选一个）
  processor:
    name: skill-processor

  # 存储底座（选一个）
  store:
    name: local-store
    config:
      path: "/Users/me/knowhub"

  # 检索引擎（可选）
  indexer:
    name: bm25-indexer

  # 消费端（可同时启用多个）
  consumers:
    - name: hugo-site
      config: { theme: synapse-default, deploy: github-pages }
    - name: obsidian-vault
      enabled: true

  # 审计（可选）
  auditor:
    name: default-auditor
```

---

## M1：基座搭建 🏗️

> **目标**：定义 Schema 规范 + 扩展点接口 + 初始化工具，让 knowhub 仓库"可用"，为生态化奠基
>
> **周期**：1 周
>
> **状态**：✅ 已完成

### 交付物

| # | 任务 | 说明 | 优先级 |
|---|------|------|--------|
| 1.1 | **Knowledge Schema 规范** | 定义 7 种页面类型的 Frontmatter 规范、`[[双向链接]]` 格式、标签分类体系 | P0 |
| 1.2 | **扩展点接口定义** | 用 Go Interface 定义 `Source` / `Processor` / `Store` / `Indexer` / `Consumer` / `Auditor` 六个扩展点 | P0 |
| 1.3 | **Schema 文件设计** | `.synapse/schema.yaml` — 知识库的"行为契约"，定义页面模板、工作流规则、质量标准 | P0 |
| 1.4 | **配置体系设计** | `.synapse/config.yaml` — 扩展点注册中心，声明各个扩展点使用哪个实现 | P0 |
| 1.5 | knowhub 仓库结构规范 | 定义 profile/topics/entities/concepts/inbox/journal/graph 各目录的规范和文件格式 | P0 |
| 1.6 | RawContent 格式规范 | 定义 Source 输出的标准原始内容格式，任何数据源都可以生成 | P0 |
| 1.7 | 用户画像（profile）格式规范 | 定义 me.md 的标准结构 | P0 |
| 1.8 | 实体页（entities）模板 | 定义人物、工具、项目、组织等实体页面模板 | P0 |
| 1.9 | 概念页（concepts）模板 | 定义技术概念、方法论、理论等概念页面模板 | P0 |
| 1.10 | 主题知识（topics）格式规范 | 定义主题目录下知识文件的标准格式，支持 `[[双向链接]]` | P0 |
| 1.11 | `synapse init` CLI 命令 | Go CLI 工具，初始化 knowhub 仓库（通过 Store 接口实现，支持幂等检查和 --force） | P0 |
| 1.11a | `synapse check` CLI 命令 | 验证全局配置文件的完整性和有效性（扩展点注册、环境变量检查等） | P0 |
| 1.11b | 全局配置模板自动创建 | 首次运行自动创建 `~/.synapse/config.yaml` 配置模板，引导用户配置扩展点 | P0 |
| 1.12 | knowhub 模板仓库 | 提供开箱即用的 knowhub 模板 | P1 |
| 1.13 | **扩展点开发指南** | 编写文档，说明如何开发自定义扩展点插件（为社区共建做准备） | P1 |
| 1.14 | README 文档 | synapse 和 knowhub 的使用说明 | P1 |

### 验收标准

- [x] 六个扩展点接口（Go Interface）定义完成，每个有清晰的输入/输出规范
- [x] 运行 `synapse init` 可以生成完整的 knowhub 目录结构（含 .synapse/schema.yaml + config.yaml）
- [x] 知识库结构规范文档完成，包含 7 种页面类型的模板定义
- [x] RawContent 格式规范清晰，第三方数据源可以据此接入
- [x] `[[双向链接]]` 格式兼容 Obsidian，可直接用 Obsidian 打开知识库
- [ ] 扩展点开发指南可指导社区开发者贡献新实现

---

## M2：Skill 集成 🧠

> **目标**：基于扩展点接口，实现第一个 Source（CodeBuddy Skill）、第一个 Processor（Skill Processor）和第一个 Store（Local Store），跑通**采集→处理→存储→反哺**闭环
>
> **周期**：2 周
>
> **状态**：✅ 已完成

### 交付物

| # | 任务 | 说明 | 优先级 |
|---|------|------|--------|
| 2.1 | **Source 参考实现：CodeBuddy Skill** | 实现 Source 接口，通过 Skill Prompt 驱动 AI 识别知识点、生成 RawContent | P0 |
| 2.2 | **Processor 参考实现：Skill Processor** | 实现 Processor 接口，通过 Skill Prompt 驱动 AI 将 RawContent 处理为实体/概念/主题/双向链接 | P0 |
| 2.3 | **Store 参考实现：Local Store** | 实现 Store 接口的本地文件系统版本（最简实现，为后续 GitHub Store 做基础） | P0 |
| 2.3a | **Store 参考实现：GitHub Store** | 实现 Store 接口的 GitHub 版本（提前完成，通过 GitHub API 实现知识库远程存储） | P0 |
| 2.4 | Skill — **反哺（Retrieve）** | AI 在对话中主动引用 knowhub 知识辅助回答，有价值的回答自动归档（知识自增长） | P0 |
| 2.5 | Skill — **审计（Audit）** | 用户说"检查知识库"时，AI 执行健康检查：孤儿页面、断链、过时内容、重复知识 | P1 |
| 2.6 | **Schema 加载机制** | Skill 启动时自动读取 `.synapse/schema.yaml`，让 AI 遵循知识库的行为契约 | P0 |
| 2.7 | `synapse install` CLI 命令 | 将 Skill 文件自动安装到目标 AI 助手的配置目录 | P1 |
| 2.8 | Skill 效果调优 | 实际使用 Skill 积累知识，根据效果迭代 Prompt 和 Schema | P1 |

### 验收标准

- [x] Source 接口有参考实现（CodeBuddy Skill），AI 能识别并采集知识到 inbox
- [x] Processor 接口有参考实现（Skill Processor），inbox 内容能被处理为实体/概念/主题/双向链接
- [x] Store 接口有参考实现（Local Store），知识文件可正确读写
- [x] Store 接口有第二个实现（GitHub Store），支持远程仓库存储
- [x] Skill Prompt 模板已实现（CodeBuddy / Claude Code / Cursor），AI 可通过 Skill 在对话中反哺知识
- [x] `synapse install` CLI 命令已实现，支持一键安装 Skill 到 AI 助手
- [x] 说"检查知识库"后，AI 能输出健康报告（审计）
- [x] 所有行为受 schema.yaml 约束，修改 schema 能改变 AI 行为

---

## M3：MCP Server 增强 + 插件管理 🔌

> **目标**：MCP Server 作为第二个 Processor + BM25 Indexer 作为第一个检索引擎 + **插件管理 CLI** 作为外部扩展基础设施，增强自动化能力并开启插件生态
>
> **注意**：GitHub Store 已在 M2 阶段提前实现
>
> **周期**：3 周
>
> **状态**：🔵 规划中

### 交付物

| # | 任务 | 说明 | 优先级 |
|---|------|------|--------|
| 3.1 | MCP Server 基础框架 | Go 实现的 MCP Server，支持 stdio 和 SSE 两种传输方式 | P0 |
| 3.2 | **Processor 实现：MCP Processor** | 通过 MCP 工具暴露核心操作：`collect` / `process` / `search` / `audit` | P0 |
| 3.3 | **Store 实现：GitHub Store** | 实现 Store 接口的 GitHub 版本，自动 git add/commit/push | P0 |
| 3.4 | **Indexer 实现：BM25 Indexer** | 实现 Indexer 接口，基于 BM25 的轻量级全文搜索 | P0 |
| 3.5 | **Auditor 实现：Default Auditor** | 实现 Auditor 接口，断链检测、孤儿页面、过时内容检查 | P1 |
| 3.6 | 采集类 MCP 工具 | `push_conversation`、`push_knowledge`：生成 RawContent 推入系统 | P0 |
| 3.7 | 反哺类 MCP 工具 | `get_profile`、`search_knowledge`、`get_topic`、`get_entity`、`get_concept` | P0 |
| 3.8 | 处理类 MCP 工具 | `process_inbox`（处理 inbox 内容）、`build_links`（构建双向链接） | P0 |
| 3.9 | 审计类 MCP 工具 | `audit_wiki`（健康检查）、`find_orphans`（孤儿页面）、`check_links`（断链检测） | P1 |
| 3.10 | MCP 配置自动生成 | `synapse mcp-config` 命令，生成各 AI 助手的 MCP 配置文件 | P1 |
| 3.11 | **插件管理 CLI** | `synapse plugin list/install/uninstall/enable/disable/update/doctor` 命令体系（参考 OpenClaw + Claude Code） | P0 |
| 3.12 | **多来源插件安装** | 支持三种安装来源：Go module（`synapse plugin install github.com/xxx/xxx`）、本地目录、Git 仓库 | P0 |
| 3.13 | **插件清单规范** | 设计 `synapse-plugin.yaml` 清单格式，声明插件元数据、扩展点类型、配置 Schema | P0 |
| 3.14 | **Layer 2 PluginAdapter** | 实现子进程 + JSON-RPC 通信的外部插件适配器，支持健康检查和能力协商 | P0 |
| 3.15 | **插件目录结构** | 安装路径 `~/.synapse/plugins/<name>/`，支持版本化缓存 | P1 |
| 3.16 | **双作用域配置** | 全局（`~/.synapse/config.yaml`）+ 项目（`.synapse/config.yaml`）两级配置，项目级优先 | P1 |
| 3.17 | **插件 Catalog** | JSON 格式的插件目录文件（`~/.synapse/plugins/catalog.json`），支持本地插件发现 | P1 |

### 验收标准

- [ ] MCP Server 可以在 Claude Code / CodeBuddy 中正常注册和调用
- [ ] AI 助手通过 MCP 可以完成知识的推送、查询、整理全流程
- [ ] 语义搜索能返回与查询语义相关的知识（而不仅仅是关键词匹配）
- [ ] `synapse plugin install` 可以从 Go module / 本地目录 / Git 仓库安装外部插件
- [ ] `synapse plugin list` 可以列出已安装的内置扩展和外部插件
- [ ] 外部插件通过 JSON-RPC 与核心通信，Engine 无需区分内置 vs 外部
- [ ] `synapse-plugin.yaml` 清单规范文档完成，社区开发者可据此开发插件

---

## M4：多平台 Source 适配 🌐

> **目标**：为更多 AI 平台实现 Source 扩展点，扩大数据源覆盖面
>
> **周期**：2 周
>
> **状态**：🔵 规划中

### 交付物

| # | 任务 | 说明 | 优先级 |
|---|------|------|--------|
| 4.1 | **Source：Claude Code** | CLAUDE.md + /commands，实现 Source 接口 | P0 |
| 4.2 | **Source：Cursor** | .cursorrules + Notepads，实现 Source 接口 | P0 |
| 4.3 | **Source：ChatGPT** | Custom Instructions / GPTs，实现 Source 接口 | P1 |
| 4.4 | **Source：Gemini** | Gems 方案，实现 Source 接口 | P2 |
| 4.5 | CLI import 命令 | `synapse import` 支持各平台的对话导出格式（通过适配器转换为 RawContent） | P1 |
| 4.6 | CLI export 命令 | `synapse export` 导出知识为指定格式 | P2 |
| 4.7 | **多 Source 编排** | 支持多个 Source 并行工作，RawContent 统一汇入处理流程 | P1 |

### 验收标准

- [ ] Claude Code、Cursor 中可以通过 Skill 完成知识采集和整理
- [ ] ChatGPT 自定义指令能引用 knowhub 知识辅助对话
- [ ] `synapse import --source chatgpt --file xxx.json` 可以正确导入

---

## M5：Consumer 实现 📊

> **目标**：实现第一批 Consumer（Hugo 网站 + Obsidian 兼容），提供知识可视化和多形态消费体验
>
> **周期**：3 周
>
> **状态**：🔵 规划中

### 交付物

| # | 任务 | 说明 | 优先级 |
|---|------|------|--------|
| 5.1 | **Consumer：Hugo 网站** | 实现 Consumer 接口，将 knowhub 渲染为静态网站 | P0 |
| 5.2 | GitHub Pages 部署 | GitHub Actions 自动构建并部署到 GitHub Pages | P0 |
| 5.3 | `[[双向链接]]` 渲染 | 将 `[[wiki-link]]` 自动渲染为可点击的超链接，支持反向链接展示 | P0 |
| 5.4 | 知识图谱可视化 | 从 `[[双向链接]]` 自动生成交互式知识关系图 | P1 |
| 5.5 | **Consumer：Obsidian 兼容** | 确保 knowhub 目录可直接作为 Obsidian Vault 打开，享受 Graph View | P1 |
| 5.6 | 时间线视图 | 按时间展示知识积累过程，可视化学习轨迹 | P1 |
| 5.7 | 搜索功能 | 静态网站内置全文搜索（复用 BM25 Indexer） | P1 |
| 5.8 | 健康仪表板 | 展示知识库健康分数、孤儿页面数、断链数等 Audit 指标 | P2 |
| 5.9 | 主题统计面板 | 展示各主题的知识量、活跃度、增长趋势 | P2 |

### 验收标准

- [ ] knowhub 仓库推送后自动生成可浏览的静态网站
- [ ] 知识图谱可以交互式地展示知识节点和关联关系
- [ ] 网站支持全文搜索

---

## M6+：社区生态共建 🚀

> **目标**：开放所有扩展点，构建完整的**插件市场（Marketplace）**，推动社区在每个扩展点贡献新实现，构建繁荣的个人知识管理生态
>
> **周期**：持续迭代
>
> **状态**：🔵 远期规划

### 插件市场（Marketplace）

借鉴 OpenClaw 和 Claude Code 的插件仓设计，Synapse 插件市场分三阶段演进：

| 阶段 | 能力 | 参考来源 |
|------|------|---------|
| **Phase 1：Catalog 目录** | JSON 格式的插件目录文件，用户可手动添加社区插件源；`synapse plugin search` 从 Catalog 搜索 | OpenClaw Channel Catalog |
| **Phase 2：Git 驱动的 Marketplace** | 官方 Marketplace 仓库（GitHub），Git 浅克隆 + 版本化缓存；`synapse plugin marketplace add/list/remove/update` 命令组 | Claude Code Marketplace |
| **Phase 3：在线 Marketplace UI** | Web 界面浏览和搜索插件，支持评分、下载统计、作者认证 | 远期规划 |

**插件市场核心机制（M6+ 逐步引入）**：

| 机制 | 说明 | 参考 |
|------|------|------|
| **Intent → State 分离** | config.yaml 声明想要的插件（Intent），`~/.synapse/plugins/cache/` 存放实际物化（State），Reconciler 负责同步 | Claude Code |
| **Reconciler 协调器** | 启动时自动对比 Intent vs State，执行安装/更新/清理 | Claude Code |
| **版本化缓存** | `~/.synapse/plugins/cache/<name>/<version>/`，旧版本延迟清理（标记后 7 天） | Claude Code |
| **Git 浅克隆** | `--depth 1`，支持 SSH/HTTPS 自动切换，Sparse Checkout 支持 monorepo | Claude Code |
| **安全防护** | checksum 校验 + 路径遍历检查 + 插件名防仿冒（保留名称集 + 正则检测） | 两者综合 |
| **社区插件提交** | PR 提交到官方 Marketplace 仓库的 `plugins/community.json` | OpenClaw |

### 方向

| # | 扩展点 | 方向 | 说明 | 优先级 |
|---|--------|------|------|--------|
| 6.1 | **Source** | 浏览器插件 | 一键采集网页版 AI 助手对话 | P1 |
| 6.2 | **Source** | RSS / Webhook | 自动采集 RSS 订阅和 Webhook 推送 | P2 |
| 6.3 | **Source** | Notion / Readwise | 同步笔记工具的内容 | P2 |
| 6.4 | **Processor** | Local LLM Processor | 集成 Ollama 等本地模型，离线处理知识 | P2 |
| 6.5 | **Processor** | GitHub Actions Processor | 定时触发知识处理 workflow | P1 |
| 6.6 | **Store** | Gitea / WebDAV | 自托管 Git 或坚果云等国内存储 | P2 |
| 6.7 | **Store** | S3 / IPFS | 云存储或去中心化存储 | P3 |
| 6.8 | **Indexer** | 向量检索 | 基于 Embedding 的语义搜索 | P2 |
| 6.9 | **Indexer** | 图谱遍历 | 基于知识图谱的关联检索 | P2 |
| 6.10 | **Consumer** | VitePress / MkDocs | 更多静态网站生成器 | P2 |
| 6.11 | **Consumer** | Anki 闪卡 | 知识导出为闪卡，间隔重复学习 | P2 |
| 6.12 | **Consumer** | Newsletter | 知识定期汇总为邮件周报 | P3 |
| 6.13 | **Consumer** | TUI 浏览器 | 终端知识浏览器 | P3 |
| 6.14 | **Consumer** | AI 反哺（MCP） | 将知识库暴露为 MCP Resource，反哺任意 AI 助手 | P1 |
| 6.15 | **Auditor** | 自定义审计规则 | 用户自定义知识质量标准和检查规则 | P3 |
| 6.16 | **跨扩展点** | 多知识库支持 | 支持工作/生活/学习等多个 knowhub 仓库 | P2 |
| 6.17 | **跨扩展点** | 团队知识库 | 支持团队共享知识库，多人协作积累 | P2 |
| 6.18 | **插件生态** | 插件市场 Phase 1 | Catalog 目录 + `synapse plugin search` | P1 |
| 6.19 | **插件生态** | 插件市场 Phase 2 | Git 驱动 Marketplace + Reconciler + 版本化缓存 | P2 |
| 6.20 | **插件生态** | 插件市场 Phase 3 | 在线 Marketplace UI + 评分/统计/认证 | P3 |
| 6.21 | **插件生态** | 插件开发 SDK | 提供 Go SDK + CLI 脚手架 `synapse plugin create`，降低社区贡献门槛 | P1 |
| 6.22 | **插件生态** | 插件热重载 | 配置变更后自动检测并重载插件，无需重启 | P2 |

---

## 🛤️ 关键路径

```
M1 基座搭建（1周）— Schema 规范 + 扩展点接口定义
 │
 ├──→ M2 Skill 集成（2周）— 第一个 Source + Processor + Store 实现
 │     │
 │     │   这是 MVP：规范 + 参考实现，跑通闭环
 │     │
 │     ├──→ M3 MCP 增强 + 插件管理（3周）— MCP Processor + GitHub Store + BM25 Indexer
 │     │     │                              + Auditor + 插件管理 CLI + Layer 2 PluginAdapter
 │     │     │
 │     │     └──→ M5 Consumer 实现（3周）— Hugo 网站 + Obsidian 兼容
 │     │
 │     └──→ M4 多平台 Source（2周）— Claude/Cursor/ChatGPT Source
 │
 └──→ M6+ 社区生态共建（持续）— 所有扩展点全面开放 + 插件市场（三阶段）
```

**MVP = M1 + M2**（约 3 周），已完成。可在 CodeBuddy 中实际使用 Synapse 积累知识。

**生态节奏**：每个里程碑产出新的扩展点参考实现，逐步验证接口设计，为社区共建铺路：

| 里程碑 | 首次实现的扩展点 | 验证的接口 | 插件生态进展 |
|--------|----------------|-----------|-------------|
| M1 ✅ | — | Schema 规范 + 六个接口定义 | — |
| M2 ✅ | Source / Processor / Store（含 GitHub Store） | 核心三件套跑通闭环 | — |
| M3 | Indexer / Auditor | 检索 + 审计能力就位 | **插件管理 CLI** + Layer 2 PluginAdapter + 多来源安装 |
| M4 | 更多 Source | 验证 Source 接口的通用性 | 外部 Source 可作为插件开发 |
| M5 | Consumer | 验证 Consumer 接口的通用性 | 外部 Consumer 可作为插件开发 |
| M6+ | 社区贡献 | 全部扩展点开放 | **Marketplace Phase 1→2→3** + 插件开发 SDK + 热重载 |

---

## 📐 技术选型

| 项 | 选择 | 理由 |
|----|------|------|
| **架构范式** | **扩展点模型（Extension Point Model）** | 星型架构，Store 为底座，六个独立扩展点按需组合 |
| 核心语言 | **Go** | 性能好，适合 CLI 和 MCP Server |
| **插件通信** | **子进程 + JSON-RPC（stdin/stdout）** | 类似 MCP stdio 模式，跨语言、安全、简单 |
| **插件管理** | **`synapse plugin *` CLI 命令体系** | 参考 OpenClaw（8 命令）+ Claude Code（11 命令），提供 list/install/uninstall/enable/disable/update/doctor |
| **插件安装来源** | **Go module + 本地目录 + Git 仓库** | 参考两者综合：OpenClaw 4 种来源 + Claude Code 7 种来源，取最实用的 3 种 |
| **插件市场** | **三阶段演进：Catalog → Git Marketplace → 在线 UI** | Phase 1 参考 OpenClaw JSON Catalog；Phase 2 参考 Claude Code Git 浅克隆 + Reconciler |
| **配置模型** | **Intent → State 分离 + 双作用域** | 参考 Claude Code Settings-First 架构，config.yaml 声明意图，`~/.synapse/plugins/cache/` 物化状态 |
| **知识规范** | **Knowledge Schema（.synapse/schema.yaml）** | 统一的知识结构定义，所有扩展点共同遵守的"语言" |
| **配置中心** | **.synapse/config.yaml** | 扩展点注册中心，声明各扩展点使用哪个实现 |
| **知识关联** | **`[[双向链接]]` + graph/relations.json** | 显式关联 > 隐式关联，兼容 Obsidian |
| Source（默认） | **Skill / Custom Instructions** | 零成本，覆盖面最广 |
| Processor（默认） | **Skill Processor（AI 助手侧处理）** | 用户零成本，不需要 API Key |
| Store（默认） | **Local Store（本地文件系统）** | 离线可用、最简、最快 |
| Store（增强） | **GitHub Store** | 版本控制、远程可用、天然 Wiki |
| Indexer（默认） | **BM25** | 不需要外部依赖，轻量高效 |
| Consumer（默认） | **Hugo + GitHub Pages** | 零成本静态网站 |
| Consumer（增强） | **Obsidian 兼容** | knowhub 可直接作为 Obsidian Vault 打开 |

---

## 📊 成功指标

### MVP 阶段（M1 + M2）

| 指标 | 目标 |
|------|------|
| 知识采集成功率 | 对话中产生的知识点 >80% 被识别和记录 |
| 整理准确率 | inbox 整理后归类到正确主题的比率 >90% |
| 使用摩擦 | 从对话产生知识到写入 knowhub，用户不需要额外操作 |
| 反哺命中率 | AI 在新对话中引用已有知识的场景 >50% 能命中 |

### 成长阶段（M3 + M4）

| 指标 | 目标 |
|------|------|
| 平台覆盖 | 支持 ≥3 个主流 AI 助手 |
| 知识库规模 | 单用户 topics 目录 >20 个主题文件 |
| MCP 响应速度 | 知识查询 <500ms |

### 成熟阶段（M5+）

| 指标 | 目标 |
|------|------|
| 日活用户 | - |
| 知识网站访问 | 用户每周至少浏览一次自己的知识库 |
| 社区贡献 | 有第三方贡献的 Skill 适配和知识模板 |

---

## ⚠️ 风险与应对

| 风险 | 可能性 | 影响 | 应对策略 |
|------|--------|------|---------|
| Skill/Prompt 处理质量不稳定 | 高 | 中 | 迭代优化 Prompt；MCP Processor 作为稳定性兜底 |
| AI 助手不支持读写本地文件 | 中 | 高 | MCP Server 方案兜底；CLI 手动操作 |
| 扩展点接口设计过早固化 | 中 | 高 | 先标记为 v0.x 实验版本，在 M1-M3 期间允许 Breaking Change |
| 插件生态冷启动难 | 中 | 中 | 官方提供足够的参考实现；开发指南要详细，降低贡献门槛；提供 `synapse plugin create` 脚手架 |
| 扩展点过度设计增加复杂度 | 中 | 中 | MVP 阶段每个接口只定义核心方法，可选能力通过接口组合渐进增加 |
| GitHub 仓库大小限制 | 低 | 中 | 知识压缩归纳，控制粒度；定期归档 |
| 各 AI 平台 Skill 机制变化 | 中 | 中 | 保持 Source 适配层薄，核心逻辑在 Processor 和 Store |
| 用户隐私担忧 | 中 | 高 | 默认 Local Store（数据不出本地）；支持敏感信息脱敏 |
| AI 平台自己做类似功能 | 中 | 低 | 差异化：开放扩展点 + 跨平台统一 + 用户拥有数据 + 社区生态 |
| 外部插件安全风险 | 中 | 高 | 子进程隔离 + checksum 校验 + 路径遍历检查 + 插件名防仿冒 |
| 插件市场维护成本 | 中 | 中 | 三阶段渐进：先 JSON Catalog（零运维），再 Git Marketplace，最后才考虑在线 UI |

---

## 📝 更新日志

| 日期 | 版本 | 变更内容 |
|------|------|---------|
| 2026-04-19 | v0.6 | M1/M2 标记为已完成；M2 增加 GitHub Store 交付物；更新里程碑总览和关键路径；M3 描述更正（GitHub Store 已提前完成） |
| 2026-04-18 | v0.5 | 插件生态增强：M3 新增插件管理 CLI + 多来源安装 + PluginAdapter；M6+ 新增三阶段 Marketplace 演进规划 + 插件开发 SDK + 热重载（基于 OpenClaw / Claude Code 插件仓调研） |
| 2026-04-18 | v0.4 | 架构再审视：从"四层协议栈"重构为"扩展点模型"，Store 为底座，五个独立扩展点按需组合 |
| 2026-04-18 | v0.3 | 生态化架构设计：四层协议模型（SCOL/SCP/SSP/SDP）、插件体系、社区共建规划 |
| 2026-04-18 | v0.2 | 融合 Karpathy LLM Wiki 思想：Schema 驱动、知识编译、四大操作、双向链接、Lint 审计 |
| 2026-04-18 | v0.1 | 初始 Roadmap，基于产品讨论整理 |
