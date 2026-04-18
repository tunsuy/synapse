<p align="center">
  <img src="assets/logo.png" alt="Synapse Logo" width="200" />
</p>

<h1 align="center">Synapse</h1>

<p align="center">
  <strong>Personal Knowledge Hub</strong><br/>
  Automatically distill, organize, and reinvest knowledge from your AI conversations, turning every chat into compounding knowledge.
</p>

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**Language: [简体中文](README.md) | English | [日本語](README_ja.md) | [한국어](README_ko.md) | [Français](README_fr.md) | [Español](README_es.md)**

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

### Initialize Knowledge Base

```bash
# Initialize a new knowledge base
synapse init ~/knowhub

# View knowledge base structure
tree ~/knowhub
```

### Knowledge Base Directory Structure

```
knowhub/
├── .synapse/
│   ├── schema.yaml       # Knowledge schema (behavior contract)
│   └── config.yaml       # Extension point configuration
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

```bash
# List installed plugins
synapse plugin list

# Install from Go module
synapse plugin install github.com/example/synapse-rss-source

# Install from Git repository
synapse plugin install --git https://github.com/example/synapse-vector-indexer.git

# Install from local directory
synapse plugin install --local ./my-custom-processor

# Enable / Disable plugins
synapse plugin enable rss-source
synapse plugin disable rss-source

# Check plugin health
synapse plugin doctor
```

---

## 📅 Roadmap

| Milestone | Content | Status |
|-----------|---------|--------|
| **M1 Foundation** | Schema spec + Extension point interfaces + CLI init | 🟡 Pending |
| **M2 Skill Integration** | First Source + Processor + Store, end-to-end pipeline | 🟡 Pending |
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
