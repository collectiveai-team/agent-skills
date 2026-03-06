---
title: Keep Logging and Print Out of Core Logic
impact: HIGH
impactDescription: Preserves pure business behavior and testability
tags: sc-principles, purity, side-effects
---

## Keep Logging and Print Out of Core Logic

**Impact: HIGH**

Do not log or print inside core calculation functions; emit observability from shell layers.

**Incorrect:**

```python
def score_order(order):
    print("scoring order")
    return 100 if order.priority else 50
```

**Correct:**

```python
def score_order(order):
    return 100 if order.priority else 50
```

Reference: rules/code-architecture/docs/sc-principles.md#side-effects-to-keep-out-of-core
