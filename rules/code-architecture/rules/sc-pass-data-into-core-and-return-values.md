---
title: Pass Data Into Core and Return Values
impact: HIGH
impactDescription: Enforces pure function contracts in business logic
tags: sc-principles, purity, core
---

## Pass Data Into Core and Return Values

**Impact: HIGH**

Inject required data into core functions and return computed outputs instead of fetching or persisting internally.

**Incorrect:**

```python
def determine_tier(user_id, repo):
    user = repo.fetch_user(user_id)
    repo.save_audit(user_id)
    return "gold" if user.spend > 1000 else "standard"
```

**Correct:**

```python
def determine_tier(total_spend: int) -> str:
    return "gold" if total_spend > 1000 else "standard"
```

Reference: rules/code-architecture/docs/sc-principles.md#purity-guidance
