# Phase 4: Structure Outline

## Purpose
Translate design decisions into a concrete change map. This is the "what and
where" without the "how exactly." Think of it as the table of contents for the
implementation plan.

## Inputs
- `02-research.md` from the artifact directory (for file:line references)
- `03-design-discussion.md` from the artifact directory (for chosen decisions)
- Artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Process

### Step 1: Load context
1. Read `03-design-discussion.md` — focus on CHOSEN decisions, constraints, and scope
2. Read `02-research.md` — focus on file paths, patterns, and code locations

### Step 2: Map decisions to file changes
For each CHOSEN design decision, determine:
- Which files need to be created (NEW)
- Which files need to be modified (MODIFY) — with approximate line ranges from research
- Which files need to be deleted (DELETE)
- Brief description of what changes in each file

### Step 3: Group into phases
Organize file changes into sequential phases:
- Each phase should be independently testable
- Each phase should produce a working (if incomplete) system
- Dependencies between phases must be explicit
- Earlier phases establish foundations (types, interfaces)
- Later phases build on those foundations (implementation, wiring)

For each phase, include:
- **Scope:** One sentence describing what this phase accomplishes
- **Files touched:** List with NEW/MODIFY/DELETE and brief description
- **Dependencies:** Which earlier phases must complete first
- **Validation:** Concrete command or test that verifies the phase is done

### Step 4: Build file impact summary
Create a table showing every file touched across all phases:

| File | Action | Phase(s) | Reason |
|------|--------|----------|--------|

### Step 5: Compile risk register
Pull forward from design discussion:
- INCOMPLETE research findings that may surface during implementation
- Decisions that depend on assumptions not fully verified
- Mitigation strategy for each risk

### Step 6: Write artifact and present

**Output file:** `04-structure-outline.md`

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
**Files touched:**
- `path/to/file.go` (NEW) — <what it contains>
- `path/to/existing.go` (MODIFY :line-range) — <what changes>

**Dependencies:** None | Phase <N>
**Validation:** <exact command and expected result>

### Phase 2: ...

## File Impact Summary
| File | Action | Phase(s) | Reason |
|------|--------|----------|--------|
| ... | ... | ... | ... |

## Risk Register
- <risk from research gap> — Mitigation: <strategy>

## What We're NOT Doing
- <explicit scope boundaries from design discussion>
```

Present the outline to the user. Ask: "Does this scope and phasing look right?
Any phases to split, merge, or reorder?"
