---
name: dw-06a-implement
description: "Use when deep-work Phase 5 plan is complete. Use when executing a written plan in a session with review checkpoints."
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

## The Process

### Step 1: Load and Review Plan
1. Read plan file
2. Review critically - identify any questions or concerns about the plan
3. If concerns: Raise them with your human partner before starting
4. If no concerns: Use TaskCreate to setup implementation tasks and proceed

### Step 1b: Bootstrap Task progress from the Plan (if needed)
The plan may have been started in another session, with progress made on some tasks.

1. Parse the plan for headers like `Phase N:` and `Task N:` to understand the structure
2. Look for progress tracking, either as an independent section or in the task and phase headers.
3. Update the task tracking system with any progress already made. Pick up from the next required task.

### Step 2: Execute Batch
**Default: First 3 tasks**

For each task:
1. Mark as in_progress
2. Follow each step exactly (plan has bite-sized steps)
3. Run verifications as specified
4. Mark as completed

### Step 3: Report
When batch complete:
- Show what was implemented
- Show verification output
- Say: "Ready for feedback."

### Step 4: Continue
Based on feedback:
- Apply changes if needed
- Execute next batch
- Repeat until complete

### Step 5: Track deviations
Note during implementation: line number shifts, unplanned tasks, materialized
risks, skipped or modified tasks.

### Step 6: Write completion artifact
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

## When to Stop and Ask for Help

**STOP executing immediately when:**
- Hit a blocker mid-batch (missing dependency, test fails, instruction unclear)
- Plan has critical gaps preventing starting
- You don't understand an instruction
- Verification fails repeatedly

**Ask for clarification rather than guessing.**

## When to Revisit Earlier Steps

**Return to Review (Step 1) when:**
- Partner updates the plan based on your feedback
- Fundamental approach needs rethinking

**Don't force through blockers** - stop and ask.

## Remember
- Review plan critically first
- Follow plan steps exactly
- Don't skip verifications
- Reference skills when plan says to
- Between batches: just report and wait
- Stop when blocked, don't guess
- Never start implementation on main/master branch without explicit user consent