---
title: Keep Functions Small and Single-Purpose
impact: HIGH
impactDescription: Improves readability and reduces bug-prone complexity
tags: architecture, kiss, complexity
---

## Keep Functions Small and Single-Purpose

**Impact: HIGH**

Split large mixed-responsibility functions into focused units.

**Incorrect:**

```python
def process_order(order, db, notifier):
    if order is None:
        raise ValueError("missing order")
    if order.total < 0:
        raise ValueError("invalid total")
    db.save(order)
    notifier.send_email(order.user_email, "order saved")
    return {"id": order.id, "status": "saved"}
```

**Correct:**

```python
def validate_order(order) -> None:
    if order is None:
        raise ValueError("missing order")
    if order.total < 0:
        raise ValueError("invalid total")


def save_order(order, db) -> None:
    db.save(order)


def process_order(order, db, notifier) -> str:
    validate_order(order)
    save_order(order, db)
    notifier.send_email(order.user_email, "order saved")
    return "saved"
```

Reference: rules/code-architecture/docs/sc-principles.md#principle-1-kiss-keep-it-simple
