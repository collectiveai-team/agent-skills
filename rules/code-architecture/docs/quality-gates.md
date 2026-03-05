# Quality Gates

## Purpose
Define mandatory quality checks and their execution order.

## Required Gate Order
1. Ruff format
2. Ruff check
3. Pyrefly check
4. Semgrep scan
5. SC principles validators (KISS, Purity, SOLID, Let It Crash)
6. Pytest

## Gate Commands (Prex)

Run gates in order:

1. `prex run "uv run ruff format ."`
2. `prex run "uv run ruff check ."`
3. `prex run "uv run pyrefly check"`
4. `prex run "uv run semgrep scan"`
5. `prex run "python rules/code-architecture/docs/scripts/validate_kiss.py --scope-root . --json"`
6. `prex run "python rules/code-architecture/docs/scripts/validate_purity.py --scope-root . --json"`
7. `prex run "python rules/code-architecture/docs/scripts/validate_solid.py --scope-root . --json"`
8. `prex run "python rules/code-architecture/docs/scripts/validate_crash.py --scope-root . --json"`
9. `prex run "uv run pytest"`

Strict principles profile:

- `prex run "python rules/code-architecture/docs/scripts/validate_kiss.py --scope-root . --strict --json"`
- `prex run "python rules/code-architecture/docs/scripts/validate_purity.py --scope-root . --strict --json"`
- `prex run "python rules/code-architecture/docs/scripts/validate_solid.py --scope-root . --strict --json"`
- `prex run "python rules/code-architecture/docs/scripts/validate_crash.py --scope-root . --strict --json"`

## SC Principles Script References

- `rules/code-architecture/docs/scripts/validate_kiss.py`
- `rules/code-architecture/docs/scripts/validate_purity.py`
- `rules/code-architecture/docs/scripts/validate_solid.py`
- `rules/code-architecture/docs/scripts/validate_crash.py`

## SC Principles Execution Order

1. `validate_kiss.py`
2. `validate_purity.py`
3. `validate_solid.py`
4. `validate_crash.py`

## Validation Rules by Script

| Script | Rule focus | Blocking conditions (default) | Key thresholds |
|---|---|---|---|
| `validate_kiss.py` | Complexity, function size, nesting, parameter count | Any `error` | complexity `>10`, cognitive `>15`, lines `>50`, depth `>4`, params warning `>5` |
| `validate_purity.py` | I/O in core logic, side effects, async I/O | Any core `error` | Core context blocks on I/O; shell context warns |
| `validate_solid.py` | SRP/OCP/LSP/ISP/DIP structural heuristics | Any `error` (`LSP` typically) | file `>300`, class methods `>5`, interface methods `>7` |
| `validate_crash.py` | Exception-handling anti-patterns | Any `error` | nested try depth warning `>2`, bare except and except:pass are errors |

Notes:

- `validate_solid.py` and `validate_crash.py` are heuristic detectors.
- Principles validators analyze changed files by default; use `--all` for full repository scans.

## Flags Reference

### Direct Script Flags (available in copied scripts)

| Script | Flags |
|---|---|
| `validate_kiss.py` | `--scope-root`, `--threshold`, `--max-lines`, `--max-depth`, `--max-params`, `--strict`, `--all`, `--json` |
| `validate_purity.py` | `--scope-root`, `--strict`, `--core-only`, `--all`, `--json` |
| `validate_solid.py` | `--scope-root`, `--max-file-lines`, `--max-class-methods`, `--max-interface-methods`, `--strict`, `--all`, `--json` |
| `validate_crash.py` | `--scope-root`, `--max-nested-try`, `--strict`, `--all`, `--json` |

### Aggregated Runner Flags (from `/sc:principles` skill contract)

These flags are defined by the upstream skill wrapper contract:

- `--kiss-only`, `--purity-only`, `--solid-only`, `--crash-only`
- `--no-kiss`, `--no-purity`, `--no-solid`, `--no-crash`
- `--max-cognitive`, `--report`, `--inline`, `--json`

Use these only where the `/sc:principles` wrapper is available; they are not direct flags on every copied script.

## How to Read Validator Output

### Text mode patterns

- KISS header: `KISS Validation: PASSED|BLOCKED`
- Purity header: `Purity Validation: PASSED|BLOCKED`
- SOLID header: `SOLID Validation: PASSED|BLOCKED`
- Crash header: `Let It Crash Validation: PASSED|BLOCKED`

Violation line formats:

- KISS: `[SEVERITY] path:line function: violation_type = value (max: threshold)`
- Purity: `[SEVERITY] path:line function: io_type - pattern (context: core|shell)`
- SOLID: `[SEVERITY] path:line [PRINCIPLE] message`
- Crash: `[SEVERITY] path:line violation_type: message`

### JSON mode fields

All four scripts emit:

- `allowed` (bool)
- `violations` (list)
- `summary` (object)
- `recommendations` (list)

Script-specific summary keys:

| Script | Summary keys |
|---|---|
| `validate_kiss.py` | `files_analyzed`, `errors`, `warnings`, `blocked` |
| `validate_purity.py` | `files_analyzed`, `core_violations`, `shell_warnings`, `errors`, `warnings`, `blocked` |
| `validate_solid.py` | `files_analyzed`, `errors`, `warnings`, `blocked`, `by_principle` |
| `validate_crash.py` | `files_analyzed`, `errors`, `warnings`, `blocked`, `by_type` |

`validate_purity.py` also emits `architecture_hint` in JSON and text mode.

### Blocking semantics

- Default mode blocks when `errors > 0`.
- `--strict` blocks when `warnings > 0` as well.
- Purity treats core findings as errors and shell findings as warnings.

## SC Principles Exit Codes

- `0`: Validation passed.
- `2`: Violations detected (blocked; refactor required).

Upstream skill documentation reserves `3` for validation execution errors. The copied scripts currently exit with `0` or `2` in their normal code paths.

## Snippets for Automation

Single-command run:

```bash
prex run "python rules/code-architecture/docs/scripts/validate_kiss.py --scope-root . --all --json"
```

Principles batch run:

```bash
prex run "python rules/code-architecture/docs/scripts/validate_kiss.py --scope-root . --json" && prex run "python rules/code-architecture/docs/scripts/validate_purity.py --scope-root . --json" && prex run "python rules/code-architecture/docs/scripts/validate_solid.py --scope-root . --json" && prex run "python rules/code-architecture/docs/scripts/validate_crash.py --scope-root . --json"
```

## Policy
- Do not skip gate steps.
- Treat gate failures as blocking.
- Fix root causes; avoid bypass-style workarounds.

## Output Expectations
- Report which gates ran.
- Report pass/fail status.
- Include concise failure context when a gate fails.
- Include validator summary fields and top violation categories (`by_principle`, `by_type`) when present.

## Related Document
- `rules/code-architecture/docs/sc-principles.md` (principles and architecture guidance)

## External References
- https://docs.astral.sh/ruff/
- https://pyrefly.org/
- https://semgrep.dev/docs/
- https://docs.pytest.org/
