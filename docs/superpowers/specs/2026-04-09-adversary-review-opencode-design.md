# Adversary Review OpenCode Skills

**Date:** 2026-04-09
**Status:** Draft

## Problem

Claude reviewing its own output (plans, code) has a structural blind spot — same model,
same training biases, same failure modes. The `opencode-consult` skill gives us a channel
to a different model (GPT-5.4 via opencode). We should use it for adversarial reviews that
provide genuinely independent perspectives.

Existing adversarial review patterns (`dw-05b-plan-review`, `spec-reviewer-prompt.md`) are
strong but run on Claude. Routing the same adversarial stance through a different model
compounds the value — different blind spots, different biases.

## Solution

Two standalone skills that compose on top of `opencode-consult`:

1. **`adversary-plan-review-opencode`** — Adversarial review of implementation plans
2. **`adversary-code-review-opencode`** — Adversarial review of code changes

Both are standalone, independently invocable, and decoupled from the deep-work pipeline.
They can be composed into pipelines later but don't require it.

## Design Decisions

- **Standalone skills, not pipeline phases.** Keeps them composable without coupling to
  deep-work phase numbering or artifact conventions.
- **Self-contained SKILL.md files.** The adversarial prompt is the skill — no separate
  prompt template files. One file per skill.
- **Lean prompts with read directives.** The subagent tells opencode which files to read
  rather than pasting code into the prompt. Keeps prompts small, leverages opencode's own
  project read access.
- **Adversarial framing from proven patterns.** Hostile default assumption, mandatory
  independent verification, strict calibration rules — all borrowed from `dw-05b` and
  `spec-reviewer-prompt.md`.

---

## Skill 1: `adversary-plan-review-opencode`

### Purpose

Send an implementation plan to opencode for adversarial review. The reviewing model
assumes the plan has gaps and looks for problems, not confirmation.

### Invocation

```
/adversary-plan-review-opencode <plan-file-path>
```

Or with a deep-work topic slug:

```
/adversary-plan-review-opencode <topic-slug>
```

When given a topic slug, resolves to `~/notes/context-engineering/<repo>/<topic-slug>/05-plan.md`
(where `<repo>` is derived from `git remote get-url origin` or `basename $(pwd)`).

### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `plan-file` or `topic-slug` | Yes | Path to plan file, or deep-work topic slug |

### Process

1. **Resolve input.**
   - **File path:** Use directly. Look for sibling artifacts in the same directory.
   - **Topic slug:** Resolve artifact directory as
     `~/notes/context-engineering/<repo>/<topic-slug>/`. The prompt must instruct
     opencode to read ALL artifacts in this directory.

2. **Discover artifacts.** The full artifact chain provides the context the reviewer
   needs to catch gaps between requirements and implementation plan:

   | Artifact | Purpose for Reviewer |
   |----------|---------------------|
   | `00-ticket.md` | Original requirements — the contract the plan must fulfill |
   | `02-research.md` | Codebase findings, patterns, constraints discovered |
   | `03-design-discussion.md` | Decided design questions, chosen approaches, scope boundaries |
   | `04-structure-outline.md` | Phase structure, risk register, scope guards |
   | `05-plan.md` | The plan under review |

   When a topic slug is provided, the prompt must include a directive like:
   > Read every file in `~/notes/context-engineering/<repo>/<topic-slug>/`. These are
   > the pipeline artifacts that led to this plan. You need all of them to understand
   > the full requirements chain — the ticket defines what must be built, the research
   > reveals codebase constraints, the design discussion records architectural decisions,
   > and the outline sets phase boundaries. Review the plan against ALL of these, not
   > just in isolation.

   If invoked with a bare file path and no companion artifacts are found, proceed with
   just the plan — the skill works standalone, but the review will be less thorough.

3. **Verify opencode server.** Run `check-server.sh`. If server is down, fail immediately.

4. **Construct adversarial prompt.** Build the prompt with:
   - Hostile default framing
   - Read directives for the plan file and companion artifacts
   - 13 review categories (from `dw-05b`)
   - Calibration rules and anti-fluff constraints
   - Structured output format

5. **Invoke opencode-consult.** Spawn subagent, send prompt via `opencode run`.

6. **Return structured result.** Present findings to user with verdict.

### Adversarial Prompt — Core Elements

The prompt sent to opencode must include these elements:

**Artifact context directive** _(included when topic slug is provided):_
> Read every file in `~/notes/context-engineering/<repo>/<topic-slug>/` before you begin
> reviewing. These artifacts form the requirements chain for this plan:
> - `00-ticket.md` — the original requirements and acceptance criteria (the contract)
> - `02-research.md` — codebase findings, patterns, and constraints
> - `03-design-discussion.md` — architectural decisions and scope boundaries (settled — do not re-litigate)
> - `04-structure-outline.md` — phase structure, risk register, scope guards
> - `05-plan.md` — the plan under review
>
> You MUST cross-reference the plan against all of these. A plan that looks fine in
> isolation may miss requirements from the ticket, contradict design decisions, or ignore
> constraints from the research.

**Framing:**
> You are performing an adversarial review of an implementation plan. Assume the plan
> has gaps until proven otherwise. Your job is to find problems, not confirm the plan
> is good.

**Independent verification mandate:**
> Read the actual plan file and any referenced source files. Do not assume file paths,
> function signatures, or module structures are correct — verify them against the
> codebase.

**Review categories:**

| Category | What to Challenge |
|----------|-------------------|
| Requirements Traceability | Every requirement maps to specific tasks. No silent scope reductions. No gold-plating. |
| Completeness | No TODOs, placeholders, or incomplete tasks. No implicit "the implementer will figure it out" gaps. |
| Task Decomposition | Tasks have clear boundaries. Steps are actionable. Dependencies are explicit. |
| Buildability | Could an engineer with zero codebase context follow this plan without getting stuck? |
| Logic Correctness | Race conditions, ordering bugs, state machine gaps, off-by-one errors, error propagation. |
| Security | Input validation, auth/authz, injection vectors, secret handling, trust boundaries. |
| Performance | N+1 queries, unbounded iterations, missing indexes, hot path allocations, missing pagination. |
| Availability & Resilience | Failure modes, retry/backoff, graceful degradation, timeout handling. |
| Durability & Data Integrity | Transaction boundaries, idempotency, migration safety, rollback path. |
| Stability & Regression Risk | Existing tests preserved, breaking changes identified, backward compatibility. |
| Code Best Practices | DRY violations across tasks. Separation of concerns. Error handling consistency. |
| Testability | Tests cover the right invariants. Missing edge case tests. Test isolation. |
| Spec Alignment | Plan implements what the spec asks — not a subset, not a superset. |

**Calibration rules:**
- Every finding MUST reference a specific task/step and explain concrete impact
- "Could be a problem" without specifics is not a finding — cut it
- "Consider adding error handling" is banned — specify WHICH error, WHERE, and WHAT happens if unhandled
- If a category has no findings, omit it — don't pad with "looks good"

**Severity levels:**

| Severity | Criteria |
|----------|----------|
| Critical | Would cause a bug, security vulnerability, data loss, or failure to meet a requirement |
| Important | Would cause performance issues, maintenance burden, fragility, or missing edge cases |
| Advisory | Would improve quality but absence won't cause failures |

**Output format:**

```
## Plan Review Verdict: APPROVED | APPROVED WITH CONDITIONS | REVISE

Verdict criteria:
- APPROVED: no Critical or Important findings
- APPROVED WITH CONDITIONS: Important findings only
- REVISE: Critical findings requiring plan changes

## Critical Issues
### [CATEGORY] Task X.Y: <title>
**What:** <specific problem>
**Impact:** <what breaks>
**Fix:** <concrete action>

## Important Issues
### [CATEGORY] Task X.Y: <title>
**What:** <specific problem>
**Impact:** <concrete consequence>
**Fix:** <concrete action>

## Advisory
- [CATEGORY] Task X.Y: <observation> — <suggested improvement>
```

---

## Skill 2: `adversary-code-review-opencode`

### Purpose

Send code changes to opencode for adversarial review. The reviewing model assumes the
implementation has defects and independently verifies all claims.

### Invocation

```
/adversary-code-review-opencode [base-sha]..[head-sha]
```

Defaults to `main..HEAD` if no range provided.

Optional spec/plan reference:

```
/adversary-code-review-opencode main..HEAD --spec path/to/spec.md
```

### Arguments

| Argument | Required | Default | Description |
|----------|----------|---------|-------------|
| `sha-range` | No | `main..HEAD` | Git range to review |
| `--spec` | No | None | Plan or spec file to verify implementation against |

### Process

1. **Resolve git range.** Parse the SHA range. Default to `main..HEAD`.

2. **Identify changed files.** Run `git diff --stat` on the range to know the scope.

3. **Verify opencode server.** Run `check-server.sh`. Fail if down.

4. **Construct adversarial prompt.** Build the prompt with:
   - Distrust-by-default framing (from `spec-reviewer-prompt.md`)
   - Read directives for git diff and changed files
   - Review checklist (correctness, security, performance, architecture, testing)
   - If `--spec` provided, add spec compliance verification
   - Calibration rules and output format

5. **Invoke opencode-consult.** Spawn subagent, send prompt via `opencode run`.

6. **Return structured result.** Present findings to user with verdict.

### Adversarial Prompt — Core Elements

**Framing:**
> You are performing an adversarial code review. The implementation may be incomplete,
> incorrect, or subtly broken. Do not assume it works — verify independently.

**Distrust mandate (from `spec-reviewer-prompt.md`):**
> DO NOT:
> - Trust commit messages about what changed
> - Assume tests prove correctness
> - Accept that "it works on my machine" means it's correct
>
> DO:
> - Read the actual diff and changed files
> - Trace logic paths through the changes
> - Check edge cases the author likely didn't consider
> - Verify error handling covers real failure modes

**Review checklist:**

| Category | What to Challenge |
|----------|-------------------|
| Correctness | Logic errors, race conditions, off-by-one, null/empty handling, state machine gaps |
| Security | Injection vectors, auth/authz gaps, secret exposure, input validation, OWASP Top 10 |
| Performance | N+1 queries, unbounded iterations, hot path allocations, missing caching, missing pagination |
| Architecture | Separation of concerns, coupling, cohesion, design pattern misuse, abstraction level |
| Error Handling | Missing error paths, swallowed errors, incorrect error propagation, missing retries |
| Testing | Missing tests for new logic, tests that test mocks not behavior, missing edge case coverage |
| Spec Compliance | _(only if --spec provided)_ Implementation matches spec. Nothing missing, nothing extra. |

**Calibration rules:**
- Every finding MUST include a file:line reference and explain why it matters
- Only report findings with high confidence — if you're uncertain, say so explicitly
- Vague findings ("improve error handling") are banned — be specific or don't report it
- If a category has no findings, omit it

**Severity levels:**

| Severity | Criteria |
|----------|----------|
| Critical | Bugs, security vulnerabilities, data loss, broken functionality |
| Important | Performance problems, missing error handling, architectural issues, test gaps |
| Advisory | Style improvements, optimization opportunities, documentation |

**Output format:**

```
## Code Review Verdict: APPROVE | REQUEST CHANGES

Verdict criteria:
- APPROVE: no Critical or Important findings
- REQUEST CHANGES: Critical or Important findings exist

## Files Reviewed
<list of files from diff>

## Critical Issues
### <file:line> — <title>
**What:** <specific problem>
**Impact:** <what breaks>
**Fix:** <concrete action>

## Important Issues
### <file:line> — <title>
**What:** <specific problem>
**Impact:** <concrete consequence>
**Fix:** <concrete action>

## Advisory
- <file:line>: <observation> — <suggested improvement>

## Strengths
- <specific positive observation with file reference>
```

---

## Interaction With opencode-consult

Both skills compose on `opencode-consult` identically:

1. Skill constructs the adversarial prompt (the unique part)
2. Skill spawns a subagent with instructions to:
   a. Verify server via `check-server.sh`
   b. Run `opencode run --attach <url> -m github-copilot/gpt-5.4 "<prompt>" 2>/dev/null`
   c. Return structured result
3. Main context receives the summary and presents to user

The subagent uses lean prompts with read directives — tells opencode which files to read
rather than pasting code. This keeps prompts manageable and lets the reviewing model form
its own understanding by reading source directly.

**Session management:** Default to new sessions (clean context per review). No need for
continue/fork unless the user explicitly wants follow-up questions.

---

## File Structure

```
.claude/skills/
├── adversary-plan-review-opencode/
│   └── SKILL.md
└── adversary-code-review-opencode/
    └── SKILL.md
```

Each SKILL.md is self-contained — no companion prompt files.

---

## What These Skills Do NOT Do

- **Don't modify the plan or code.** They report findings. The user decides what to fix.
- **Don't integrate with deep-work phase state.** No `.state.json` updates, no artifact
  directory assumptions (though they can resolve topic slugs as a convenience).
- **Don't replace Claude-native reviews.** They provide a second, independent perspective.
  Running both Claude and opencode reviews is valid and encouraged for critical work.
- **Don't auto-run.** They're invoked explicitly by the user or by other skills that choose
  to compose them.
