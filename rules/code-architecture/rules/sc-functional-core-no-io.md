---
title: Keep Domain Logic Free of I/O
impact: HIGH
impactDescription: Makes domain behavior deterministic and testable
tags: architecture, purity, domain
---

## Keep Domain Logic Free of I/O

**Impact: HIGH**

Do not perform network, file, or database operations inside core computation logic.

**Incorrect:**

```python
def calculate_discount(user_id, session):
    user = session.get_user(user_id)
    return 0.2 if user.is_premium else 0.1
```

**Correct:**

```python
def calculate_discount(is_premium: bool) -> float:
    return 0.2 if is_premium else 0.1
```

Reference: rules/code-architecture/docs/sc-principles.md#principle-2-functional-core-imperative-shell
