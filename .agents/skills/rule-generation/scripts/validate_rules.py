#!/usr/bin/env python3
"""Validate generated rule files for strict template conformance."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


REQUIRED_FRONTMATTER_ORDER = ["title", "impact", "impactDescription", "tags"]
VALID_IMPACTS = {"CRITICAL", "HIGH", "MEDIUM-HIGH", "MEDIUM", "LOW-MEDIUM", "LOW"}
FILENAME_PATTERN = re.compile(r"^([a-z0-9]+)-([a-z0-9][a-z0-9-]*)\.md$")
TOC_SECTION_PATTERN = re.compile(r"^##\s+\d+\.\s+.+\(([^)]+)\)\s*$")
TOC_RULE_LINK_PATTERN = re.compile(r"\[[^\]]+\]\(([^)]+\.md)\)")
DOC_REFERENCE_PATTERN = re.compile(r"^rules/[a-z0-9-]+/docs/.+\.md#.+$")
PLACEHOLDER_LINE_PATTERN = re.compile(
    r"^(pass|\.\.\.|raise\s+NotImplementedError(?:\(.*\))?)$"
)
COMMENT_PREFIXES = ("#", "//", "/*", "*", "<!--")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Validate generated rule markdown files."
    )
    parser.add_argument("rules_dir", help="Path to rules/<package>/rules directory")
    parser.add_argument(
        "--validate-toc-mapping",
        action="store_true",
        help="Validate _toc.md exists and check section-prefix mapping.",
    )
    return parser.parse_args()


def load_toc_sections(rules_dir: Path) -> tuple[set[str], set[str], list[str]]:
    toc_path = rules_dir / "_toc.md"
    if not toc_path.exists():
        return set(), set(), ["Missing required _toc.md"]

    errors: list[str] = []
    section_ids: set[str] = set()
    rule_files: set[str] = set()
    for line in toc_path.read_text(encoding="utf-8").splitlines():
        match = TOC_SECTION_PATTERN.match(line)
        if match:
            section_ids.add(match.group(1).strip())
        for rule_match in TOC_RULE_LINK_PATTERN.finditer(line):
            rule_files.add(Path(rule_match.group(1)).name)

    if not section_ids:
        errors.append("No section IDs found in _toc.md")
    if not rule_files:
        errors.append("No rule links found in _toc.md")

    return section_ids, rule_files, errors


def split_frontmatter(content: str) -> tuple[list[str], str] | tuple[None, None]:
    if not content.startswith("---\n"):
        return None, None

    end_marker = "\n---\n"
    end_idx = content.find(end_marker, 4)
    if end_idx == -1:
        return None, None

    frontmatter_block = content[4:end_idx]
    body = content[end_idx + len(end_marker) :]
    return frontmatter_block.splitlines(), body


def parse_frontmatter_lines(
    lines: list[str],
) -> tuple[dict[str, str], list[str], list[str]]:
    errors: list[str] = []
    keys_in_order: list[str] = []
    data: dict[str, str] = {}

    for line in lines:
        if not line.strip():
            continue
        if ":" not in line:
            errors.append(f"Invalid frontmatter line: {line}")
            continue
        key, value = line.split(":", 1)
        key = key.strip()
        value = value.strip()
        keys_in_order.append(key)
        data[key] = value

    if keys_in_order != REQUIRED_FRONTMATTER_ORDER:
        errors.append(
            "Frontmatter keys/order must be exactly: "
            + ", ".join(REQUIRED_FRONTMATTER_ORDER)
        )

    for required in REQUIRED_FRONTMATTER_ORDER:
        if required not in data:
            errors.append(f"Missing frontmatter key: {required}")
        elif not data[required]:
            errors.append(f"Empty frontmatter value: {required}")

    impact = data.get("impact", "")
    if impact and impact not in VALID_IMPACTS:
        errors.append(f"Invalid impact '{impact}'")

    return data, keys_in_order, errors


def extract_section_code(body: str, label: str) -> str | None:
    pattern = re.compile(rf"\*\*{re.escape(label)}:\*\*\n\n```[\w+-]*\n([\s\S]*?)\n```")
    match = pattern.search(body)
    if not match:
        return None
    return match.group(1)


def has_actionable_code(code: str) -> bool:
    executable_lines: list[str] = []
    for raw_line in code.splitlines():
        line = raw_line.strip()
        if not line:
            continue
        if line.startswith(COMMENT_PREFIXES):
            continue
        executable_lines.append(line)

    if not executable_lines:
        return False

    if all(PLACEHOLDER_LINE_PATTERN.match(line) for line in executable_lines):
        return False

    return True


def validate_body(body: str, title: str, impact: str) -> list[str]:
    errors: list[str] = []
    expected_heading = f"## {title}"
    expected_impact = f"**Impact: {impact}**"

    indices: dict[str, int] = {}
    for label, token in (
        ("heading", expected_heading),
        ("impact", expected_impact),
        ("incorrect", "**Incorrect:**"),
        ("correct", "**Correct:**"),
        ("reference", "Reference: "),
    ):
        idx = body.find(token)
        if idx == -1:
            errors.append(f"Missing body token: {token}")
        indices[label] = idx

    order = [
        indices["heading"],
        indices["impact"],
        indices["incorrect"],
        indices["correct"],
        indices["reference"],
    ]
    if all(idx >= 0 for idx in order) and order != sorted(order):
        errors.append("Body sections are out of required order")

    incorrect_code = extract_section_code(body, "Incorrect")
    if incorrect_code is None:
        errors.append("Incorrect section must include a fenced code block")
    elif not has_actionable_code(incorrect_code):
        errors.append(
            "Incorrect code block must include executable, non-placeholder code"
        )

    correct_code = extract_section_code(body, "Correct")
    if correct_code is None:
        errors.append("Correct section must include a fenced code block")
    elif not has_actionable_code(correct_code):
        errors.append(
            "Correct code block must include executable, non-placeholder code"
        )

    if (
        incorrect_code is not None
        and correct_code is not None
        and incorrect_code.strip() == correct_code.strip()
    ):
        errors.append("Incorrect and Correct code blocks must differ")

    reference_line = None
    for line in body.splitlines():
        if line.startswith("Reference: "):
            reference_line = line
            break
    if reference_line is None:
        errors.append("Missing Reference line")
    else:
        ref_value = reference_line.replace("Reference: ", "", 1).strip()
        if not DOC_REFERENCE_PATTERN.match(ref_value):
            errors.append(
                "Reference must use docs anchor format: rules/<package>/docs/<doc-file>.md#<heading-anchor>"
            )

    return errors


def validate_rule_file(file_path: Path, allowed_sections: set[str] | None) -> list[str]:
    errors: list[str] = []

    match = FILENAME_PATTERN.match(file_path.name)
    if not match:
        return [
            "Filename must follow <section>-<title>.md with lowercase kebab-case title"
        ]

    section = match.group(1)
    if allowed_sections is not None and section not in allowed_sections:
        errors.append(f"Section prefix '{section}' is not declared in _toc.md")

    content = file_path.read_text(encoding="utf-8")
    parsed = split_frontmatter(content)
    if parsed == (None, None):
        return errors + ["Invalid or missing frontmatter block"]

    frontmatter_lines, body = parsed
    if frontmatter_lines is None or body is None:
        return errors + ["Invalid or missing frontmatter block"]

    data, _, fm_errors = parse_frontmatter_lines(frontmatter_lines)
    errors.extend(fm_errors)
    if fm_errors:
        return errors

    title = data["title"]
    impact = data["impact"]
    errors.extend(validate_body(body, title, impact))
    return errors


def discover_rule_files(rules_dir: Path) -> list[Path]:
    files = []
    for path in sorted(rules_dir.glob("*.md")):
        if path.name in {"_toc.md", "_template.md", "_sections.md"}:
            continue
        if path.name.startswith("_"):
            continue
        files.append(path)
    return files


def main() -> int:
    args = parse_args()
    rules_dir = Path(args.rules_dir)

    if not rules_dir.exists() or not rules_dir.is_dir():
        print(f"Rules directory does not exist: {rules_dir}")
        return 1

    allowed_sections: set[str] | None = None
    toc_rule_files: set[str] | None = None
    if args.validate_toc_mapping:
        allowed_sections, toc_rule_files, toc_errors = load_toc_sections(rules_dir)
        if toc_errors:
            for error in toc_errors:
                print(f"ERROR: {error}")
            return 1

    rule_files = discover_rule_files(rules_dir)
    if not rule_files:
        print("No rule files found to validate")
        return 1

    has_errors = False
    for rule_file in rule_files:
        errors = validate_rule_file(rule_file, allowed_sections)
        if errors:
            has_errors = True
            print(f"FAIL: {rule_file}")
            for error in errors:
                print(f"  - {error}")

        if (
            args.validate_toc_mapping
            and toc_rule_files is not None
            and rule_file.name not in toc_rule_files
        ):
            has_errors = True
            print(f"FAIL: {rule_file}")
            print("  - Rule file is not listed in _toc.md")

    if has_errors:
        return 1

    print(f"Validation passed for {len(rule_files)} rule file(s)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
