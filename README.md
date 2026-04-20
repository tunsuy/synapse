<p align="center">
  <img src="assets/logo.png" alt="Synapse Logo" width="200" />
</p>

<h1 align="center">Synapse</h1>

<p align="center">
  <strong>Personal Knowledge Hub</strong><br/>
  Automatically distill, organize, and reinvest knowledge from your AI conversations, turning every chat into compounding knowledge.
</p>

[![CI](https://github.com/tunsuy/synapse/actions/workflows/ci.yml/badge.svg)](https://github.com/tunsuy/synapse/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tunsuy/synapse)](https://goreportcard.com/report/github.com/tunsuy/synapse)
[![codecov](https://codecov.io/gh/tunsuy/synapse/branch/main/graph/badge.svg)](https://codecov.io/gh/tunsuy/synapse)
[![Go Reference](https://pkg.go.dev/badge/github.com/tunsuy/synapse.svg)](https://pkg.go.dev/github.com/tunsuy/synapse)
[![Release](https://img.shields.io/github/v/release/tunsuy/synapse?include_prereleases)](https://github.com/tunsuy/synapse/releases)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**Language: [简体中文](README_zh.md) | English | [日本語](README_ja.md) | [한국어](README_ko.md) | [Français](README_fr.md) | [Español](README_es.md)**

---

## 🎯 Why Synapse?

We use various AI assistants (ChatGPT, Claude, CodeBuddy, Gemini, etc.) in our daily work and learning. Every conversation is essentially an accumulation of knowledge. But the reality is:

- **Fragmented Knowledge** — Scattered across different AI assistants, hard to revisit
- **Cognitive Isolation** — AI's understanding of you is fragmented; every conversation starts from scratch
- **Dark Assets** — Valuable conversation outputs are used once and forgotten

**Synapse's Goal**: Turn every AI conversation into knowledge assets that can be retained, retrieved, and reinvested.

> "Wikis are persistent, compound-growth knowledge products." — Andrej Karpathy

---

## ✨ Core Features

- 🔌 **Extension Point Model** — Six independent extension points (Source / Processor / Store / Indexer / Consumer / Auditor), composable and independently replaceable
- 📥 **Multi-Source Ingestion** — Zero-friction content acquisition from any AI assistant, RSS, Notion, podcasts, etc.
- 🧠 **Intelligent Processing** — AI-driven knowledge extraction, classification, and correlation; automatically compiles raw conversations into structured knowledge
- 💾 **Storage Sovereignty** — Data resides in any backend you choose (Local / GitHub / S3 / WebDAV), fully self-controlled
- 🔍 **Flexible Retrieval** — Pluggable retrieval engines (BM25 / Vector Search / Graph Traversal)
- 📊 **Multi-Format Consumption** — Knowledge output as static sites, Obsidian Vaults, Anki flashcards, email digests, etc.
- 🔗 **Bidirectional Links** — `[[wiki-link]]` format, Obsidian-compatible, building your personal knowledge graph
- 📋 **Schema-Driven** — Define AI behavior contracts via Schema files; modify the Schema to change all AI assistants' behavior
- 🧩 **Plugin Ecosystem** — Complete plugin management CLI, multi-source installation, community-contributed extension point implementations

---

## 🏗️ Architecture Overview

Synapse adopts an **Extension Point Model** — a star-shaped architecture with Store as the foundation and six independent extension points composed on demand:

```
                    ┌─────────────┐
                    │   Source     │  Data Sources (AI Chat / RSS / Notion / ...)
                    └──────┬──────┘
                           │ RawContent
                           ▼
                    ┌─────────────┐
                    │  Processor  │  Processing Engines (Skill / MCP / LocalLLM / ...)
                    └──────┬──────┘
                           │ KnowledgeFile
                           ▼
┌──────────────────────────────────────────────────────┐
│                  Store (Storage Layer)                │
│        Local FS / GitHub / S3 / WebDAV / ...         │
└────────┬──────────────────┬──────────────────┬───────┘
         │                  │                  │
         ▼                  ▼                  ▼
  ┌─────────────┐   ┌─────────────┐   ┌───────────────┐
  │   Indexer    │   │   Auditor   │   │   Consumer    │
  │  Retrieval   │   │   Quality   │   │   Output      │
  └─────────────┘   └─────────────┘   └───────────────┘
```

> For detailed architecture documentation, see [ARCHITECTURE.md](ARCHITECTURE.md).

---

## 🚀 Quick Start

### Requirements

- Go >= 1.21

### Installation

```bash
go install github.com/tunsuy/synapse@latest
```

On the first run of any synapse command, a global configuration template will be automatically created at `~/.synapse/config.yaml`:

```bash
# Trigger auto-creation of config template
synapse --version

# Output:
# 📝 Created global config template: /Users/you/.synapse/config.yaml
#    Please edit this file to configure your store and extensions.
#    Then run 'synapse check' to verify your configuration.
```

### Step 1: Configure Extensions

Edit the global config file `~/.synapse/config.yaml` to select your storage backend and other extensions.

#### Option A: Local Filesystem Store (Recommended for Beginners)

```yaml
synapse:
  version: "1.0"

  sources:
    - name: "skill-source"
      enabled: true

  processor:
    name: "skill-processor"

  # Local storage
  store:
    name: "local-store"
    config:
      path: "~/knowhub"        # Local knowledge base path
```

#### Option B: GitHub Repository Store (For Cloud Sync)

```yaml
synapse:
  version: "1.0"

  sources:
    - name: "skill-source"
      enabled: true

  processor:
    name: "skill-processor"

  # GitHub storage
  store:
    name: "github-store"
    config:
      owner: "${GITHUB_OWNER}"   # Your GitHub username
      repo: "${GITHUB_REPO}"     # Knowledge base repository name
      token: "${GITHUB_TOKEN}"   # GitHub Personal Access Token
      branch: "main"
```

> 💡 **Tip**: Use `${ENV_VAR}` format to reference environment variables, avoiding hardcoded sensitive information in config files.

### Step 2: Verify Configuration

```bash
synapse check
```

Sample output:

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

The `check` command validates the following:

| Check | Description |
|-------|-------------|
| Config file existence | Whether `~/.synapse/config.yaml` exists |
| YAML validity | Whether the file is valid YAML |
| Required fields | Whether `synapse.version` and `synapse.store.name` are set |
| Extension registration | Whether configured Store/Source/Processor are registered in the Registry |
| Environment variables | Whether `${ENV_VAR}` placeholders have corresponding env vars set |

### Step 3: Initialize Knowledge Base

```bash
# Initialize using global config
synapse init

# Specify knowledge base owner name
synapse init --name "Your Name"

# Use a specific config file
synapse init --config /path/to/config.yaml

# Force re-initialization (existing data will NOT be deleted)
synapse init --force
```

The `init` command automatically performs initialization based on the Store backend specified in your config:

| Store | Initialization Behavior |
|-------|------------------------|
| `local-store` | Creates knowledge base directory structure and template files locally |
| `github-store` | Creates knowledge base skeleton files in the repository via GitHub API |

Knowledge base directory structure after initialization:

```
knowhub/
├── .synapse/
│   └── schema.yaml       # Knowledge schema (behavior contract)
├── profile/
│   └── me.md             # User profile
├── topics/               # Topic knowledge
│   ├── golang/
│   ├── architecture/
│   └── ...
├── entities/             # Entity pages (people, tools, projects)
├── concepts/             # Concept pages (tech concepts, methodologies)
├── inbox/                # Pending items
├── journal/              # Timeline journal
└── graph/
    └── relations.json    # Knowledge relation graph
```

> ⚠️ **Idempotency**: If the knowledge base has already been initialized, `init` will display a warning and skip. Use `--force` to force re-initialization.

### Step 4: Install Skill to AI Assistants

A Skill is a pre-configured Prompt instruction file. Once installed, your AI assistant will automatically help you collect, organize, and retrieve knowledge during conversations.

```bash
# Install to CodeBuddy (recommended)
synapse install codebuddy

# Install to Claude Code
synapse install claude --target /path/to/project

# Install to Cursor
synapse install cursor

# List all supported AI assistants
synapse install --list
```

After installation, you can use the following trigger phrases in your AI assistant:

| You say | AI will do |
|---------|------------|
| "remember this" / "save to knowhub" | Immediately collect knowledge from the current conversation |
| "check knowledge base" / "audit" | Run a knowledge base health check |
| "what do I know about X" | Retrieve relevant content from the knowledge base |
| "organize inbox" | Help you organize pending items |

### Step 5: Daily Usage

#### Manual Knowledge Collection

In addition to automatic collection via Skill, you can also use the CLI to collect manually:

```bash
# Pass content directly
synapse collect --content "Go interfaces are implicitly implemented" --title "Go Interfaces" \
  --topics "Go" --concepts "Duck Typing"

# Pipe input
echo "Learning notes..." | synapse collect --topics "Distributed Systems" --entities "Raft"
```

#### Search the Knowledge Base

```bash
# Keyword search
synapse search goroutine

# Filter by type
synapse search --type topic "concurrency model"

# Limit results
synapse search --limit 5 golang
```

#### Audit the Knowledge Base

```bash
synapse audit
```

The audit report includes:

| Check | Description |
|-------|-------------|
| Health Score | Overall score (out of 100) |
| Frontmatter Completeness | Whether required fields like title and type are present |
| Broken Links | Whether `[[wiki-links]]` point to existing pages |
| Orphan Pages | Pages not linked from any other page |
| Statistics | File counts, link counts, distribution by type |

#### Manage Extension Plugins

```bash
# List all registered extensions
synapse plugin list
```

---

## 📖 Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `synapse init` | Initialize knowledge base | `synapse init --name "Your Name"` |
| `synapse check` | Verify config validity | `synapse check` |
| `synapse collect` | Collect knowledge | `synapse collect --content "..." --topics "Go"` |
| `synapse search` | Search knowledge base | `synapse search goroutine` |
| `synapse audit` | Audit knowledge base health | `synapse audit` |
| `synapse install` | Install Skill to AI assistant | `synapse install codebuddy` |
| `synapse plugin list` | List registered plugins | `synapse plugin list` |

### Global Flags

| Flag | Description |
|------|-------------|
| `--config`, `-c` | Config file path (default `~/.synapse/config.yaml`) |
| `--version`, `-v` | Show version number |
| `--help`, `-h` | Show help information |

---

## 🔌 Extension Points

| Extension Point | Responsibility | Default Implementation | Community Contributions |
|----------------|---------------|----------------------|------------------------|
| **Source** | Fetch raw content from external sources | CodeBuddy Skill | RSS / Notion / Twitter / Podcast / WeChat... |
| **Processor** | Raw content → Structured knowledge | Skill Processor | Local LLM / Rule Engine / Hybrid... |
| **Store** | CRUD + version control for knowledge files | Local Store | GitHub / S3 / WebDAV / SQLite / IPFS... |
| **Indexer** | Knowledge base retrieval | BM25 Indexer | Vector Search / Graph Traversal / Elasticsearch... |
| **Consumer** | Output knowledge in various formats | Hugo Site | VitePress / Anki / Email / TUI... |
| **Auditor** | Quality checks and repairs | Default Auditor | Custom audit rules... |

---

## 🧩 Plugin Management

Synapse extends functionality through a plugin system. Full plugin management will be available in M3:

```bash
# List registered extension plugins (available now)
synapse plugin list

# The following commands will be implemented in M3:
# synapse plugin install github.com/example/synapse-rss-source  # Go module
# synapse plugin install --git https://github.com/example/xxx.git  # Git repo
# synapse plugin install --local ./my-custom-processor  # Local directory
# synapse plugin enable rss-source   # Enable plugin
# synapse plugin disable rss-source  # Disable plugin
# synapse plugin doctor              # Check plugin health
```

---

## 📅 Roadmap

| Milestone | Content | Status |
|-----------|---------|--------|
| **M1 Foundation** | Schema spec + Extension point interfaces + CLI init/check | ✅ Done |
| **M2 Skill Integration** | First Source + Processor + Store, end-to-end pipeline | ✅ Done |
| **M3 MCP + Plugin Mgmt** | MCP Server + GitHub Store + BM25 Indexer + Plugin CLI | 🔵 Planned |
| **M4 Multi-Platform** | Claude Code / Cursor / ChatGPT Source | 🔵 Planned |
| **M5 Consumer Impl** | Hugo site + Obsidian compat + Knowledge graph | 🔵 Planned |
| **M6+ Community** | Plugin marketplace + Full extension points + Community | 🔵 Long-term |

> For the detailed roadmap, see [docs/roadmap.md](docs/roadmap.md).

---

## 🤝 Contributing

We welcome all forms of contributions! Whether it's submitting bug reports, suggesting new features, or contributing code directly.

- 📖 Read the [Contributing Guide](CONTRIBUTING.md) to learn how to participate
- 🏛️ Read the [Architecture Guide](ARCHITECTURE.md) to understand the technical design
- 📋 Read the [Code of Conduct](CODE_OF_CONDUCT.md) for community standards
- 🗺️ Read the [Roadmap](docs/roadmap.md) for project planning
- 🔐 Read the [Security Policy](SECURITY.md) for vulnerability reporting

### Contribution Areas

Every extension point welcomes community-contributed implementations:

- 🔌 **Source Plugins**: Connect more data sources (RSS, Notion, WeChat, Podcasts...)
- ⚙️ **Processor Plugins**: Support more processing engines (Local LLM, Rule Engine...)
- 💾 **Store Plugins**: Support more storage backends (S3, WebDAV, IPFS...)
- 🔍 **Indexer Plugins**: Support more retrieval engines (Vector Search, Graph Traversal...)
- 📊 **Consumer Plugins**: Support more output formats (VitePress, Anki, TUI...)

---

## 📄 License

This project is licensed under the [Apache License 2.0](LICENSE).

---

## 💬 Contact

- **Issues**: [GitHub Issues](https://github.com/tunsuy/synapse/issues)
- **Discussions**: [GitHub Discussions](https://github.com/tunsuy/synapse/discussions)

---

> *Synapse — Turn every AI conversation into compounding knowledge.*
