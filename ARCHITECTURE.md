# Synapse 架构说明

本文档描述 Synapse 的整体架构设计、核心概念和技术决策，帮助开发者快速理解项目的技术全貌。

---

## 📐 架构范式：扩展点模型（Extension Point Model）

Synapse 采用**星型架构**，以 Store 为底座，六个独立扩展点按需组合：

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
                                          └───────────────┘

    贯穿所有模块：Schema（知识规范）— 统一的"语言"
```

### 为什么是扩展点模型？

| 特性 | 说明 |
|------|------|
| **不强制分层** | Source 和 Processor 可以合并为一个实现，也可以分开 |
| **不强制线性** | Consumer 直接读 Store，不必等前面所有步骤完成 |
| **按需组合** | 最小集合只需 Source + Processor + Store，其他扩展点可选 |
| **正交替换** | 换 Indexer 不影响 Processor，换 Store 不影响 Source |

---

## 🧩 核心概念

### Knowledge Schema

`.synapse/schema.yaml` 是知识库的"行为契约"，定义：

- **页面类型**：7 种（profile / topic / entity / concept / inbox / journal / graph）
- **Frontmatter 规范**：每种页面类型的元数据结构
- **工作流规则**：知识处理的自动化规则
- **质量标准**：知识库健康检查的基准

所有扩展点共同遵守 Schema 定义的规范，修改 Schema 即修改所有 AI 助手的行为。

### RawContent

Source 输出的标准原始内容格式，是连接数据源与处理引擎的统一中间表示。任何数据源都必须将其原始数据转换为 RawContent 格式。

### KnowledgeFile

Processor 输出的结构化知识文件，符合 Knowledge Schema 规范，包含：

- Frontmatter 元数据（标题、标签、关联、创建时间等）
- Markdown 正文
- `[[双向链接]]`（兼容 Obsidian）

---

## 📁 项目结构

> 💡 以下为当前已实现的项目结构。标记 `(planned)` 的部分为规划中尚未实现。

```
synapse/
├── cmd/                    # 应用入口
│   └── synapse/
│       ├── main.go         # CLI 主程序（所有命令定义）
│       └── plugins.go      # 内置扩展点注册
├── internal/               # 核心逻辑（不对外暴露）
│   ├── engine/             # 编排引擎：协调各扩展点
│   ├── schema/             # Schema 加载与校验
│   ├── config/             # 配置加载与管理
│   │   ├── config.go       # 配置文件解析（YAML → SynapseConfig）
│   │   └── global.go       # 全局配置模板自动创建与管理
│   ├── initializer/        # 知识库初始化逻辑（通过 Store 接口实现）
│   ├── source/             # Source 扩展点
│   │   └── skill/          # 参考实现：CodeBuddy Skill Source
│   ├── processor/          # Processor 扩展点
│   │   └── skill/          # 参考实现：Skill Processor
│   ├── store/              # Store 扩展点
│   │   ├── local/          # 参考实现：Local Store（本地文件系统）
│   │   └── github/         # 参考实现：GitHub Store（远程仓库）
│   ├── indexer/            # Indexer 扩展点
│   │   └── bm25/           # 参考实现：BM25 Indexer (planned)
│   ├── consumer/           # Consumer 扩展点
│   │   ├── hugo/           # 参考实现：Hugo Consumer (planned)
│   │   └── obsidian/       # 参考实现：Obsidian Consumer (planned)
│   └── auditor/            # Auditor 扩展点
│       └── defaultauditor/ # 参考实现：Default Auditor (planned)
├── pkg/                    # 共享工具包
│   ├── model/              # 领域模型（RawContent, KnowledgeFile, SearchResult 等）
│   ├── link/               # 双向链接解析工具
│   └── extension/          # 扩展点接口定义与 Registry（注册中心）
├── skills/                 # AI 助手 Skill 文件
│   ├── codebuddy/          # CodeBuddy Skill Prompt
│   ├── claude-code/        # Claude Code Skill（SYNAPSE.md）
│   └── cursor/             # Cursor Skill（.cursorrules）
├── configs/                # 配置模板（预留目录）
├── docs/                   # 项目文档
│   ├── roadmap.md          # 路线图
│   └── decision.md         # 方案决策记录
└── test/                   # 测试工具与集成测试
    ├── fixtures/           # 测试数据
    └── integration/        # 集成测试
```

### 规划中的目录（M3 阶段）

```
├── api/                    # MCP Server 定义 (M3)
│   └── mcp/                # MCP 工具定义与处理器
├── internal/
│   └── plugin/             # 插件管理 (M3)
│       ├── adapter/        # Layer 2 PluginAdapter（JSON-RPC 通信）
│       ├── catalog/        # 插件目录管理
│       └── reconciler/     # Intent → State 协调器
└── pkg/
    └── jsonrpc/            # JSON-RPC 通信工具 (M3)
```

---

## 🔌 六大扩展点接口

> 所有扩展点接口统一定义在 `pkg/extension/extension.go` 中，并通过 `pkg/extension/registry.go` 提供注册中心。

### Source — 数据源

```go
// Source 从外部数据源获取原始内容
type Source interface {
    // Name 返回数据源名称
    Name() string

    // Fetch 获取原始内容
    Fetch(ctx context.Context, opts FetchOptions) ([]RawContent, error)
}
```

### Processor — 处理引擎

```go
// Processor 将原始内容处理为结构化知识
type Processor interface {
    // Name 返回处理引擎名称
    Name() string

    // Process 将原始内容转换为知识文件
    Process(ctx context.Context, raw RawContent) ([]KnowledgeFile, error)
}
```

### Store — 存储底座

```go
// Store 提供知识文件的持久化存储
type Store interface {
    // Name 返回存储后端名称
    Name() string

    // Read 读取知识文件
    Read(ctx context.Context, path string) (KnowledgeFile, error)

    // Write 写入知识文件
    Write(ctx context.Context, file KnowledgeFile) error

    // Delete 删除知识文件
    Delete(ctx context.Context, path string) error

    // List 列出指定目录下的知识文件
    List(ctx context.Context, dir string) ([]KnowledgeFile, error)
}
```

### Indexer — 检索引擎

```go
// Indexer 提供知识库检索能力
type Indexer interface {
    // Name 返回检索引擎名称
    Name() string

    // Index 索引知识文件
    Index(ctx context.Context, file KnowledgeFile) error

    // Search 搜索知识库
    Search(ctx context.Context, query string, opts SearchOptions) ([]SearchResult, error)
}
```

### Consumer — 消费端

```go
// Consumer 将知识输出为特定消费形式
type Consumer interface {
    // Name 返回消费端名称
    Name() string

    // Consume 消费知识库，生成输出产物
    Consume(ctx context.Context, store Store, opts ConsumeOptions) error
}
```

### Auditor — 质量审计

```go
// Auditor 检查知识库健康状态
type Auditor interface {
    // Name 返回审计器名称
    Name() string

    // Audit 执行知识库审计
    Audit(ctx context.Context, store Store) (AuditReport, error)
}
```

---

## 🧩 插件系统

> ⚠️ 插件系统（外部插件的安装/卸载/通信）规划在 M3 阶段实现。当前版本仅支持内置扩展点。

### 内置扩展 vs 外部插件

| 维度 | 内置扩展 | 外部插件 |
|------|---------|---------|
| **位置** | `internal/` 目录下 | `~/.synapse/plugins/` |
| **通信** | Go 函数调用 | 子进程 + JSON-RPC（stdin/stdout） |
| **语言** | Go | 任意（Go/Python/Node/Rust...） |
| **安装** | 编译时包含 | `synapse plugin install` 运行时安装 |

### 插件通信

外部插件通过 **子进程 + JSON-RPC** 通信（类似 MCP stdio 模式）：

```
synapse CLI ←→ [JSON-RPC over stdin/stdout] ←→ 插件子进程
```

### 插件管理流程

```
Intent（config.yaml 声明）
       │
       ▼
Reconciler（协调器）
       │
       ├── 对比已安装插件（State）
       │
       ├── 缺少 → 自动安装
       ├── 多余 → 标记清理
       └── 版本不匹配 → 自动更新
       │
       ▼
State（~/.synapse/plugins/cache/）
```

---

## 🔧 配置模型

### 全局配置

首次运行任意 `synapse` 命令时，自动在 `~/.synapse/config.yaml` 创建全局配置模板。用户编辑此文件选择扩展点实现：

```yaml
synapse:
  version: "1.0"
  sources:
    - name: "skill-source"
      enabled: true
  processor:
    name: "skill-processor"
  store:
    name: "local-store"
    config:
      path: "~/knowhub"
```

配置中使用 `${ENV_VAR}` 格式引用环境变量，避免硬编码敏感信息（如 GitHub Token）。

### 作用域

| 作用域 | 路径 | 用途 | 状态 |
|--------|------|------|------|
| 全局 | `~/.synapse/config.yaml` | 用户级默认配置 | ✅ 已实现 |
| 项目 | `.synapse/config.yaml` | 知识库级配置（优先级更高） | 🔵 M3 规划 |

> 当前版本只支持全局配置。M3 阶段将引入双作用域配置，项目级配置优先于全局配置。

### Extension Registry（扩展点注册中心）

`pkg/extension/` 提供统一的扩展点注册机制。所有内置扩展点在 `cmd/synapse/plugins.go` 中注册到 Registry，配置文件中通过 `name` 字段引用已注册的扩展点实现。`synapse check` 命令会验证配置中引用的扩展点是否已注册。

### Intent → State 分离（M3 规划）

- **Intent**：`config.yaml` 中声明的期望状态（"我想使用哪些扩展点和插件"）
- **State**：`~/.synapse/plugins/cache/` 中的实际物化状态（"实际安装了什么"）
- **Reconciler**：启动时自动对比 Intent vs State，执行安装/更新/清理

---

## 🛤️ 数据流

### 采集 → 处理 → 存储

```
外部数据源 → Source.Fetch() → RawContent → Processor.Process() → KnowledgeFile → Store.Write()
```

### 检索

```
用户查询 → Indexer.Search() → SearchResult[] → 知识文件路径 → Store.Read() → KnowledgeFile
```

### 消费

```
Store → Consumer.Consume() → 输出产物（静态网站 / Obsidian Vault / Anki 闪卡 / ...）
```

### 审计

```
Store → Auditor.Audit() → AuditReport（断链 / 孤儿页面 / 过时内容 / 重复知识）
```

---

## 📊 技术选型

| 项 | 选择 | 理由 |
|----|------|------|
| 核心语言 | **Go** | 性能好，适合 CLI 和 MCP Server，单二进制分发 |
| 插件通信 | **子进程 + JSON-RPC** | 跨语言、安全隔离、简单，类似 MCP stdio 模式 |
| 知识关联 | **`[[双向链接]]`** | 显式关联，兼容 Obsidian 生态 |
| 配置格式 | **YAML** | 可读性好，社区接受度高 |
| 默认存储 | **本地文件系统** | 离线可用、最简、隐私安全 |
| 默认检索 | **BM25** | 不需要外部依赖，轻量高效 |
| MCP Server | **Go + stdio/SSE** | 兼容主流 AI 助手 |

---

## 🔐 安全考虑

- **数据本地优先**：默认使用 Local Store，数据不出本地
- **插件沙箱**：外部插件运行在子进程中，通过 JSON-RPC 通信，天然隔离
- **Checksum 校验**：插件安装时校验完整性
- **路径遍历防护**：阻止插件访问其目录之外的文件
- **插件名防仿冒**：保留名称集 + 正则检测，防止恶意插件仿冒官方

---

## 📚 相关文档

- [路线图](docs/roadmap.md) — 项目里程碑和计划
- [方案决策记录](docs/decision.md) — 关键设计决策和讨论过程
- [贡献指南](CONTRIBUTING.md) — 如何参与项目开发
