# Skill Installer

Interactive TUI to copy skill folders into supported coding agent locations.

## Run

```bash
cd skill-installer
go run ./cmd/installer
```

## What it does

- Clones the latest skills from `https://github.com/collectiveai-team/agent-skills.git`
- Discovers skills under `skills/*/SKILL.md` in the cloned repo
- Lets you select agents and scope (Project or Global)
- Copies each skill folder into the agent-specific skills directory
- Prompts before overwriting existing skills

## Requirements

- `git` available on PATH
