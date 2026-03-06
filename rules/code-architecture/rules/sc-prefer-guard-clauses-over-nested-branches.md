---
title: Prefer Guard Clauses Over Nested Branches
impact: MEDIUM
impactDescription: Improves readability by flattening control flow
tags: sc-principles, kiss, readability
---

## Prefer Guard Clauses Over Nested Branches

**Impact: MEDIUM**

Use early returns to keep the main path obvious instead of stacking nested conditionals.

**Incorrect:**

```python
def compute_total(order):
    if order is not None:
        if order.items:
            if order.currency == "USD":
                return sum(item.price for item in order.items)
    return 0
```

**Correct:**

```python
def compute_total(order):
    if order is None:
        return 0
    if not order.items:
        return 0
    if order.currency != "USD":
        return 0
    return sum(item.price for item in order.items)
```

Reference: rules/code-architecture/docs/sc-principles.md#kiss-guidance
