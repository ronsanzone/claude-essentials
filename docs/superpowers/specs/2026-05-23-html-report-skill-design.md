# `html-report` skill — design spec

**Date:** 2026-05-23
**Status:** Draft, pending user review
**Branch:** `html-report-skill`

## Motivation

`code-tour` currently does three jobs in one file: research orchestration, content structuring, and HTML rendering. The rendering job (Step 6) is encoded as a prose brief that leans on `frontend-design` for styling guidance. This has two problems:

1. `frontend-design`'s philosophy is "commit to a BOLD direction, never converge." That's right for product UI but wrong for long-form technical reading documents, where consistency is the brand.
2. Every future skill that wants to produce an HTML report (perf reports, retrospectives, postmortems, design docs as HTML) would re-state the same structural UX requirements (sticky TOC, scroll-spy, marginalia, collapsible sections, SVG cards) and re-invoke `frontend-design` from scratch.

We extract HTML rendering into its own skill, `html-report`, that owns both the structural shell and a curated set of visual design languages. Callers delegate the entire HTML phase with a content payload plus a design-language identifier.

## Philosophy: brand guidelines, not page templates

`html-report` is the opposite of `frontend-design`. Where `frontend-design` says "be distinctive, never converge," `html-report` says "converge on one of N curated languages, vary freely within."

The mental model is **brand guidelines**:

- **Locked elements** carry the brand: font stack, anchor palette, voice rules, anti-patterns.
- **Steered elements** vary per report: extra accent colors, exact spacing rhythm, hero composition, which components appear, diagram composition.
- **Plumbing elements** are fully canonical: scroll-spy JS, sticky-TOC HTML pattern, progress-bar wiring, skip-link. Nobody reinvents these per report.

Per-report variation is the goal, not the bug. Two reports in the same language should feel like two articles in the same publication: visibly the same voice, visibly not xeroxed.

## Skill identity

- **Name:** `html-report`
- **Location:** `~/.claude/skills/html-report/`
- **Trigger description:** "Use when you need to render a self-contained HTML report from compiled content. Triggers when a caller has structured content (sections, prose, code, diagrams) and needs styled HTML output. Used by `code-tour` and intended for future report-producing skills."
- **In scope:** structural shell rendering (TOC, scroll-spy, progress bar, layout grid), visual design-language application (typography, color, components, diagram palette), single-file HTML output.
- **Out of scope:** content compilation, research, dimension selection, audience profiling. The caller supplies all of that.

## Caller contract

The caller invokes `html-report` with:

| Input | Description |
|-------|-------------|
| `content` | Compiled outline. Sections with headings, prose, code snippets (with `file:line` citations), tables, callouts, and diagram briefs (what each diagram should show — the renderer composes the SVG). |
| `audience` | One of: `engineer-internal-ramp-up`, `engineer-internal-pr-review`, `stakeholder-external`. Used to tune density and tone within the chosen language. |
| `design-language` | Identifier matching a file in `design-languages/`. e.g. `editorial-parchment`. |
| `hero` | Title, lede/dek, eyebrow text, meta pills (date, branch, reading time, scope). |
| `output-path` | Absolute path for the resulting HTML file. |

The skill writes a single self-contained HTML file (all CSS and SVG inline, fonts via Google Fonts CDN) and prints the absolute path on completion. It does not open the file; callers decide whether to invoke `open`.

## Directory layout

```
~/.claude/skills/html-report/
├── SKILL.md
├── references/
│   ├── structural-shell.md      # Canonical plumbing (locked)
│   ├── content-components.md    # Generic HTML structure for hero/section/card/etc.
│   └── diagram-kit.md           # Generic SVG patterns + role taxonomy
└── design-languages/
    ├── README.md                # Index + selection guide
    └── editorial-parchment.md   # First language (built from atlas/admin-backup reports)
```

The two-axis split is the load-bearing design choice:

- **`references/`** — structure. What every report has. Language-agnostic.
- **`design-languages/`** — visual identity. Which paint job. Steering briefs, not templates.

This lets us add a new language without touching the structural shell, and fix a structural bug once across all languages.

## `SKILL.md` process

The SKILL.md stays small:

```
1. Receive content + audience + design-language + hero + output-path from caller.
2. Read references/structural-shell.md, references/content-components.md,
   references/diagram-kit.md, design-languages/<chosen>.md.
3. Read the exemplar report(s) listed in the chosen design-language file —
   absorb the feel.
4. Dispatch a fresh `general-purpose` subagent with: all four reference files
   inline, the exemplar paths, the content payload, the audience, and the
   output path. The work is composition (writing HTML/CSS/SVG informed by
   provided references), not codebase analysis.
5. Subagent commits to the language, writes one fresh HTML file in a single
   pass, returns the absolute path.
6. Main skill prints the path and exits.
```

**Why a fresh subagent for rendering:** the caller (e.g. `code-tour`) has spent its context on research findings from 5 parallel agents. That context is wrong for an aesthetic commitment. A fresh subagent loads only the design context (references + language + exemplar + content), commits cleanly, and produces more distinctive output. This matches how `frontend-design` works best.

## What goes in `references/structural-shell.md`

This file is **canon, not steering**. The plumbing is consistent across all reports:

- Sticky-TOC HTML pattern (the `<aside class="toc">` with `<nav>` and `<ol>`).
- Scroll-spy JS (vanilla, no dependencies; updates `.active` class on TOC links as sections enter viewport).
- Reading-progress bar wiring (CSS custom property `--scroll` updated on scroll event).
- Skip-link for accessibility.
- Layout grid: a default 2-column (TOC + main) with media-query collapse to single-column on narrow viewports. Languages may override the grid template entirely (e.g. a future `field-report` language uses a 4-column named-grid with named asides); the shell exposes the grid as a CSS variable seam, not a hardcoded template.
- Marginalia gutter pattern (`.section { position: relative; }` + `.num { position: absolute; left: -64px; ... }`).
- Document `<head>` boilerplate (charset, viewport, font preconnect).

CSS that *paints* these elements (TOC link color, active-state indicator, progress-bar color, marginalia number font) belongs to the design-language file. The HTML structure and JS behavior live here.

## What goes in `references/content-components.md`

Generic HTML structure for each report component, with comments explaining slots. Language-agnostic — no colors, no specific fonts.

Components covered:

- **Hero** — `header.hero` with eyebrow row, h1 title, lede paragraph, meta-row.
- **Section** — `<section>` with marginalia number, h2 heading, body content.
- **Card** — bordered content block.
- **Callout** — variants for note / warn / success.
- **Pull quote** — for the 1-2 stickiest insights.
- **Details block** — collapsible `<details>` styled as elegant chapter sections.
- **Table** — clean technical-doc table pattern.
- **Code block** — `<pre>` with manual `<span class="kw|str|com|num|fn|type">` highlighting.
- **Ticket pill** — colored chip for ticket IDs.
- **Diagram card** — `<figure>` wrapping inline `<svg>` with italic `<figcaption>`.

## What goes in `references/diagram-kit.md`

Generic SVG patterns plus a **role taxonomy** that design languages map to colors.

Diagram types covered (the *which*):

- Linear flow (left → right boxes + arrows).
- Swim-lane (e.g. cloud-agnostic vs. provider-specific).
- Before/after comparison.
- Sequence diagram.
- State machine.
- Hypothesis ledger (table-as-diagram for outcomes).
- System-boundary diagram.

Role taxonomy (the *which slots*):

- `default` — neutral box, the workhorse.
- `hot` — active / current focus / hot path.
- `store` — persistence / data at rest.
- `external` — third-party API or cross-system boundary.
- `warn` — caution / failure path.
- `future` — not yet wired / planned work.

Each design language fills these slots with its own palette positions. The diagram-kit file doesn't specify colors — it specifies the slots and gives generic SVG markup that uses CSS variable hooks.

It also covers:

- Arrow marker definitions (`<defs><marker>`).
- Inline `<style>` block pattern inside `<svg>` for typography.
- Caption styling reference (caption itself comes from the language).

## What goes in a `design-languages/<name>.md` file

This is a **steering brief**, not a template. Six sections:

### 1. Personality + when to pick

Two to four sentences on the aesthetic feel, plus a bullet list of pick-me signals (audience, content type, length, density).

### 2. Locked brand stack

Prescriptive. The renderer must honor these:

- Exact Google Fonts URL.
- Font stack with fallbacks for `--serif`, `--sans`, `--mono`.
- Body font choice (this is the personality knob — Editorial Parchment uses Geist sans for body; a future Field Report language would use Newsreader serif for body).
- Voice rules — e.g. "italic serif for ledes and pull quotes; mono caps for eyebrows; never use bold for emphasis in body prose."

### 3. Anchor palette

The 4–6 colors the renderer **must** use as foundation, with semantic intent. Plus an explicit "you may introduce 1–2 additional accents if the content calls for it; here's the kind of accents that fit the language."

Each anchor token includes:

- Hex value.
- Semantic name (`--bg`, `--ink`, `--accent-primary`, etc.).
- "Use for" guidance.

### 4. Anti-patterns

What *breaks* the language. e.g. "No purple gradients. No Inter/Arial. No shadows heavier than `0 10px 30px rgba(60,45,20,.04)` — kills the paper feel. No full-bleed photographs; this is a typographic language."

### 5. Patterns to draw from

Brief CSS sketches showing the *flavor* of each component in this language. Not copy-paste blocks. e.g. "Callouts use a 4px left border in the accent color over a wash background; pull quotes are large serif (≥24px), italicized, set off by a thin rule."

The renderer writes its own CSS, but in the right register.

### 6. Diagram philosophy

What kind of diagrams look right in this language. The role taxonomy from `diagram-kit.md` filled in with this language's palette positions:

| Role | Fill | Stroke |
|------|------|--------|
| default | `--bg-card` | `--rule` |
| hot | `--accent-soft` | `--accent-primary` |
| store | `--bg-soft` | `--rule-strong` |
| external | `--secondary-soft` | `--secondary` |
| warn | `--warn-soft` | `--warn` |
| future | dashed `--rule-soft` | dashed `--ink-ghost` |

Plus typography rules inside SVG (which font, which weight, which sizes).

### 7. Exemplar reports

File paths to 1–2 existing reports as ground truth. Treated as **evidence of the language**, not a template to clone. The SKILL.md instructs the rendering subagent to open and read at least one exemplar before writing.

## First language: `editorial-parchment`

Built from:

- `~/notes/04_Research/atlas-on-alibaba-cloud-tour.html` (denser, ramp-up flavor).
- `~/notes/04_Research/admin-backup-snapshots-overview.html` (shorter, operator reference flavor).

**Locked stack:**

- Fonts: Fraunces (serif display, variable opsz/wght) + Geist (sans body) + JetBrains Mono.
- Body font: Geist sans.
- Voice: italic serif for ledes / pull quotes / "em" inside h1 title; mono caps for eyebrows and meta-row labels.

**Anchor palette:**

- `--bg: #faf7f0` (warm parchment).
- `--bg-soft: #f4efe2` (deeper parchment for hero gradient, table headers).
- `--bg-card: #fdfcf8` (card surface).
- `--ink: #1c1a17` (near-black with warmth).
- `--ink-soft: #4d4842`, `--ink-faint: #8a8175`, `--ink-ghost: #b4ac9e`.
- `--rule: #d9d2c1`, `--rule-soft: #e8e2d2`.
- `--teal: #0a5c4f` (primary accent — italic-em in title, links, active TOC, primary diagram nodes).
- `--terracotta: #a84e1a` (secondary accent — eyebrow dot, marginalia number, ticket pills).

**Permitted additional accents** (if content calls for them):

- `--warn: #9a6300` (deprecated / failure path).
- `--success` from the same warm family if outcomes need positive signaling.

**Anti-patterns:**

- No Inter / Arial / Roboto.
- No purple gradients.
- No dark themes — this language is light only.
- No heavy shadows.
- No full-bleed photographs.

**Diagram role table:** filled in with the palette above.

**Exemplars:** the two files listed.

## Integration with `code-tour`

`code-tour` Step 6 (lines 142–173 of its current SKILL.md) shrinks from a prose brief plus `frontend-design` invocation to a thin delegation:

```
### Step 6: Render

Invoke the `html-report` skill with:
- content: the compiled outline from Step 5
- audience: <reader-profile from Step 3 mapped to {engineer-internal-ramp-up |
  engineer-internal-pr-review | stakeholder-external}>
- design-language: editorial-parchment (default; promote others as we add them)
- hero: title, lede, eyebrow, meta pills
- output-path: ~/notes/04_Research/<slug>-tour.html

html-report returns the absolute path; print it and run `open <path>`. If
`open` fails (worktree cleanup, sandbox, headless host), leave the path
printed.
```

The `frontend-design` reference is removed from `code-tour`. The Step 6 content (~30 lines of prose brief) is replaced with the ~10 lines above.

## Future languages (seeded but not built in this iteration)

For awareness only — these are not in scope for the implementation plan that follows:

- **`field-report`** — built from `beating_adf_report.html` + `goal_loop_journey.html`. Newsreader serif body, burnt-clay accent (`#8c3a2c`), 4-column named-grid with real asides, heavier paper grain. For outcomes / journey / postmortem reports.
- **`library-doc`** — built from `mongodb-server-extensions-tour.html`. Playfair Display, Source Sans, saddle-brown accent, fixed 260px sidebar layout, gold gradient progress bar. For "scholarly" technical docs.

Adding either language is a pure additive change: write the file in `design-languages/`, update `design-languages/README.md`, no changes to skill code or structural shell.

## Decisions log

| Decision | Choice | Reason |
|----------|--------|--------|
| Skill scope | Full HTML phase (structure + design) | Future report skills reuse the structural shell without re-stating it. |
| Design-language file format | Steering brief, not template | Per-report variation is the goal; convergence is achieved via locked elements + exemplars, not blanks-to-fill. |
| Structural shell rigidity | Locked plumbing, free presentation | Scroll-spy/TOC HTML/JS is canonical; CSS painting is steered. |
| Rendering execution | Fresh subagent | Caller's research context is wrong for aesthetic commitment; clean context produces more distinctive output. |
| Diagram styling | Generic kit + per-language role-to-palette map | One taxonomy of node roles across all languages; each language fills the slots its way. |
| First language | `editorial-parchment` only | Two existing reports already use it; lowest-risk place to land the architecture. |

## Out of scope for this design

- Building `field-report` or `library-doc` languages.
- Refactoring other skills that produce HTML (e.g. any future perf-report skill).
- Changing the contents of existing reports in `~/notes/`.
- Removing `frontend-design` (still useful for product UI work).

## Open questions for the implementation plan

- Whether `references/content-components.md` should include the manual span-class syntax highlighting recipe or defer to the language file. Leaning: keep span classes generic in components, color them in the language.
- Whether to add a smoke test (write a tiny report, diff against a golden) or rely on visual inspection of `code-tour` runs.
