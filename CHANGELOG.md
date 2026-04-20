# Changelog

本项目的所有重要变更都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，
版本号遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

---

## [Unreleased]

### Added

- 项目初始化：README、LICENSE、CONTRIBUTING、CODE_OF_CONDUCT、ARCHITECTURE
- 产品构思与方案讨论（docs/decision.md）
- 项目路线图（docs/roadmap.md）
  - 扩展点模型架构设计（6 大扩展点：Source / Processor / Store / Indexer / Consumer / Auditor）
  - 6 个里程碑规划（M1-M6+）
  - 插件生态设计（三阶段 Marketplace 演进）

- **Schema 增强**
  - `PageTypeDefinition` 新增 `Emoji` 字段（目录展示用的 emoji 图标）
  - `PageTypeDefinition` 新增 `Example` 字段（示例文件名，用于 README 表格展示）
  - Default Schema 为所有 7 种页面类型填充了默认 Emoji 和 Example

- **README 动态生成**
  - Store `tmpl/readme.go` 支持从 Schema 的 `PageTypeDefinition` 动态生成知识库目录结构表格
  - 知识库 README.md 中的目录表格不再硬编码，由 Schema 驱动生成

- **Skill Prompt 模板**（`skills/`）
  - CodeBuddy Skill（`skills/codebuddy/synapse-knowledge.md`）：最详细版本，包含页面类型规范、采集决策规则、触发词、示例命令
  - Claude Code Skill（`skills/claude/SYNAPSE.md`）：适配 CLAUDE.md 引用机制的简洁版
  - Cursor Skill（`skills/cursor/.cursorrules`）：适配 .cursorrules 的精简版
  - 通用模板（`skills/common/`）：共享知识规范定义
  - 三个模板均定义了 Retrieve（反哺）和 Collect（采集）双职责工作流

- **SECURITY.md** — 安全策略与漏洞报告流程
- **GitHub 模板** — Issue 模板（Bug Report / Feature Request）和 PR 模板

- **CLI 命令：`synapse install`**（`cmd/synapse/main.go`）
  - 支持三个目标平台：`codebuddy`、`claude`/`claude-code`、`cursor`
  - `--target` 参数指定目标项目目录（默认当前目录）
  - `--list` 列出所有支持的 AI 助手
  - `findSkillTemplate()` 多路径模板发现（可执行文件目录 → 当前目录 → `~/.synapse/skills/`）
  - Cursor 安装支持追加模式（不覆盖已有 .cursorrules）

---

## [v0.2.0] — M2 Skill 集成 🧠

### Added

- **Source 参考实现：Skill Source**（`internal/source/skill`）
  - 通过 `init()` 自注册为 `skill-source`
  - 支持 `FetchOptions.Config` 传入 content/title/session_id/suggested_topics/entities/concepts/key_points
  - 自动提取标题（取前 50 字符）
  - 扩展元数据透传

- **Processor 参考实现：Skill Processor**（`internal/processor/skill`）
  - 通过 `init()` 自注册为 `skill-processor`
  - 基于规则的知识分类：SuggestedTopics → topic 文件，SuggestedEntities → entity 文件，SuggestedConcepts → concept 文件
  - 无分类建议时自动放入 inbox
  - 交叉链接：自动为相关知识生成 `[[wiki-links]]`
  - 实体分类猜测：基于关键词判定 tool/person/organization/project
  - slug 生成：标题转 URL 友好文件名

- **CLI 命令：`synapse collect`**
  - 支持 `--content`、`--title`、`--topics`、`--entities`、`--concepts`、`--key-points` 参数
  - 支持 stdin 管道输入（`echo "content" | synapse collect`）
  - 运行完整 Source.Fetch → Processor.Process → Store.Write 管道

- **CLI 命令：`synapse search`**
  - M2 阶段基于文件遍历 + 文本匹配的简单搜索
  - 支持按页面类型过滤（`--type`）、结果数量限制（`--limit`）
  - 标题/标签/正文加权匹配

- **CLI 命令：`synapse audit`**
  - 基础审计：Frontmatter 完整性、断链检测、孤儿页面检测
  - 知识库健康评分（0-100 分）
  - 分类统计报告

- **Schema 加载机制**
  - Engine 支持 `WithSchema` Option 模式
  - 自动从 config.yaml 同目录加载 schema.yaml
  - 回退到默认 Schema

### Changed

- **Engine 重构**
  - `NewFromConfig` 支持 `Option` 可变参数
  - `Collect` 方法支持 `CollectOptions` 传递 FetchOptions
  - 新增 `Schema()` 和 `Config()` 访问器方法

- **Local Store 增强**
  - Frontmatter 解析从手动字符串匹配升级为 YAML 反序列化
  - 新增 `parseFrontmatterSimple()` 回退机制

### Tests

- Engine 单元测试（`internal/engine/engine_test.go`）
  - 测试 NewFromConfig、WithSchema、Collect 管道、Search/Audit/Publish 错误处理
- 集成测试（`test/integration/pipeline_test.go`）
  - 完整管道测试：Topic/Entity/Concept/Inbox 采集
  - 交叉链接验证
  - Store 读回验证

---

## [v0.1.0] — M1 基座搭建 🏗️

### Added

- **Knowledge Schema 规范**（`internal/schema/schema.go`）
  - 7 种页面类型定义：profile/topic/entity/concept/inbox/journal/graph
  - Frontmatter 字段规范（必填 + 可选）
  - 质量标准定义
  - Schema 加载、校验和默认值

- **六大扩展点接口**（`pkg/extension/extension.go`）
  - Source / Processor / Store / Indexer / Consumer / Auditor
  - VersionedStore 可选能力接口

- **全局注册表**（`pkg/extension/registry.go`）
  - Factory 模式 + init() 自注册
  - Register*/Get*/List* 系列 API
  - 线程安全（sync.RWMutex）

- **配置体系**（`internal/config/config.go`）
  - .synapse/config.yaml 扩展点注册中心
  - 支持启用/禁用扩展点、自定义配置

- **Core 领域模型**（`pkg/model/`）
  - KnowledgeFile + Frontmatter + PageType
  - RawContent（Source 输出格式）
  - SearchResult / AuditReport / ListOptions / FetchOptions

- **引擎编排**（`internal/engine/engine.go`）
  - 从配置实例化所有扩展点
  - Collect / Search / Audit / Publish 流程编排

- **CLI 框架**（`cmd/synapse/main.go`）
  - `synapse init` — 初始化 knowhub 仓库
  - `synapse plugin list` — 列出已注册插件

- **Store 参考实现：Local Store**（`internal/store/local/local.go`）
  - 基于本地文件系统的 CRUD
  - Frontmatter 解析和序列化
  - Markdown 文件列表、遍历

- **Wiki Links 解析工具**（`pkg/link/`）

---

## 版本规划

| 版本 | 里程碑 | 预计时间 |
|------|--------|---------|
| v0.1.0 | M1 基座搭建 — Schema 规范 + 扩展点接口 + CLI init | 2026 Q2 ✅ |
| v0.2.0 | M2 Skill 集成 — 第一个 Source + Processor + Store | 2026 Q2 ✅ |
| v0.3.0 | M3 MCP + 插件管理 — MCP Server + GitHub Store + 插件 CLI | 2026 Q3 |
| v0.4.0 | M4 多平台适配 — Claude Code / Cursor / ChatGPT Source | 2026 Q3 |
| v0.5.0 | M5 Consumer 实现 — Hugo 网站 + Obsidian 兼容 | 2026 Q4 |
| v1.0.0 | M6+ 社区生态 — 插件市场 + 全扩展点开放 | 2026 Q4+ |

[Unreleased]: https://github.com/tunsuy/synapse/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/tunsuy/synapse/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/tunsuy/synapse/releases/tag/v0.1.0
