# Agent Skills Repository

This repository has two main parts:

- `registry/`: packaged resource definitions (skills, rules, MCP servers, subagents, profiles) for coding agents.
- `installer/`: Go TUI (`agent-setup`) that installs resources into supported providers.

## Resource types

| Type | Marker file | Location |
|------|-------------|----------|
| Skill | `SKILL.md` | `registry/skills/*/` |
| Rule | `RULE.md` | `registry/rules/*/` |
| MCP Server | `MCP.yaml` | `registry/mcp-servers/*/` |
| Subagent | `SUBAGENT.md` | `registry/subagents/*/` |
| Profile | `PROFILE.yaml` | `registry/profiles/*/` |

Profiles bundle multiple resources into a curated set that can be installed together.

## Add a new skill

1. Create a new skill directory under `registry/skills/` using the template script:

```bash
python3 registry/skills/skill-creator-0.1.0/scripts/init_skill.py <skill-name>-0.1.0 --path registry/skills
```

2. Edit the generated `SKILL.md` and delete any unused example resources.
3. Validate the skill (recommended):

```bash
python3 registry/skills/skill-creator-0.1.0/scripts/quick_validate.py registry/skills/<skill-name>-0.1.0
```

## Contribute to the installer

The installer lives in `installer/`.

Common tasks:

- Run locally:

```bash
cd installer
go run ./cmd/agent-setup
```

- UI logic: `installer/internal/ui/`
- Provider definitions: `installer/internal/provider/`
- Resource model and discovery: `installer/internal/resource/`
- Installation adapters: `installer/internal/adapter/`
- Orchestration: `installer/internal/installer/`
- Filesystem helpers: `installer/internal/fs/`

## Supported providers

Claude Code, Cursor, Windsurf, Antigravity, Gemini, OpenCode, Codex.

## Releases and Installation

You can download the latest pre-compiled binaries from the [GitHub Releases page](https://github.com/collectiveai-team/agent-skills/releases/latest).

Tag pushes like `v0.1.0` trigger GitHub Actions to build binaries and attach them to the release.

### macOS (Apple Silicon)
1. Download `agent-setup-darwin-arm64`.
2. Make it executable and clear the macOS quarantine attribute (to bypass the "malware" warning):
   ```bash
   chmod +x agent-setup-darwin-arm64
   xattr -c agent-setup-darwin-arm64
   ```
3. (Optional) Move to a directory in your PATH to run globally:
   ```bash
   sudo mv agent-setup-darwin-arm64 /usr/local/bin/agent-setup
   ```
4. Run the installer:
   ```bash
   agent-setup
   ```

### Linux (amd64)
1. Download `agent-setup-linux-amd64`.
2. Make it executable:
   ```bash
   chmod +x agent-setup-linux-amd64
   ```
3. (Optional) Move to a directory in your PATH to run globally:
   ```bash
   sudo mv agent-setup-linux-amd64 /usr/local/bin/agent-setup
   ```
4. Run the installer:
   ```bash
   agent-setup
   ```

### Windows (amd64)
1. Download `agent-setup-windows-amd64.exe`.
2. Double-click to run, or run from Command Prompt/PowerShell:
   ```powershell
   .\agent-setup-windows-amd64.exe
   ```
