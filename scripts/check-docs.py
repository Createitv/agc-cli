#!/usr/bin/env python3
"""Check published Markdown files, local links, anchors, and image paths."""
from __future__ import annotations

import argparse
import re
import sys
import unicodedata
from html.parser import HTMLParser
from pathlib import Path
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[1]


def without_fences(text: str) -> str:
    lines = []
    fence = None
    for line in text.splitlines():
        match = re.match(r"^\s*(`{3,}|~{3,})", line)
        if match:
            marker = match[1]
            if fence is None:
                fence = marker
            elif marker[0] == fence[0] and len(marker) >= len(fence):
                fence = None
            continue
        if fence is None:
            lines.append(line)
    if fence is not None:
        raise ValueError("unclosed code fence")
    return "\n".join(lines)


def anchors(text: str) -> set[str]:
    result = set()
    counts: dict[str, int] = {}
    for heading in re.findall(r"^#{1,6}\s+(.+?)(?:\s+#+)?$", without_fences(text), re.M):
        heading = re.sub(r"<[^>]+>", "", heading)
        heading = re.sub(r"\[([^\]]+)\]\([^)]*\)", r"\1", heading).lower()
        slug = "".join(c for c in heading if c in "_- " or unicodedata.category(c)[0] not in "PS")
        slug = slug.replace(" ", "-")
        count = counts.get(slug, 0)
        counts[slug] = count + 1
        result.add(slug if not count else f"{slug}-{count}")
    return result


class HTMLLinks(HTMLParser):
    def __init__(self) -> None:
        super().__init__()
        self.links: list[str] = []

    def handle_starttag(self, tag, attrs):
        for name, value in attrs:
            if value and ((tag == "img" and name == "src") or (tag == "a" and name == "href")):
                self.links.append(value)


def check(files: list[Path]) -> list[str]:
    errors = []
    for file in files:
        try:
            text = without_fences(file.read_text())
        except (OSError, ValueError) as exc:
            errors.append(f"{file}: {exc}")
            continue
        links = re.findall(r'\]\(\s*(<[^>]*>|[^\s)]+)(?:\s+"[^\"]*")?\s*\)', text)
        html = HTMLLinks()
        html.feed(text)
        for link in links + html.links:
            parsed = urlsplit(link.strip("<>"))
            if parsed.scheme or parsed.netloc:
                continue
            target = (file.parent / unquote(parsed.path)).resolve() if parsed.path else file
            if not target.exists():
                errors.append(f"{file}: missing target {link}")
            elif parsed.fragment and target.suffix == ".md":
                try:
                    if unquote(parsed.fragment) not in anchors(target.read_text()):
                        errors.append(f"{file}: missing anchor {link}")
                except (OSError, ValueError) as exc:
                    errors.append(f"{file}: cannot inspect {link}: {exc}")
    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("files", nargs="*", help="Optional files; defaults to published project documentation")
    args = parser.parse_args()
    if args.files:
        files = [Path(name).resolve() for name in args.files]
    else:
        patterns = ("README*.md", "CONTRIBUTING.md", "CHANGELOG.md", "docs/README.md", "docs/CLI_USAGE*.md", "docs/AGC_CLI_FULL_PLAN.md", "docs/features/**/*.md", ".github/pull_request_template.md")
        files = sorted({file for pattern in patterns for file in ROOT.glob(pattern)})
    errors = check(files)
    for error in errors:
        print(error, file=sys.stderr)
    if errors:
        return 1
    print(f"Documentation checks passed ({len(files)} files): local links, anchors, images, and code fences.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
