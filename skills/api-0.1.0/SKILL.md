---
name: api
description: API Skill for FastAPI. Use when creating or editing FastAPI API code, routers, main app setup, lifecycle, or logging.
---

# API Skill

Follow these steps when creating or editing API code.

## Review Checklist

1. Use Pydantic models for request and response payloads.
2. Group related endpoints under router modules.
3. Use main.py to include routers, add time middleware, add lifecycle events, and include a __main__ block to run the API.
4. If a logger is not created, create one following the logger skill.

## Example Structure

```
api/
    __init__.py
    main.py
    router/
        __init__.py
        users/
            __init__.py
            router.py
            models.py
            services.py
```
