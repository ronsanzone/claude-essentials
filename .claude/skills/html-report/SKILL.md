---
name: html-report
description: Use when a caller has compiled content (sections, prose, code, diagrams) and needs to render it as a self-contained HTML technical report. Owns structural UX (TOC, scroll-spy, progress bar, marginalia, collapsible sections, SVG diagram cards) and visual design languages. Used by code-tour; intended for future report-producing skills.
---

# html-report

Render a self-contained HTML report from compiled content, applying a curated design language. Brand-guidelines philosophy: locked font stack and anchor palette per language; free composition within.

## Usage

Invoked by a calling skill (e.g. `code-tour`) with:

- **content** — compiled outline (sections, prose, code with `file:line` citations, tables, callouts, diagram briefs)
- **audience** — one of: `engineer-internal-ramp-up`, `engineer-internal-pr-review`, `stakeholder-external`
- **design-language** — identifier matching a file in `design-languages/` (e.g. `editorial-parchment`)
- **hero** — title, lede, eyebrow text, meta pills
- **output-path** — absolute path for the resulting HTML file

## Process

(See later task — finalized after references are written.)
