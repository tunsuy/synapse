# 贡献指南

感谢你对 Synapse 项目的关注！我们欢迎任何形式的贡献，包括但不限于：

- 🐛 报告 Bug
- 💡 提出新功能建议
- 📖 改进文档
- 🔌 开发扩展点插件
- 🧪 补充测试用例
- 🔧 优化代码质量

---

## 📋 目录

- [行为准则](#行为准则)
- [如何开始](#如何开始)
- [开发环境搭建](#开发环境搭建)
- [开发工作流](#开发工作流)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [Pull Request 流程](#pull-request-流程)
- [扩展点插件开发](#扩展点插件开发)
- [Issue 指南](#issue-指南)

---

## 行为准则

参与本项目即表示你同意遵守我们的 [行为准则](CODE_OF_CONDUCT.md)。请确保在所有互动中保持尊重和包容。

---

## 如何开始

1. **Fork** 本仓库
2. **Clone** 你的 Fork：
   ```bash
   git clone https://github.com/<your-username>/synapse.git
   cd synapse
   ```
3. **添加上游远程仓库**：
   ```bash
   git remote add upstream https://github.com/tunsuy/synapse.git
   ```
4. **创建功能分支**：
   ```bash
   git checkout -b feature/your-feature-name
   ```

---

## 开发环境搭建

### 环境要求

- Go >= 1.24
- Git >= 2.30
- Make（可选，用于执行常用命令）

### 安装依赖

```bash
# 下载 Go 模块依赖
go mod download

# 验证依赖完整性
go mod verify
```

### 常用命令

```bash
# 运行所有测试
go test ./...

# 运行测试并查看覆盖率
go test -cover ./...

# 运行代码检查
go vet ./...

# 格式化代码
gofmt -w .
goimports -w .

# 构建项目
go build ./cmd/synapse
```

---

## 开发工作流

1. **同步上游代码**：
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **在功能分支上开发**：每个功能或修复对应一个分支

3. **编写测试**：所有新功能必须有对应的测试用例

4. **本地验证**：
   ```bash
   go test ./...
   go vet ./...
   ```

5. **提交代码**：遵循 [提交规范](#提交规范)

6. **发起 Pull Request**：

---

## 代码规范

本项目严格遵循 Go 社区的编码规范，详细要求如下：

### 格式化

- **必须**使用 `gofmt` 格式化所有 Go 代码
- **必须**使用 `goimports` 管理 import 语句
- 建议一行代码不超过 120 列

### 命名

- 文件名：小写 + 下划线分割，如 `local_store.go`
- 包名：小写单词，与目录名一致
- 变量 / 函数 / 结构体：驼峰命名，首字母按可见性大小写
- 接口：单方法接口以 `er` 结尾（如 `Reader`、`Processor`）

### 注释

- 所有导出的类型、函数、方法、常量、变量必须有 GoDoc 注释
- 格式：`// TypeName description`
- 每个包必须有包注释（`// Package xxx ...`）

### 错误处理

- `error` 必须是函数的最后一个返回参数
- 所有 `error` 必须显式处理或明确忽略
- 使用 `fmt.Errorf("context: %w", err)` 包装错误

### 项目结构

```
synapse/
├── cmd/            # 应用入口
├── internal/       # 核心逻辑（不对外暴露）
├── pkg/            # 共享工具包
├── api/            # gRPC/REST 定义
├── configs/        # 配置模板
├── docs/           # 项目文档
└── test/           # 测试工具和集成测试
```

### 测试

- 测试文件命名：`xxx_test.go`
- 使用表驱动测试（Table-Driven Tests）
- 测试函数命名：`TestXxx` 或 `Test_Xxx`
- 每个导出函数必须有对应的测试用例

---

## 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Type 类型

| Type | 说明 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `docs` | 文档变更 |
| `style` | 代码格式（不影响逻辑） |
| `refactor` | 重构（既不修复 bug 也不添加功能） |
| `perf` | 性能优化 |
| `test` | 测试相关 |
| `chore` | 构建过程或辅助工具变更 |

### Scope 范围

| Scope | 说明 |
|-------|------|
| `source` | Source 扩展点 |
| `processor` | Processor 扩展点 |
| `store` | Store 扩展点 |
| `indexer` | Indexer 扩展点 |
| `consumer` | Consumer 扩展点 |
| `auditor` | Auditor 扩展点 |
| `cli` | CLI 命令 |
| `schema` | Schema 规范 |
| `plugin` | 插件管理 |
| `mcp` | MCP Server |

### 示例

```
feat(source): add RSS source implementation

Implement Source interface for RSS feeds, supporting
Atom and RSS 2.0 formats.

Closes #42
```

```
fix(store): handle concurrent file write conflicts

Use file-level locking to prevent data corruption
when multiple goroutines write to the same knowledge file.
```

---

## Pull Request 流程

### PR 要求

1. **标题**：遵循提交规范格式
2. **描述**：清楚说明改动内容、动机和影响
3. **测试**：新功能必须包含测试，Bug 修复建议补充回归测试
4. **文档**：涉及 API 或行为变更时同步更新文档
5. **单一职责**：一个 PR 只做一件事

### PR 模板

```markdown
## 改动说明

简要描述本次改动的内容和动机。

## 改动类型

- [ ] 新功能（feat）
- [ ] Bug 修复（fix）
- [ ] 文档更新（docs）
- [ ] 重构（refactor）
- [ ] 测试（test）

## 测试方式

描述如何验证本次改动。

## 关联 Issue

Closes #xxx
```

### Review 流程

1. 提交 PR 后自动触发 CI（lint + test）
2. 至少需要 1 位 Maintainer 审核通过
3. CI 通过 + Review 通过后可合并
4. 使用 **Squash and Merge** 合并方式

---

## 扩展点插件开发

Synapse 的核心价值在于其扩展性。我们鼓励社区为任意扩展点开发插件。

### 插件类型

| 扩展点 | Go 接口 | 示例 |
|--------|---------|------|
| Source | `Source` | RSS Reader、Notion Sync |
| Processor | `Processor` | Local LLM、Rules Engine |
| Store | `Store` | S3 Store、WebDAV Store |
| Indexer | `Indexer` | Vector Search、Graph Traversal |
| Consumer | `Consumer` | VitePress、Anki Export |
| Auditor | `Auditor` | Custom Lint Rules |

### 开发步骤

1. 实现对应的 Go 接口
2. 创建 `synapse-plugin.yaml` 清单文件
3. 编写测试用例
4. 编写 README 说明文档
5. 提交 PR 或发布为独立 Go module

### 插件清单示例

```yaml
# synapse-plugin.yaml
name: rss-source
version: 1.0.0
description: RSS feed source for Synapse
author: your-name
extension_point: source
config_schema:
  feeds:
    type: array
    description: RSS feed URLs
    required: true
```

> 详细的插件开发指南将在 M3 里程碑完成后发布。

---

## Issue 指南

### Bug 报告

请包含以下信息：

- **环境**：Go 版本、操作系统、Synapse 版本
- **复现步骤**：最小化的复现步骤
- **预期行为**：你期望发生什么
- **实际行为**：实际发生了什么
- **日志/截图**：如有相关错误日志或截图

### 功能建议

请说明：

- **使用场景**：你在什么场景下需要这个功能
- **期望行为**：你希望它如何工作
- **替代方案**：你考虑过哪些替代方案

### 标签说明

| 标签 | 说明 |
|------|------|
| `good first issue` | 适合新贡献者的简单任务 |
| `help wanted` | 需要社区帮助 |
| `bug` | Bug 报告 |
| `enhancement` | 功能增强 |
| `documentation` | 文档改进 |
| `plugin` | 插件相关 |

---

## 🙏 致谢

感谢所有贡献者让 Synapse 变得更好！

---

> 如果你有任何问题，欢迎在 [GitHub Discussions](https://github.com/tunsuy/synapse/discussions) 中提问。
