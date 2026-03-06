---
title: Avoid Nested Try Fallback Ladders
impact: MEDIUM
impactDescription: Keeps failure paths understandable and debuggable
tags: sc-principles, let-it-crash, reliability
---

## Avoid Nested Try Fallback Ladders

**Impact: MEDIUM**

Avoid deeply nested `try/except` fallback trees that hide the original error source.

**Incorrect:**

```python
def load_config(path):
    try:
        return read_primary(path)
    except Exception:
        try:
            return read_backup(path)
        except Exception:
            return {}
```

**Correct:**

```python
def load_config(path):
    return read_primary(path)
```

Reference: rules/code-architecture/docs/sc-principles.md#prohibited-patterns
