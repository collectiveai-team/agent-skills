---
title: Classify Core and Shell Layers Explicitly
impact: HIGH
impactDescription: Prevents boundary drift between business logic and adapters
tags: sc-principles, purity, architecture
---

## Classify Core and Shell Layers Explicitly

**Impact: HIGH**

Treat domain and logic paths as core and keep I/O in shell paths such as handlers, adapters, API, CLI, and scripts.

**Incorrect:**

```python
# domain/pricing.py
import requests

def compute_price(product_id: str) -> float:
    return requests.get(f"https://svc/{product_id}").json()["price"]
```

**Correct:**

```python
# domain/pricing.py
def compute_price(base_price: float, discount: float) -> float:
    return base_price - discount
```

Reference: rules/code-architecture/docs/sc-principles.md#layer-model
