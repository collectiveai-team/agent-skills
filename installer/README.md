# Agent Setup

Interactive TUI to install skills, rules, MCP servers, and subagents into supported coding agent providers.

## Run

```bash
cd installer
go run ./cmd/agent-setup
```

## What it does

- Clones the latest registry from `https://github.com/collectiveai-team/agent-skills.git`
- Discovers all resource types (skills, rules, MCP servers, subagents, profiles)
- Lets you choose: install a profile (curated bundle) or pick individual resources
- Select target providers and scope (Project or Global)
- Installs resources using provider-specific formatting:
  - Skills/subagents: copied as directories
  - Rules: merged into provider config files (idempotent with section markers)
  - MCP servers: merged into provider JSON config under `mcpServers` key

## Requirements

- `git` available on PATH
