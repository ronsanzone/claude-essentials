# /write-plan and /implement-plan — Design

**Date:** 2026-05-01
**Status:** Approved (brainstorming)
**Repo:** claude-essentials

## Problem

The deep-work pipeline (Phases 1-6) is the right tool for tickets that need design discussion, structured outline, and a multi-session paper trail. Smaller changes — a localized refactor, a bug fix with an obvious shape, a small feature where the design is uncontroversial — don't earn the upstream phases. Existing alternatives don't fit:

- `investigate-and-fix` produces a transient `EnterPlanMode` plan and implements in-process. No durable artifact, no fresh-subagent-per-task, no two-stage review.
- `superpowers:writing-plans` produces a durable plan, but `superpowers:executing-plans` is loose ("review plan, execute tasks") and lacks dw-06's review discipline.

Result: small changes lose the implementation quality discipline of dw-06 (fresh subagent per task, spec-compliance review, code-quality review, final session review).

## Goal

Two new skills that pair to give small changes a durable plan + dw-06-grade implementation, without the four upstream dw phases:

- `/write-plan <slug> [<brief|file|jira-key>]` — light inline research, then drafts a full dw-05-format plan to disk.
- `/implement-plan <slug>` — executes the plan using the dw-06 fresh-subagent-per-task review loop.

## Non-goals

- Replacing `investigate-and-fix`. It stays — different shape (transient plan, in-process implementation, Jira-driven flow).
- Replacing the deep-work pipeline. The new skills are for tasks that explicitly do not need Phases 1-4.
- Supporting tasks that need design discussion. If a brief implies real design questions, the skill should suggest `/dw-01-research-questions` instead.

## Architecture

```
/write-plan <slug> <brief>
   1. dw-setup.sh -> REPO, SLUG, ARTIFACT_DIR
   2. Parse input: free-form text | file path | Jira key
   3. Light research (parallel subagents):
        codebase-locator
        codebase-analyzer
        codebase-pattern-finder (conditional)
   4. Draft plan.md with dw-05 task format + inline Research Context section
   5. Write to <ARTIFACT_DIR>/plan.md
   6. Present summary; suggest /implement-plan <slug>

/implement-plan <slug>
   1. dw-setup.sh -> REPO, SLUG, ARTIFACT_DIR
   2. Pre-flight: <ARTIFACT_DIR>/plan.md exists
   3. Extract tasks from plan, create TaskCreate items
   4. For each task:
        implementer subagent (sonnet)
        -> spec-compliance reviewer (sonnet)
        -> code-quality reviewer (sonnet)
        Re-loop on rejection. Update Task Completion table.
   5. Final session review via /quick-review
   6. AskUserQuestion on Critical/Significant findings
```

Both skills share `~/.claude/skills/deep-work/dw-setup.sh` for repo + slug + artifact-dir derivation, including its `MISSING_SLUG` exit code 2 contract.

Plans live at `~/notes/context-engineering/<repo>/<slug>/plan.md` — same root as deep-work artifacts, flat single-file layout. The flat layout (no `05-` prefix) is the signal that "this is a one-shot plan, not a phase artifact." A small task that grows in scope can be promoted to the full pipeline by writing `00-ticket.md` + `02-research.md` + `03-design-discussion.md` + `04-structure-outline.md` and resuming with `/dw-03-design-discussion <slug>`.

## Components

### New files

| Path | Purpose |
|------|---------|
| `.claude/skills/write-plan/SKILL.md` | Input parsing, research dispatch, plan drafting, artifact write |
| `.claude/skills/implement-plan/SKILL.md` | Pre-flight, task extraction, per-task review loop, session review |
| `.claude/skills/implement-plan/implementer-prompt.md` | Copy of `dw-06-implement/implementer-prompt.md` |
| `.claude/skills/implement-plan/spec-reviewer-prompt.md` | Copy of `dw-06-implement/spec-reviewer-prompt.md` |
| `.claude/skills/implement-plan/code-quality-reviewer-prompt.md` | Copy of `dw-06-implement/code-quality-reviewer-prompt.md` |

The three reviewer prompts are **copied**, not symlinked or referenced cross-skill. This buys independence at the cost of drift risk; if the dw-06 prompts evolve, `/implement-plan` will not auto-inherit. Acceptable because the lighter skill may legitimately want to diverge (e.g., adjust review tone for smaller scopes) without coupling.

### Reused (no changes)

- `~/.claude/skills/deep-work/dw-setup.sh` — repo + slug + artifact-dir derivation.
- `codebase-locator`, `codebase-analyzer`, `codebase-pattern-finder` agents — for `/write-plan` light research.
- `/quick-review` skill — for `/implement-plan` final session review.

## Data flow

### `/write-plan`

1. Parse `$ARGUMENTS` into `slug` and `input`. If `slug` missing → `dw-setup.sh` exits 2 → `AskUserQuestion` for slug → re-run.
2. Resolve `input`:
   - Looks like a Jira key (`PROJ-12345`): fetch via `mcp__glean_default__search`.
   - Existing file path: `Read` the file.
   - Otherwise: treat as free-form brief.
3. Dispatch research subagents in parallel:
   - `codebase-locator`: "Find files and components related to: <key nouns from brief>."
   - `codebase-analyzer`: based on locator hits, analyze the most relevant component (current behavior, data flow, error handling, test coverage, file:line refs).
   - `codebase-pattern-finder` (conditional): only if the brief implies a known pattern (e.g., "add another handler like X").
4. Synthesize research findings inline in the plan header as a `## Research Context` section: brief restated, file:line refs, patterns to follow, constraints surfaced.
5. Draft phases and tasks per the dw-05 task format:
   - Each task: Files (NEW/MODIFY + paths + line range), Pattern (ref to research finding), What to create/modify (signatures, fields, exact names), Tests (function names, inputs, expected outputs), Validation (command + expected result), Commit (files + message).
   - TDD pattern: failing test → run (expect fail) → implement → run (expect pass) → commit.
   - Phase Progress table, Task Completion table, Deviation Log, frontmatter — all per dw-05.
6. Write `plan.md` to `<ARTIFACT_DIR>`.
7. Present summary; instruct: "Plan ready at `<path>`. Run `/implement-plan <slug>` to execute."

### `/implement-plan`

Direct port of `dw-06-implement` Process section, with three substitutions:

1. Read `plan.md` (not `05-plan.md`).
2. No `.state.json` updates (no pipeline state to track).
3. Implementer's reference document is the `## Research Context` section of `plan.md` itself, not separate `00-ticket.md` / `02-research.md` files. The implementer subagent prompt passes the relevant Research Context excerpt alongside the task text.

Otherwise identical: `TaskCreate` per task, dispatch implementer → spec reviewer → code-quality reviewer per task, re-loop on rejection, update Task Completion table, final `/quick-review`, `AskUserQuestion` on Critical/Significant findings.

## Plan format

Same as `dw-05-plan` output, with one addition and zero removals.

```markdown
# <Topic> Implementation Plan

**Goal:** <one sentence>
**Architecture:** <key decisions>
**Tech Stack:** <relevant tech>

## Research Context

### Brief
<original input, normalized>

### Files in scope
- `path/to/file.ext:LINES` — <what it does>
- ...

### Patterns to follow
- `path/to/example.ext:LINES` — <pattern name and shape>
- ...

### Constraints
- <surfaced constraint 1>
- ...

## Execution Progress
### Phase Progress
| # | Phase | Status | Validation Command | Result |
| ... |

### Task Completion
| Task | Description | Status | Committed | Deviations |
| ... |

### Deviation Log
> Record deviations here. Format: **Task X.Y:** <what> — <why> — <impact>.
_No deviations recorded._

## Phase 1: <name>
### Task 1.1: <name>
**Files:** ...
**Step 1: Write the failing test** ...
**Step 2: Run test (expect FAIL)** ...
**Step 3: Implement** ...
**Step 4: Run test (expect PASS)** ...
**Step 5: Commit** ...

...

---
phase: plan
date: 2026-05-01
topic: <slug>
repo: <repo>
git_sha: <HEAD>
total_phases: <N>
total_tasks: <N>
status: complete
---
```

## Error handling

| Scenario | Behavior |
|----------|----------|
| `slug` missing | `dw-setup.sh` exits 2 → `AskUserQuestion` → re-run |
| `plan.md` already exists in `/write-plan` | `AskUserQuestion`: overwrite / new slug / abort |
| `plan.md` missing in `/implement-plan` | "Plan not found at `<path>`. Run `/write-plan <slug>` first." Stop. |
| Implementer subagent asks questions | Answer; re-dispatch (per dw-06) |
| Reviewer rejects | Same implementer fixes, re-review (per dw-06) |
| Resume mid-implementation | `/implement-plan <slug>` reads Task Completion table, skips `[x]` tasks |
| Brief implies real design questions | `/write-plan` recognizes ambiguity, suggests `/dw-01-research-questions <slug>` instead, stops |

## Relationship to existing skills

| Skill | When to use |
|-------|-------------|
| `/dw-01-research-questions` ... `/dw-06-implement` | Tickets needing design discussion, multi-session paper trail, or large scope |
| `/deep-work-pipeline` | Same as above, end-to-end in one session |
| `/investigate-and-fix` | Quick Jira-driven fix; transient plan via `EnterPlanMode`; in-process implementation |
| `/write-plan` + `/implement-plan` | **(new)** Small change with durable plan + dw-06 implementation discipline; no upstream phases |
| `/refine-ticket` | Upstream of any of the above; produces a refined `ticket.md` |

## Testing

Skills cannot be unit-tested. Validation is dogfooding:

1. Pick a small recent fix or feature from the repo's git log (one with 3-5 commits).
2. Run `/write-plan <slug> <brief>` from the original brief.
3. Compare the generated plan against the actual commits — does the plan's task decomposition resemble what was actually done?
4. Run `/implement-plan <slug>` on a fresh worktree of the parent commit; compare the resulting code to the historical commits.
5. Iterate skill prompts until the dogfood tests produce comparable quality.

## Open questions

- **Dogfood candidate**: pick a recent commit range from claude-essentials or another repo to validate against. Defer until skills are written.
- **Glean key detection regex**: same regex `investigate-and-fix` uses, or relaxed? Defer to implementation; check `investigate-and-fix` SKILL.md.
- **Worktree integration**: should `/implement-plan` auto-create a worktree like the dw pipeline? Current call: no — keep it minimal; user can `EnterWorktree` themselves if desired. Revisit after first dogfood pass.

## Migration / rollout

No migration needed — purely additive. Existing skills unchanged. Once written, add the two new skill names to any internal skill index docs.
