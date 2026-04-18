<p align="center">
  <img src="assets/logo.png" alt="Synapse Logo" width="200" />
</p>

<h1 align="center">Synapse</h1>

<p align="center">
  <strong>Hub de Conocimiento Personal (Personal Knowledge Hub)</strong><br/>
  Destila, organiza y reinvierte automáticamente el conocimiento de tus conversaciones con IA, convirtiendo cada interacción en capital de conocimiento compuesto.
</p>

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**Idioma: [简体中文](README.md) | [English](README_en.md) | [日本語](README_ja.md) | [한국어](README_ko.md) | [Français](README_fr.md) | Español**

---

## 🎯 ¿Por qué Synapse?

Usamos diversos asistentes de IA (ChatGPT, Claude, CodeBuddy, Gemini, etc.) en nuestro trabajo y aprendizaje diario. Cada conversación es esencialmente una acumulación de conocimiento. Pero la realidad es:

- **Conocimiento fragmentado** — Disperso entre diferentes asistentes de IA, difícil de revisar
- **Aislamiento cognitivo de la IA** — La comprensión que la IA tiene de ti es fragmentada; cada conversación comienza desde cero
- **Activos oscuros** — Productos valiosos de conversación que se usan una vez y se olvidan

**El objetivo de Synapse**: Convertir cada conversación con IA en activos de conocimiento que pueden ser retenidos, buscados y reinvertidos.

> "Los wikis son productos de conocimiento persistentes con crecimiento compuesto." — Andrej Karpathy

---

## ✨ Características Principales

- 🔌 **Modelo de Puntos de Extensión** — Seis puntos de extensión independientes (Source / Processor / Store / Indexer / Consumer / Auditor), componibles e independientemente reemplazables
- 📥 **Ingesta Multi-Fuente** — Adquisición de contenido sin fricción desde cualquier asistente IA, RSS, Notion, podcasts, etc.
- 🧠 **Procesamiento Inteligente** — Extracción, clasificación y correlación de conocimiento impulsadas por IA; compilación automática de conversaciones brutas en conocimiento estructurado
- 💾 **Soberanía de Almacenamiento** — Los datos residen en cualquier backend que elijas (Local / GitHub / S3 / WebDAV), totalmente auto-controlado
- 🔍 **Búsqueda Flexible** — Motores de búsqueda conectables (BM25 / Búsqueda vectorial / Recorrido de grafos)
- 📊 **Consumo Multi-Formato** — Salida de conocimiento como sitios estáticos, Obsidian Vaults, tarjetas Anki, resúmenes por email, etc.
- 🔗 **Enlaces Bidireccionales** — Formato `[[wiki-link]]`, compatible con Obsidian, construcción de tu grafo de conocimiento personal
- 📋 **Dirigido por Schema** — Define contratos de comportamiento de IA mediante archivos Schema; modificar el Schema cambia el comportamiento de todos los asistentes IA
- 🧩 **Ecosistema de Plugins** — CLI completo de gestión de plugins, instalación multi-fuente, implementaciones de puntos de extensión contribuidas por la comunidad

---

## 🏗️ Visión General de la Arquitectura

Synapse adopta un **Modelo de Puntos de Extensión (Extension Point Model)** — una arquitectura en estrella con Store como base y seis puntos de extensión independientes compuestos según demanda:

```
                    ┌─────────────┐
                    │   Source     │  Fuentes de datos (Chat IA / RSS / Notion / ...)
                    └──────┬──────┘
                           │ RawContent
                           ▼
                    ┌─────────────┐
                    │  Processor  │  Motores de procesamiento (Skill / MCP / LocalLLM / ...)
                    └──────┬──────┘
                           │ KnowledgeFile
                           ▼
┌──────────────────────────────────────────────────────┐
│                  Store (Capa de almacenamiento)       │
│        Local FS / GitHub / S3 / WebDAV / ...         │
└────────┬──────────────────┬──────────────────┬───────┘
         │                  │                  │
         ▼                  ▼                  ▼
  ┌─────────────┐   ┌─────────────┐   ┌───────────────┐
  │   Indexer    │   │   Auditor   │   │   Consumer    │
  │  Búsqueda   │   │   Calidad   │   │   Salida      │
  └─────────────┘   └─────────────┘   └───────────────┘
```

> Para documentación detallada de la arquitectura, consulta [ARCHITECTURE.md](ARCHITECTURE.md).

---

## 🚀 Inicio Rápido

### Requisitos

- Go >= 1.21

### Instalación

```bash
go install github.com/tunsuy/synapse@latest
```

### Inicializar la Base de Conocimiento

```bash
# Inicializar una nueva base de conocimiento
synapse init ~/knowhub

# Ver la estructura de la base de conocimiento
tree ~/knowhub
```

### Estructura de Directorios

```
knowhub/
├── .synapse/
│   ├── schema.yaml       # Schema de conocimiento (contrato de comportamiento)
│   └── config.yaml       # Configuración de puntos de extensión
├── profile/
│   └── me.md             # Perfil de usuario
├── topics/               # Conocimiento por temas
│   ├── golang/
│   ├── architecture/
│   └── ...
├── entities/             # Páginas de entidades (personas, herramientas, proyectos)
├── concepts/             # Páginas de conceptos (conceptos técnicos, metodologías)
├── inbox/                # Elementos pendientes
├── journal/              # Diario cronológico
└── graph/
    └── relations.json    # Grafo de relaciones de conocimiento
```

---

## 🔌 Puntos de Extensión

| Punto de Extensión | Responsabilidad | Implementación por Defecto | Contribuciones de la Comunidad |
|-------------------|----------------|---------------------------|-------------------------------|
| **Source** | Obtener contenido bruto de fuentes externas | CodeBuddy Skill | RSS / Notion / Twitter / Podcast / WeChat... |
| **Processor** | Contenido bruto → Conocimiento estructurado | Skill Processor | LLM local / Motor de reglas / Híbrido... |
| **Store** | CRUD + control de versiones para archivos de conocimiento | Local Store | GitHub / S3 / WebDAV / SQLite / IPFS... |
| **Indexer** | Búsqueda en la base de conocimiento | BM25 Indexer | Búsqueda vectorial / Recorrido de grafos / Elasticsearch... |
| **Consumer** | Salida de conocimiento en diversos formatos | Sitio Hugo | VitePress / Anki / Email / TUI... |
| **Auditor** | Verificación de calidad y reparaciones | Default Auditor | Reglas de auditoría personalizadas... |

---

## 🧩 Gestión de Plugins

```bash
# Listar plugins instalados
synapse plugin list

# Instalar desde módulo Go
synapse plugin install github.com/example/synapse-rss-source

# Instalar desde repositorio Git
synapse plugin install --git https://github.com/example/synapse-vector-indexer.git

# Instalar desde directorio local
synapse plugin install --local ./my-custom-processor

# Habilitar / Deshabilitar plugins
synapse plugin enable rss-source
synapse plugin disable rss-source

# Verificar salud de plugins
synapse plugin doctor
```

---

## 📅 Hoja de Ruta

| Hito | Contenido | Estado |
|------|-----------|--------|
| **M1 Fundación** | Espec. Schema + Interfaces de puntos de extensión + CLI init | 🟡 Pendiente |
| **M2 Integración Skill** | Primer Source + Processor + Store, pipeline E2E | 🟡 Pendiente |
| **M3 MCP + Gestión Plugins** | MCP Server + GitHub Store + BM25 Indexer + Plugin CLI | 🔵 Planificado |
| **M4 Multi-Plataforma** | Claude Code / Cursor / ChatGPT Source | 🔵 Planificado |
| **M5 Impl. Consumer** | Sitio Hugo + Compat. Obsidian + Grafo de conocimiento | 🔵 Planificado |
| **M6+ Comunidad** | Mercado de plugins + Puntos de extensión completos + Comunidad | 🔵 Largo plazo |

> Para la hoja de ruta detallada, consulta [docs/roadmap.md](docs/roadmap.md).

---

## 🤝 Contribuir

¡Damos la bienvenida a todas las formas de contribución! Ya sea enviando reportes de bugs, sugiriendo nuevas funcionalidades o contribuyendo código directamente.

- 📖 Lee la [Guía de Contribución](CONTRIBUTING.md) para saber cómo participar
- 🏛️ Lee la [Guía de Arquitectura](ARCHITECTURE.md) para entender el diseño técnico
- 📋 Lee el [Código de Conducta](CODE_OF_CONDUCT.md) para las normas de la comunidad
- 🗺️ Lee la [Hoja de Ruta](docs/roadmap.md) para la planificación del proyecto

### Áreas de Contribución

Cada punto de extensión da la bienvenida a implementaciones de la comunidad:

- 🔌 **Plugins Source**: Conectar más fuentes de datos (RSS, Notion, WeChat, Podcasts...)
- ⚙️ **Plugins Processor**: Soportar más motores de procesamiento (LLM local, Motor de reglas...)
- 💾 **Plugins Store**: Soportar más backends de almacenamiento (S3, WebDAV, IPFS...)
- 🔍 **Plugins Indexer**: Soportar más motores de búsqueda (Búsqueda vectorial, Recorrido de grafos...)
- 📊 **Plugins Consumer**: Soportar más formatos de salida (VitePress, Anki, TUI...)

---

## 📄 Licencia

Este proyecto está licenciado bajo la [Apache License 2.0](LICENSE).

---

## 💬 Contacto

- **Issues**: [GitHub Issues](https://github.com/tunsuy/synapse/issues)
- **Discussions**: [GitHub Discussions](https://github.com/tunsuy/synapse/discussions)

---

> *Synapse — Convierte cada conversación con IA en capital de conocimiento compuesto.*
