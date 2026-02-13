---
name: dw-06-implement
description: "Use when deep-work Phase 5 plan is complete. Executes the implementation plan with progress tracking, deviation logging, and completion verification."
---

# Phase 6: Implement

Execute the implementation plan. Delegates to existing execution skills.

**Announce at start:** "Starting deep-work Phase 6: Implementation."

## Setup

1. Parse `$ARGUMENTS` as `<topic-slug>`
   - If empty, ask user via AskUserQuestion
2. Derive repo: `basename $(git remote get-url origin 2>/dev/null | sed 's/.git$//') 2>/dev/null || basename $(pwd)`
3. Set artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Pre-flight Validation

- `05-plan.md` exists → if not: "Plan not found. Complete Phases 1-5 first." **Stop.**

## Process

### Step 1: Load plan
Read `05-plan.md` completely. Identify total phases, tasks, and risk items.

### Step 2: Choose execution mode
Ask via AskUserQuestion:

| Option | Description |
|--------|-------------|
| Subagent-driven (this session) | Fresh subagent per task, review between. **REQUIRED:** superpowers:subagent-driven-development |
| Parallel session | New session with executing-plans. **REQUIRED:** superpowers:executing-plans |
| Manual | User implements, skill tracks progress |

### Step 3: Execute
Delegate to chosen execution skill with the plan file path as input.

### Step 4: Track deviations
Note during implementation: line number shifts, unplanned tasks, materialized
risks, skipped or modified tasks.

### Step 5: Write completion artifact
Write `06-completion.md` to the artifact directory:
```yaml
---
phase: implementation
date: <today>
topic: <topic-slug>
repo: <repo>
git_sha_start: <HEAD at start>
git_sha_end: <HEAD at end>
input_artifacts: [05-plan.md]
status: complete
---

## Completion Summary
- Phase 1: complete (<N>/<N> tasks)
...

## Deviations from Plan
- <task>: <what changed and why>

## Verification
- <command> — PASS|FAIL
```

## Completion

1. Present completion summary
2. Update `.state.json` with `current_phase: 6, completed_phases: [1, 2, 3, 4, 5, 6]`
3. Suggest: "Implementation complete. Use superpowers:finishing-a-development-branch to wrap up."
