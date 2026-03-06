---
title: Never Use Except Pass
impact: CRITICAL
impactDescription: Eliminates silent failures and hidden defects
tags: sc-principles, let-it-crash, error-handling
---

## Never Use Except Pass

**Impact: CRITICAL**

Do not swallow exceptions with `except: pass`; surface or translate failures intentionally.

**Incorrect:**

```python
def parse_payload(raw):
    try:
        return int(raw)
    except:
        pass
```

**Correct:**

```python
def parse_payload(raw):
    try:
        return int(raw)
    except ValueError as err:
        raise ValueError("invalid payload") from err
```

Reference: rules/code-architecture/docs/sc-principles.md#prohibited-patterns
