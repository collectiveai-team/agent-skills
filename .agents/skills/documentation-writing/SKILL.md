---
name: documentation-writing
description: Write skill-like documentation from a set of reference documents. Use when the user asks to synthesize a topic into imperative, readable, maintainable docs and save output to rules/<package>/docs/<topic>/<subtopic>.md, including package inference when package is omitted.
metadata:
  version: 1.0.0
---

# Documentation Writing

Write one documentation file from reference input.
Keep output imperative, readable, maintainable, and self-consistent.

## Use This Skill When

- The user wants documentation generated from several source documents.
- The user asks for a topic/subtopic guide in skill-like structure.
- The output must be written to `rules/<package>/docs/<topic>/<subtopic>.md`.
- The package may be omitted and must be inferred.

## Inputs

Collect and normalize these inputs:

- `references`: list of source files or source texts
- `topic`: domain topic for destination folder
- `subtopic`: final markdown filename stem
- `package` (optional): explicit package override
- `constraints` (optional): style, exclusions, quality checks

If `topic` or `subtopic` is missing, infer from the request and keep names in kebab-case.

## Output Path Rule

Write exactly one output file:

- `rules/<package>/docs/<topic>/<subtopic>.md`

Do not create companion `references/` or `assets/` directories for this task.

## Package Inference Rule

Resolve `<package>` in this order:

1. Use `package` if explicitly provided.
2. Match topic intent to the closest existing package.
3. Prefer the narrowest package that already contains related docs.
4. Use `code-architecture` only as final fallback.

Inference hints:

- testing topics -> `python-testing-pytest`
- API/auth/dependencies/background jobs -> `python-apis-fastapi`
- database/sqlmodel/migrations -> `python-databases`
- typing/logging/reliability/quality gates/project setup -> `python-project-foundations`
- cross-cutting architecture principles -> `code-architecture`

## Required Output Format

Generate documentation that starts with YAML frontmatter containing exactly:

- `name`
- `description`
- `tags`

Immediately after frontmatter, include this block verbatim:

```markdown
> This document is mainly for agents and LLMs to follow when maintaining,
> generating, or refactoring ${self.description}. Humans
> may also find it useful, but guidance here is optimized for automation
> and consistency by AI-assisted workflows.
```

Then write the body using this exact structure:

```markdown
# <Title>

## Purpose

<2-4 imperative sentences that define scope and intent>

## Standard

- <imperative rule 1>
- <imperative rule 2>
- <imperative rule 3>

## Guidance

1. <step 1>
2. <step 2>
3. <step 3>

## Quality Gates

- <check 1>
- <check 2>
- <check 3>

## External References

- https://<official-source-1>
- https://<official-source-2>
```

## Authoring Process

Follow this sequence every time:

1. Read references and extract stable facts.
2. Resolve conflicts by preferring precise and current facts.
3. Create a concise outline focused on execution.
4. Draft in imperative voice.
5. Remove provenance language about input references.
6. Add only authoritative external references.
7. Run quality checks.
8. Write final file to target path.

## Content Rules

Enforce all rules in the produced documentation:

- Use imperative instructions.
- Keep terms consistent across sections.
- Avoid contradictory requirements.
- Keep scope aligned to topic and subtopic.
- Keep examples and commands consistent with stated rules.
- Never mention that the content came from provided references.

## Citation Rules

Allowed references:

- Official language docs
- Official framework or library docs
- Official standards/specifications

Do not cite low-authority sources unless the user explicitly asks for them.

## Exact Body Template For Generated Docs

Use this exact sequence for generated documents:

1. `# <Title>`
2. `## Purpose`
3. `## Standard`
4. `## Guidance`
5. `## Quality Gates`
6. `## External References`

## Quality Gates

Run or emulate these checks before finishing:

- Frontmatter exists and includes `name`, `description`, `tags`
- Required notice block appears exactly once, immediately after frontmatter
- No provenance phrases like "based on provided references"
- Markdown links in `External References` are valid
- Imperative phrasing dominates procedural sections

Script-based gates are allowed. Example command pattern:

```bash
python scripts/validate_docs.py "rules/<package>/docs/<topic>/<subtopic>.md"
```

## Done Criteria

Finish only when all are true:

- Exactly one output file is produced at `rules/<package>/docs/<topic>/<subtopic>.md`
- Output matches required frontmatter + notice block format
- Documentation is imperative, readable, maintainable, self-consistent
- Package inference is correct when package input is omitted
- External references are authoritative
