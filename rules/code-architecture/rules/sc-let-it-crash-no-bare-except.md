---
title: Never Use Bare Except
impact: CRITICAL
impactDescription: Prevents silent failure and hidden defects
tags: error-handling, let-it-crash, reliability
---

## Never Use Bare Except

**Impact: CRITICAL**

Catch explicit exception types and handle them intentionally.

**Incorrect:**

```python
def load_data(path: str):
    try:
        return open(path).read()
    except:
        return ""
```

**Correct:**

```python
def load_data(path: str) -> str:
    try:
        with open(path, encoding="utf-8") as file:
            return file.read()
    except FileNotFoundError as err:
        raise ValueError(f"missing file: {path}") from err
```

Reference: rules/code-architecture/docs/sc-principles.md#prohibited-patterns
