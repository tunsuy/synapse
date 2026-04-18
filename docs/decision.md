# Synapse 方案决策记录

## 2026-04-17 产品构思与方案讨论

### 背景与问题

我们在日常生活中会使用各种 AI 助手（ChatGPT、Claude、CodeBuddy、Gemini 等），与 AI 助手的每一次交流本质上都是一次个人知识的积累过程。但当前存在以下问题：

1. **知识碎片化**：这些零散的知识散落在各个 AI 助手中
2. **难以系统回顾**：即使在同一个 AI 助手中，也很难系统性地看到曾经讨论过的知识
3. **AI 认知不足**：AI 助手对用户的认识不够系统，记忆是割裂的
4. **知识是"暗资产"**：与 AI 对话产生的大量知识是一次性消费品，用完就散落了

### 产品定位

**Synapse — 个人知识中枢（Personal Knowledge Hub）**

核心价值链：
```
各种AI助手对话 → 采集 → 整理归纳 → 知识Wiki → 反哺回AI助手
```

它能：
- 收集任何 AI 助手的交流过程
- 进行系统性的整理、归纳、展示
- 挂载到各个 AI 助手中，让 AI 助手了解用户的一切

### 设计原则

- **轻量化**：不使用任何云端数据库
- **用户拥有数据**：直接使用 GitHub 仓库存储
- **工具与数据解耦**：产品代码和用户知识库分离

### 仓库规划

| 仓库 | 用途 | 说明 |
|------|------|------|
| synapse | 产品代码仓库 | 工具本身的源代码 |
| knowhub | 用户知识库仓库 | 用户使用产品时提供的知识存储仓库（每个用户一个） |

### 架构设计

```
┌─────────────────────────────────────────────────┐
│                   采集层（Input）                  │
├─────────────────────────────────────────────────┤
│  ChatGPT   Claude   CodeBuddy   Gemini   ...    │
│     │         │         │          │             │
│     ▼         ▼         ▼          ▼             │
│  ┌─────────────────────────────────────┐        │
│  │  统一采集器（Collector）              │        │
│  │  - 浏览器插件（网页版AI助手）          │        │
│  │  - CLI hook（终端AI助手）             │        │
│  │  - API导出（支持导出的平台）           │        │
│  │  - 手动导入（复制粘贴/文件上传）       │        │
│  └──────────────┬──────────────────────┘        │
│                 ▼                                │
├─────────────────────────────────────────────────┤
│              处理层（Process）                     │
│  ┌─────────────────────────────────────┐        │
│  │  知识提取引擎（用AI做）               │        │
│  │  - 主题识别与分类                     │        │
│  │  - 关键知识点提取                     │        │
│  │  - 用户画像更新                       │        │
│  │  - 知识去重与合并                     │        │
│  │  - 知识关联与图谱构建                 │        │
│  └──────────────┬──────────────────────┘        │
│                 ▼                                │
├─────────────────────────────────────────────────┤
│             存储层（Storage）= GitHub Repo        │
│  ┌─────────────────────────────────────┐        │
│  │  knowhub/（用户知识库仓库）            │        │
│  │  ├── profile/                        │        │
│  │  │   └── me.md          # 用户画像    │        │
│  │  ├── topics/                         │        │
│  │  │   ├── golang/                     │        │
│  │  │   │   ├── _index.md  # 主题概览    │        │
│  │  │   │   ├── concurrency.md          │        │
│  │  │   │   └── error-handling.md       │        │
│  │  │   ├── architecture/               │        │
│  │  │   └── ai/                         │        │
│  │  ├── journal/           # 时间线      │        │
│  │  │   └── 2026-04-17.md               │        │
│  │  ├── graph/             # 知识图谱    │        │
│  │  │   └── relations.json              │        │
│  │  └── .mcp/              # MCP配置     │        │
│  │      └── knowledge-server.json       │        │
│  └─────────────────────────────────────┘        │
│                 ▼                                │
├─────────────────────────────────────────────────┤
│            输出层（Output）                        │
│  ┌────────┐ ┌────────┐ ┌──────────────┐        │
│  │GitHub   │ │MCP     │ │ 静态网站      │        │
│  │Pages    │ │Server  │ │ (Hugo/ViteP) │        │
│  │Wiki浏览 │ │反哺AI  │ │ 知识浏览      │        │
│  └────────┘ └────────┘ └──────────────┘        │
└─────────────────────────────────────────────────┘
```

### 使用 GitHub 存储的分析

**优势：**
- 免费、无限（公开仓库）/ 私有仓库也免费
- 天然版本控制，知识的每次变更都有记录
- Markdown 原生渲染，就是一个天然的 Wiki
- 不需要运维，零成本
- 可以被各种工具直接读取（MCP、CLI 等）
- 支持 GitHub Pages，可以直接生成静态网站浏览
- 支持 GitHub Actions，可以自动化处理流程

**劣势与应对：**

| 劣势 | 应对方案 |
|------|---------|
| 搜索能力弱 | 本地构建索引 + 向量检索，或利用 AI 做语义搜索 |
| 频繁提交会很嘈杂 | 批量合并提交，或用 orphan branch 专门存储 |
| 大文件/大量文件性能下降 | 知识做好压缩归纳，控制粒度 |
| 私密性 | 用私有仓库，敏感信息脱敏处理 |

### 核心模块

#### 1. 采集器（Collector）— 最关键模块

零摩擦地把对话内容收集进来：

| 优先级 | 方式 | 说明 |
|--------|------|------|
| P0 | 手动导入 | 支持各平台的对话导出格式（ChatGPT JSON、Claude JSON等） |
| P0 | CLI工具 | 命令行直接提交一段对话或知识 |
| P1 | 浏览器插件 | 一键采集当前AI对话页面 |
| P1 | MCP协议 | AI助手主动写入知识库 |
| P2 | Webhook | 接收各平台的导出回调 |

#### 2. 知识引擎（AI驱动）

用 AI 对原始对话做结构化处理：
- 主题识别与分类
- 关键知识点提取
- 用户画像更新
- 知识去重与合并
- 知识关联与图谱构建

#### 3. MCP Server（反哺AI助手）

暴露为 MCP Server，任何支持 MCP 的 AI 助手都能读取知识库：
- search_knowledge：语义搜索知识库
- get_profile：获取用户画像
- get_topic：获取某个主题的完整知识
- get_recent：获取最近学习的知识

### 技术选型

| 项 | 选择 | 理由 |
|----|------|------|
| 核心语言 | Go | 擅长、性能好、适合做CLI工具 |
| 存储 | GitHub Repo（Markdown + JSON） | 轻量、免费、版本控制 |
| 知识处理 | LLM API（OpenAI/Claude/本地模型） | AI驱动的知识提取 |
| 展示层 | GitHub Pages + Hugo/MkDocs | 零成本静态网站 |
| 分发 | MCP Server + CLI 工具 | 反哺AI助手 + 命令行交互 |

### MVP 路径

```
Phase 1（2周）- 最小可用：
  ├── CLI工具：手动导入对话 + AI提取知识 + 写入GitHub
  ├── 知识库结构：Markdown文件 + 目录分类
  └── 基本浏览：GitHub 自带的 Markdown 渲染

Phase 2（2周）- 反哺能力：
  ├── MCP Server：让AI助手能读取知识库
  ├── 用户画像：自动维护 profile
  └── 知识关联：简单的标签和引用关系

Phase 3（持续）- 体验提升：
  ├── 浏览器插件：一键采集
  ├── 静态网站：更好的知识浏览体验
  ├── 知识图谱可视化
  └── 自动化：GitHub Actions 定期整理
```

### 产品价值

1. **知识复利**：每次与AI对话都是在给自己的知识库做增量
2. **跨平台统一**：打破各AI平台的数据孤岛
3. **用户拥有数据**：数据存在用户自己的GitHub仓库，不被平台绑架
4. **AI认知增强**：通过MCP让所有AI助手都能系统性地了解用户

---

## 2026-04-18 方案深化讨论

### 讨论一：AI 整理能力放在哪一侧？

#### 问题背景

原方案中 Synapse 需要集成 AI（调用 LLM API）来整理用户知识库，但很多普通用户可能没有专门的大模型 API Key。

#### 三种方案对比

**方案 A：Synapse 侧集成 AI（原方案）**
```
AI助手 → 原始对话 → Synapse CLI → Synapse调用LLM API整理 → 写入knowhub
```
- 优点：整理逻辑集中、标准统一、知识质量可控
- 缺点：用户需要 API Key（成本门槛），Synapse 变重

**方案 B：AI 助手侧整理（用户提出）**
```
AI助手 ←→ knowhub（通过MCP拉取+推送）
AI助手负责：1.推送原始记录 2.拉取知识库 3.整理后推回
```
- 优点：零成本（复用 AI 助手算力），Synapse 极度轻量
- 缺点：整理标准不统一，依赖 AI 助手的 MCP 能力

**方案 C：混合方案（最终方案）**

把"整理"拆成两级：

| | 实时增量整理（AI助手侧） | 全局深度整理（可选） |
|---|---|---|
| **谁来做** | 用户正在使用的 AI 助手 | AI 助手 / Synapse CLI / GitHub Actions |
| **什么时候做** | 每次对话结束时 | 定期 / 手动触发 |
| **做什么** | 提取本次知识摘要，推到 inbox | 合并 inbox，去重，更新全局结构 |
| **算力成本** | 零（复用 AI 助手） | 零（让 AI 助手做）或低（自有 API Key） |

核心设计：
1. **inbox 缓冲区**：不同 AI 助手推送的增量包统一存放在 inbox，格式由协议约束
2. **全局整理也由 AI 助手完成**：用户在任意 AI 助手中说"帮我整理下知识库"即可，全程不需要 API Key
3. **渐进增强**：最基础只有 inbox 堆积 → 进阶由 AI 整理 → 高级有 API Key 自动整理

#### 知识库结构调整

新增 inbox 目录作为增量缓冲区：

```
knowhub/
├── profile/
│   └── me.md                    # 用户画像
├── topics/                      # 整理后的知识（结构化）
│   ├── golang/
│   ├── architecture/
│   └── ai/
├── inbox/                       # 新增：增量缓冲区
│   ├── 2026-04-18-chatgpt-01.md  # AI助手推送的知识增量包
│   ├── 2026-04-18-claude-01.md
│   └── 2026-04-18-codebuddy-01.md
├── journal/                     # 时间线
├── graph/                       # 知识图谱
└── .synapse/                    # Synapse配置
    ├── config.yaml              # 知识库配置
    └── mcp-server.json          # MCP服务配置
```

#### 增量包格式设计

```markdown
---
source: chatgpt
timestamp: 2026-04-18T08:30:00+08:00
session_id: abc123
suggested_topics: ["golang", "concurrency"]
profile_updates:
  - type: skill
    content: "熟悉Go的channel和select机制"
---

## 知识摘要

本次讨论了Go语言中goroutine泄漏的排查方法...

## 关键知识点

1. 使用 `runtime.NumGoroutine()` 监控goroutine数量
2. ...

## 原始对话摘要

用户询问了如何排查线上goroutine泄漏问题...
```

---

### 讨论二：集成方式分析

如何把 Synapse 集成到各个 AI 助手中，是产品能不能跑通的关键。

#### 方式一：MCP Server

```
AI助手 ←→ Synapse MCP Server ←→ knowhub（GitHub仓库）
```

- 覆盖范围：Claude Code、CodeBuddy、Cursor、Windsurf 等本地/IDE 类 AI 助手
- 局限：网页版 AI 助手（ChatGPT Web、Gemini Web）不支持 MCP
- 能力：双向交互，支持复杂操作

MCP Server 暴露的工具设计：

```
📥 采集类：
  - push_conversation: AI助手推送原始对话/知识增量到inbox
  - push_knowledge:    AI助手推送整理后的知识到topics

📤 反哺类：
  - get_profile:       获取用户画像
  - search_knowledge:  搜索知识库
  - get_topic:         获取某主题的知识
  - get_recent:        获取最近的知识动态

🔧 整理类：
  - get_inbox:         获取待整理的增量包列表
  - organize_inbox:    AI助手整理inbox（拉取→整理→推回）
  - update_profile:    更新用户画像
```

#### 方式二：浏览器插件

覆盖 MCP 触达不了的网页版 AI 助手（ChatGPT、Claude Web、Gemini）。

工作模式：
1. 被动采集：监听页面 DOM，提取对话内容
2. 主动注入：在用户发送消息前，自动追加相关知识库上下文

#### 方式三：Rules 文件

各 AI 助手支持的项目级配置文件（`.claude/CLAUDE.md`、`.codebuddy/rules.md`、`.cursorrules`）。

Synapse CLI 可以自动生成/更新这些配置文件，把用户画像和相关知识写入。

- 优点：极其简单，AI 助手原生支持
- 缺点：静态、不能实时交互、容量有限

#### 方式四：CLI 管道

最基础的兜底方案，手动导入导出。

```bash
synapse import --source chatgpt --file exported-conversations.json
synapse import --source claude-code --dir ~/.claude/conversations/
echo "今天学到了..." | synapse add --topic golang
```

#### 方式五：Skill / Custom Instructions（用户提出，关键突破）

**核心洞察**：Skill 本质是一段预置的 Prompt 指令，当用户触发时自动加载到 AI 上下文中。不需要启动任何服务，不需要 MCP Server，纯文本驱动。

各 AI 助手的对应机制：

| AI 助手 | 对应机制 | 形式 |
|---------|---------|------|
| CodeBuddy | Skill | `.codebuddy/skills/` 目录下的 prompt 文件 |
| Claude Code | CLAUDE.md + /commands | 项目级指令 + 自定义斜杠命令 |
| Cursor | .cursorrules + Notepads | 规则文件 + 笔记本 |
| ChatGPT | Custom Instructions / GPTs | 自定义指令 / 自定义 GPT |
| Gemini | Gems | 自定义角色 |

Skill 方式 vs MCP 对比：

| | MCP Server | Skill |
|---|---|---|
| 需要启动服务 | 是 | 否，纯文本 |
| 安装成本 | 中等 | 极低（复制一个文件） |
| AI助手覆盖面 | 支持MCP的助手 | 几乎所有AI助手 |
| 能力 | 强（自定义工具、逻辑） | 依赖AI助手自身能力 |
| 维护成本 | 需要维护Server代码 | 维护Prompt文本即可 |
| 可靠性 | 高（代码逻辑确定） | 中（AI可能理解偏差） |

**在 Skill 方案下 Synapse 变成了：一套 Skill/Prompt 模板 + 一个轻量 CLI**

```
synapse/
├── skills/                          # 各平台的 Skill 模板
│   ├── codebuddy/
│   │   └── synapse-knowledge.md     # CodeBuddy Skill
│   ├── claude-code/
│   │   └── synapse-commands.md      # Claude Code 自定义命令
│   ├── cursor/
│   │   └── .cursorrules             # Cursor Rules
│   └── chatgpt/
│       └── custom-instructions.md   # ChatGPT 自定义指令
├── cli/                             # 轻量CLI（可选）
│   └── synapse
│       ├── init                     # 初始化 knowhub 仓库结构
│       ├── install                  # 安装 Skill 到各 AI 助手
│       └── import                   # 手动导入对话
└── templates/                       # 知识库模板
    ├── profile/me.md
    ├── inbox/.gitkeep
    └── topics/.gitkeep
```

#### 最终集成策略（三层架构）

```
第一层：Skill（覆盖面最广，零成本）— MVP核心
  所有支持自定义指令的 AI 助手
  能力：采集、整理、反哺（基础版）
  原理：AI直接读写本地knowhub文件

第二层：MCP Server（可选增强）
  支持MCP的AI助手
  增强：语义搜索、自动触发、复杂操作
  原理：标准化工具接口

第三层：浏览器插件（后续扩展）
  网页版AI助手
  能力：对话采集、上下文注入
```

---

### MVP 路径（修订版）

```
MVP（1周）：
├── 1. 设计 knowhub 仓库标准结构
├── 2. 编写 Synapse Skill Prompt（先做 CodeBuddy 版）
├── 3. 写一个简单的 CLI：synapse init（初始化 knowhub 结构）
└── 4. 实际使用：在 CodeBuddy 中加载 Skill，开始积累知识

后续迭代：
├── 适配更多 AI 助手的 Skill 版本
├── MCP Server（增强搜索和自动化能力）
├── CLI 丰富（import、export、search）
└── 浏览器插件
```

### 技术选型（修订版）

| 项 | 选择 | 理由 |
|----|------|------|
| 核心产出 | Skill Prompt 模板 | 零成本集成，覆盖面最广 |
| 辅助工具 | Go CLI | 初始化、安装、导入等辅助操作 |
| 存储 | GitHub Repo（Markdown + JSON） | 轻量、免费、版本控制 |
| AI 整理 | 由 AI 助手自身完成 | 用户零成本，不需要 API Key |
| 展示层 | GitHub Pages + Hugo/MkDocs | 零成本静态网站 |
| 增强层 | MCP Server（可选） | 语义搜索、自动化 |

---

### Roadmap 制定

基于以上所有讨论，制定了完整的产品 Roadmap，详见 [roadmap.md](./roadmap.md)。

六大里程碑：

| 里程碑 | 内容 | 周期 |
|--------|------|------|
| M1 基座搭建 | knowhub 结构规范 + `synapse init` CLI | 1 周 |
| M2 Skill 集成 | CodeBuddy Skill（采集/整理/反哺） | 2 周 |
| M3 MCP Server 增强 | 标准化工具接口 + 语义搜索 | 3 周 |
| M4 多平台适配 | Claude Code / Cursor / ChatGPT Skill | 2 周 |
| M5 知识可视化 | 静态网站 + 知识图谱 + 时间线 | 3 周 |
| M6+ 生态扩展 | 浏览器插件、团队知识库、本地LLM等 | 持续 |

**MVP = M1 + M2**（约 3 周），完成即可在 CodeBuddy 中实际使用。

关键路径：M1 → M2（MVP） → M3/M4 并行 → M5 → M6+

---

## 2026-04-18 融合 Karpathy LLM Wiki 思想

### 背景

Andrej Karpathy 在 2025 年 4 月提出了 LLM Wiki 架构——一种完全不同于 RAG 的知识管理范式。核心思想是**"知识编译"**：

> "不要让 LLM 在查询时去理解原始文档，而是提前让 LLM 把文档'编译'成结构化的知识。"

我们之前已经在 `llm-wiki` 项目中实践了这套架构（三层架构、三大操作、Schema 契约），现在要把这些思想融入 Synapse。

### 核心对齐

| 维度 | Karpathy LLM Wiki | Synapse（原方案） | Synapse（融合后） |
|------|-------------------|------------------|------------------|
| 核心理念 | 知识编译（预处理 > 实时检索） | 知识复利（对话 → 积累） | **知识编译 + 知识复利** |
| LLM 角色 | Worker（执行 Ingest/Query/Lint） | 采集/整理/反哺 | **Worker + Personal Assistant** |
| 行为定义 | Schema 文件控制 | Skill Prompt 控制 | **Schema + Skill 双层驱动** |
| 存储 | Markdown + Git | Markdown + GitHub Repo | 一致 |
| 知识关联 | `[[双向链接]]` | graph/relations.json | **`[[双向链接]]` + 知识图谱** |
| 页面类型 | 5 类（Entity/Concept/Summary/Synthesis/Query） | 4 类（Profile/Topics/Inbox/Journal） | **7 类（融合）** |
| 健康维护 | Lint 操作 | 无 | **新增 Lint/Health Check** |
| 知识增长 | Query 归档 → Wiki 自增长 | 反哺是被动引用 | **主动编译 + 自我增长** |

### 关键融合点

#### 1. Schema 驱动（最重要的借鉴）

Karpathy 的核心洞察：**Schema 是人类和 LLM 之间的"行为契约"**。修改 Schema 就能修改 LLM 行为，不需要改代码。

**融合方案**：在 knowhub 仓库中引入 `.synapse/schema.md` 文件

```
knowhub/
├── .synapse/
│   ├── config.yaml        # 知识库配置
│   ├── schema.md          # 知识库行为契约（新增！）
│   └── mcp-server.json    # MCP 服务配置
```

`schema.md` 定义：
- 页面模板（每种知识类型的 Markdown 结构）
- 工作流规则（Ingest/Compile/Query/Lint 的步骤）
- 质量标准（页面质量、链接质量、标签分类）
- 冲突解决策略（新旧知识矛盾时的处理规则）

**价值**：任何 AI 助手加载 Skill 时，都会读取 schema.md 并遵守规范。用户只需修改 schema.md 就能调整所有 AI 助手的行为——这比改 Skill Prompt 更优雅。

#### 2. 知识编译流程

Karpathy 的 Ingest 流程：原始文档 → 摘要 → 实体提取 → 概念提取 → 双向链接 → 索引

**融合到 Synapse 的"整理"流程**：

```
inbox 增量包 → AI "编译"：
  1. 提取知识摘要 → 写入 topics/
  2. 提取实体（人物/工具/项目） → 写入 entities/（新增）
  3. 提取概念（技术概念/方法论） → 写入 concepts/（新增）
  4. 建立 [[双向链接]] → 更新 graph/
  5. 更新用户画像 → profile/me.md
  6. 更新索引 → topics/_index.md
  7. 记录日志 → journal/
```

这比原方案的"inbox → topics"更结构化，知识的颗粒度更细、关联更丰富。

#### 3. 页面类型丰富化

原方案 4 类 → 融合后 7 类：

| 类型 | 目录 | 来源 | 说明 |
|------|------|------|------|
| **Profile** | `profile/` | Synapse 原有 | 用户画像 |
| **Topics** | `topics/` | Synapse 原有 | 主题知识（类似 Synthesis） |
| **Entities** | `entities/` | Karpathy 借鉴 | 人物、工具、项目、组织 |
| **Concepts** | `concepts/` | Karpathy 借鉴 | 技术概念、方法论、理论 |
| **Inbox** | `inbox/` | Synapse 原有 | 增量缓冲区 |
| **Journal** | `journal/` | Synapse 原有 | 时间线 |
| **Graph** | `graph/` | 两者共有 | 知识关联图谱 |

#### 4. 新增 Lint 操作

借鉴 Karpathy 的 Lint 健康检查，为 Synapse 新增"知识库体检"能力：

```
用户说"帮我检查下知识库" → AI 执行 Lint：
- 孤儿页面：没有任何入链的知识文件
- 断链检测：[[链接]] 指向不存在的页面
- 过时内容：超过 N 天未更新的知识
- 重复知识：不同文件中的重复内容
- 画像偏差：profile 与实际知识积累不匹配
- 标签缺失：没有标签的知识文件
```

输出健康报告，并给出修复建议。

#### 5. 双向链接机制

Karpathy 的 `[[wiki-link]]` 方式比 Synapse 原方案的 JSON 关系图更自然：

```markdown
# Go 并发编程

## 核心概念
Go 的并发模型基于 [[CSP]] 理论，核心原语是 [[goroutine]] 和 [[channel]]。

## 常见问题
- [[goroutine-leak]]: 使用 `runtime.NumGoroutine()` 监控
- [[race-condition]]: 使用 `-race` flag 检测

## 相关实体
- [[Rob-Pike]]: Go 语言设计者之一
- [[Google]]: Go 语言的背后组织
```

**优势**：
- 知识关联内嵌在内容中，更自然
- 可以自动从 `[[链接]]` 解析生成 `graph/relations.json`
- 兼容 Obsidian，可以直接用 Obsidian 打开知识库浏览图谱

#### 6. 查询自增长

Karpathy 的 Query 操作会把有价值的问答自动归档为新的 Wiki 页面。

**融合方案**：Synapse 的"反哺"过程中，如果 AI 生成了有价值的综合性回答（比如跨多个主题的分析），可以自动归档到 topics/ 或新建一个 synthesis 文件。

```
用户问"Go 的错误处理和 Rust 的错误处理有什么区别？"
→ AI 从 knowhub 检索 Go 和 Rust 的知识
→ 生成综合对比分析
→ 这个回答本身就是有价值的知识
→ 自动归档到 topics/error-handling-comparison.md
→ Wiki 自增长 🔄
```

### 更新后的知识库结构

```
knowhub/
├── .synapse/
│   ├── config.yaml            # 知识库配置
│   ├── schema.md              # 行为契约（Karpathy Schema 思想）
│   └── mcp-server.json        # MCP 服务配置
├── profile/
│   └── me.md                  # 用户画像
├── topics/                    # 主题知识（深度整理后的知识）
│   ├── golang/
│   │   ├── _index.md
│   │   ├── concurrency.md     # 含 [[双向链接]]
│   │   └── error-handling.md
│   └── architecture/
├── entities/                  # 实体页（新增，借鉴 Karpathy）
│   ├── rob-pike.md
│   ├── google.md
│   └── kubernetes.md
├── concepts/                  # 概念页（新增，借鉴 Karpathy）
│   ├── csp.md
│   ├── goroutine.md
│   └── clean-architecture.md
├── inbox/                     # 增量缓冲区（原始对话的"编译"入口）
│   ├── 2026-04-18-codebuddy-01.md
│   └── 2026-04-18-claude-01.md
├── journal/                   # 时间线
│   └── 2026-04-18.md
└── graph/                     # 知识图谱
    └── relations.json         # 从 [[双向链接]] 自动生成
```

### 更新后的操作模型

从 Karpathy 的三大操作扩展为 Synapse 的四大操作：

| 操作 | Karpathy 对应 | 触发方式 | 说明 |
|------|-------------|---------|------|
| **Capture**（采集） | Ingest 的前半部分 | AI 对话中自动触发 | 识别知识点，写入 inbox 增量包 |
| **Compile**（编译） | Ingest 的后半部分 | 用户说"整理知识库" | 从 inbox 编译出实体、概念、主题知识、双向链接 |
| **Retrieve**（反哺） | Query | AI 对话中自动触发 | 检索知识库辅助回答，有价值的回答自动归档 |
| **Audit**（审计） | Lint | 定期 / 用户说"检查知识库" | 健康检查：孤儿页、断链、过时、重复、画像偏差 |

### 融合后的四字箴言

**采集 → 编译 → 反哺 → 审计**（Capture → Compile → Retrieve → Audit）

对比原方案的"采集 → 整理 → 反哺"，新增了：
1. "整理"升级为"编译"——更强调结构化处理，借鉴 Karpathy 的知识编译思想
2. 新增"审计"——借鉴 Karpathy 的 Lint 操作，持续维护知识库健康度

### 对 Roadmap 的影响

这次融合**不改变 Roadmap 的大结构**，但丰富了每个里程碑的内容：

| 里程碑 | 新增/变更内容 |
|--------|-------------|
| **M1** | knowhub 结构新增 entities/ concepts/ .synapse/schema.md |
| **M2** | Skill 拆为 4 个（采集/编译/反哺/审计），引入 Schema 驱动 |
| **M3** | MCP 工具对应四大操作，新增 lint 工具 |
| **M4** | 各平台 Skill 都需要遵循统一 schema.md |
| **M5** | 知识图谱从 `[[双向链接]]` 自动生成，兼容 Obsidian |

---

## 2026-04-18 生态化架构设计：分层协议 + 插件体系

### 背景

在融合了 Karpathy LLM Wiki 思想之后，我们进一步思考一个核心问题：**如何让 Synapse 从第一天就具备生态化的能力？**

关键洞察：**Synapse 不应该是一个"产品"，而应该是一套"协议 + 参考实现"。** 就像 HTTP 协议之于 Web 生态，LSP 协议之于 IDE 生态——真正有持久力的不是具体的实现，而是协议本身。

### 核心架构原则

#### 原则一：协议先行，实现可替换

每一层都先定义**协议/规范**，然后提供**参考实现**。第三方可以基于协议开发自己的实现。

#### 原则二：分层独立，正交解耦

各层之间只通过协议交互，层内实现可以独立替换，不影响其他层。

#### 原则三：渐进增强，核心极简

核心协议极简（一个 JSON Schema 就够），增强能力通过 Provider 插件渐进添加。

### 生态化架构：四层协议模型

```
┌───────────────────────────────────────────────────────────────────┐
│                    Synapse 四层协议架构                             │
│                                                                   │
│  ┌───────────────────────────────────────────────────────────┐    │
│  │  Layer 4: 展示层（Presentation Protocol）                   │    │
│  │                                                           │    │
│  │  协议：Synapse Display Protocol (SDP)                      │    │
│  │  定义：如何将知识库渲染为可浏览、可交互的形式                     │    │
│  │                                                           │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐    │    │
│  │  │ Hugo     │ │ Obsidian │ │ Web App  │ │ 社区共建  │    │    │
│  │  │ (静态站) │ │ (本地)   │ │ (动态)   │ │ ...      │    │    │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘    │    │
│  └──────────────────────────┬────────────────────────────────┘    │
│                             │ 读取                                │
│  ┌──────────────────────────┴────────────────────────────────┐    │
│  │  Layer 3: 存储层（Storage Protocol）                        │    │
│  │                                                           │    │
│  │  协议：Synapse Storage Protocol (SSP)                      │    │
│  │  定义：知识如何持久化存储，提供统一的 CRUD 接口                   │    │
│  │                                                           │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐    │    │
│  │  │ GitHub   │ │ Local FS │ │ S3/OSS   │ │ 社区共建  │    │    │
│  │  │ (Git)    │ │ (本地)   │ │ (云存储) │ │ ...      │    │    │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘    │    │
│  └──────────────────────────┬────────────────────────────────┘    │
│                             │ 存/取                               │
│  ┌──────────────────────────┴────────────────────────────────┐    │
│  │  Layer 2: 编译层（Compilation Protocol）                    │    │
│  │                                                           │    │
│  │  协议：Synapse Compilation Protocol (SCP)                  │    │
│  │  定义：原始知识如何被编译为结构化知识                             │    │
│  │  四大操作：Capture / Compile / Retrieve / Audit             │    │
│  │                                                           │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐    │    │
│  │  │ Skill    │ │ MCP      │ │ GitHub   │ │ 社区共建  │    │    │
│  │  │ Prompt   │ │ Server   │ │ Actions  │ │ ...      │    │    │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘    │    │
│  └──────────────────────────┬────────────────────────────────┘    │
│                             │ 输入                                │
│  ┌──────────────────────────┴────────────────────────────────┐    │
│  │  Layer 1: 采集层（Collection Protocol）                     │    │
│  │                                                           │    │
│  │  协议：Synapse Collection Protocol (SCOL)                  │    │
│  │  定义：各种数据源如何将原始知识送入系统                           │    │
│  │                                                           │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐    │    │
│  │  │ AI 对话  │ │ 浏览器   │ │ RSS/     │ │ 社区共建  │    │    │
│  │  │ (Skill)  │ │ 插件     │ │ Webhook  │ │ ...      │    │    │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘    │    │
│  └───────────────────────────────────────────────────────────┘    │
│                                                                   │
│  ┌───────────────────────────────────────────────────────────┐    │
│  │  Layer 0: 知识模型层（Knowledge Model）                     │    │
│  │                                                           │    │
│  │  协议：Synapse Knowledge Schema (SKS)                      │    │
│  │  定义：知识的标准数据结构（7 种页面类型 + Frontmatter 规范）     │    │
│  │  这是所有层的共同"语言"                                       │    │
│  └───────────────────────────────────────────────────────────┘    │
└───────────────────────────────────────────────────────────────────┘
```

### Layer 0：知识模型层（Synapse Knowledge Schema, SKS）

**这是整个生态的基石。** 所有层都基于统一的知识模型进行交互。

```yaml
# .synapse/schema.yaml — 知识模型定义
version: "1.0"

# 知识类型定义
page_types:
  - name: profile
    directory: profile/
    template: templates/profile.md
    description: "用户画像"
  - name: topic
    directory: topics/
    template: templates/topic.md
    description: "主题知识"
  - name: entity
    directory: entities/
    template: templates/entity.md
    description: "实体（人物/工具/项目/组织）"
  - name: concept
    directory: concepts/
    template: templates/concept.md
    description: "概念（技术概念/方法论/理论）"
  - name: inbox
    directory: inbox/
    template: templates/inbox.md
    description: "增量缓冲区"
  - name: journal
    directory: journal/
    template: templates/journal.md
    description: "时间线"
  - name: graph
    directory: graph/
    template: null
    description: "知识关联图谱"

# Frontmatter 标准字段
frontmatter:
  required: [type, title, created, updated]
  optional: [tags, links, source, confidence]

# 双向链接格式
link_format: "[[page-id]]"

# 操作定义
operations: [capture, compile, retrieve, audit]
```

**生态价值**：任何人只要遵循 SKS 规范，就可以构建自己的采集器、编译器、存储后端、展示前端。知识库本身是**可移植的**——从 GitHub 迁移到本地文件系统，或者从 Hugo 切换到 Obsidian，知识不变、协议不变。

### Layer 1：采集层协议（Synapse Collection Protocol, SCOL）

**定义：数据源如何将原始知识送入 Synapse 系统。**

#### 协议规范

采集层协议定义的是一个**标准的增量包格式（Increment Package）**。任何数据源只要能生成符合此格式的增量包，就可以成为 Synapse 的数据源。

```yaml
# SCOL 增量包标准格式
---
synapse_version: "1.0"
type: increment
source:
  provider: "codebuddy"       # 数据源提供者标识
  provider_version: "1.0.0"   # 提供者版本
  session_id: "abc123"        # 会话标识
  timestamp: "2026-04-18T08:30:00+08:00"

# 提取的知识元素（由数据源预处理）
knowledge:
  summary: "讨论了 Go 并发编程中的 goroutine 泄漏排查方法"
  key_points:
    - "使用 runtime.NumGoroutine() 监控"
    - "使用 pprof 定位泄漏点"
  suggested_topics: ["golang", "concurrency"]
  suggested_entities: ["Go", "pprof"]
  suggested_concepts: ["goroutine", "goroutine-leak"]
  profile_updates:
    - type: skill
      content: "熟悉 Go 的 goroutine 泄漏排查"

# 原始内容（可选，用于全局编译时深度处理）
raw_content: |
  用户: 如何排查 goroutine 泄漏？
  AI: ...
---
```

#### Provider 接口（Go 接口定义）

```go
// Provider 是采集层的插件接口
// 任何数据源只要实现此接口，就可以接入 Synapse 生态
type Provider interface {
    // Name 返回 Provider 的唯一标识
    Name() string

    // Collect 从数据源采集知识，返回标准增量包
    Collect(ctx context.Context, opts CollectOptions) ([]*IncrementPackage, error)

    // Watch 持续监听数据源的新知识（可选，用于实时采集）
    Watch(ctx context.Context, handler IncrementHandler) error
}
```

#### 可能的 Provider 实现

| Provider | 数据源 | 实现方式 | 谁来做 |
|----------|--------|---------|--------|
| `codebuddy-skill` | CodeBuddy 对话 | Skill Prompt 驱动 AI 生成增量包 | Synapse 官方 |
| `claude-skill` | Claude Code 对话 | CLAUDE.md 驱动 | Synapse 官方 |
| `browser-extension` | 网页版 AI 助手 | 浏览器插件截取 DOM | Synapse 官方 |
| `chatgpt-export` | ChatGPT 导出 JSON | CLI 解析导出文件 | Synapse 官方 |
| `rss-reader` | RSS 订阅 | 定期拉取并生成增量包 | **社区共建** |
| `notion-sync` | Notion 笔记 | Notion API 同步 | **社区共建** |
| `readwise-sync` | Readwise 高亮 | Readwise API 同步 | **社区共建** |
| `twitter-bookmark` | Twitter 收藏 | Twitter API | **社区共建** |
| `podcast-transcript` | 播客转录 | Whisper 转文字后生成增量包 | **社区共建** |
| `github-star` | GitHub Star 项目 | GitHub API | **社区共建** |
| `wechat-article` | 微信公众号文章 | 微信 API / 剪藏工具 | **社区共建** |

**生态价值**：社区不需要理解 Synapse 的全部架构，只需要**实现一个 Provider 接口 + 生成标准增量包**，就可以把任何数据源接入 Synapse 生态。这是最容易产生社区贡献的一层。

### Layer 2：编译层协议（Synapse Compilation Protocol, SCP）

**定义：原始知识（增量包）如何被编译为结构化的 Wiki 知识。**

#### 协议规范

编译层协议定义的是**四大操作的标准输入输出格式**。

```yaml
# SCP 操作定义
operations:
  capture:
    input: "原始对话/内容"
    output: "SCOL 增量包"
    description: "从原始内容中提取知识元素"

  compile:
    input: "SCOL 增量包集合"
    output:
      - "topics/*.md（主题知识）"
      - "entities/*.md（实体页）"
      - "concepts/*.md（概念页）"
      - "graph/relations.json（关联图谱）"
      - "profile/me.md（画像更新）"
    description: "将增量包编译为结构化知识"

  retrieve:
    input: "查询请求（关键词/语义/图谱遍历）"
    output: "匹配的知识列表 + 相关性评分"
    description: "从知识库中检索相关知识"

  audit:
    input: "知识库当前状态"
    output: "健康报告（问题列表 + 修复建议）"
    description: "检查知识库健康度"
```

#### Compiler 接口（Go 接口定义）

```go
// Compiler 是编译层的插件接口
// 不同的编译器可以有不同的编译策略
type Compiler interface {
    // Name 返回编译器标识
    Name() string

    // Compile 将增量包编译为结构化知识
    Compile(ctx context.Context, packages []*IncrementPackage, store Storage) (*CompileResult, error)

    // Retrieve 从知识库中检索知识
    Retrieve(ctx context.Context, query *RetrieveQuery, store Storage) (*RetrieveResult, error)

    // Audit 对知识库进行健康检查
    Audit(ctx context.Context, store Storage) (*AuditReport, error)
}
```

#### 可能的 Compiler 实现

| Compiler | 实现方式 | 特点 | 谁来做 |
|----------|---------|------|--------|
| `skill-compiler` | AI 助手 Skill 驱动 | 零成本，复用 AI 助手算力 | Synapse 官方 |
| `mcp-compiler` | MCP Server 调用 AI | 更可控，需要 MCP 支持 | Synapse 官方 |
| `actions-compiler` | GitHub Actions + LLM API | 自动化，需要 API Key | Synapse 官方 |
| `local-llm-compiler` | Ollama / llama.cpp | 离线、隐私友好 | **社区共建** |
| `rule-compiler` | 基于规则的编译（无 AI） | 轻量、确定性高 | **社区共建** |
| `hybrid-compiler` | 规则预处理 + AI 精炼 | 平衡性能和质量 | **社区共建** |

**生态价值**：编译策略可以多样化。有人追求隐私用本地 LLM，有人追求质量用 GPT-4，有人不想用 AI 用纯规则。协议统一，实现自由选择。

### Layer 3：存储层协议（Synapse Storage Protocol, SSP）

**定义：知识如何被持久化存储和访问。**

#### 协议规范

存储层协议定义的是**知识库的标准 CRUD 操作**。

```go
// Storage 是存储层的插件接口
// 不同的存储后端只要实现此接口，知识库就可以存在任何地方
type Storage interface {
    // --- 基本 CRUD ---

    // Read 读取指定路径的知识文件
    Read(ctx context.Context, path string) (*KnowledgeFile, error)

    // Write 写入/更新知识文件
    Write(ctx context.Context, path string, file *KnowledgeFile) error

    // Delete 删除知识文件
    Delete(ctx context.Context, path string) error

    // List 列出指定目录下的知识文件
    List(ctx context.Context, dir string, opts ListOptions) ([]*KnowledgeFile, error)

    // --- 搜索 ---

    // Search 全文搜索知识库
    Search(ctx context.Context, query string, opts SearchOptions) ([]*SearchResult, error)

    // --- 版本控制（可选） ---

    // Commit 提交当前变更
    Commit(ctx context.Context, message string) error

    // History 获取文件变更历史
    History(ctx context.Context, path string) ([]*ChangeRecord, error)

    // --- 元数据 ---

    // Stats 获取知识库统计信息
    Stats(ctx context.Context) (*StorageStats, error)
}
```

#### 可能的 Storage 实现

| Storage Backend | 存储介质 | 特点 | 谁来做 |
|-----------------|---------|------|--------|
| `github-storage` | GitHub Repo | 免费、版本控制、天然 Wiki | Synapse 官方 |
| `local-storage` | 本地文件系统 | 离线可用、最快 | Synapse 官方 |
| `gitea-storage` | Gitea / Gogs | 自托管 Git，隐私友好 | **社区共建** |
| `s3-storage` | AWS S3 / 阿里 OSS | 云端存储、大容量 | **社区共建** |
| `webdav-storage` | WebDAV（坚果云等） | 国内用户友好 | **社区共建** |
| `notion-storage` | Notion Database | Notion 用户零迁移 | **社区共建** |
| `sqlite-storage` | SQLite | 单文件、高性能查询 | **社区共建** |
| `ipfs-storage` | IPFS | 去中心化、永久存储 | **社区共建** |

**生态价值**：用户不被绑定在 GitHub 上。企业用户可能需要内部 Git 服务器，国内用户可能更喜欢坚果云 WebDAV，极客用户可能选择 IPFS。存储可以自由切换，知识不丢失。

### Layer 4：展示层协议（Synapse Display Protocol, SDP）

**定义：知识库如何被渲染为可浏览、可交互的形式。**

#### 协议规范

展示层协议定义的是**知识库的标准展示接口**。

```go
// Renderer 是展示层的插件接口
// 不同的渲染器将知识库渲染为不同的展示形式
type Renderer interface {
    // Name 返回渲染器标识
    Name() string

    // Render 将知识库渲染为目标格式
    Render(ctx context.Context, store Storage, opts RenderOptions) error

    // Serve 启动展示服务（可选，用于动态渲染）
    Serve(ctx context.Context, store Storage, addr string) error
}
```

#### 可能的 Renderer 实现

| Renderer | 展示形式 | 特点 | 谁来做 |
|----------|---------|------|--------|
| `hugo-renderer` | 静态网站（Hugo） | GitHub Pages 部署 | Synapse 官方 |
| `obsidian-renderer` | Obsidian Vault | 本地图谱浏览 | Synapse 官方（兼容即可） |
| `mkdocs-renderer` | 静态网站（MkDocs） | 文档风格 | **社区共建** |
| `vitepress-renderer` | 静态网站（VitePress） | 现代感、Vue 生态 | **社区共建** |
| `web-app-renderer` | 动态 Web App | 实时搜索、图谱交互 | **社区共建** |
| `tui-renderer` | 终端 UI（TUI） | 命令行极客 | **社区共建** |
| `pdf-renderer` | PDF 导出 | 打印/离线阅读 | **社区共建** |
| `anki-renderer` | Anki 闪卡 | 知识复习/间隔重复 | **社区共建** |
| `mindmap-renderer` | 思维导图 | 可视化知识结构 | **社区共建** |
| `newsletter-renderer` | 邮件周报 | 定期知识回顾 | **社区共建** |

**生态价值**：知识的"消费方式"可以无限多样。有人喜欢网站浏览，有人喜欢 Obsidian 图谱，有人喜欢 Anki 复习，有人喜欢收到邮件周报。展示层完全开放。

### 协议之间的交互模型

```
                    ┌─────────────────┐
                    │  用户 / AI 助手  │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │  Layer 1: SCOL   │ ← 各种 Provider 生成增量包
                    │  (采集协议)       │
                    └────────┬────────┘
                             │ IncrementPackage
                    ┌────────▼────────┐
                    │  Layer 2: SCP    │ ← 各种 Compiler 编译知识
                    │  (编译协议)       │
                    └────────┬────────┘
                             │ KnowledgeFile
                    ┌────────▼────────┐
                    │  Layer 3: SSP    │ ← 各种 Storage 持久化
                    │  (存储协议)       │
                    └────────┬────────┘
                             │ Read/Query
                    ┌────────▼────────┐
                    │  Layer 4: SDP    │ ← 各种 Renderer 展示
                    │  (展示协议)       │
                    └─────────────────┘

    贯穿所有层的基石：Layer 0 (SKS) — 统一的知识模型
```

### 关键设计决策

#### 1. 先协议后实现

每一层的开发顺序：
1. 定义协议（Go Interface + 数据结构）
2. 编写参考实现（官方默认实现）
3. 编写 Provider/Compiler/Storage/Renderer 开发指南
4. 社区基于协议贡献新实现

#### 2. Registry（注册中心）

```yaml
# .synapse/config.yaml — 配置使用哪些 Provider
synapse:
  version: "1.0"

  # 采集层：使用哪些 Provider
  collection:
    providers:
      - name: codebuddy-skill
        enabled: true
      - name: browser-extension
        enabled: true

  # 编译层：使用哪个 Compiler
  compilation:
    compiler: skill-compiler
    # compiler: local-llm-compiler  # 可切换

  # 存储层：使用哪个 Storage Backend
  storage:
    backend: github-storage
    config:
      repo: "username/knowhub"
      branch: "main"
    # backend: local-storage  # 可切换
    # config:
    #   path: "/Users/me/knowhub"

  # 展示层：使用哪些 Renderer
  display:
    renderers:
      - name: hugo-renderer
        config:
          theme: "synapse-default"
          deploy: "github-pages"
      - name: obsidian-renderer
        enabled: true
```

#### 3. CLI 作为协调器

`synapse` CLI 不再是"做所有事情的工具"，而是**协调各层 Provider 的编排器**：

```bash
# synapse 只做协调，具体实现由各层 Provider 完成
synapse init              # 初始化知识库 → 调用 Storage.Init()
synapse capture           # 触发采集 → 调用 Provider.Collect()
synapse compile           # 触发编译 → 调用 Compiler.Compile()
synapse retrieve "query"  # 检索知识 → 调用 Compiler.Retrieve()
synapse audit             # 健康检查 → 调用 Compiler.Audit()
synapse render            # 生成展示 → 调用 Renderer.Render()

# Provider 管理
synapse provider list                          # 列出所有已注册的 Provider
synapse provider install <name>                # 安装社区 Provider
synapse provider enable/disable <name>         # 启用/禁用 Provider
```

#### 4. 插件分发

社区 Provider 的分发方式：

| 方式 | 说明 | 适用场景 |
|------|------|---------|
| Go Plugin | 编译为 `.so`，动态加载 | 性能敏感的 Provider |
| 子进程 | 独立二进制，通过 stdin/stdout 通信 | 跨语言 Provider |
| WASM | WebAssembly 沙箱执行 | 安全性要求高的场景 |
| Git Repo | 直接引用 Git 仓库 | 开发阶段、社区共享 |

**推荐方案**：子进程 + JSON-RPC 通信（类似 MCP 的 stdio 模式），简单、跨语言、安全。

```
synapse CLI ←→ [JSON-RPC over stdin/stdout] ←→ Provider 子进程
```

### 对比分析：当前方案 vs 生态化方案

| 维度 | 当前方案 | 生态化方案 |
|------|---------|-----------|
| 数据源 | 只有 AI 助手对话 | **任何数据源**（RSS/Notion/Twitter/播客...） |
| 存储 | 只有 GitHub | **任何存储后端**（本地/Git/S3/WebDAV/IPFS...） |
| 编译 | 只有 Skill Prompt | **任何编译器**（Skill/MCP/本地LLM/规则...） |
| 展示 | 只有 Hugo + Obsidian | **任何展示形式**（Web/TUI/PDF/Anki/邮件...） |
| 扩展方式 | 改 Synapse 代码 | **实现接口即可**，不需要改 Synapse |
| 社区参与 | 贡献 Skill Prompt | **四个层面都可以贡献** |
| 迁移成本 | 被锁定在 GitHub | **知识可移植**，任意切换后端 |

### 对 Roadmap 的影响

生态化架构**不增加 MVP 的工作量**，但需要在 M1 阶段就把协议定义好：

| 里程碑 | 新增内容 |
|--------|---------|
| **M1** | 定义四层协议接口（Go Interface）；知识模型 Schema（SKS） |
| **M2** | Skill 作为第一个 Provider（SCOL）+ 第一个 Compiler（SCP）实现 |
| **M3** | MCP Server 作为第二个 Compiler 实现；GitHub 作为第一个 Storage（SSP）实现 |
| **M4** | 各平台 Skill 作为更多 Provider 实现 |
| **M5** | Hugo/Obsidian 作为第一批 Renderer（SDP）实现 |
| **M6+** | 开放社区共建，发布 Provider 开发指南 |

### 类比

这个架构像什么？

- **LSP（Language Server Protocol）**：定义了编辑器和语言服务之间的协议，任何编辑器 + 任何语言都可以组合
- **CSI（Container Storage Interface）**：定义了容器和存储之间的协议，任何存储都可以接入 K8s
- **MCP（Model Context Protocol）**：定义了 AI 助手和外部工具之间的协议，任何工具都可以接入 AI

**Synapse 的四层协议 = 个人知识管理领域的 LSP/CSI/MCP**

---

## 2026-04-18 架构再审视：从第一性原理重新推导分层

### 背景

上一轮讨论中，用户提出了"各层级可以生态化"的核心想法，我快速设计了一个"采集-编译-存储-展示"四层协议模型。但用户指出：**"不一定是这四层，你需要自己思考整体架构"**。

这是一个非常关键的提醒。我之前的做法是**从用户的举例出发去设计**，而不是**从 Synapse 的本质需求出发去推导**。这次要回到第一性原理。

### 对上一版四层架构的反思

#### 问题一：采集和编译的边界模糊

上一版把"采集"和"编译"分为两个独立的层（SCOL + SCP），但实际上：

- 在 Skill 场景下，AI 助手**既做采集又做编译**，是同一个 Prompt 驱动的同一个过程
- 采集的输出（增量包）和编译的输入（增量包）是同一个东西
- 更根本的问题：**"从原始内容到结构化知识"是一个连续的处理过程**，按哪里切分取决于实际需要，不应该预设为固定的两层

真正需要**独立、可替换**的是：
- **数据源**：原始内容从哪里来？（AI 对话 / RSS / Notion / 播客...）
- **处理引擎**：谁来做知识提取和结构化？（AI / 规则引擎 / 本地 LLM...）

这两个才是真正正交的扩展点。

#### 问题二：存储不是数据流上的一个"层"

上一版把存储放在采集→编译之后、展示之前，暗示数据流是 `采集→编译→存储→展示` 的线性流。但实际上：

- **处理引擎**需要读写存储（编译时要读已有知识、写入新知识）
- **检索引擎**需要读存储
- **展示**需要读存储
- **审计**需要读存储

存储更像是**所有模块都依赖的基础设施**，是一个**底座**，而不是流水线上的一个节点。

#### 问题三：编译层职责过重

上一版的 `Compiler` 接口包含了 `Compile()` + `Retrieve()` + `Audit()` 三个完全不同的操作：

- `Compile` 是**写入**操作（增量包 → 结构化知识）
- `Retrieve` 是**读取**操作（查询 → 匹配结果）
- `Audit` 是**运维**操作（检查健康度）

把这三个塞到一个接口里违反了单一职责原则，也让社区贡献变得困难——想贡献一个新的检索引擎，不应该被迫实现编译逻辑。

#### 问题四：展示层概念太窄

"展示"暗示的是**人类浏览**，但知识被消费的方式远不止于此：

- 反哺 AI 助手（通过 MCP / Skill）—— 消费者是 AI，不是人
- 导出为 Anki 闪卡 —— 不是"展示"，更像是"转换/分发"
- 推送邮件周报 —— 是"推送"，不是被动的"展示"
- 同步到 Notion —— 是另一个系统的"导入"

更准确的概念应该是**"消费"或"输出"**——知识库的内容以各种形式被各种消费者使用。

### 从第一性原理重新推导

#### 核心问题：Synapse 生态化到底需要什么？

回到最根本的问题：**Synapse 需要哪些能力是独立、可替换、可由社区共建的？**

让我列出所有职责，然后判断每个职责的**正交性**（即它们之间是否可以独立变化）：

| # | 职责 | 描述 | 需要可替换吗？ | 为什么？ |
|---|------|------|--------------|---------|
| 1 | 知识结构定义 | 页面类型、Frontmatter、链接格式 | 需要可**扩展** | 不同领域（编程/医学/法律）需要不同的知识结构 |
| 2 | 数据获取 | 从外部拿到原始内容 | **必须**可替换 | 数据源千差万别，这是最容易社区贡献的 |
| 3 | 知识处理 | 原始内容 → 结构化知识 | **必须**可替换 | AI / 规则 / 本地LLM，策略完全不同 |
| 4 | 持久化 | 知识文件的读写 | **必须**可替换 | GitHub / 本地 / S3 / WebDAV，场景差异大 |
| 5 | 检索 | 从知识库中找相关内容 | 需要可替换 | BM25 / 向量 / 图谱遍历，算法差异大 |
| 6 | 质量维护 | 健康检查、去重、修复 | 可选可替换 | 不同用户对质量标准不同 |
| 7 | 知识输出 | 以各种形式被消费 | **必须**可替换 | 网站 / Obsidian / Anki / 邮件 / AI反哺，形态各异 |

**正交性分析**（两个职责是否可以独立变化）：

```
          数据获取  处理引擎  持久化  检索  质量维护  输出
数据获取    -       独立      独立    无关   无关     无关
处理引擎    独立     -        依赖    无关   无关     无关
持久化      独立     被依赖    -      被依赖  被依赖   被依赖
检索        无关     无关      依赖    -     无关     无关
质量维护    无关     无关      依赖    无关    -      无关
输出        无关     无关      依赖    可选    无关     -
```

关键发现：
1. **持久化**被所有其他模块依赖 → 它是底座，不是流水线上的一个节点
2. **数据获取**和**处理引擎**是独立的（同一个数据源可以用不同的处理引擎）
3. **检索**和**处理引擎**是独立的（同一个知识库可以用不同的检索算法）
4. **输出**和其他所有模块都是独立的（除了依赖持久化读数据）
5. **质量维护**是相对独立的横切能力

### 新架构：扩展点模型（Extension Point Model）

基于以上分析，我认为 Synapse 不应该是一个"固定N层的协议栈"，而应该是一个**以知识规范为中心、以存储为底座、围绕多个独立扩展点的星型架构**。

```
                        ┌─────────────┐
                        │   Source     │ 数据源（原始内容从哪来）
                        │  ┌───────┐  │
                        │  │AI对话  │  │  可替换：AI 对话 / RSS / Notion / 播客 / ...
                        │  │RSS    │  │
                        │  │Notion │  │
                        │  └───┬───┘  │
                        └─────┼───────┘
                              │ RawContent
                              ▼
                        ┌─────────────┐
                        │  Processor  │ 处理引擎（原始内容 → 结构化知识）
                        │  ┌───────┐  │
                        │  │Skill  │  │  可替换：Skill / MCP / LocalLLM / Rules / ...
                        │  │MCP    │  │
                        │  │Ollama │  │
                        │  └───┬───┘  │
                        └─────┼───────┘
                              │ KnowledgeFile
                              ▼
    ┌──────────────────────────────────────────────────────┐
    │                                                      │
    │               Store（存储底座）                        │
    │                                                      │
    │  可替换：GitHub / Local FS / S3 / WebDAV / IPFS / ...│
    │                                                      │
    │  职责：知识文件的 CRUD + 版本控制                       │
    │                                                      │
    └───────┬──────────────┬───────────────┬───────────────┘
            │              │               │
            ▼              ▼               ▼
    ┌───────────┐  ┌───────────┐   ┌───────────────┐
    │  Indexer   │  │  Auditor  │   │   Consumer    │
    │  检索引擎  │  │  质量审计  │   │   消费端      │
    │           │  │           │   │              │
    │ BM25     │  │ 断链检测  │   │ Hugo 网站    │
    │ 向量检索  │  │ 孤儿页面  │   │ Obsidian    │
    │ 图谱遍历  │  │ 过时检测  │   │ Anki 闪卡   │
    │           │  │ 去重      │   │ 邮件周报     │
    │           │  │           │   │ AI 反哺(MCP) │
    └───────────┘  └───────────┘   │ TUI 浏览     │
                                   └───────────────┘

    贯穿所有模块：Schema（知识规范）— 统一的"语言"
```

#### 与上一版的核心区别

| 维度 | 上一版（四层协议栈） | 新版（扩展点模型） |
|------|-------------------|------------------|
| **架构形态** | 线性流水线：采集→编译→存储→展示 | **星型**：存储为底座，多个扩展点围绕 |
| **存储定位** | 流水线上的一个节点（Layer 3） | **底座**（所有模块都依赖它） |
| **采集+编译** | 两个独立的层（SCOL + SCP） | **两个扩展点**（Source + Processor），但不强制分层 |
| **检索** | 塞在 Compiler 接口里 | **独立扩展点**（Indexer），可以单独替换 |
| **审计** | 塞在 Compiler 接口里 | **独立扩展点**（Auditor），可以单独替换 |
| **展示** | 一个层（SDP） | **Consumer 扩展点**，不仅是展示，还包括 AI 反哺、导出等 |
| **层间关系** | 严格的上下级 | **按需组合**，扩展点之间不强制依赖 |

#### 五个扩展点（Extension Points）

每个扩展点都是一个**独立的 Go 接口**，社区可以针对任意一个扩展点贡献新实现，不需要理解其他扩展点的逻辑。

**1. Source（数据源）**

```go
// Source 定义了如何从外部获取原始内容
type Source interface {
    Name() string
    Collect(ctx context.Context, opts CollectOptions) ([]*RawContent, error)
    Watch(ctx context.Context, handler func(*RawContent)) error // 可选：实时监听
}
```

**社区共建场景**：想接入 RSS 订阅？实现 Source 接口即可。想同步 Notion？实现 Source 接口即可。

**2. Processor（处理引擎）**

```go
// Processor 定义了如何将原始内容处理为结构化知识
type Processor interface {
    Name() string
    Process(ctx context.Context, raw []*RawContent, store Store) (*ProcessResult, error)
}
```

**社区共建场景**：想用本地 LLM 处理？实现 Processor 接口。想用纯规则引擎？实现 Processor 接口。

**3. Store（存储底座）**

```go
// Store 定义了知识的持久化方式
type Store interface {
    Read(ctx context.Context, path string) (*KnowledgeFile, error)
    Write(ctx context.Context, path string, file *KnowledgeFile) error
    Delete(ctx context.Context, path string) error
    List(ctx context.Context, dir string, opts ListOptions) ([]*FileInfo, error)
    Exists(ctx context.Context, path string) (bool, error)
    
    // 可选能力（通过接口组合实现）
    // VersionedStore: Commit() / History()
    // SearchableStore: Search()
}
```

**4. Indexer（检索引擎）**

```go
// Indexer 定义了如何从知识库中检索相关内容
type Indexer interface {
    Name() string
    Build(ctx context.Context, store Store) error                           // 构建/更新索引
    Search(ctx context.Context, query string, opts SearchOptions) ([]*SearchResult, error)
}
```

**社区共建场景**：想加向量检索？实现 Indexer 接口。想做图谱遍历？实现 Indexer 接口。

**5. Consumer（消费端）**

```go
// Consumer 定义了知识如何被外部消费
type Consumer interface {
    Name() string
    Consume(ctx context.Context, store Store, opts ConsumeOptions) error
}
```

这个接口足够泛化，可以涵盖：
- Hugo 渲染静态网站（`HugoConsumer`）
- Obsidian Vault 兼容（`ObsidianConsumer`）
- Anki 闪卡导出（`AnkiConsumer`）
- 邮件周报推送（`NewsletterConsumer`）
- AI 反哺（`MCPConsumer`）
- TUI 终端浏览（`TUIConsumer`）

**6. Auditor（质量审计）**— 可选扩展点

```go
// Auditor 定义了知识库的质量检查策略
type Auditor interface {
    Name() string
    Audit(ctx context.Context, store Store) (*AuditReport, error)
    Fix(ctx context.Context, store Store, issues []*Issue) (*FixResult, error) // 可选：自动修复
}
```

#### 为什么是"扩展点"而不是"层"？

1. **不强制分层**：Source 和 Processor 可以合并为一个实现（比如 Skill 同时做采集和处理），也可以分开实现
2. **不强制线性**：Consumer 不必等前面所有步骤完成，它只需要能读 Store
3. **按需组合**：用户可以只用 Source + Processor + Store（最小集合），不用 Indexer、Auditor、Consumer
4. **正交替换**：换一个 Indexer 不影响 Processor，换一个 Store 不影响 Source

#### Schema 的定位

Schema 不是一个"扩展点"，它是**所有扩展点遵循的规范**。就像 HTTP 规范不是 TCP/IP 协议栈的一层，而是所有 Web 参与者共同遵守的约定。

```yaml
# .synapse/schema.yaml
version: "1.0"

# 知识类型定义（所有扩展点都基于这个结构工作）
page_types:
  - name: profile
    directory: profile/
  - name: topic
    directory: topics/
  - name: entity
    directory: entities/
  - name: concept
    directory: concepts/
  - name: inbox
    directory: inbox/
  - name: journal
    directory: journal/

# Frontmatter 规范
frontmatter:
  required: [type, title, created, updated]
  optional: [tags, links, source, confidence]

# 双向链接格式
link_format: "[[page-id]]"
```

#### 配置模型

```yaml
# .synapse/config.yaml
synapse:
  version: "1.0"

  # 数据源（可以同时启用多个）
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
    # name: local-llm-processor
    # name: rule-processor

  # 存储底座（选一个）
  store:
    name: local-store
    config:
      path: "/Users/me/knowhub"
    # name: github-store
    # config:
    #   repo: "user/knowhub"

  # 检索引擎（可选）
  indexer:
    name: bm25-indexer

  # 消费端（可以同时启用多个）
  consumers:
    - name: hugo-site
      config:
        theme: synapse-default
        deploy: github-pages
    - name: obsidian-vault
      enabled: true

  # 审计（可选）
  auditor:
    name: default-auditor
```

#### CLI 作为编排器

```bash
# 核心流程
synapse collect          # Source.Collect() → 获取原始内容
synapse process          # Processor.Process() → 结构化处理
synapse run              # collect + process 一条龙

# 检索
synapse search "query"   # Indexer.Search()
synapse index            # Indexer.Build() — 重建索引

# 审计
synapse audit            # Auditor.Audit()
synapse fix              # Auditor.Fix() — 自动修复

# 消费/输出
synapse publish          # Consumer.Consume() — 触发所有已启用的消费端
synapse publish --only hugo-site  # 只触发特定消费端

# 扩展点管理
synapse plugin list
synapse plugin install <name>
synapse plugin info <name>
```

### 总结

新架构的核心思想：

1. **存储是底座，不是流水线节点** — 所有模块都围绕 Store 工作
2. **五个独立扩展点，不是固定分层** — 按需组合，正交替换
3. **Schema 是共同语言** — 不是一个"层"，而是所有扩展点遵循的规范
4. **星型而非线性** — 不强制数据经过每一层，Consumer 直接读 Store

---

## 2026-04-18 扩展点集成机制设计

### 背景

上一轮确定了"扩展点模型"——六个独立扩展点（Source / Processor / Store / Indexer / Consumer / Auditor），每个可独立替换、社区共建。现在要回答一个关键问题：**这些扩展点到底怎么集成到 Synapse 里？**

这个问题涉及：
- 插件的**发现**：系统怎么知道有哪些可用的插件？
- 插件的**注册**：用户怎么选择使用哪个插件？
- 插件的**加载**：运行时怎么把插件实例化？
- 插件的**通信**：Synapse 核心和插件之间怎么交互？

### 业界方案调研

| 项目 | 插件机制 | 通信方式 | 优缺点 |
|------|---------|---------|--------|
| **Terraform** | Go Provider + gRPC（go-plugin） | 子进程 + gRPC over stdin/stdout | ✅ 跨语言、稳定隔离 ❌ 部署重、性能开销 |
| **Caddy** | Go Module（编译时注入） | 进程内直接调用 | ✅ 零开销、类型安全 ❌ 必须用 Go、需要重新编译 |
| **Docker / K8s CSI** | gRPC over Unix Socket | 独立进程 + gRPC | ✅ 完全解耦 ❌ 运维成本高 |
| **MCP** | JSON-RPC over stdin/stdout | 子进程 + JSON-RPC | ✅ 极简、跨语言 ❌ 性能一般 |
| **VS Code** | Node.js Extension Host | 独立进程 + JSON-RPC | ✅ 隔离好 ❌ 限定 JS/TS |
| **Hugo** | Go 模板 + 配置 | 无插件系统，全部内置 | ✅ 简单 ❌ 不可扩展 |
| **Grafana** | Go Plugin SDK + gRPC | HashiCorp go-plugin | ✅ 成熟方案 ❌ 复杂度高 |

### Synapse 的约束条件

在选择方案之前，先明确 Synapse 的特殊约束：

1. **Synapse 是 CLI 工具，不是长驻服务**：每次执行一个命令，完成即退出。不像 Docker/K8s 那样需要 daemon 管理插件生命周期
2. **MVP 阶段要极简**：不能一上来就引入 gRPC/WASM 等重量级方案
3. **社区贡献门槛要低**：理想情况下，贡献一个 Source 只需要实现一个 Go 接口
4. **要支持跨语言（长期）**：Python/Node 开发者也应该能贡献插件
5. **性能不是瓶颈**：知识处理是低频操作，不需要纳秒级通信

### 三层集成模型

基于以上约束，我设计一个**渐进式的三层集成模型**——从简单到复杂，从内置到外部：

```
┌───────────────────────────────────────────────────────────────────────┐
│                                                                       │
│  Layer 1: 内置扩展（Built-in）                                         │
│                                                                       │
│  Go 接口 + init() 注册                                                │
│  零开销、类型安全、Go 代码直接调用                                        │
│                                                                       │
│  适用：官方实现 + Go 社区贡献                                            │
│  示例：local-store, skill-processor, bm25-indexer                      │
│                                                                       │
├───────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  Layer 2: 本地插件（Local Plugin）                                      │
│                                                                       │
│  独立可执行文件 + JSON-RPC over stdin/stdout                            │
│  跨语言、进程隔离、MCP 风格通信                                          │
│                                                                       │
│  适用：非 Go 语言的社区贡献、需要进程隔离的场景                             │
│  示例：rss-source (Python), notion-source (Node.js)                    │
│                                                                       │
├───────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  Layer 3: 远程插件（Remote Plugin）— 远期                               │
│                                                                       │
│  HTTP/gRPC over network                                               │
│  云端运行、SaaS 集成                                                    │
│                                                                       │
│  适用：需要持久运行的服务（如 Webhook 接收器）                              │
│  示例：webhook-source (云端服务), vector-indexer (GPU 服务器)             │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

**MVP 只需要 Layer 1（内置扩展）**。Layer 2 在 M3-M4 引入，Layer 3 在 M6+ 远期考虑。

### Layer 1：内置扩展（Built-in）— MVP 核心

#### 核心思路

借鉴 **Caddy** 和 **Go database/sql** 的模式：

- 每个扩展点定义为 Go Interface
- 每个实现是一个独立的 Go Package
- 通过 `init()` 函数自注册到全局 Registry
- 主程序通过 `_ "github.com/xxx/synapse/plugins/xxx"` 导入即激活
- 运行时通过配置文件选择使用哪个实现

#### 目录结构

```
synapse/
├── pkg/
│   └── extension/                    # 扩展点框架
│       ├── registry.go               # 全局注册表
│       ├── source.go                 # Source 接口定义
│       ├── processor.go              # Processor 接口定义
│       ├── store.go                  # Store 接口定义
│       ├── indexer.go                # Indexer 接口定义
│       ├── consumer.go              # Consumer 接口定义
│       └── auditor.go               # Auditor 接口定义
│
├── internal/
│   ├── engine/                       # 核心编排引擎
│   │   └── engine.go                 # 读取配置，实例化扩展点，编排执行
│   │
│   ├── source/                       # Source 扩展点实现
│   │   ├── skill/                    # 内置：CodeBuddy Skill Source
│   │   │   └── skill.go
│   │   └── cli_import/               # 内置：CLI 导入 Source
│   │       └── import.go
│   │
│   ├── processor/                    # Processor 扩展点实现
│   │   ├── skill/                    # 内置：Skill Processor
│   │   │   └── skill.go
│   │   └── mcp/                      # 内置：MCP Processor（M3）
│   │       └── mcp.go
│   │
│   ├── store/                        # Store 扩展点实现
│   │   ├── local/                    # 内置：本地文件系统
│   │   │   └── local.go
│   │   └── github/                   # 内置：GitHub Store（M3）
│   │       └── github.go
│   │
│   ├── indexer/                      # Indexer 扩展点实现
│   │   └── bm25/                     # 内置：BM25 检索
│   │       └── bm25.go
│   │
│   ├── consumer/                     # Consumer 扩展点实现
│   │   ├── hugo/                     # 内置：Hugo 网站
│   │   │   └── hugo.go
│   │   └── obsidian/                 # 内置：Obsidian 兼容
│   │       └── obsidian.go
│   │
│   └── auditor/                      # Auditor 扩展点实现
│       └── default/                  # 内置：默认审计器
│           └── default.go
│
├── cmd/
│   └── synapse/
│       ├── main.go                   # CLI 入口
│       └── plugins.go                # 插件注册（import 所有内置扩展）
│
└── go.mod
```

#### 注册机制（Registry Pattern）

```go
// pkg/extension/registry.go

// Registry 是扩展点的全局注册表
type Registry struct {
    mu         sync.RWMutex
    sources    map[string]SourceFactory
    processors map[string]ProcessorFactory
    stores     map[string]StoreFactory
    indexers   map[string]IndexerFactory
    consumers  map[string]ConsumerFactory
    auditors   map[string]AuditorFactory
}

// 全局单例
var globalRegistry = &Registry{
    sources:    make(map[string]SourceFactory),
    processors: make(map[string]ProcessorFactory),
    stores:     make(map[string]StoreFactory),
    indexers:   make(map[string]IndexerFactory),
    consumers:  make(map[string]ConsumerFactory),
    auditors:   make(map[string]AuditorFactory),
}

// SourceFactory 是 Source 的工厂函数
// config 来自 .synapse/config.yaml 中该 Source 的 config 字段
type SourceFactory func(config map[string]any) (Source, error)

// RegisterSource 注册一个 Source 实现
func RegisterSource(name string, factory SourceFactory) {
    globalRegistry.mu.Lock()
    defer globalRegistry.mu.Unlock()
    if _, exists := globalRegistry.sources[name]; exists {
        panic(fmt.Sprintf("source %q already registered", name))
    }
    globalRegistry.sources[name] = factory
}

// GetSource 根据名称获取 Source 实例
func GetSource(name string, config map[string]any) (Source, error) {
    globalRegistry.mu.RLock()
    factory, ok := globalRegistry.sources[name]
    globalRegistry.mu.RUnlock()
    if !ok {
        return nil, fmt.Errorf("unknown source: %s", name)
    }
    return factory(config)
}

// ListSources 列出所有已注册的 Source
func ListSources() []string {
    globalRegistry.mu.RLock()
    defer globalRegistry.mu.RUnlock()
    names := make([]string, 0, len(globalRegistry.sources))
    for name := range globalRegistry.sources {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}

// 其他扩展点的 Register/Get/List 方法类似...
// RegisterProcessor, RegisterStore, RegisterIndexer, RegisterConsumer, RegisterAuditor
```

#### 扩展点实现示例

```go
// internal/store/local/local.go

package local

import (
    "context"
    "os"
    "path/filepath"

    "github.com/xxx/synapse/pkg/extension"
)

func init() {
    // 自注册到全局 Registry
    extension.RegisterStore("local-store", New)
}

// LocalStore 是基于本地文件系统的 Store 实现
type LocalStore struct {
    basePath string
}

// New 创建一个新的 LocalStore 实例
func New(config map[string]any) (extension.Store, error) {
    path, ok := config["path"].(string)
    if !ok || path == "" {
        return nil, fmt.Errorf("local-store requires 'path' config")
    }
    absPath, err := filepath.Abs(path)
    if err != nil {
        return nil, fmt.Errorf("invalid path %q: %w", path, err)
    }
    return &LocalStore{basePath: absPath}, nil
}

func (s *LocalStore) Read(ctx context.Context, path string) (*extension.KnowledgeFile, error) {
    fullPath := filepath.Join(s.basePath, path)
    data, err := os.ReadFile(fullPath)
    if err != nil {
        return nil, fmt.Errorf("read %s: %w", path, err)
    }
    return extension.ParseKnowledgeFile(data)
}

func (s *LocalStore) Write(ctx context.Context, path string, file *extension.KnowledgeFile) error {
    fullPath := filepath.Join(s.basePath, path)
    if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
        return fmt.Errorf("create directory for %s: %w", path, err)
    }
    data, err := file.Marshal()
    if err != nil {
        return fmt.Errorf("marshal %s: %w", path, err)
    }
    return os.WriteFile(fullPath, data, 0o644)
}

// ... Delete, List, Exists 等方法
```

#### 插件激活

```go
// cmd/synapse/plugins.go

package main

// 通过 blank import 激活所有内置扩展点实现
// 每个 package 的 init() 会自动注册到全局 Registry
import (
    // Source 实现
    _ "github.com/xxx/synapse/internal/source/skill"
    _ "github.com/xxx/synapse/internal/source/cli_import"

    // Processor 实现
    _ "github.com/xxx/synapse/internal/processor/skill"

    // Store 实现
    _ "github.com/xxx/synapse/internal/store/local"

    // Indexer 实现
    _ "github.com/xxx/synapse/internal/indexer/bm25"

    // Consumer 实现
    _ "github.com/xxx/synapse/internal/consumer/hugo"
    _ "github.com/xxx/synapse/internal/consumer/obsidian"

    // Auditor 实现
    _ "github.com/xxx/synapse/internal/auditor/default_auditor"
)
```

#### 编排引擎

```go
// internal/engine/engine.go

// Engine 是核心编排器
// 它读取配置，实例化各扩展点，协调它们的工作
type Engine struct {
    config    *Config
    sources   []extension.Source
    processor extension.Processor
    store     extension.Store
    indexer   extension.Indexer    // 可选
    consumers []extension.Consumer // 可选
    auditor   extension.Auditor   // 可选
}

// New 从配置文件创建引擎实例
func New(configPath string) (*Engine, error) {
    cfg, err := LoadConfig(configPath)
    if err != nil {
        return nil, fmt.Errorf("load config: %w", err)
    }

    e := &Engine{config: cfg}

    // 实例化 Store（底座，最先初始化）
    e.store, err = extension.GetStore(cfg.Store.Name, cfg.Store.Config)
    if err != nil {
        return nil, fmt.Errorf("init store %q: %w", cfg.Store.Name, err)
    }

    // 实例化 Sources（可多个）
    for _, srcCfg := range cfg.Sources {
        if !srcCfg.Enabled {
            continue
        }
        src, err := extension.GetSource(srcCfg.Name, srcCfg.Config)
        if err != nil {
            return nil, fmt.Errorf("init source %q: %w", srcCfg.Name, err)
        }
        e.sources = append(e.sources, src)
    }

    // 实例化 Processor
    e.processor, err = extension.GetProcessor(cfg.Processor.Name, cfg.Processor.Config)
    if err != nil {
        return nil, fmt.Errorf("init processor %q: %w", cfg.Processor.Name, err)
    }

    // 实例化可选扩展点
    if cfg.Indexer != nil {
        e.indexer, err = extension.GetIndexer(cfg.Indexer.Name, cfg.Indexer.Config)
        if err != nil {
            return nil, fmt.Errorf("init indexer %q: %w", cfg.Indexer.Name, err)
        }
    }
    // ... consumers, auditor 类似

    return e, nil
}

// Collect 执行采集流程：Source.Collect → Processor.Process → Store.Write
func (e *Engine) Collect(ctx context.Context) error {
    // 1. 从所有启用的 Source 采集原始内容
    var allRaw []*extension.RawContent
    for _, src := range e.sources {
        raw, err := src.Collect(ctx, extension.CollectOptions{})
        if err != nil {
            // 单个 Source 失败不阻断整体流程
            log.Printf("WARN: source %s collect failed: %v", src.Name(), err)
            continue
        }
        allRaw = append(allRaw, raw...)
    }
    if len(allRaw) == 0 {
        log.Println("no raw content collected")
        return nil
    }

    // 2. 用 Processor 处理原始内容为结构化知识
    result, err := e.processor.Process(ctx, allRaw, e.store)
    if err != nil {
        return fmt.Errorf("process: %w", err)
    }

    // 3. 如果有 Indexer，更新索引
    if e.indexer != nil {
        if err := e.indexer.Build(ctx, e.store); err != nil {
            log.Printf("WARN: index build failed: %v", err)
        }
    }

    log.Printf("collected %d raw items, produced %d knowledge files",
        len(allRaw), len(result.Files))
    return nil
}

// Search 执行检索
func (e *Engine) Search(ctx context.Context, query string) ([]*extension.SearchResult, error) {
    if e.indexer == nil {
        return nil, fmt.Errorf("no indexer configured")
    }
    return e.indexer.Search(ctx, query, extension.SearchOptions{})
}

// Audit 执行审计
func (e *Engine) Audit(ctx context.Context) (*extension.AuditReport, error) {
    if e.auditor == nil {
        return nil, fmt.Errorf("no auditor configured")
    }
    return e.auditor.Audit(ctx, e.store)
}

// Publish 触发所有消费端
func (e *Engine) Publish(ctx context.Context) error {
    for _, c := range e.consumers {
        if err := c.Consume(ctx, e.store, extension.ConsumeOptions{}); err != nil {
            log.Printf("WARN: consumer %s failed: %v", c.Name(), err)
        }
    }
    return nil
}
```

#### CLI 与引擎的连接

```go
// cmd/synapse/main.go

func main() {
    app := &cli.App{
        Name:  "synapse",
        Usage: "Personal Knowledge Hub — 个人知识中枢",
        Commands: []*cli.Command{
            {
                Name:  "collect",
                Usage: "从所有启用的 Source 采集知识并处理",
                Action: func(c *cli.Context) error {
                    eng, err := engine.New(".synapse/config.yaml")
                    if err != nil {
                        return err
                    }
                    return eng.Collect(c.Context)
                },
            },
            {
                Name:  "search",
                Usage: "检索知识库",
                Action: func(c *cli.Context) error {
                    eng, err := engine.New(".synapse/config.yaml")
                    if err != nil {
                        return err
                    }
                    results, err := eng.Search(c.Context, c.Args().First())
                    if err != nil {
                        return err
                    }
                    // 输出结果...
                    return nil
                },
            },
            {
                Name:  "audit",
                Usage: "检查知识库健康度",
                Action: func(c *cli.Context) error {
                    eng, err := engine.New(".synapse/config.yaml")
                    if err != nil {
                        return err
                    }
                    report, err := eng.Audit(c.Context)
                    if err != nil {
                        return err
                    }
                    // 输出报告...
                    return nil
                },
            },
            {
                Name:  "publish",
                Usage: "触发所有消费端输出",
                Action: func(c *cli.Context) error {
                    eng, err := engine.New(".synapse/config.yaml")
                    if err != nil {
                        return err
                    }
                    return eng.Publish(c.Context)
                },
            },
            {
                Name:  "plugin",
                Usage: "管理扩展点插件",
                Subcommands: []*cli.Command{
                    {
                        Name:  "list",
                        Usage: "列出所有已注册的插件",
                        Action: pluginListAction,
                    },
                },
            },
        },
    }
    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
```

#### 配置与插件的对应关系

```yaml
# .synapse/config.yaml
synapse:
  version: "1.0"

  sources:
    - name: skill-source          # → 匹配 extension.GetSource("skill-source", config)
      enabled: true
    - name: rss-source            # → 如果是外部插件，走 Layer 2 通信
      enabled: true
      plugin: "./plugins/rss-source"  # 指向外部可执行文件
      config:
        feeds: ["https://..."]

  processor:
    name: skill-processor         # → 匹配 extension.GetProcessor("skill-processor", config)

  store:
    name: local-store             # → 匹配 extension.GetStore("local-store", config)
    config:
      path: "/Users/me/knowhub"

  indexer:
    name: bm25-indexer

  consumers:
    - name: hugo-site
      config: { theme: synapse-default }
    - name: obsidian-vault
```

### Layer 2：本地插件（Local Plugin）— M3+ 引入

当社区用 Python/Node 等非 Go 语言开发插件时，采用**子进程 + JSON-RPC** 通信。

#### 通信协议

```
synapse CLI                          Plugin Process
    │                                     │
    │  ──── Launch subprocess ────>       │
    │                                     │
    │  ──── {"jsonrpc":"2.0",             │
    │        "method":"initialize",       │
    │        "params":{...}} ────>        │
    │                                     │
    │  <──── {"jsonrpc":"2.0",            │
    │         "result":{                  │
    │           "name":"rss-source",      │
    │           "type":"source",          │
    │           "version":"1.0.0"         │
    │         }} ─────                    │
    │                                     │
    │  ──── {"jsonrpc":"2.0",             │
    │        "method":"collect",          │
    │        "params":{...}} ────>        │
    │                                     │
    │  <──── {"jsonrpc":"2.0",            │
    │         "result":[                  │
    │           {"title":"...",           │
    │            "content":"..."}         │
    │         ]} ─────                    │
    │                                     │
    │  ──── {"jsonrpc":"2.0",             │
    │        "method":"shutdown"} ────>   │
    │                                     │
    │  <──── process exits ─────          │
    │                                     │
```

#### JSON-RPC 方法映射

每个扩展点的 Go 接口方法映射为 JSON-RPC method：

| 扩展点 | Go 方法 | JSON-RPC method | 请求 params | 响应 result |
|--------|---------|-----------------|-------------|-------------|
| Source | `Name()` | `initialize` | `{}` | `{name, type, version}` |
| Source | `Collect()` | `collect` | `{options}` | `[RawContent...]` |
| Source | `Watch()` | `watch` | `{}` | stream of RawContent |
| Processor | `Process()` | `process` | `{raw_contents, ...}` | `{files, links}` |
| Store | `Read()` | `read` | `{path}` | `{file}` |
| Store | `Write()` | `write` | `{path, file}` | `{}` |
| Indexer | `Build()` | `build` | `{}` | `{}` |
| Indexer | `Search()` | `search` | `{query, options}` | `[SearchResult...]` |
| Consumer | `Consume()` | `consume` | `{options}` | `{}` |
| Auditor | `Audit()` | `audit` | `{}` | `{report}` |
| Auditor | `Fix()` | `fix` | `{issues}` | `{result}` |
| 通用 | — | `shutdown` | `{}` | `{}` |

#### 外部插件适配器

Synapse 核心提供一个 `PluginAdapter`，它实现了 Go 接口但内部通过 JSON-RPC 与子进程通信：

```go
// pkg/extension/plugin/adapter.go

// PluginSource 将外部子进程适配为 Source 接口
type PluginSource struct {
    name    string
    cmd     *exec.Cmd
    client  *jsonrpc.Client
    stdin   io.WriteCloser
    stdout  io.ReadCloser
}

// NewPluginSource 启动外部插件进程并建立 JSON-RPC 通信
func NewPluginSource(execPath string, config map[string]any) (*PluginSource, error) {
    cmd := exec.Command(execPath)
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, fmt.Errorf("create stdin pipe: %w", err)
    }
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, fmt.Errorf("create stdout pipe: %w", err)
    }

    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("start plugin %s: %w", execPath, err)
    }

    client := jsonrpc.NewClient(stdout, stdin)

    // 初始化握手
    var initResult InitializeResult
    if err := client.Call("initialize", config, &initResult); err != nil {
        cmd.Process.Kill()
        return nil, fmt.Errorf("initialize plugin %s: %w", execPath, err)
    }

    return &PluginSource{
        name:   initResult.Name,
        cmd:    cmd,
        client: client,
        stdin:  stdin,
        stdout: stdout,
    }, nil
}

func (p *PluginSource) Name() string { return p.name }

func (p *PluginSource) Collect(ctx context.Context, opts extension.CollectOptions) ([]*extension.RawContent, error) {
    var result []*extension.RawContent
    if err := p.client.Call("collect", opts, &result); err != nil {
        return nil, fmt.Errorf("plugin %s collect: %w", p.name, err)
    }
    return result, nil
}

// Close 优雅关闭插件进程
func (p *PluginSource) Close() error {
    p.client.Call("shutdown", nil, nil)
    p.stdin.Close()
    return p.cmd.Wait()
}
```

#### Python 插件示例

```python
#!/usr/bin/env python3
"""
synapse-rss-source — RSS 数据源插件
通过 stdin/stdout JSON-RPC 与 Synapse 通信
"""

import json
import sys
import feedparser
from datetime import datetime

def handle_initialize(params):
    return {
        "name": "rss-source",
        "type": "source",
        "version": "1.0.0"
    }

def handle_collect(params):
    config = params.get("config", {})
    feeds = config.get("feeds", [])

    raw_contents = []
    for feed_url in feeds:
        feed = feedparser.parse(feed_url)
        for entry in feed.entries[:10]:  # 最近 10 条
            raw_contents.append({
                "title": entry.title,
                "content": entry.get("summary", ""),
                "source": feed_url,
                "timestamp": datetime.now().isoformat(),
                "metadata": {
                    "url": entry.link,
                    "author": entry.get("author", ""),
                }
            })

    return raw_contents

# JSON-RPC 主循环
handlers = {
    "initialize": handle_initialize,
    "collect": handle_collect,
}

for line in sys.stdin:
    request = json.loads(line.strip())
    method = request.get("method")

    if method == "shutdown":
        break

    handler = handlers.get(method)
    if handler:
        result = handler(request.get("params", {}))
        response = {"jsonrpc": "2.0", "id": request.get("id"), "result": result}
    else:
        response = {"jsonrpc": "2.0", "id": request.get("id"),
                     "error": {"code": -32601, "message": f"unknown method: {method}"}}

    sys.stdout.write(json.dumps(response) + "\n")
    sys.stdout.flush()
```

#### 引擎对外部插件的透明支持

Engine 在初始化时，如果发现 config 中有 `plugin` 字段，自动切换为外部插件模式：

```go
// internal/engine/engine.go

func (e *Engine) initSource(cfg SourceConfig) (extension.Source, error) {
    if cfg.Plugin != "" {
        // 外部插件：通过子进程 + JSON-RPC 通信
        return plugin.NewPluginSource(cfg.Plugin, cfg.Config)
    }
    // 内置扩展：直接从 Registry 获取
    return extension.GetSource(cfg.Name, cfg.Config)
}
```

**对用户完全透明**——用户只需要在 config.yaml 中声明，不需要知道插件是内置的还是外部的。

### 插件发现与安装

#### 内置插件

开箱即用，无需安装。`synapse plugin list` 即可查看：

```bash
$ synapse plugin list

Sources:
  ✅ skill-source      (built-in)  CodeBuddy Skill 数据源
  ✅ cli-import-source (built-in)  CLI 手动导入

Processors:
  ✅ skill-processor   (built-in)  Skill Prompt 处理引擎

Stores:
  ✅ local-store       (built-in)  本地文件系统存储

Indexers:
  ✅ bm25-indexer      (built-in)  BM25 全文检索

Consumers:
  ✅ hugo-consumer     (built-in)  Hugo 静态网站
  ✅ obsidian-consumer (built-in)  Obsidian Vault 兼容

Auditors:
  ✅ default-auditor   (built-in)  默认审计策略
```

#### 外部插件安装

```bash
# 从 GitHub 安装社区插件
$ synapse plugin install github.com/community/synapse-rss-source

Installing rss-source v1.0.0...
  → Downloaded to ~/.synapse/plugins/rss-source
  → Verified checksum: ok
  → Plugin type: source
  → Ready to use!

Add to your .synapse/config.yaml:
  sources:
    - name: rss-source
      plugin: "~/.synapse/plugins/rss-source"
      config:
        feeds: ["https://..."]
```

插件安装目录结构：

```
~/.synapse/
├── plugins/                        # 外部插件安装目录
│   ├── rss-source                  # 可执行文件
│   ├── notion-source               # 可执行文件
│   └── vector-indexer              # 可执行文件
└── registry.json                   # 已安装插件的元数据
```

### 扩展点接口的设计约束

为了保证内置和外部插件的一致性，接口设计有以下约束：

#### 1. 所有数据结构必须 JSON 可序列化

```go
// ✅ 好：所有字段都是 JSON 友好的
type RawContent struct {
    Title     string            `json:"title"`
    Content   string            `json:"content"`
    Source    string            `json:"source"`
    Timestamp time.Time         `json:"timestamp"`
    Metadata  map[string]any    `json:"metadata,omitempty"`
}

// ❌ 坏：包含不可序列化的类型
type RawContent struct {
    Reader io.Reader  // 不能 JSON 序列化
}
```

#### 2. 接口方法签名统一模式

```go
// 统一模式：(ctx, input) → (output, error)
// 这样可以直接映射为 JSON-RPC call
type Source interface {
    Name() string
    Collect(ctx context.Context, opts CollectOptions) ([]*RawContent, error)
}
```

#### 3. 可选能力通过接口组合

```go
// 基础接口（所有实现必须满足）
type Store interface {
    Read(ctx context.Context, path string) (*KnowledgeFile, error)
    Write(ctx context.Context, path string, file *KnowledgeFile) error
    Delete(ctx context.Context, path string) error
    List(ctx context.Context, dir string, opts ListOptions) ([]*FileInfo, error)
    Exists(ctx context.Context, path string) (bool, error)
}

// 可选能力接口（通过类型断言检查）
type VersionedStore interface {
    Store
    Commit(ctx context.Context, message string) error
    History(ctx context.Context, path string) ([]*ChangeRecord, error)
}

type SearchableStore interface {
    Store
    Search(ctx context.Context, query string, opts SearchOptions) ([]*SearchResult, error)
}

// 使用时通过类型断言
func doSomething(s extension.Store) {
    if vs, ok := s.(extension.VersionedStore); ok {
        // 支持版本控制
        vs.Commit(ctx, "update knowledge")
    }
}
```

### 总结：三层集成模型

```
┌──────────────────────────────────────────────────────────────────┐
│                   扩展点集成全景                                    │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  配置层：.synapse/config.yaml                               │  │
│  │  用户声明使用哪些扩展点、哪个实现                                │  │
│  └────────────────────────┬───────────────────────────────────┘  │
│                           │                                      │
│  ┌────────────────────────▼───────────────────────────────────┐  │
│  │  编排层：Engine                                              │  │
│  │  读取配置 → 实例化扩展点 → 协调执行                             │  │
│  └──────┬──────────────────────────────────┬─────────────────┘  │
│         │                                  │                     │
│  ┌──────▼────────────┐           ┌─────────▼──────────────┐     │
│  │  内置扩展           │           │  外部插件                │     │
│  │  (Built-in)        │           │  (Local Plugin)        │     │
│  │                    │           │                        │     │
│  │  Go Interface      │           │  子进程 + JSON-RPC     │     │
│  │  + init() 注册     │           │  + PluginAdapter       │     │
│  │  + Registry 查找   │           │                        │     │
│  │                    │           │  Go / Python / Node    │     │
│  │  进程内直接调用      │           │  / Rust / ...          │     │
│  │  零开销             │           │                        │     │
│  └────────────────────┘           └────────────────────────┘     │
│                                                                  │
│  对 Engine 和用户来说，内置 vs 外部 完全透明                          │
│  唯一的区别：config.yaml 中有没有 plugin 字段                        │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 落地节奏

| 阶段 | 集成方式 | 内容 | 插件生态 |
|------|---------|------|---------|
| **M1** | 定义接口 | 六个扩展点 Go Interface + Registry 框架 + Engine 骨架 | — |
| **M2** | 内置扩展 | skill-source + skill-processor + local-store（全部内置，零外部依赖） | — |
| **M3** | 内置 + 外部 | 新增 mcp-processor + github-store + bm25-indexer（内置）；引入 PluginAdapter + JSON-RPC 协议 | **插件管理 CLI** + 多来源安装 + 插件清单规范 |
| **M4** | 外部为主 | 各平台 Source 可以是内置或外部插件 | 外部 Source 可作为插件 |
| **M5** | 混合 | Hugo/Obsidian Consumer（内置）；社区 Consumer 可以是外部插件 | 外部 Consumer 可作为插件 |
| **M6+** | 插件市场 | `synapse plugin install` + 插件注册中心 | **Marketplace 三阶段演进** + 插件开发 SDK + 热重载 |

### 与之前设计的对齐

这个集成机制和之前的"扩展点模型"完全兼容——只是把之前**概念上的扩展点**落实为**具体的代码结构和通信机制**：

| 之前（概念） | 现在（实现） |
|-------------|-------------|
| "Source 是一个 Go 接口" | `pkg/extension/source.go` 定义接口 |
| "社区可以贡献 Source" | 内置：实现接口 + init() 注册；外部：实现 JSON-RPC 协议 |
| "配置文件选择实现" | Engine 读 config.yaml，按 name 从 Registry 查找或启动外部进程 |
| "正交替换" | 改 config.yaml 一个字段即可切换实现，Engine 无需修改 |

---

## 2026-04-18 插件机制调研：OpenClaw vs Claude Code vs Synapse 综合对比

### 调研目的

为 Synapse 的三层集成模型（内置扩展 / 本地插件 / 远程插件）提供设计参考，通过深度分析两个成熟项目的插件机制，提取可借鉴的模式和需要规避的复杂度。

---

### 一、三者架构概览

| 维度 | OpenClaw | Claude Code | Synapse（规划） |
|------|----------|-------------|-----------------|
| **语言** | TypeScript / Node.js | TypeScript / Node.js | Go |
| **插件运行模型** | 进程内代码注入（IoC） | 声明式清单 + 外部进程执行 | 混合：进程内接口 + 子进程 JSON-RPC |
| **插件与核心关系** | 共享进程，直接调用 API | 隔离进程，通过 Shell/HTTP/Prompt/MCP 通信 | Layer 1 共享进程；Layer 2 隔离进程 |
| **注册机制** | `register(api)` 回调 + 10 种注册通道 | `plugin.json` 清单 + 目录约定 | `init()` 自注册 + Registry + config.yaml 声明 |
| **插件发现** | 四层目录扫描（config > workspace > global > bundled） | 三种来源（Marketplace > Session > Builtin） | 两层（内置 + `~/.synapse/plugins/`） |
| **配置驱动** | `openclaw.plugin.json` + JSON Schema | `plugin.json` + Zod Schema + 多层设置优先级 | `.synapse/config.yaml` + Go struct 验证 |
| **Hook 系统** | 24 种生命周期钩子 | 26 种事件 + 4 种执行类型 | 暂无规划 |

---

### 二、核心维度对比

#### 2.1 插件注册与发现

##### OpenClaw — IoC 控制反转

```
核心 → 扫描发现 → 加载模块 → 调用 register(api) → 插件自主注册
```

- **四层优先级发现**：config 目录 > workspace 目录 > 全局目录 > 内置目录
- **Manifest 清单**：`openclaw.plugin.json` 声明元数据，JSON Schema 校验配置
- **10 种注册通道**：Tool / Hook / Channel / Provider / Service / Command / HttpRoute / GatewayMethod / CLI / ContextEngine
- **插件通过 `register(api)` 回调主动注册**，核心提供 `OpenClawPluginApi` 对象
- **典型代码**：
  ```typescript
  export const register: OpenClawPluginDefinition["register"] = (api) => {
    api.registerTool(myTool);
    api.registerHook("onAgentStart", myHandler);
  };
  ```

##### Claude Code — 声明式清单

```
插件目录 → plugin.json 清单 → 目录约定扫描 → Markdown/JSON 文件加载
```

- **三种来源**：Marketplace（Git 仓库浅克隆/NPM 安装）、Session（`--plugin-dir`）、Builtin（`name@builtin`）
- **`plugin.json` 清单**声明 6 种组件路径（commands/agents/skills/hooks/output-styles/MCP/LSP servers）
- **Markdown Frontmatter 接口**：commands/skills/agents 通过 YAML frontmatter 声明元数据
- **没有进程内 API 注入**，插件纯粹是内容+配置，核心系统负责解释和执行
- **典型结构**：
  ```
  .claude-plugin/
  ├── plugin.json       # 清单：声明所有组件
  ├── commands/         # 斜杠命令 (.md)
  ├── agents/           # AI Agent (.md)
  ├── skills/           # 技能 (SKILL.md)
  ├── hooks/hooks.json  # Hook 配置
  └── .mcp.json         # MCP 服务器配置
  ```

##### Synapse — Registry Pattern

```
init() 注册工厂 → Engine 读 config.yaml → Registry 查找/子进程启动
```

- **内置扩展**：Go Interface + `init()` 自注册 + Registry 按名称查找
- **外部插件**：子进程 + JSON-RPC over stdin/stdout + PluginAdapter 适配
- **配置驱动**：`config.yaml` 中有 `plugin` 字段则启动外部进程，否则走 Registry
- **对 Engine 透明**：内置和外部实现统一接口，Engine 无需区分

#### 2.2 插件-核心通信机制

| 通信方式 | OpenClaw | Claude Code | Synapse |
|---------|----------|-------------|---------|
| **进程内直接调用** | ✅ 主要方式（共享 V8） | ❌ 不支持 | ✅ Layer 1 内置扩展 |
| **子进程 stdin/stdout** | ❌ 不支持 | ✅ Hook command 类型 | ✅ Layer 2 JSON-RPC |
| **HTTP 请求** | ❌ 不支持 | ✅ Hook http 类型 | 🔮 Layer 3（远期） |
| **MCP 协议** | ❌ 不支持 | ✅ 集成 MCP Server | 🔮 远期 |
| **LSP 协议** | ❌ 不支持 | ✅ 集成 LSP Server | ❌ 不适用 |
| **LLM Prompt 评估** | ❌ 不支持 | ✅ Hook prompt 类型 | ❌ 不适用 |
| **Symbol 全局共享** | ✅ `Symbol.for()` 跨模块单例 | ❌ 不使用 | ❌ 不适用（Go 用 sync.Once） |

**关键差异**：
- OpenClaw 走的是**进程内强耦合**路线，性能好但安全性依赖信任
- Claude Code 走的是**进程外声明式**路线，安全隔离好但通信开销大
- Synapse 走**混合路线**：Layer 1 进程内（零开销），Layer 2 进程外（隔离性好）

#### 2.3 Hook 系统

| 特性 | OpenClaw | Claude Code | Synapse |
|------|----------|-------------|---------|
| **Hook 数量** | 24 种 | 26 种事件 | 暂无 |
| **执行模型** | 进程内函数调用 | 外部进程（Shell/HTTP/Prompt/Agent） |  — |
| **并行/顺序** | `runVoidHook`（并行）<br>`runModifyingHook`（顺序） | 同步（阻塞等待）<br>异步（后台运行） | — |
| **Hook 过滤** | 按 Hook 名称注册 | Matcher 模式（toolName/command 过滤） | — |
| **Hook 修改能力** | 可修改传入参数和返回值 | 可修改（approve/reject/modify 输入输出） | — |
| **排他性** | 支持 Exclusive Slot | 不支持 | — |

**OpenClaw 的 Hook 覆盖 6 大类**：
- Agent（start/end/think/response/error）
- Message（before/after）
- Tool（before/after/added）
- Session（init/end）
- Subagent（before/after）
- Gateway（request/response）

**Claude Code 的 Hook 覆盖更广泛**：
- Tool（PreToolUse/PostToolUse/PostToolUseFailure）
- Session（SessionStart/SessionEnd）
- User（UserPromptSubmit/PermissionRequest/PermissionDenied）
- Agent（SubagentStart/SubagentStop/Stop/StopFailure）
- System（Notification/ConfigChange/CwdChanged/FileChanged）
- Task（TaskCreated/TaskCompleted/TeammateIdle）
- Context（PreCompact/PostCompact/InstructionsLoaded）
- Setup/Elicitation/WorktreeCreate/WorktreeRemove

#### 2.4 安全模型

| 安全特性 | OpenClaw | Claude Code | Synapse |
|---------|----------|-------------|---------|
| **进程隔离** | ❌ 共享进程 | ✅ 外部进程 | Layer 2 ✅ 隔离 |
| **路径安全** | ✅ 路径逃逸检测 + 世界可写检查 + 所有权验证 | ✅ `checkPathTraversal()` | 🔮 需设计 |
| **名称防冒充** | ❌ 无 | ✅ `isBlockedOfficialName()` + 同形文字攻击检测 | 🔮 需设计 |
| **企业策略** | ❌ 无 | ✅ `isPluginBlockedByPolicy()` + `strictPluginOnlyCustomization` | ❌ 不适用（个人工具） |
| **敏感数据保护** | ❌ 无明确机制 | ✅ Keychain 分级存储 + `sensitive` 字段 | 🔮 config.yaml 密钥引用 |
| **Hook 安全** | ✅ Prompt 注入控制 | ✅ Agent 权限限制（不能设置 permissionMode/hooks/mcpServers） | — |
| **自动脱列** | ❌ 无 | ✅ `detectDelistedPlugins()` 自动卸载 | ❌ 不适用 |

#### 2.5 配置与生命周期

| 维度 | OpenClaw | Claude Code | Synapse |
|------|----------|-------------|---------|
| **配置格式** | JSON（openclaw.plugin.json） | JSON（plugin.json）+ Markdown frontmatter | YAML（config.yaml） |
| **配置验证** | JSON Schema | Zod Schema（TypeScript 运行时验证） | Go struct tag（静态验证） |
| **热重载** | ❌ 需重启 | ✅ `/reload-plugins` 命令 + 设置变更监听 | 🔮 远期（watch config） |
| **缓存策略** | TTL 缓存（Manifest Registry） | `lodash memoize()` + `clearPluginCache()` | 🔮 需设计 |
| **版本化** | 无版本化缓存 | ✅ `~/.claude/plugins/cache/{marketplace}/{plugin}/{version}/` | 🔮 `~/.synapse/plugins/{name}@{version}` |
| **自动更新** | ❌ 无 | ✅ 官方 Marketplace 默认自动更新 | 🔮 `synapse plugin update` |
| **启用/禁用** | 白名单/黑名单 + 自动启用 | 多层设置优先级（6层合并） | config.yaml 显式声明 |

---

### 三、设计模式提炼

#### 3.1 共性模式（三者共有）

| 模式 | 说明 | 三者的实现 |
|------|------|-----------|
| **注册表模式** | 全局注册表管理扩展实现 | OC: `createPluginRegistry()`<br>CC: `BUILTIN_PLUGINS` Map<br>Syn: `Registry` + `init()` |
| **清单/配置驱动** | 插件通过声明式配置描述自身 | OC: `openclaw.plugin.json`<br>CC: `plugin.json`<br>Syn: `config.yaml` |
| **统一接口抽象** | 不同实现统一到相同接口 | OC: 10种注册通道<br>CC: 6种组件类型<br>Syn: 6个 Go Interface |
| **路径安全检查** | 防止路径遍历攻击 | OC: 路径逃逸检测<br>CC: `checkPathTraversal()`<br>Syn: 需实现 |

#### 3.2 OpenClaw 特有模式

| 模式 | 说明 | Synapse 可借鉴度 |
|------|------|-----------------|
| **IoC API 注入** | 核心向插件注入 API 对象，插件主动注册 | ⭐⭐⭐ Layer 1 内置扩展已类似（init() 注册） |
| **排他性槽位** | 某些类型只允许一个活跃实现 | ⭐⭐⭐ Synapse 每个扩展点本就只选一个实现 |
| **Core Bridge 延迟导入** | 外部插件通过动态 import 松耦合访问核心 | ⭐⭐ 不适用 Go（编译时链接），但 JSON-RPC 已实现松耦合 |
| **策略模式 Provider** | 如 VoiceCallProvider 多实现切换 | ⭐⭐⭐ 直接映射到 Synapse 的 Store/Processor 等接口 |
| **适配器模式 Channel** | 20+ 个通道适配器 | ⭐⭐ Source 扩展点本质相同 |

#### 3.3 Claude Code 特有模式

| 模式 | 说明 | Synapse 可借鉴度 |
|------|------|-----------------|
| **声明式 Markdown 接口** | 用 Markdown + frontmatter 定义能力 | ⭐ 不适用（Synapse 不是 AI Agent 框架） |
| **多层设置合并** | 6 层优先级配置合并 | ⭐⭐ config.yaml 可支持 default + user + project 三层 |
| **Marketplace 系统** | Git 浅克隆 + NPM 安装 + 版本化缓存 | ⭐⭐⭐ M6+ 阶段的插件市场直接参考 |
| **Memoize 缓存** | 统一缓存 + 统一清理 | ⭐⭐ Go 中用 `sync.Once` + 缓存失效机制 |
| **Hook stdin/stdout 协议** | Shell 命令通过 JSON stdin 输入 / JSON stdout 输出 | ⭐⭐⭐ 与 Synapse Layer 2 JSON-RPC 高度相似！ |
| **模板变量替换** | `${PLUGIN_ROOT}` / `${user_config.KEY}` 等 | ⭐⭐ config.yaml 支持 `${env.XXX}` 变量替换 |
| **脱列检测 + 自动卸载** | 防止恶意/失效插件残留 | ⭐⭐ 远期插件市场需要 |
| **Seed 目录预置** | `PLUGIN_SEED_DIR` 容器化部署 | ⭐ 个人工具暂不需要 |
| **命名空间隔离** | MCP Server 添加 `plugin:{name}:` 前缀 | ⭐⭐⭐ 外部插件需要命名空间隔离 |

---

### 四、Synapse 设计决策建议

基于 OpenClaw 和 Claude Code 的分析，对 Synapse 三层集成模型提出以下设计建议：

#### 4.1 Layer 1 内置扩展 — 保持现有设计，借鉴排他性槽位

**现有设计已经很好**，与 OpenClaw 的 IoC + Registry 模式高度吻合：

```go
// 已有：init() 自注册 + Registry 按名查找
func init() {
    extension.RegisterSource("skill-source", NewSkillSource)
}
```

**补充建议**：
- ✅ 借鉴 OpenClaw 的**排他性槽位**思想 — Synapse 每个扩展点在 config.yaml 中只选一个实现，天然满足
- ✅ 借鉴 Claude Code 的**内置插件启用/禁用** — `BUILTIN_PLUGINS` Map 按 settings 过滤

#### 4.2 Layer 2 本地插件 — 强化 JSON-RPC 协议，借鉴 Claude Code 的 stdin/stdout 协议

**现有 JSON-RPC over stdin/stdout 设计**与 Claude Code 的 Hook command 类型高度相似：

| Synapse Layer 2 | Claude Code Hook (command) |
|-----------------|---------------------------|
| JSON-RPC request → stdin | `$ARGUMENTS` env + stdin |
| JSON-RPC response ← stdout | stdout JSON lines |
| `initialize` / `collect` / `shutdown` | 自由 shell 命令 |

**补充建议**：

1. **添加健康检查协议**：
   ```go
   // 借鉴 Claude Code 的 Hook 超时机制
   type PluginHealth struct {
       Status    string `json:"status"`    // "healthy" | "degraded" | "unhealthy"
       Latency   int64  `json:"latency"`   // 最近一次调用的延迟 (ms)
       ErrorRate float64 `json:"errorRate"` // 近 N 次调用的错误率
   }
   ```

2. **添加能力协商**（借鉴 Claude Code 的组件类型声明）：
   ```go
   // initialize 响应中声明支持的能力
   type InitializeResult struct {
       Name         string   `json:"name"`
       Type         string   `json:"type"`         // "source" | "processor" | ...
       Version      string   `json:"version"`
       Capabilities []string `json:"capabilities"` // 可选能力列表
   }
   ```

3. **添加 Graceful Shutdown**（借鉴 Claude Code 的 SessionEnd hook）：
   ```
   Engine → "shutdown" → Plugin → cleanup → exit(0)
   Engine 等待 5s，超时则 SIGKILL
   ```

#### 4.3 插件清单格式 — 借鉴 Claude Code 的 plugin.json

为 Layer 2 外部插件设计清单格式：

```yaml
# ~/.synapse/plugins/rss-source/synapse-plugin.yaml
name: rss-source
version: 1.0.0
type: source                              # 扩展点类型
description: "RSS/Atom 数据源"
author: "community"
license: "MIT"
homepage: "https://github.com/xxx/synapse-rss-source"

# 运行时配置
runtime:
  command: "./rss-source"                 # 可执行文件
  protocol: "jsonrpc"                     # 通信协议
  timeout: 30s                            # 单次调用超时

# 插件配置 Schema（借鉴 OpenClaw 的 JSON Schema 验证）
config_schema:
  type: object
  properties:
    feeds:
      type: array
      items:
        type: string
        format: uri
      description: "RSS 订阅源列表"
  required: ["feeds"]
```

#### 4.4 插件市场（M6+）— 直接参考 Claude Code Marketplace

Claude Code 的 Marketplace 设计非常成熟，Synapse 可直接参考：

| Claude Code Marketplace | Synapse Plugin Market（远期） |
|------------------------|------------------------------|
| `marketplace.json` 描述多个插件 | `synapse-registry.json` 或 GitHub Topics 发现 |
| Git 浅克隆 / NPM 安装 | `go install` 或 GitHub Release 下载二进制 |
| 版本化缓存 `~/.claude/plugins/cache/` | 版本化缓存 `~/.synapse/plugins/{name}@{version}/` |
| 自动更新（官方 marketplace） | `synapse plugin update` 命令 |
| 脱列检测 + 自动卸载 | `synapse plugin audit` 检查 |
| 名称防冒充（同形文字检测） | 远期需要 |

#### 4.5 安全机制 — 融合两者的最佳实践

| 阶段 | 安全措施 | 参考来源 |
|------|---------|---------|
| **安装时** | 校验 checksum / 签名验证 | Claude Code（Git 克隆验证） |
| **加载时** | 路径遍历检查 | OpenClaw + Claude Code |
| **运行时** | 子进程隔离 + 超时控制 | Claude Code（Hook 超时） |
| **配置** | JSON Schema 校验插件配置 | OpenClaw（Manifest Schema） |
| **密钥** | config.yaml 支持 `${env.API_KEY}` 引用，不硬编码 | Claude Code（Keychain 分级） |

#### 4.6 暂不引入的复杂度

以下特性虽然在 OpenClaw/Claude Code 中存在，但对 Synapse 当前阶段**过度复杂**，建议远期再考虑：

| 特性 | 原因 |
|------|------|
| **Hook 系统** | Synapse 是数据流水线，不是 Agent 框架，暂不需要 24+ 种生命周期钩子 |
| **多层设置优先级** | 个人工具单 config.yaml 足够，不需要 6 层合并 |
| **Markdown 接口** | 不适用于 Go 数据处理管道 |
| **企业策略控制** | 个人知识中枢不需要 |
| **Seed 目录** | 不需要容器化部署场景 |
| **热重载** | 初期可重启，远期可加 watch |

---

### 五、结论：Synapse 三层模型的定位

```
                OpenClaw                    Claude Code                 Synapse
              ┌──────────────┐           ┌──────────────┐          ┌──────────────┐
              │  进程内 IoC    │           │  声明式清单    │          │  混合模型      │
              │  API 注入      │           │  外部进程执行  │          │              │
              │              │           │              │          │  Layer 1:    │
通信方式       │  共享 V8 进程  │           │  Shell/HTTP/  │          │  进程内接口    │
              │  直接函数调用   │           │  Prompt/Agent│          │              │
              │              │           │  MCP/LSP     │          │  Layer 2:    │
              │              │           │              │          │  子进程 RPC   │
              └──────────────┘           └──────────────┘          │              │
                                                                   │  Layer 3:    │
              更高性能                     更强隔离                    │  远程 HTTP   │
              更强能力                     更安全                     └──────────────┘
              更紧耦合                     更松耦合                    兼顾两者优势

适用场景       Node.js 全栈应用             AI Agent 编排               Go 数据管道
              需要深度集成的插件             安全敏感的插件生态            个人知识处理工具
```

**Synapse 的混合模型恰好取了两者的长处**：

1. **Layer 1（像 OpenClaw）**：内置扩展享受进程内零开销调用，适合核心功能
2. **Layer 2（像 Claude Code）**：外部插件通过子进程隔离，适合社区贡献的不受信任代码
3. **Layer 3（远期）**：远程 HTTP/gRPC，适合微服务化部署

**最重要的设计原则保持不变**：对 Engine 和用户完全透明——无论是哪一层的实现，都统一到同一个 Go Interface，改 `config.yaml` 一个字段即可切换。

---

### 六、插件仓与管理命令调研

在上一节对比了两者的插件机制（注册/发现/Hook/通信）之后，进一步调研了 **插件仓（Marketplace / Registry）** 和 **插件管理命令（Plugin CLI）** 两个维度。

---

#### 6.1 OpenClaw 插件仓

OpenClaw **没有独立的 Marketplace 浏览界面**，采用 **npm + 本地 JSON 目录** 的分散式方案：

| 机制 | 说明 |
|------|------|
| **npm 分发** | 官方插件以 `@openclaw/*` 发布到 npm，安装时直接 `npm install --ignore-scripts` |
| **外部插件目录（Channel Catalog）** | 支持从 JSON 文件加载插件元数据，路径：`~/.openclaw/mpm/plugins.json`、`~/.openclaw/mpm/catalog.json` 等，也可通过 `OPENCLAW_PLUGIN_CATALOG_PATHS` 环境变量指定 |
| **社区插件列表** | 通过 PR 提交到 `docs/plugins/community.md`（目前仅 1 个：`@icesword760/openclaw-wechat`） |
| **40 个内置扩展** | 在 `extensions/` 目录下，作为 PNPM workspace 包管理 |
| **Bundled 插件** | 裸名安装时优先使用本地 bundled 版本（`bundled-sources.ts` + `plugin-install-plan.ts`） |

配置结构（`PluginsConfig`）：`enabled` / `allow` / `deny` / `load.paths` / `slots` / `entries` / `installs`

排他性槽位：`slots.memory`（memory-core / memory-lancedb / none）、`slots.contextEngine`

#### 6.2 OpenClaw 插件管理命令

完整 CLI 命令集（`src/cli/plugins-cli.ts`，827 行）：

```bash
openclaw plugins list           # 列出已安装/可用插件
openclaw plugins info <id>      # 查看插件详情
openclaw plugins install <spec> # 安装（支持 4 种来源）
openclaw plugins uninstall <id> # 卸载
openclaw plugins enable <id>    # 启用
openclaw plugins disable <id>   # 禁用
openclaw plugins update [id]    # 更新（仅 npm 来源）
openclaw plugins doctor         # 诊断插件健康状态
```

**四种安装来源**：

| 来源 | 示例 | 实现文件 |
|------|------|---------|
| npm spec | `@openclaw/memory-core` | `install.ts` → `installPluginFromNpmSpec()` |
| 本地目录 | `./my-plugin` | `install.ts` → `installPluginFromDir()` |
| 归档文件 | `plugin.tgz` (.zip/.tgz/.tar.gz/.tar) | `install.ts` → `installPluginFromArchive()` |
| 单文件 | `handler.ts` (.ts/.js) | `install.ts` → `installPluginFromFile()` |

安装路径：`~/.openclaw/extensions/<pluginId>/`（CONFIG_DIR/extensions）

安装安全：`npm install --ignore-scripts` + 代码安全扫描（warn-only）+ integrity drift 检测

卸载逻辑（`uninstall.ts`）：清理 entries/installs/allow/load.paths/slots.memory → 删除目录（链接类 source=path 不删源目录）

更新逻辑（`update.ts`）：仅支持 npm 来源，有 integrity drift 检测

#### 6.3 Claude Code 插件仓

Claude Code 有一套 **成熟的、Git 驱动的 Marketplace 体系**（核心文件 `marketplaceManager.ts`，2644 行）：

| 机制 | 说明 |
|------|------|
| **官方 Marketplace** | `github:anthropics/claude-plugins-official`，优先从 GCS 镜像获取 |
| **多来源类型** | 支持 `github` / `git` / `url` / `npm` / `file` / `directory` / `settings` 7 种来源 |
| **两层架构（Intent → State）** | Intent 层 = `settings.json` 中的 `extraKnownMarketplaces` + `enabledPlugins`；State 层 = `~/.claude/plugins/known_marketplaces.json` + `~/.claude/plugins/cache/` |
| **Reconciler 协调器** | `reconciler.ts` 对比 declared intent 与 materialized state，执行安装/更新/跳过 |
| **版本化缓存** | `~/.claude/plugins/cache/<marketplace>/<plugin>/<version>/` |
| **版本优先级** | `plugin.json` version > marketplace entry version > Git commit SHA（前 12 位） > `'unknown'` |
| **延迟 GC** | 旧版本标记 `.orphaned_at` 后 7 天才删除（`cacheUtils.ts`） |
| **企业策略** | `strictKnownMarketplaces`（白名单）+ `blockedMarketplaces`（黑名单），支持 `hostPattern` + `pathPattern` 正则匹配 |
| **防仿冒** | 保留名称集合 `ALLOWED_OFFICIAL_MARKETPLACE_NAMES` + 正则检测 + 非 ASCII 检测 + GitHub org 验证 |
| **Seed 目录** | `CLAUDE_CODE_PLUGIN_SEED_DIR` 预置只读插件（容器化部署场景） |
| **Settings-First 架构** | 安装/启用操作先写 settings（声明意图），再执行物化（缓存/克隆） |
| **后台安装管理** | `PluginInstallationManager.ts` 启动时自动安装/更新 Marketplace（`reconcileMarketplaces`） |

Git 浅克隆机制：`--depth 1`、SSH/HTTPS 自动切换（`isGitHubSshLikelyConfigured()`）、Sparse Checkout 支持（monorepo）、120s 超时

#### 6.4 Claude Code 插件管理命令

**外部 CLI 命令**（`src/main.tsx` 注册 + `src/cli/handlers/plugins.ts` 实现）：

```bash
# 插件管理
claude plugin validate <path>       # 验证插件格式
claude plugin list                  # 列出插件
claude plugin install <source>      # 安装
claude plugin uninstall <id>        # 卸载
claude plugin enable <id>           # 启用
claude plugin disable <id>          # 禁用
claude plugin update [id]           # 更新

# Marketplace 管理（独立子命令组）
claude plugin marketplace add <source>     # 添加 Marketplace 源
claude plugin marketplace list             # 列出已添加的 Marketplace
claude plugin marketplace remove <source>  # 移除 Marketplace 源
claude plugin marketplace update [source]  # 更新 Marketplace
```

**REPL 内交互式命令**（`/plugin`，别名 `/plugins`、`/marketplace`）：

```
/plugin install <source>          # 安装
/plugin manage                    # 交互式管理菜单
/plugin uninstall <id>            # 卸载
/plugin enable <id>               # 启用
/plugin disable <id>              # 禁用
/plugin validate <path>           # 验证
/plugin marketplace add <source>  # 添加 Marketplace
/plugin marketplace remove        # 移除 Marketplace
/plugin marketplace update        # 更新 Marketplace
/plugin marketplace list          # 列出 Marketplace
/plugin help                      # 帮助
/reload-plugins                   # 热重载所有插件
```

**三作用域系统**：安装/启用操作支持 `--scope user|project|local`

**核心操作库**（`pluginOperations.ts`，1089 行）：`installPluginOp()` / `uninstallPluginOp()` / `setPluginEnabledOp()` / `enablePluginOp()` / `disablePluginOp()` / `updatePluginOp()`

Settings-First 流程：操作 → 先写 settings 声明意图 → 再执行物化（Git 克隆 / NPM 安装 / 缓存）

---

#### 6.5 核心差异对比

| 维度 | OpenClaw | Claude Code |
|------|----------|-------------|
| **插件仓类型** | npm + 本地 JSON 目录，无在线 UI | 完整 Marketplace 系统，Git 驱动 |
| **官方仓库** | npm `@openclaw/*` | `github:anthropics/claude-plugins-official` |
| **来源类型** | 4 种（npm / 目录 / 归档 / 单文件） | 7 种（github / git / url / npm / file / directory / settings） |
| **CLI 命令数** | 8 个子命令 | 11 个子命令（含 marketplace 子组） |
| **REPL 命令** | ❌ 无 | ✅ `/plugin` 交互式命令 + `/reload-plugins` 热重载 |
| **作用域** | 单一全局配置 | 三层（user / project / local） |
| **版本管理** | npm semver | 自定义版本计算 + Git SHA fallback |
| **缓存策略** | npm 缓存 | 版本化目录 + 延迟 GC（7 天） |
| **企业管控** | allow / deny list | 白名单 / 黑名单 + 正则匹配 + 防仿冒 |
| **安全机制** | `--ignore-scripts` + 代码扫描 | Settings-First + Reconciler 校验 |
| **后台管理** | ❌ 无 | ✅ `PluginInstallationManager` 启动时自动安装/更新 |
| **一致性模型** | 即时生效 | Intent → State 两层架构 + Reconciler 协调 |

---

#### 6.6 对 Synapse 的设计启示

| 启示 | 说明 | 参考来源 | 落地阶段 |
|------|------|---------|---------|
| **Intent → State 分离** | 配置声明（config.yaml）与物化状态（~/.synapse/plugins/cache/）分离，更安全可控 | Claude Code | M3+ |
| **多来源支持** | 至少支持 Go module + 本地目录 + Git 仓库 三种来源 | 两者综合 | M3+ |
| **CLI 命令体系** | `synapse plugin list/install/uninstall/enable/disable/update/doctor` | 两者综合 | M3+ |
| **版本化缓存** | `~/.synapse/plugins/cache/<name>/<version>/`，带延迟 GC | Claude Code | M6+ |
| **轻量级 Catalog** | 初期用 JSON 文件做插件目录，远期演进为在线 Marketplace | OpenClaw | M3+ |
| **安全安装** | 校验 checksum + 路径遍历检查 + 子进程隔离 | 两者综合 | M3+ |
| **Reconciler 机制** | 启动时自动对比配置与实际状态，执行安装/更新/清理 | Claude Code | M6+ |
| **作用域** | 支持全局（~/.synapse/config.yaml）+ 项目（.synapse/config.yaml）两级 | Claude Code 简化 | M3+ |

---

### 七、后续行动项

- [ ] 设计 `synapse-plugin.yaml` 清单格式规范
- [ ] 实现 Layer 2 PluginAdapter 的健康检查和能力协商
- [ ] 设计 config.yaml 的变量替换机制（`${env.XXX}`）
- [ ] 设计 `synapse plugin *` CLI 命令体系（参考 OpenClaw + Claude Code）
- [ ] 设计插件安装来源机制（Go module + 本地目录 + Git 仓库）
- [ ] 设计 Intent → State 配置分离模型
- [ ] M6+ 阶段：参考 Claude Code Marketplace 设计插件市场

---

## 待讨论

（下一轮讨论内容待补充）
