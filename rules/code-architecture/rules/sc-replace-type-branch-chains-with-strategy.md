---
title: Replace Type Branch Chains With Strategy
impact: MEDIUM
impactDescription: Improves extensibility and reduces branch growth
tags: sc-principles, solid, ocp
---

## Replace Type Branch Chains With Strategy

**Impact: MEDIUM**

Replace repeated `if/elif` type switches with strategy dispatch or a registry.

**Incorrect:**

```python
def ship(order):
    if order.kind == "air":
        return ship_air(order)
    elif order.kind == "ground":
        return ship_ground(order)
    elif order.kind == "sea":
        return ship_sea(order)
    raise ValueError("unsupported")
```

**Correct:**

```python
SHIPPERS = {
    "air": ship_air,
    "ground": ship_ground,
    "sea": ship_sea,
}


def ship(order):
    return SHIPPERS[order.kind](order)
```

Reference: rules/code-architecture/docs/sc-principles.md#solid-smell-matrix
