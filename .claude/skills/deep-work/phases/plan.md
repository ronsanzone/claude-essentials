# Phase 5: Plan

## Purpose
Expand the structure outline into a fully detailed implementation plan. Every
task has exact file paths, function signatures, code patterns to follow, test
cases, and validation commands. The implementing agent should execute
mechanically — no architectural decisions remain.

## Inputs
- `04-structure-outline.md` from the artifact directory (phase structure, file map)
- `02-research.md` from the artifact directory (code patterns, line numbers)
- Artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Process

### Step 1: Load context
1. Read `04-structure-outline.md` completely — this is your phase structure
2. Read `02-research.md` completely — this has the code patterns and file:line refs

### Step 2: Expand each phase into tasks
For each phase in the outline, create detailed tasks. Each task covers ONE file
change (or a tightly coupled pair like implementation + test).

**Every task MUST include:**

1. **File:** Exact path, action (NEW/MODIFY), line range for modifications
2. **Pattern:** Which research finding to follow, with file:line reference
   (e.g., "Follow `pkg/handlers/user.go:30-55` pattern from Q2")
3. **What to create/modify:** Exact function/type names, signatures, struct fields.
   Include enough code detail that the implementer doesn't need to make design
   decisions. For modifications, specify what changes and what stays.
4. **Tests:** Specific test function names, test cases with inputs and expected
   outputs. Reference test patterns from research (e.g., "Follow table-driven
   pattern at `pkg/handlers/user_test.go:15-80` from Q2")
5. **Validation:** Exact command to run after the task, and expected result.
   (e.g., "`go test ./pkg/handlers/... -run TestHandleWidget` — PASS")
6. **Commit:** What to include in the commit and suggested message.

**Task granularity:** Each task should be completable in 2-5 minutes. If a task
would take longer, split it. The pattern is:
- Write failing test → Run (expect fail) → Write implementation → Run (expect pass) → Commit

### Step 3: Add phase success criteria
After all tasks in a phase, add:
- **Automated criteria:** Commands that must pass (go test, go vet, go build, etc.)
- **Manual criteria:** What to verify manually (if applicable)

### Step 4: Include scope guards
For each phase, explicitly state:
- "This phase does NOT include [X]" — prevent scope creep
- "Do NOT modify [file] in this phase" — prevent premature changes

### Step 5: Address risks
For each risk in the outline's risk register, include a specific task or
checkpoint that addresses it.

### Step 6: Write artifact

**Output file:** `05-plan.md`

The plan uses this structure:

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

## Goal
<from outline's change summary>

## Architecture
<key design decisions, restated briefly>

## Phase 1: <name>

### Overview
<what this phase accomplishes, what the system looks like after>

### Task 1.1: <descriptive name>

**File:** `exact/path/to/file.ext` (NEW|MODIFY :line-range)
**Pattern:** Follow `exact/path/to/example.ext:30-55` (Research Q<N>)

<exact details: struct fields, function signatures, logic description>

**Test:** `exact/path/to/file_test.ext`
- `TestFunctionName`: <input> → <expected output>
- `TestOtherCase`: <input> → <expected output>

**Validation:**
Run: `<exact command>`
Expected: <exact result>

**Commit:**
```
git add <files>
git commit -m "<message>"
```

### Task 1.2: ...

### Phase 1 Success Criteria
**Automated:**
- [ ] `<command>` passes
- [ ] `<command>` passes

**Manual:**
- [ ] <what to verify>

**Scope guard:** This phase does NOT include <X>.

---

## Phase 2: ...
```

Present the full plan to the user. Ask: "Full plan ready. Review and approve
to begin implementation."
