<p align="center">
  <img src="assets/logo.png" alt="Synapse Logo" width="200" />
</p>

<h1 align="center">Synapse</h1>

<p align="center">
  <strong>Hub de Connaissances Personnel (Personal Knowledge Hub)</strong><br/>
  Distillez, organisez et réinvestissez automatiquement les connaissances de vos conversations IA, transformant chaque échange en capital de connaissances cumulé.
</p>

<p align="center">
  <img src="assets/hero-banner.png" alt="Synapse Flux de Travail" width="800" />
</p>

[![CI](https://github.com/tunsuy/synapse/actions/workflows/ci.yml/badge.svg)](https://github.com/tunsuy/synapse/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tunsuy/synapse)](https://goreportcard.com/report/github.com/tunsuy/synapse)
[![codecov](https://codecov.io/gh/tunsuy/synapse/branch/main/graph/badge.svg)](https://codecov.io/gh/tunsuy/synapse)
[![Go Reference](https://pkg.go.dev/badge/github.com/tunsuy/synapse.svg)](https://pkg.go.dev/github.com/tunsuy/synapse)
[![Release](https://img.shields.io/github/v/release/tunsuy/synapse?include_prereleases)](https://github.com/tunsuy/synapse/releases)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.24-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**Langue : [简体中文](README_zh.md) | [English](README.md) | [日本語](README_ja.md) | [한국어](README_ko.md) | Français | [Español](README_es.md)**

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

## 📸 Démo

### Initialiser la Base de Connaissances

<p align="center">
  <img src="assets/demo-init.png" alt="synapse init" width="700" />
</p>

### Collecter des Connaissances

<p align="center">
  <img src="assets/demo-collect.png" alt="synapse collect" width="700" />
</p>

### Rechercher dans la Base de Connaissances

<p align="center">
  <img src="assets/demo-search.png" alt="synapse search" width="700" />
</p>

### Skill dans l'Assistant IA

Utilisation du skill synapse-knowledge dans CodeBuddy IDE pour une gestion intelligente des connaissances :

<p align="center">
  <img src="assets/demo-skill-collect.png" alt="skill collecter connaissances" width="700" />
</p>

<p align="center">
  <img src="assets/demo-skill-search.png" alt="skill rechercher connaissances" width="700" />
</p>

---

## 🚀 Démarrage Rapide

### Prérequis

- Go >= 1.24

### Installation

```bash
go install github.com/tunsuy/synapse@latest
```

Lors de la première exécution d'une commande synapse, un modèle de configuration globale sera automatiquement créé dans `~/.synapse/config.yaml` :

```bash
# Déclencher la création automatique du modèle de configuration
synapse --version

# Sortie :
# 📝 Created global config template: /Users/you/.synapse/config.yaml
#    Please edit this file to configure your store and extensions.
#    Then run 'synapse check' to verify your configuration.
```

### Étape 1 : Configurer les Extensions

Éditez le fichier de configuration globale `~/.synapse/config.yaml` pour sélectionner votre backend de stockage et autres extensions.

#### Option A : Stockage Fichier Local (Recommandé pour les débutants)

```yaml
synapse:
  version: "1.0"

  sources:
    - name: "skill-source"
      enabled: true

  processor:
    name: "skill-processor"

  # Stockage local
  store:
    name: "local-store"
    config:
      path: "~/knowhub"        # Chemin local de la base de connaissances
```

#### Option B : Stockage Dépôt GitHub (Pour la synchronisation cloud)

```yaml
synapse:
  version: "1.0"

  sources:
    - name: "skill-source"
      enabled: true

  processor:
    name: "skill-processor"

  # Stockage GitHub
  store:
    name: "github-store"
    config:
      owner: "${GITHUB_OWNER}"   # Votre nom d'utilisateur GitHub
      repo: "${GITHUB_REPO}"     # Nom du dépôt de base de connaissances
      token: "${GITHUB_TOKEN}"   # GitHub Personal Access Token
      branch: "main"
```

> 💡 **Conseil** : Utilisez le format `${ENV_VAR}` pour référencer les variables d'environnement, évitant ainsi de coder en dur des informations sensibles dans les fichiers de configuration.

### Étape 2 : Vérifier la Configuration

```bash
synapse check
```

Exemple de sortie :

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

La commande `check` valide les éléments suivants :

| Vérification | Description |
|-------------|-------------|
| Existence du fichier de config | Si `~/.synapse/config.yaml` existe |
| Validité YAML | Si le fichier est un YAML valide |
| Champs requis | Si `synapse.version` et `synapse.store.name` sont définis |
| Enregistrement des extensions | Si les Store/Source/Processor configurés sont enregistrés dans le Registry |
| Variables d'environnement | Si les espaces réservés `${ENV_VAR}` ont des variables d'environnement correspondantes |

### Étape 3 : Initialiser la Base de Connaissances

```bash
# Initialiser avec la configuration globale
synapse init

# Spécifier le nom du propriétaire
synapse init --name "Votre Nom"

# Utiliser un fichier de configuration spécifique
synapse init --config /path/to/config.yaml

# Forcer la ré-initialisation (les données existantes ne seront PAS supprimées)
synapse init --force
```

La commande `init` effectue automatiquement l'initialisation en fonction du backend Store spécifié dans votre configuration :

| Store | Comportement d'initialisation |
|-------|------------------------------|
| `local-store` | Crée la structure de répertoires et les fichiers modèles localement |
| `github-store` | Crée les fichiers squelettes dans le dépôt via l'API GitHub |

Structure des répertoires après initialisation :

```
knowhub/
├── .synapse/
│   └── schema.yaml       # Schéma de connaissances (contrat de comportement)
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

> ⚠️ **Idempotence** : Si la base de connaissances a déjà été initialisée, `init` affichera un avertissement et passera. Utilisez `--force` pour forcer la ré-initialisation.

### Étape 4 : Installer le Skill dans les Assistants IA

Un Skill est un fichier d'instructions Prompt pré-configuré. Une fois installé, votre assistant IA vous aidera automatiquement à collecter, organiser et réutiliser les connaissances pendant les conversations.

```bash
# Installer dans CodeBuddy (recommandé)
synapse install codebuddy

# Installer dans Claude Code
synapse install claude --target /path/to/project

# Installer dans Cursor
synapse install cursor

# Lister tous les assistants IA supportés
synapse install --list
```

Après installation, vous pouvez utiliser les phrases déclencheurs suivantes :

| Vous dites | L'IA fait |
|-----------|----------|
| « retiens ça » / « sauvegarde dans la base » | Collecte immédiatement les connaissances de la conversation |
| « vérifie la base de connaissances » / « audit » | Exécute un bilan de santé |
| « qu'est-ce que je sais sur X » | Recherche le contenu pertinent |
| « organise l'inbox » | Aide à organiser les éléments en attente |

### Étape 5 : Utilisation Quotidienne

#### Collecte Manuelle de Connaissances

En plus de la collecte automatique via Skill, vous pouvez aussi utiliser le CLI manuellement :

```bash
# Passer le contenu directement
synapse collect --content "Les interfaces Go sont implémentées implicitement" --title "Go Interfaces" \
  --topics "Go" --concepts "Duck Typing"

# Entrée par pipe
echo "Notes d'apprentissage..." | synapse collect --topics "Systèmes Distribués" --entities "Raft"
```

#### Rechercher dans la Base de Connaissances

```bash
# Recherche par mot-clé
synapse search goroutine

# Filtrer par type
synapse search --type topic "modèle de concurrence"

# Limiter les résultats
synapse search --limit 5 golang
```

#### Auditer la Base de Connaissances

```bash
synapse audit
```

Le rapport d'audit comprend :

| Vérification | Description |
|-------------|-------------|
| Score de santé | Score global (sur 100) |
| Complétude Frontmatter | Si les champs requis comme titre et type sont présents |
| Liens cassés | Si les `[[wiki-links]]` pointent vers des pages existantes |
| Pages orphelines | Pages non liées depuis d'autres pages |
| Statistiques | Nombre de fichiers, de liens, distribution par type |

#### Gérer les Plugins d'Extension

```bash
# Lister toutes les extensions enregistrées
synapse plugin list
```

---

## 📖 Référence des Commandes

| Commande | Description | Exemple |
|----------|-------------|---------|
| `synapse init` | Initialiser la base de connaissances | `synapse init --name "Jean"` |
| `synapse check` | Vérifier la validité de la config | `synapse check` |
| `synapse collect` | Collecter des connaissances | `synapse collect --content "..." --topics "Go"` |
| `synapse search` | Rechercher dans la base | `synapse search goroutine` |
| `synapse audit` | Auditer la santé de la base | `synapse audit` |
| `synapse install` | Installer le Skill dans un assistant IA | `synapse install codebuddy` |
| `synapse plugin list` | Lister les plugins enregistrés | `synapse plugin list` |

### Options Globales

| Option | Description |
|--------|-------------|
| `--config`, `-c` | Chemin du fichier de config (défaut `~/.synapse/config.yaml`) |
| `--version`, `-v` | Afficher le numéro de version |
| `--help`, `-h` | Afficher l'aide |

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
