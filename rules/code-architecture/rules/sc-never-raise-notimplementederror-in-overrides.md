---
title: Never Raise Not Implemented Error in Overrides
impact: CRITICAL
impactDescription: Prevents runtime contract breaks under polymorphism
tags: sc-principles, solid, lsp
---

## Never Raise Not Implemented Error in Overrides

**Impact: CRITICAL**

Do not ship subtype overrides that raise `NotImplementedError`; fix the contract or inheritance model.

**Incorrect:**

```python
class BaseExporter:
    def export(self, payload):
        raise NotImplementedError


class CsvExporter(BaseExporter):
    def export(self, payload):
        raise NotImplementedError
```

**Correct:**

```python
class BaseExporter:
    def export(self, payload):
        raise ValueError("contract requires concrete behavior")


class CsvExporter(BaseExporter):
    def export(self, payload):
        return ",".join(payload)
```

Reference: rules/code-architecture/docs/sc-principles.md#detection-heuristics
