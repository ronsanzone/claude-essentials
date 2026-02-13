---
name: dw-04-outline
description: "Use when deep-work Phase 3 design decisions are complete. Maps chosen decisions to concrete file changes organized into implementable phases."
---

# Phase 4: Structure Outline

Translate design decisions into a concrete change map — the "what and where"
without the "how exactly." A table of contents for the implementation plan.

**Announce at start:** "Starting deep-work Phase 4: Structure Outline."

## Setup

1. Parse `$ARGUMENTS` as `<topic-slug>`
   - If empty, ask user via AskUserQuestion
2. Derive repo: `basename $(git remote get-url origin 2>/dev/null | sed 's/.git$//') 2>/dev/null || basename $(pwd)`
3. Set artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Pre-flight Validation

- `03-design-discussion.md` exists → if not: "Design decisions not found. Complete Phases 1-3 first." **Stop.**
- `02-research.md` exists → if not: "Research not found. Complete Phases 1-2 first." **Stop.**

## Process

### Step 1: Load context
1. Read `03-design-discussion.md` — focus on CHOSEN decisions, constraints, scope
2. Read `02-research.md` — focus on file paths, patterns, code locations

### Step 2: Map decisions to file changes
For each CHOSEN decision, determine:
- Files to create (NEW), modify (MODIFY with line ranges from research), or delete (DELETE)
- Brief description of what changes in each file

### Step 3: Group into phases
Organize changes into sequential phases where:
- Each phase is independently testable
- Each produces a working (if incomplete) system
- Dependencies between phases are explicit
- Earlier phases establish foundations; later phases build on them

Per phase include: **Scope**, **Files touched**, **Dependencies**, **Validation command**.

### Step 4: Build file impact summary

| File | Action | Phase(s) | Reason |
|------|--------|----------|--------|

### Step 5: Compile risk register
Pull forward from design discussion: incomplete research, unverified assumptions, mitigations.

### Step 6: Write artifact
Write `04-structure-outline.md` to the artifact directory:
```yaml
---
phase: structure-outline
date: <today>
topic: <topic-slug>
repo: <repo>
git_sha: <HEAD>
input_artifacts: [02-research.md, 03-design-discussion.md]
phases_count: <N>
total_files_touched: <N>
status: complete
---

## Change Summary
<2-3 sentences describing the total change set>

## Phases

### Phase 1: <name> — <one-line goal>
**Scope:** <what this phase accomplishes>
**Files:**
- `path/to/file` (NEW) — <what it contains>
- `path/to/existing` (MODIFY :line-range) — <what changes>
**Dependencies:** None | Phase <N>
**Validation:** <exact command and expected result>

## File Impact Summary
## Risk Register
## What We're NOT Doing
```

Present outline. Ask: "Does this scope and phasing look right?"

## Completion

1. Present outline summary
2. Update `.state.json` with `current_phase: 4, completed_phases: [1, 2, 3, 4]`
3. Instruct: "Run `/dw-05-plan <topic-slug>` in a **fresh conversation** to continue."
