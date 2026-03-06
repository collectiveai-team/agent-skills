---
title: Avoid Defensive Null Guard Chains
impact: MEDIUM
impactDescription: Keeps contracts explicit and avoids masking defects
tags: sc-principles, let-it-crash, kiss
---

## Avoid Defensive Null Guard Chains

**Impact: MEDIUM**

Avoid layered `if not x` fallback chains in core logic; validate once at the boundary and keep core contracts explicit.

**Incorrect:**

```python
def calculate_invoice_total(order):
    if not order:
        return 0
    if not order.items:
        return 0
    if not order.customer:
        return 0
    return sum(item.price for item in order.items)
```

**Correct:**

```python
def validate_order(order) -> None:
    if order is None or not order.items or order.customer is None:
        raise ValueError("invalid order payload")


def calculate_invoice_total(order):
    validate_order(order)
    return sum(item.price for item in order.items)
```

Reference: rules/code-architecture/docs/sc-principles.md#prohibited-patterns
