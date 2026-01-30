# Agent Skills Repository

This repository has two main parts:

- `skills/`: packaged skill definitions for coding agents.
- `skill-installer/`: Go TUI that installs skills into supported agents.

## Add a new skill

1. Create a new skill directory under `skills/` using the template script:

```bash
python3 skills/skill-creator-0.1.0/scripts/init_skill.py <skill-name>-0.1.0 --path skills
```

2. Edit the generated `SKILL.md` and delete any unused example resources.
3. Validate the skill (recommended):

```bash
python3 skills/skill-creator-0.1.0/scripts/quick_validate.py skills/<skill-name>-0.1.0
```

## Contribute to the installer

The installer lives in `skill-installer/`.

Common tasks:

- Run locally:

```bash
cd skill-installer
go run ./cmd/installer
```

- Add or fix features in the UI under `skill-installer/internal/ui/`.
- Update agent paths in `skill-installer/internal/agents/`.
- Update file operations in `skill-installer/internal/skills/` and `skill-installer/internal/fs/`.

## Releases

Tag pushes like `v0.1.0` trigger GitHub Actions to build binaries and attach them to the release.
