---
name: dw-05-plan
description: "Use when deep-work Phase 4 structure outline is complete. Expands the outline into a detailed implementation plan with exact file paths, code patterns, tests, and validation commands."
---

# Phase 5: Plan

Expand the structure outline into a fully detailed implementation plan. Every
task has exact file paths, function signatures, code patterns, test cases, and
validation commands. The implementing agent executes mechanically — no
architectural decisions remain.

Write comprehensive implementation plans assuming the engineer has zero context for our codebase and questionable taste. Document everything they need to know: which files to touch for each task, code, testing, docs they might need to check, how to test it. Give them the whole plan as bite-sized tasks. DRY. YAGNI. TDD. Frequent commits.

Assume they are a skilled developer, but know almost nothing about our toolset or problem domain. Assume they don't know good test design very well.

**Announce at start:** "Starting deep-work Phase 5: Plan."

## Setup

1. Parse `$ARGUMENTS` as `<topic-slug>`
   - If empty, ask user via AskUserQuestion
2. Derive repo: `basename $(git remote get-url origin 2>/dev/null | sed 's/.git$//') 2>/dev/null || basename $(pwd)`
3. Set artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Pre-flight Validation

- `04-structure-outline.md` exists → if not: "Outline not found. Complete Phases 1-4 first." **Stop.**
- `02-research.md` exists → if not: "Research not found. Complete Phases 1-2 first." **Stop.**
- `03-design-discussion.md` exists → if not: "Design decisions not found. Complete Phases 1-3 first." **Stop.**

## Process

### Step 1: Load context
1. Read `04-structure-outline.md` — phase structure and file map
2. Read `02-research.md` — code patterns and file:line references
3. Read `03-design-discussion.md` — chosen design decisions

### Step 2: Expand phases into tasks
For each phase in the outline, create tasks covering ONE file change (or tightly coupled pair).

**Every task MUST include:**
1. **File:** Exact path, action (NEW/MODIFY), line range for modifications
2. **Pattern:** Research finding to follow with file:line ref
   (e.g., "Follow `pkg/handlers/user.go:30-55` pattern from Q2")
3. **What to create/modify:** Exact names, signatures, fields — enough detail
   that the implementer makes no design decisions
4. **Tests:** Test function names, cases with inputs/expected outputs, reference
   test patterns from research
5. **Validation:** Exact command and expected result
6. **Commit:** Files to include and suggested message

**Task granularity:** 2-5 minutes each. Pattern: write failing test → run
(expect fail) → implement → run (expect pass) → commit.

**You MUST follow the outline**

#### Task Structure
## Task Structure

```markdown

### Task N: [Component Name]

**Files:**
- Create: `exact/path/to/file.py`
- Modify: `exact/path/to/existing.py:123-145`
- Test: `tests/exact/path/to/test.py`

**Step 1: Write the failing test**

```python
def test_specific_behavior():
    result = function(input)
    assert result == expected
```

**Step 2: Run test to verify it fails**

Run: `pytest tests/path/test.py::test_name -v`
Expected: FAIL with "function not defined"

**Step 3: Write minimal implementation**

```python
def function(input):
    return expected
```

**Step 4: Run test to verify it passes**

Run: `pytest tests/path/test.py::test_name -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tests/path/test.py src/path/file.py
git commit -m "feat: add specific feature"
```
```

### Step 3: Phase success criteria
Per phase: automated criteria (commands that must pass) + manual criteria.

### Step 4: Scope guards
Per phase: "This phase does NOT include [X]" and "Do NOT modify [file] in this phase."

### Step 5: Address risks
Include specific task or checkpoint for each risk from the outline's register.

### Step 6: Write artifact
Write `05-plan.md` to the artifact directory. Include plan header:

```markdown
# <Topic> Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** <from outline>
**Architecture:** <key decisions>
**Tech Stack:** <relevant tech>
```

Followed by an **Execution Progress** section (before the phase details), then full phase/task detail in standard plan format.

### Step 7: Add execution progress tracker

Insert an `## Execution Progress` section between the plan header and Phase 1. This
enables independent phase execution — an implementing agent in a fresh conversation
reads this section first to know exactly where to resume. Include three subsections:

#### 7a: Phase Progress table

```markdown
### Phase Progress

| # | Phase | Status | Validation Command | Result |
|---|-------|--------|--------------------|--------|
| 1 | <phase name> | `[ ] NOT STARTED` | `<exact validation command from phase>` | — |
| ... | ... | ... | ... | ... |

**Status legend:** `[ ] NOT STARTED` | `[~] IN PROGRESS` | `[x] DONE` | `[!] BLOCKED`
```

One row per phase. The validation command is the same command from the phase's success
criteria. The `Result` column is updated with PASS/FAIL + timestamp when validation runs.

#### 7b: Task Completion table

```markdown
### Task Completion

| Task | Description | Status | Committed | Deviations |
|------|-------------|--------|-----------|------------|
| **Phase 1** | | | | |
| 1.1 | <short description> | `[ ]` | — | |
| ... | ... | ... | ... | ... |

**Task status legend:** `[ ]` pending | `[~]` in progress | `[x]` done | `[!]` blocked | `[-]` skipped
```

One row per task, grouped under phase header rows. Description is a terse summary
(e.g., "`ApiDateCriteriaView` record"). Committed column gets the short SHA when
committed. Deviations column notes any changes from the plan.

#### 7c: Deviation Log

```markdown
### Deviation Log

> Record any deviations from the plan here. Include: task ID, what changed, why, and
> impact on downstream tasks. This is critical for maintaining plan integrity across
> sessions.

_No deviations recorded._
```

This is a free-form section. The implementing agent appends entries here when a task
requires changes from the plan. Format: `**Task X.Y:** <what changed> — <why> — <downstream impact>`.

---

End the artifact with the YAML frontmatter block:

```yaml
---
phase: plan
date: <today>
topic: <topic-slug>
repo: <repo>
git_sha: <HEAD>
input_artifacts: [02-research.md, 04-structure-outline.md]
total_phases: <N>
total_tasks: <N>
status: complete
---
```

## Completion

1. Present full plan to user for review
2. Update `.state.json` with `current_phase: 5, completed_phases: [1, 2, 3, 4, 5]`
3. Instruct: "Plan ready. Run `/dw-06-implement <topic-slug>` in a **fresh conversation** to execute."
