---
name: logger
description: Logger setup and usage guidance. Use when logging is needed and no logger is present.
---

# Logger Skill

Create a logger instance if one does not exist.

## Layout

```
logger/
    __init__.py
    logger.py
```

## logger.py

```python
import logging
from logging import Logger
from rich.console import Console
from rich.logging import RichHandler

_console = Console()


def get_logger(name: str | None = None, level: int = logging.INFO) -> Logger:
    logger = logging.getLogger(name)
    if not logger.handlers:
        handler = RichHandler(console=_console, rich_tracebacks=True, markup=True)
        formatter = logging.Formatter("%(message)s")
        handler.setFormatter(formatter)
        logger.addHandler(handler)
        logger.setLevel(level)
        logger.propagate = False
    return logger
```
