---
name: dw-plan
description: "Use when deep-work Phase 4 structure outline is complete. Expands the outline into a detailed implementation plan with exact file paths, code patterns, tests, and validation commands."
---

# Phase 5: Plan

Expand the structure outline into a fully detailed implementation plan. Every
task has exact file paths, function signatures, code patterns, test cases, and
validation commands. The implementing agent executes mechanically — no
architectural decisions remain.

**Announce at start:** "Starting deep-work Phase 5: Plan."

## Setup

1. Parse `$ARGUMENTS` as `<topic-slug>`
   - If empty, ask user via AskUserQuestion
2. Derive repo: `basename $(git remote get-url origin 2>/dev/null | sed 's/.git$//') 2>/dev/null || basename $(pwd)`
3. Set artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Pre-flight Validation

- `04-structure-outline.md` exists → if not: "Outline not found. Complete Phases 1-4 first." **Stop.**
- `02-research.md` exists → if not: "Research not found. Complete Phases 1-2 first." **Stop.**

## Process

### Step 1: Load context
1. Read `04-structure-outline.md` — phase structure and file map
2. Read `02-research.md` — code patterns and file:line references

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

Followed by full phase/task detail in standard plan format.

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
3. Instruct: "Plan ready. Run `/dw-implement <topic-slug>` in a **fresh conversation** to execute."
