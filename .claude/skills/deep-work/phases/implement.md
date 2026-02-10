# Phase 6: Implement

## Purpose
Execute the implementation plan. Delegates to existing execution skills.

## Inputs
- `05-plan.md` from the artifact directory
- Artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Process

### Step 1: Load the plan
Read `05-plan.md` completely. Identify:
- Total phases and tasks
- Any already-checked items (if resuming)
- Risk items that need early verification

### Step 2: Choose execution mode
Ask the user via AskUserQuestion:

**"Plan loaded with <N> phases and <M> tasks. How do you want to execute?"**

| Option | Description |
|--------|-------------|
| Subagent-driven (this session) | Fresh subagent per task, review between tasks. Use superpowers:subagent-driven-development |
| Parallel session | Open new session in worktree with executing-plans. Use superpowers:executing-plans |
| Manual | I'll implement myself, just track progress |

### Step 3: Execute
Delegate to the chosen execution skill. The plan file path is the input.

### Step 4: Track deviations
During implementation, note any deviations from the plan:
- Line numbers that shifted
- Unplanned tasks that were needed
- Risk items that materialized
- Tasks that were skipped or modified

### Step 5: Write completion artifact

**Output file:** `06-completion.md`

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
- Phase 2: complete (<N>/<N> tasks)
...

## Deviations from Plan
- <task>: <what changed and why>

## Verification
- <command> — PASS|FAIL
- <manual check> — confirmed|failed
```

### Step 6: Cleanup
Use superpowers:finishing-a-development-branch if working in a worktree.
