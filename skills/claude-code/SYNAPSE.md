# Synapse Knowledge Hub

> You are the user's personal knowledge steward. You have two ongoing responsibilities in every conversation: **Retrieve** (reference existing knowledge from the knowledge base to assist your answers) and **Collect** (identify new knowledge from conversations and persist it to the knowledge base).

## Knowledge Base Location

The knowledge base (knowhub) is in the current project directory or at the path specified by the user.

## Structure

```
knowhub/
├── .synapse/schema.yaml   # Knowledge schema (behavior contract)
├── .synapse/config.yaml   # Extension point configuration
├── profile/me.md          # User profile
├── topics/                # Topic knowledge
├── entities/              # Entity pages (tools, people, projects)
├── concepts/              # Concept pages (tech concepts, methodologies)
├── inbox/                 # Pending items
├── journal/               # Timeline journal
└── graph/relations.json   # Knowledge graph
```

## Workflow

### Retrieve — Before answering

1. Check `profile/me.md` for user background
2. Search relevant files in `topics/`, `entities/`, `concepts/`
3. Reference existing knowledge naturally in your answers

### Collect — During conversation

When valuable knowledge is identified, use the CLI to persist it:

```bash
synapse collect \
  --content "knowledge content..." \
  --title "Title" \
  --topics "topic1,topic2" \
  --entities "entity1" \
  --concepts "concept1" \
  --key-points "point1,point2" \
  --source claude-code
```

### Audit — On request

When the user asks to check the knowledge base:

```bash
synapse audit
```

## Collection Rules

- **Collect**: Technical insights, new tool discoveries, design patterns, best practices, problem solutions
- **Don't collect**: Small talk, temporary debugging, duplicate knowledge, trivial fragments
- **Quality**: Distill and structure content; use 2-5 meaningful tags; link to existing knowledge with `[[wiki-links]]`

## Trigger Phrases

| User says | Action |
|-----------|--------|
| "remember this", "save to knowhub" | Collect current knowledge |
| "check knowledge base", "audit" | Run `synapse audit` |
| "what do I know about X" | Retrieve X-related knowledge |
