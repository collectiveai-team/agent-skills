---
name: python-project-setup
description: Use for Python project setup and structure. Enforce uv with pyproject.toml (use workspaces when applicable) and avoid src/ layouts by placing the package at the repository root using the package name.
---

# Python Project Setup Skill

Follow these requirements when setting up a Python project.

## Requirements

1. Use uv for dependency and project management.
2. Define project configuration in pyproject.toml.
3. Use uv workspaces when the project has multiple packages.
4. Avoid a src/ folder. Place the package directory at the repository root, named exactly after the package.
