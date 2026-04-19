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

**Idioma: [简体中文](README_zh.md) | [English](README.md) | [日本語](README_ja.md) | [한국어](README_ko.md) | [Français](README_fr.md) | Español**

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

En la primera ejecución de cualquier comando synapse, se creará automáticamente una plantilla de configuración global en `~/.synapse/config.yaml`:

```bash
# Activar la creación automática de la plantilla de configuración
synapse --version

# Salida:
# 📝 Created global config template: /Users/you/.synapse/config.yaml
#    Please edit this file to configure your store and extensions.
#    Then run 'synapse check' to verify your configuration.
```

### Paso 1: Configurar Extensiones

Edita el archivo de configuración global `~/.synapse/config.yaml` para seleccionar tu backend de almacenamiento y otras extensiones.

#### Opción A: Almacenamiento en Sistema de Archivos Local (Recomendado para principiantes)

```yaml
synapse:
  version: "1.0"

  sources:
    - name: "skill-source"
      enabled: true

  processor:
    name: "skill-processor"

  # Almacenamiento local
  store:
    name: "local-store"
    config:
      path: "~/knowhub"        # Ruta local de la base de conocimiento
```

#### Opción B: Almacenamiento en Repositorio GitHub (Para sincronización en la nube)

```yaml
synapse:
  version: "1.0"

  sources:
    - name: "skill-source"
      enabled: true

  processor:
    name: "skill-processor"

  # Almacenamiento GitHub
  store:
    name: "github-store"
    config:
      owner: "${GITHUB_OWNER}"   # Tu nombre de usuario de GitHub
      repo: "${GITHUB_REPO}"     # Nombre del repositorio de la base de conocimiento
      token: "${GITHUB_TOKEN}"   # GitHub Personal Access Token
      branch: "main"
```

> 💡 **Consejo**: Usa el formato `${ENV_VAR}` para referenciar variables de entorno, evitando codificar información sensible en los archivos de configuración.

### Paso 2: Verificar la Configuración

```bash
synapse check
```

Ejemplo de salida:

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

El comando `check` valida lo siguiente:

| Verificación | Descripción |
|-------------|-------------|
| Existencia del archivo de config | Si `~/.synapse/config.yaml` existe |
| Validez YAML | Si el archivo es un YAML válido |
| Campos requeridos | Si `synapse.version` y `synapse.store.name` están definidos |
| Registro de extensiones | Si los Store/Source/Processor configurados están registrados en el Registry |
| Variables de entorno | Si los marcadores `${ENV_VAR}` tienen variables de entorno correspondientes |

### Paso 3: Inicializar la Base de Conocimiento

```bash
# Inicializar usando la configuración global
synapse init

# Especificar el nombre del propietario
synapse init --name "Tu Nombre"

# Usar un archivo de configuración específico
synapse init --config /path/to/config.yaml

# Forzar la re-inicialización (los datos existentes NO se eliminarán)
synapse init --force
```

El comando `init` realiza automáticamente la inicialización basándose en el backend Store especificado en tu configuración:

| Store | Comportamiento de inicialización |
|-------|--------------------------------|
| `local-store` | Crea la estructura de directorios y archivos plantilla localmente |
| `github-store` | Crea los archivos esqueleto en el repositorio vía API de GitHub |

Estructura de directorios después de la inicialización:

```
knowhub/
├── .synapse/
│   └── schema.yaml       # Schema de conocimiento (contrato de comportamiento)
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

> ⚠️ **Idempotencia**: Si la base de conocimiento ya ha sido inicializada, `init` mostrará una advertencia y se saltará. Usa `--force` para forzar la re-inicialización.

### Paso 4: Instalar Skill en Asistentes IA

Un Skill es un archivo de instrucciones Prompt pre-configurado. Una vez instalado, tu asistente IA te ayudará automáticamente a recopilar, organizar y reutilizar conocimiento durante las conversaciones.

```bash
# Instalar en CodeBuddy (recomendado)
synapse install codebuddy

# Instalar en Claude Code
synapse install claude --target /path/to/project

# Instalar en Cursor
synapse install cursor

# Listar todos los asistentes IA soportados
synapse install --list
```

Después de la instalación, puedes usar las siguientes frases activadoras:

| Tú dices | La IA hace |
|----------|-----------|
| "recuerda esto" / "guarda en la base" | Recopila inmediatamente conocimiento de la conversación |
| "revisa la base de conocimiento" / "auditoría" | Ejecuta una verificación de salud |
| "¿qué sé sobre X?" | Busca contenido relevante |
| "organiza el inbox" | Ayuda a organizar elementos pendientes |

### Paso 5: Uso Diario

#### Recopilación Manual de Conocimiento

Además de la recopilación automática vía Skill, también puedes usar el CLI manualmente:

```bash
# Pasar contenido directamente
synapse collect --content "Las interfaces de Go se implementan implícitamente" --title "Go Interfaces" \
  --topics "Go" --concepts "Duck Typing"

# Entrada por pipe
echo "Notas de aprendizaje..." | synapse collect --topics "Sistemas Distribuidos" --entities "Raft"
```

#### Buscar en la Base de Conocimiento

```bash
# Búsqueda por palabra clave
synapse search goroutine

# Filtrar por tipo
synapse search --type topic "modelo de concurrencia"

# Limitar resultados
synapse search --limit 5 golang
```

#### Auditar la Base de Conocimiento

```bash
synapse audit
```

El informe de auditoría incluye:

| Verificación | Descripción |
|-------------|-------------|
| Puntuación de salud | Puntuación global (sobre 100) |
| Completitud Frontmatter | Si los campos requeridos como título y tipo están presentes |
| Enlaces rotos | Si los `[[wiki-links]]` apuntan a páginas existentes |
| Páginas huérfanas | Páginas no enlazadas desde otras páginas |
| Estadísticas | Cantidad de archivos, enlaces, distribución por tipo |

#### Gestionar Plugins de Extensión

```bash
# Listar todas las extensiones registradas
synapse plugin list
```

---

## 📖 Referencia de Comandos

| Comando | Descripción | Ejemplo |
|---------|-------------|---------|
| `synapse init` | Inicializar base de conocimiento | `synapse init --name "Juan"` |
| `synapse check` | Verificar validez de la config | `synapse check` |
| `synapse collect` | Recopilar conocimiento | `synapse collect --content "..." --topics "Go"` |
| `synapse search` | Buscar en la base | `synapse search goroutine` |
| `synapse audit` | Auditar salud de la base | `synapse audit` |
| `synapse install` | Instalar Skill en asistente IA | `synapse install codebuddy` |
| `synapse plugin list` | Listar plugins registrados | `synapse plugin list` |

### Opciones Globales

| Opción | Descripción |
|--------|-------------|
| `--config`, `-c` | Ruta del archivo de config (por defecto `~/.synapse/config.yaml`) |
| `--version`, `-v` | Mostrar número de versión |
| `--help`, `-h` | Mostrar información de ayuda |

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
