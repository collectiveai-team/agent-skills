---
title: Keep One Responsibility per Module
impact: MEDIUM
impactDescription: Reduces coupling and change blast radius
tags: architecture, solid, srp
---

## Keep One Responsibility per Module

**Impact: MEDIUM**

Do not combine unrelated concerns in a single module.

**Incorrect:**

```python
# users.py
def create_user(...): ...
def send_marketing_email(...): ...
def generate_pdf_report(...): ...
```

**Correct:**

```python
# users_service.py
def create_user(...): ...

# marketing_service.py
def send_marketing_email(...): ...

# report_service.py
def generate_pdf_report(...): ...
```

Reference: rules/code-architecture/docs/sc-principles.md#detection-heuristics
