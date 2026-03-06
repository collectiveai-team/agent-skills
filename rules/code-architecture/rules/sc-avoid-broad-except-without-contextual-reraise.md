---
title: Avoid Broad Except Without Contextual Re-raise
impact: HIGH
impactDescription: Preserves root cause and prevents hidden failures
tags: sc-principles, let-it-crash, error-handling
---

## Avoid Broad Except Without Contextual Re-raise

**Impact: HIGH**

Catch specific exception types and re-raise with clear context when translating failures.

**Incorrect:**

```python
def load_user_profile(user_id: str, repo):
    try:
        return repo.fetch_user(user_id)
    except Exception:
        return {"id": user_id, "status": "unknown"}
```

**Correct:**

```python
def load_user_profile(user_id: str, repo):
    try:
        return repo.fetch_user(user_id)
    except ConnectionError as err:
        raise RuntimeError(f"profile backend unavailable for user {user_id}") from err
```

Reference: rules/code-architecture/docs/sc-principles.md#prohibited-patterns
