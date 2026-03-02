# Python Code Quality - Comprehensive Reference

Complete guide to Python code quality with ruff and pyrefly.

## Ruff (Linter & Formatter)

Ruff is an extremely fast Python linter and code formatter, written in Rust. It replaces black, isort, flake8, and many plugins.

### Installation

```bash
uv add --dev ruff
```

### Linting

```bash
# Check all files
ruff check .

# Auto-fix
ruff check --fix .

# Watch mode
ruff check --watch .
```

### Formatting

```bash
# Format code
ruff format .

# Check formatting
ruff format --check .
```

### Configuration

```toml
[tool.ruff]
line-length = 88
target-version = "py311"
src = ["src", "tests"]

[tool.ruff.lint]
select = ["E", "F", "I", "N", "UP", "B"]
ignore = ["E501"]

[tool.ruff.format]
quote-style = "double"
indent-style = "space"
```

---

## pyrefly (Type Checking)

Fast Python type checker and language server from Meta, written in Rust. Successor to Pyre, capable of checking 1.8M lines/sec.

### Installation

```bash
uv add --dev pyrefly
```

### Usage

```bash
# Initialize config (creates/updates pyproject.toml or pyrefly.toml)
pyrefly init

# Type check project
pyrefly check

# Type check with error summary
pyrefly check --summarize-errors

# Suppress existing errors with ignore comments
pyrefly suppress
```

### Configuration

```toml
[tool.pyrefly]
python-version = "3.12"
project-includes = ["src"]
project-excludes = ["**/.[!/.]*", "**/tests"]

[tool.pyrefly.errors]
bad-assignment = false
```

---

## Best Practices

1. Run ruff before committing
2. Enable type hints on all functions
3. Use pre-commit hooks
4. Configure in pyproject.toml
5. Run in CI/CD

---

## References

- **Ruff**: https://docs.astral.sh/ruff/
- **pyrefly**: https://pyrefly.org/
- **Type Hints**: https://docs.python.org/3/library/typing.html
