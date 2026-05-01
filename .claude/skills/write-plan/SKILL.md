---
name: write-plan
description: "Use when you have a small, well-scoped change that needs an implementation plan but doesn't warrant the deep-work pipeline's research/design phases. Use instead of investigate-and-fix when you want a durable plan artifact rather than a transient EnterPlanMode plan."
---

# /write-plan

Drafts a full dw-05-format implementation plan from a brief, a file, or a Jira key. Light inline research via codebase agents fills in file:line references and patterns; the plan lands at `~/notes/context-engineering/<repo>/<slug>/plan.md` ready for `/implement-plan <slug>`.

**Use this when:** the change is small enough that Phases 1-4 of the deep-work pipeline (research questions, research, design discussion, outline) would be overkill, but you still want a durable plan artifact, fresh-subagent-per-task implementation, and two-stage review.

**Do NOT use this when:** the change has real design questions, needs cross-session paper trail, or has a scope you can't articulate in a paragraph. Use `/dw-01-research-questions <slug>` instead.

**Announce at start:** "Starting /write-plan."

## Setup

1. Run `~/.claude/skills/deep-work/dw-setup.sh "<slug>"` (extract `<slug>` from `$ARGUMENTS`; everything after the slug is the brief input). Parse stdout for `REPO`, `TOPIC_SLUG`, `ARTIFACT_DIR`.
   - If the script exits 2 (`MISSING_SLUG` on stderr), use `AskUserQuestion` to ask the user for a topic slug, then re-run with the slug.

## Pre-flight Validation

- `<ARTIFACT_DIR>/plan.md` does NOT already exist → if it DOES, use `AskUserQuestion`:
  - **Overwrite** — proceed and clobber the existing plan
  - **New slug** — ask for a different slug, re-run Setup
  - **Abort** — stop the skill
- The brief is non-empty → if `$ARGUMENTS` after the slug is empty AND no file path was provided, use `AskUserQuestion` to ask for the brief inline.

## Process

### Step 1: Parse input

`$ARGUMENTS` after the slug is the input. Resolve it to a brief string in this order:

| Input shape | Action |
|-------------|--------|
| Matches `^[A-Z]+-[0-9]+$` (Jira key, e.g. `PROJ-12345`) | Fetch via `mcp__glean_default__search` with the key. Extract problem statement, acceptance criteria, linked context. If linked docs are referenced, optionally `mcp__glean_default__read_document` for the most relevant. |
| Existing file path (use `Read`) | Read the file. Treat full contents as the brief. |
| Otherwise | Treat as free-form pasted text. |

If Glean returns nothing useful, note the gap and proceed with what's available — do NOT block.
