---
title: Inject Infrastructure Through Abstractions
impact: HIGH
impactDescription: Keeps core policy independent from frameworks
tags: sc-principles, solid, dip
---

## Inject Infrastructure Through Abstractions

**Impact: HIGH**

Depend on interfaces in business logic and inject concrete adapters from the shell.

**Incorrect:**

```python
from psycopg import Connection


def load_profile(user_id: str):
    conn = Connection.connect("postgres://...")
    return conn.execute("select * from users where id=%s", [user_id])
```

**Correct:**

```python
class UserReader:
    def get_user(self, user_id: str):
        raise NotImplementedError


def load_profile(user_id: str, reader: UserReader):
    return reader.get_user(user_id)
```

Reference: rules/code-architecture/docs/sc-principles.md#solid-guidance
