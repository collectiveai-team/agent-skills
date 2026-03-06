#!/usr/bin/env python3
"""Generate _toc.md from validated rule files, including all rules."""

from __future__ import annotations

import argparse
import re
import sys
from dataclasses import dataclass
from pathlib import Path


IMPACT_RANK = {
    "CRITICAL": 6,
    "HIGH": 5,
    "MEDIUM-HIGH": 4,
    "MEDIUM": 3,
    "LOW-MEDIUM": 2,
    "LOW": 1,
}

SECTION_HEADING_PATTERN = re.compile(r"^##\s+\d+\.\s+(.+)\s+\(([^)]+)\)\s*$")

DEFAULT_SECTION_METADATA = {
    "gates": (
        "Quality Gates",
        "Mandatory gate order, blocking semantics, and validator reporting rules.",
    ),
    "sc": (
        "SC Principles",
        "Core simplicity, purity, SOLID, and fail-fast error-handling rules.",
    ),
}


@dataclass
class RuleMeta:
    file_path: Path
    section_id: str
    title: str
    impact: str


@dataclass
class SectionMeta:
    section_id: str
    title: str
    description: str


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Generate _toc.md from rule files.")
    parser.add_argument("rules_dir", help="Path to rules/<package>/rules directory")
    parser.add_argument(
        "--output",
        help="Output path for _toc.md (defaults to <rules_dir>/_toc.md)",
    )
    return parser.parse_args()


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


def parse_frontmatter(lines: list[str]) -> dict[str, str]:
    data: dict[str, str] = {}
    for line in lines:
        if not line.strip() or ":" not in line:
            continue
        key, value = line.split(":", 1)
        data[key.strip()] = value.strip()
    return data


def discover_rule_files(rules_dir: Path) -> list[Path]:
    files: list[Path] = []
    for path in sorted(rules_dir.glob("*.md")):
        if path.name in {"_toc.md", "_template.md", "_sections.md"}:
            continue
        if path.name.startswith("_"):
            continue
        files.append(path)
    return files


def parse_sections_metadata(rules_dir: Path) -> dict[str, SectionMeta]:
    sections_path = rules_dir / "_sections.md"
    if not sections_path.exists():
        return {}

    lines = sections_path.read_text(encoding="utf-8").splitlines()
    metadata: dict[str, SectionMeta] = {}
    idx = 0

    while idx < len(lines):
        heading_match = SECTION_HEADING_PATTERN.match(lines[idx].strip())
        if not heading_match:
            idx += 1
            continue

        title = heading_match.group(1).strip()
        section_id = heading_match.group(2).strip()
        description = ""

        scan = idx + 1
        while scan < len(lines):
            raw = lines[scan]
            stripped = raw.strip()
            if stripped.startswith("## "):
                break
            if stripped.startswith("**Description:**"):
                description = stripped.replace("**Description:**", "", 1).strip()
                continuation = scan + 1
                while continuation < len(lines):
                    next_line = lines[continuation].strip()
                    if not next_line:
                        break
                    if next_line.startswith("## ") or next_line.startswith(
                        "**Impact:**"
                    ):
                        break
                    description += f" {next_line}"
                    continuation += 1
                break
            scan += 1

        metadata[section_id] = SectionMeta(
            section_id=section_id,
            title=title,
            description=description,
        )
        idx = scan + 1

    return metadata


def parse_section_id(filename: str, known_sections: set[str]) -> str | None:
    stem = filename.removesuffix(".md")

    matched = [section for section in known_sections if stem.startswith(f"{section}-")]
    if matched:
        return max(matched, key=len)

    if "-" not in stem:
        return None
    section = stem.split("-", 1)[0].strip()
    return section or None


def load_rule_metadata(rule_file: Path, known_sections: set[str]) -> RuleMeta:
    section_id = parse_section_id(rule_file.name, known_sections)
    if not section_id:
        raise ValueError(f"Invalid filename format: {rule_file.name}")

    parsed = split_frontmatter(rule_file.read_text(encoding="utf-8"))
    if parsed == (None, None):
        raise ValueError(f"Missing frontmatter: {rule_file}")

    frontmatter_lines, _ = parsed
    if frontmatter_lines is None:
        raise ValueError(f"Missing frontmatter: {rule_file}")

    data = parse_frontmatter(frontmatter_lines)
    title = data.get("title", "").strip()
    impact = data.get("impact", "").strip()
    if not title:
        raise ValueError(f"Missing 'title' in frontmatter: {rule_file}")
    if impact not in IMPACT_RANK:
        raise ValueError(f"Invalid or missing 'impact' in frontmatter: {rule_file}")

    return RuleMeta(
        file_path=rule_file, section_id=section_id, title=title, impact=impact
    )


def humanize_section_id(section_id: str) -> str:
    return " ".join(
        part.upper() if part == "sc" else part.capitalize()
        for part in section_id.split("-")
    )


def default_section_meta(section_id: str) -> SectionMeta:
    default = DEFAULT_SECTION_METADATA.get(section_id)
    if default:
        return SectionMeta(
            section_id=section_id,
            title=default[0],
            description=default[1],
        )

    title = humanize_section_id(section_id)
    return SectionMeta(
        section_id=section_id,
        title=title,
        description=f"Rules grouped under {title}.",
    )


def compute_section_impact(impacts: list[str]) -> str:
    return max(impacts, key=lambda impact: IMPACT_RANK[impact])


def build_toc_content(
    section_rows: list[tuple[str, str, str, str, list[RuleMeta]]],
) -> str:
    lines = [
        "# Sections",
        "",
        "This file defines all sections, their ordering, impact levels, and descriptions.",
        "The section ID (in parentheses) is the filename prefix used to group rules.",
        "",
        "---",
        "",
    ]

    for index, (section_id, section_title, impact, description, rules) in enumerate(
        section_rows, start=1
    ):
        lines.append(f"## {index}. {section_title} ({section_id})")
        lines.append("")
        lines.append(f"**Impact:** {impact}")
        lines.append(f"**Description:** {description}")
        lines.append("")

        for rule in rules:
            lines.append(f"- [{rule.title}]({rule.file_path.name}) - {rule.impact}")
        lines.append("")

    return "\n".join(lines).rstrip() + "\n"


def main() -> int:
    args = parse_args()
    rules_dir = Path(args.rules_dir)
    if not rules_dir.exists() or not rules_dir.is_dir():
        print(f"Rules directory does not exist: {rules_dir}")
        return 1

    section_metadata = parse_sections_metadata(rules_dir)
    known_sections = set(section_metadata.keys())

    rule_files = discover_rule_files(rules_dir)
    if not rule_files:
        print("No rule files found")
        return 1

    by_section: dict[str, list[RuleMeta]] = {}
    for rule_file in rule_files:
        try:
            meta = load_rule_metadata(rule_file, known_sections)
        except ValueError as error:
            print(f"ERROR: {error}")
            return 1
        by_section.setdefault(meta.section_id, []).append(meta)

    section_rows: list[tuple[str, str, str, str, list[RuleMeta]]] = []
    for section_id, rules in by_section.items():
        section_impact = compute_section_impact([rule.impact for rule in rules])
        selected_section_meta = section_metadata.get(
            section_id, default_section_meta(section_id)
        )
        section_title = selected_section_meta.title
        section_description = selected_section_meta.description

        sorted_rules = sorted(
            rules,
            key=lambda rule: (
                -IMPACT_RANK[rule.impact],
                rule.title.lower(),
                rule.file_path.name,
            ),
        )

        section_rows.append(
            (
                section_id,
                section_title,
                section_impact,
                section_description,
                sorted_rules,
            )
        )

    section_rows.sort(key=lambda row: (-IMPACT_RANK[row[2]], row[0]))

    toc_content = build_toc_content(section_rows)
    output_path = Path(args.output) if args.output else rules_dir / "_toc.md"
    output_path.write_text(toc_content, encoding="utf-8")

    total_rules = sum(len(row[4]) for row in section_rows)
    print(
        f"Generated {output_path} with {len(section_rows)} section(s) and {total_rules} rule(s)"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
