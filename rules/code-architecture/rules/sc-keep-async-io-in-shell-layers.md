---
title: Keep Async IO in Shell Layers
impact: HIGH
impactDescription: Avoids leaking transport concerns into core logic
tags: sc-principles, purity, async
---

## Keep Async IO in Shell Layers

**Impact: HIGH**

Keep `async` and `await` out of core business modules and confine async I/O to boundaries.

**Incorrect:**

```python
async def compute_risk(order_id, client):
    data = await client.fetch(order_id)
    return data["risk"]
```

**Correct:**

```python
def compute_risk(score: int, threshold: int) -> str:
    return "high" if score >= threshold else "low"
```

Reference: rules/code-architecture/docs/sc-principles.md#side-effects-to-keep-out-of-core
