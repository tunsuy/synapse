# Issue #001: Skill 无法自动触发知识收集（被动采集机制缺失）

> **状态**：🟡 待处理
>
> **优先级**：P1
>
> **关联里程碑**：M2（Skill 集成）/ M3（MCP Server 增强）
>
> **发现日期**：2026-04-20

---

## 问题描述

当前 synapse-knowledge Skill 的知识收集（Collect）依赖 **用户显式触发**（如说 "remember this"、"save to knowhub"），
无法在对话中自动识别有价值的内容并主动收集。这与 Skill Prompt 中描述的 "continuously identifies valuable knowledge" 设计意图不符。

### 根因分析

这是一个 **"鸡生蛋"问题**：

1. **Skill 加载依赖显式触发**：CodeBuddy 的 Skill 框架根据 `description` 中的 Trigger 词来决定是否加载 Skill。
   当前 `metadata.yaml` 中的 Triggers 全部是显式关键词：
   ```
   Triggers: 'remember this', 'save to knowhub', 'check knowledge base',
             'what do I know about X', 'review X', 'what did I learn today'
   ```

2. **缺少隐式/被动触发条件**：description 中虽然写了 "automatically" 和 "continuously"，但这只是描述性文字，
   不构成实际的触发机制。如果用户没有说出 Trigger 词，Skill 就不会被加载。

3. **Skill 未加载 = 自动收集不存在**：即使 Skill Prompt 里写了 "Continuously identify valuable knowledge during conversation"，
   但如果 Skill 本身没有被加载到上下文中，这段指令根本不会被 AI 看到。

4. **首次使用者无感知**：对于一个首次使用 synapse-knowledge Skill 的 AI 助手，它完全不知道需要自动收集知识，
   因为 Skill 只在匹配到 Trigger 词时才加载。

### 影响范围

- 直接影响 MVP 成功指标中的 **"知识采集成功率 > 80%"** — 当前实际采集率远低于此
- 直接影响 **"使用摩擦：用户不需要额外操作"** — 当前用户必须主动要求收集
- 违背产品愿景 **"让你的每一次 AI 对话都成为知识复利"** — 如果不主动触发就无法积累

---

## 复现步骤

1. 在 CodeBuddy 中安装 synapse-knowledge Skill
2. 正常进行技术讨论（如讨论项目架构、解决技术问题）
3. 不说任何 Trigger 词
4. 观察：Skill 不会被加载，不会有任何知识被自动收集

---

## 优化方向

### 方案 A：扩展 Skill Description 中的隐式触发场景（短期可行）

在 `metadata.yaml` 的 `description` 中补充更多隐式触发场景关键词，
让更多对话场景能触发 Skill 加载：

```yaml
description: >
  Personal knowledge steward. ...
  Triggers: 'remember this', 'save to knowhub', ...
  Also activates during: any technical discussion, problem-solving session,
  architecture review, new tool discovery, code review, learning new concept,
  best practice sharing, debugging session with reusable solution.
```

**优点**：改动最小，只需修改一个文件
**缺点**：触发范围仍然有限，无法覆盖所有有价值的对话场景

### 方案 B：探索 "always-on" 机制（中期理想）

调研 CodeBuddy Skill 框架是否支持类似 `alwaysApply: true` 的配置，
让 Skill 在每次对话开始时自动加载到上下文中。

```yaml
# 理想的配置
activation:
  mode: always-on    # 始终加载，不需要 Trigger 词
```

**优点**：彻底解决问题，AI 始终带着"收集知识"的意识
**缺点**：依赖平台框架支持，可能增加上下文占用

### 方案 C：MCP Server 驱动的自动采集（M3 阶段）

在 M3 实现 MCP Server 后，利用 MCP 的 Resource 机制将 "自动采集" 逻辑注入到 AI 上下文中：

- MCP Resource 可以在每次对话时自动提供知识库状态和采集规则
- AI 助手侧无需 Skill 加载，MCP Server 始终在后台运行

**优点**：不依赖 Skill 框架的触发机制
**缺点**：需要等 M3 完成

### 方案 D：Automation 定时回顾（补充方案）

设置一个 CodeBuddy Automation，定时回顾对话历史并提取有价值的知识：

```
每天结束时自动回顾当天的对话，提取有价值的知识点并通过 synapse collect 收集
```

**优点**：不遗漏任何对话中的知识
**缺点**：非实时，有延迟；依赖对话历史的可访问性

---

## 建议行动

1. **短期（本周）**：实施方案 A，扩展 description 中的触发场景
2. **中期（M3）**：方案 B + C 并行调研，在 MCP Server 实现时一并解决
3. **补充**：方案 D 作为兜底方案，确保不遗漏

---

## 相关文件

- `skills/common/metadata.yaml` — Skill 元数据（包含 description 和 triggers）
- `skills/codebuddy/synapse-knowledge.md` — CodeBuddy Skill Prompt 模板
- `docs/roadmap.md` — 项目路线图（M2/M3 相关）
