# Synapse

> **Hub de Connaissances Personnel (Personal Knowledge Hub)** — Distillez, organisez et réinvestissez automatiquement les connaissances de vos conversations IA, transformant chaque échange en capital de connaissances cumulé.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**Langue : [简体中文](README.md) | [English](README_en.md) | [日本語](README_ja.md) | [한국어](README_ko.md) | Français | [Español](README_es.md)**

---

## 🎯 Pourquoi Synapse ?

Nous utilisons quotidiennement divers assistants IA (ChatGPT, Claude, CodeBuddy, Gemini, etc.) dans notre travail et nos études. Chaque conversation est essentiellement une accumulation de connaissances. Mais la réalité est :

- **Connaissances fragmentées** — Éparpillées entre différents assistants IA, difficiles à retrouver
- **Isolement cognitif de l'IA** — La compréhension que l'IA a de vous est fragmentée ; chaque conversation repart de zéro
- **Actifs obscurs** — Des productions conversationnelles précieuses utilisées une fois puis oubliées

**L'objectif de Synapse** : Transformer chaque conversation IA en un actif de connaissances pouvant être conservé, recherché et réinvesti.

> « Les wikis sont des produits de connaissances persistants à croissance composée. » — Andrej Karpathy

---

## ✨ Fonctionnalités Clés

- 🔌 **Modèle à Points d'Extension** — Six points d'extension indépendants (Source / Processor / Store / Indexer / Consumer / Auditor), composables et remplaçables indépendamment
- 📥 **Ingestion Multi-Sources** — Acquisition de contenu sans friction depuis tout assistant IA, RSS, Notion, podcasts, etc.
- 🧠 **Traitement Intelligent** — Extraction, classification et corrélation de connaissances pilotées par l'IA ; compilation automatique des conversations brutes en connaissances structurées
- 💾 **Souveraineté du Stockage** — Les données résident dans le backend de votre choix (Local / GitHub / S3 / WebDAV), entièrement auto-contrôlé
- 🔍 **Recherche Flexible** — Moteurs de recherche enfichables (BM25 / Recherche vectorielle / Parcours de graphe)
- 📊 **Consommation Multi-Format** — Sortie des connaissances en sites statiques, Obsidian Vaults, flashcards Anki, digests email, etc.
- 🔗 **Liens Bidirectionnels** — Format `[[wiki-link]]`, compatible Obsidian, construction de votre graphe de connaissances personnel
- 📋 **Piloté par Schéma** — Définition des contrats de comportement IA via des fichiers Schema ; modifier le Schema modifie le comportement de tous les assistants IA
- 🧩 **Écosystème de Plugins** — CLI complet de gestion des plugins, installation multi-sources, implémentations de points d'extension contribuées par la communauté

---

## 🏗️ Aperçu de l'Architecture

Synapse adopte un **Modèle à Points d'Extension (Extension Point Model)** — une architecture en étoile avec Store comme fondation et six points d'extension indépendants composés à la demande :

```
                    ┌─────────────┐
                    │   Source     │  Sources de données (Chat IA / RSS / Notion / ...)
                    └──────┬──────┘
                           │ RawContent
                           ▼
                    ┌─────────────┐
                    │  Processor  │  Moteurs de traitement (Skill / MCP / LocalLLM / ...)
                    └──────┬──────┘
                           │ KnowledgeFile
                           ▼
┌──────────────────────────────────────────────────────┐
│                  Store (Couche de stockage)           │
│        Local FS / GitHub / S3 / WebDAV / ...         │
└────────┬──────────────────┬──────────────────┬───────┘
         │                  │                  │
         ▼                  ▼                  ▼
  ┌─────────────┐   ┌─────────────┐   ┌───────────────┐
  │   Indexer    │   │   Auditor   │   │   Consumer    │
  │  Recherche   │   │   Qualité   │   │   Sortie      │
  └─────────────┘   └─────────────┘   └───────────────┘
```

> Pour la documentation détaillée de l'architecture, voir [ARCHITECTURE.md](ARCHITECTURE.md).

---

## 🚀 Démarrage Rapide

### Prérequis

- Go >= 1.21

### Installation

```bash
go install github.com/tunsuy/synapse@latest
```

### Initialiser la Base de Connaissances

```bash
# Initialiser une nouvelle base de connaissances
synapse init ~/knowhub

# Voir la structure de la base de connaissances
tree ~/knowhub
```

### Structure des Répertoires

```
knowhub/
├── .synapse/
│   ├── schema.yaml       # Schéma de connaissances (contrat de comportement)
│   └── config.yaml       # Configuration des points d'extension
├── profile/
│   └── me.md             # Profil utilisateur
├── topics/               # Connaissances par thème
│   ├── golang/
│   ├── architecture/
│   └── ...
├── entities/             # Pages d'entités (personnes, outils, projets)
├── concepts/             # Pages de concepts (concepts tech, méthodologies)
├── inbox/                # Éléments en attente
├── journal/              # Journal chronologique
└── graph/
    └── relations.json    # Graphe de relations de connaissances
```

---

## 🔌 Points d'Extension

| Point d'Extension | Responsabilité | Implémentation par Défaut | Contributions Communautaires |
|-------------------|---------------|--------------------------|------------------------------|
| **Source** | Récupérer le contenu brut depuis des sources externes | CodeBuddy Skill | RSS / Notion / Twitter / Podcast / WeChat... |
| **Processor** | Contenu brut → Connaissances structurées | Skill Processor | LLM local / Moteur de règles / Hybride... |
| **Store** | CRUD + contrôle de version pour les fichiers de connaissances | Local Store | GitHub / S3 / WebDAV / SQLite / IPFS... |
| **Indexer** | Recherche dans la base de connaissances | BM25 Indexer | Recherche vectorielle / Parcours de graphe / Elasticsearch... |
| **Consumer** | Sortie des connaissances dans divers formats | Site Hugo | VitePress / Anki / Email / TUI... |
| **Auditor** | Vérification de qualité et réparations | Default Auditor | Règles d'audit personnalisées... |

---

## 🧩 Gestion des Plugins

```bash
# Lister les plugins installés
synapse plugin list

# Installer depuis un module Go
synapse plugin install github.com/example/synapse-rss-source

# Installer depuis un dépôt Git
synapse plugin install --git https://github.com/example/synapse-vector-indexer.git

# Installer depuis un répertoire local
synapse plugin install --local ./my-custom-processor

# Activer / Désactiver des plugins
synapse plugin enable rss-source
synapse plugin disable rss-source

# Vérifier la santé des plugins
synapse plugin doctor
```

---

## 📅 Feuille de Route

| Jalon | Contenu | Statut |
|-------|---------|--------|
| **M1 Fondation** | Spéc. Schema + Interfaces des points d'extension + CLI init | 🟡 En attente |
| **M2 Intégration Skill** | Premier Source + Processor + Store, pipeline E2E | 🟡 En attente |
| **M3 MCP + Gestion Plugins** | MCP Server + GitHub Store + BM25 Indexer + Plugin CLI | 🔵 Planifié |
| **M4 Multi-Plateforme** | Claude Code / Cursor / ChatGPT Source | 🔵 Planifié |
| **M5 Impl. Consumer** | Site Hugo + Compat. Obsidian + Graphe de connaissances | 🔵 Planifié |
| **M6+ Communauté** | Marché de plugins + Points d'extension complets + Communauté | 🔵 Long terme |

> Pour la feuille de route détaillée, voir [docs/roadmap.md](docs/roadmap.md).

---

## 🤝 Contribuer

Nous accueillons toutes les formes de contributions ! Que ce soit des rapports de bugs, des suggestions de fonctionnalités ou des contributions directes au code.

- 📖 Lisez le [Guide de Contribution](CONTRIBUTING.md) pour savoir comment participer
- 🏛️ Lisez le [Guide d'Architecture](ARCHITECTURE.md) pour comprendre la conception technique
- 📋 Lisez le [Code de Conduite](CODE_OF_CONDUCT.md) pour les normes communautaires
- 🗺️ Lisez la [Feuille de Route](docs/roadmap.md) pour la planification du projet

### Domaines de Contribution

Chaque point d'extension accueille les implémentations de la communauté :

- 🔌 **Plugins Source** : Connecter plus de sources de données (RSS, Notion, WeChat, Podcasts...)
- ⚙️ **Plugins Processor** : Supporter plus de moteurs de traitement (LLM local, Moteur de règles...)
- 💾 **Plugins Store** : Supporter plus de backends de stockage (S3, WebDAV, IPFS...)
- 🔍 **Plugins Indexer** : Supporter plus de moteurs de recherche (Recherche vectorielle, Parcours de graphe...)
- 📊 **Plugins Consumer** : Supporter plus de formats de sortie (VitePress, Anki, TUI...)

---

## 📄 Licence

Ce projet est sous licence [Apache License 2.0](LICENSE).

---

## 💬 Contact

- **Issues** : [GitHub Issues](https://github.com/tunsuy/synapse/issues)
- **Discussions** : [GitHub Discussions](https://github.com/tunsuy/synapse/discussions)

---

> *Synapse — Transformez chaque conversation IA en capital de connaissances cumulé.*
