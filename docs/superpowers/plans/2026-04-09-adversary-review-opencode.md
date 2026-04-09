# Adversary Review OpenCode Skills Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create two standalone skills that route adversarial reviews through opencode for independent model perspective — one for plans, one for code.

**Architecture:** Each skill is a single self-contained SKILL.md that constructs an adversarial prompt and dispatches it through opencode-consult's subagent pattern. The adversarial framing, review categories, calibration rules, and output format are all baked into the prompt sent to opencode.

**Tech Stack:** Markdown skill files, bash (opencode CLI, check-server.sh), opencode-consult building block skill.

---

## File Structure

```
.claude/skills/
├── adversary-plan-review-opencode/
│   └── SKILL.md          # Create: adversarial plan review skill
└── adversary-code-review-opencode/
    └── SKILL.md          # Create: adversarial code review skill
```

---

### Task 1: RED — Baseline test (no skill)

**Files:** None (testing only)

**Prerequisite:** opencode server must be running on port 4096.

This task establishes what happens when a subagent tries to do an adversarial plan
review via opencode WITHOUT skill guidance. We need to see what's weak so the skill
addresses real gaps.

- [ ] **Step 1: Verify opencode server is running**

Run: `~/.claude/skills/opencode-consult/check-server.sh`

Expected: prints `http://localhost:4096`. If it fails, start the server before proceeding.

- [ ] **Step 2: Run baseline plan review (no skill)**

Spawn a `general-purpose` subagent with this prompt — deliberately vague, mimicking
what an agent would do without skill guidance:

```
Review this implementation plan for problems. Be critical.

Verify the opencode server is running by executing:
  ~/.claude/skills/opencode-consult/check-server.sh
If it fails, stop and report that the opencode server is not available.

Then run:
  opencode run --attach "http://localhost:4096" -m github-copilot/gpt-5.4 "Review the implementation plan at docs/superpowers/plans/2026-04-09-adversary-review-opencode.md for this project. Be adversarial — find problems, don't confirm it's good. Report findings with severity." 2>/dev/null

Return the full response.
```

- [ ] **Step 3: Document baseline weaknesses**

Record what was weak about the unguided review. Expected failure modes:
- Vague findings without task/step references ("error handling could be improved")
- Missing review categories (security? performance? testability?)
- No structured output format (hard to act on)
- No calibration (everything marked same severity, or rubber-stamped)
- No cross-referencing against requirements artifacts
- No independent codebase verification

Save observations — these are what the skill must fix.

- [ ] **Step 4: Run baseline code review (no skill)**

Spawn a `general-purpose` subagent with this prompt:

```
Verify the opencode server is running by executing:
  ~/.claude/skills/opencode-consult/check-server.sh
If it fails, stop and report that the opencode server is not available.

Then run:
  opencode run --attach "http://localhost:4096" -m github-copilot/gpt-5.4 "Review the code changes in git range main~3..HEAD for this project. Be adversarial — assume the code has bugs. Report findings with file:line references." 2>/dev/null

Return the full response.
```

- [ ] **Step 5: Document baseline code review weaknesses**

Record what was weak. Expected failure modes:
- No structured DO/DON'T distrust framing
- Missing categories (only checks obvious bugs, misses security/performance)
- No calibration (vague "could be better" findings)
- No verdict (APPROVE/REQUEST CHANGES)
- No spec compliance checking

Save observations alongside plan review observations.

---

### Task 2: GREEN — Create adversary-plan-review-opencode skill

**Files:**
- Create: `.claude/skills/adversary-plan-review-opencode/SKILL.md`

- [ ] **Step 1: Create the skill directory**

Run: `mkdir -p .claude/skills/adversary-plan-review-opencode`

- [ ] **Step 2: Write SKILL.md**

Create `.claude/skills/adversary-plan-review-opencode/SKILL.md` with the following content.
This skill must address every baseline weakness documented in Task 1.

````markdown
---
name: adversary-plan-review-opencode
description: >
  Use when an implementation plan needs independent adversarial review from a different
  model before implementation begins. Accepts plan file path or deep-work topic slug.
---

# Adversary Plan Review (OpenCode)

Route an implementation plan to opencode for adversarial review. The reviewing model
assumes the plan has gaps and actively looks for problems — not confirmation.

**This skill composes on `opencode-consult`.** It constructs the adversarial prompt;
opencode-consult handles the transport.

## When to Use

- Before implementing a plan — especially high-stakes or complex plans
- When you want a genuinely independent perspective (different model, different biases)
- When Claude already reviewed the plan (via `dw-05b`) and you want a second opinion

**When NOT to use:**
- For quick sanity checks — use `dw-05b-plan-review` (faster, no server dependency)
- When there's no plan file yet — this reviews existing plans, not creates them
- For code review — use `adversary-code-review-opencode` instead

## Arguments

Parse `$ARGUMENTS` as one of:
- **File path** — a path to a plan file (e.g., `path/to/plan.md`)
- **Topic slug** — a deep-work topic slug (e.g., `adversary-review-opencode`)

Detection: if the argument contains `/` or ends in `.md`, treat as file path. Otherwise
treat as topic slug.

If empty, ask user via AskUserQuestion.

## Process

### Step 1: Resolve input

**If topic slug:**

```bash
REPO=$(basename "$(git remote get-url origin 2>/dev/null | sed 's/.git$//')" 2>/dev/null || basename "$(pwd)")
ARTIFACT_DIR="$HOME/notes/context-engineering/$REPO/$SLUG"
PLAN_FILE="$ARTIFACT_DIR/05-plan.md"
```

Verify `$PLAN_FILE` exists. If not: "Plan not found at `$PLAN_FILE`. Check the topic slug." **Stop.**

**If file path:**

Verify the file exists. If not: "Plan file not found at `$ARGUMENTS`." **Stop.**

Set `ARTIFACT_DIR` to the parent directory. Check for sibling artifacts (`00-ticket.md`,
`02-research.md`, `03-design-discussion.md`, `04-structure-outline.md`).

### Step 2: Build the prompt

Construct the adversarial review prompt. The prompt has two forms depending on whether
a full artifact directory was found.

**When artifact directory exists (topic slug or sibling artifacts found):**

Use the full artifact-aware prompt:

```
You are performing an adversarial review of an implementation plan. Assume the plan
has gaps until proven otherwise. Your job is to find problems, not confirm the plan
is good.

FIRST: Read every file in <ARTIFACT_DIR>. These artifacts form the requirements chain
for this plan:
- 00-ticket.md — the original requirements and acceptance criteria (the contract)
- 02-research.md — codebase findings, patterns, and constraints
- 03-design-discussion.md — architectural decisions and scope boundaries (settled — do NOT re-litigate these)
- 04-structure-outline.md — phase structure, risk register, scope guards
- 05-plan.md — the plan under review

You MUST cross-reference the plan against ALL of these. A plan that looks fine in
isolation may miss requirements from the ticket, contradict design decisions, or ignore
constraints from the research.

THEN: Read any codebase source files referenced in the plan. Do not assume file paths,
function signatures, or module structures mentioned in the plan are correct — verify
them against the actual codebase.

Review the plan against EACH of these categories. For each finding, reference the
specific task number (e.g., "Task 2.3") and explain the concrete impact:

| Category | What to Challenge |
|----------|-------------------|
| Requirements Traceability | Every requirement in 00-ticket.md maps to specific tasks. No silent scope reductions. No gold-plating beyond requirements. |
| Completeness | No TODOs, placeholders, or incomplete tasks. No implicit "the implementer will figure it out" gaps. |
| Spec Alignment | Plan implements what the spec and design decisions ask for — not a subset, not a superset. |
| Task Decomposition | Tasks have clear boundaries. Steps are actionable. Each task is independently executable. Dependencies are explicit. |
| Buildability | Could an engineer with zero codebase context follow this plan without getting stuck? Are file paths, signatures, and commands correct? |
| Logic Correctness | Race conditions, ordering bugs, state machine gaps, off-by-one errors, null/empty handling, error propagation paths. |
| Security | Input validation, auth/authz checks, injection vectors, secret handling, OWASP Top 10 relevance, trust boundary violations. |
| Performance | N+1 queries, unbounded iterations, missing indexes, large payload handling, hot path allocations, missing pagination. |
| Availability & Resilience | Failure modes, retry/backoff strategy, graceful degradation, timeout handling, dependency failure cascading. |
| Durability & Data Integrity | Transaction boundaries, idempotency, data migration safety, rollback path, schema evolution strategy. |
| Stability & Regression Risk | Existing tests preserved, breaking changes identified, backward compatibility, shared module impact. |
| Code Best Practices | DRY violations across tasks. Separation of concerns. Error handling consistency. |
| Testability | Planned tests cover the right invariants. Missing edge case tests. Integration test coverage for failure modes. Test isolation. |

CALIBRATION RULES:
- Every finding MUST reference a specific task/step and explain concrete impact
- "Could be a problem" without specifics is not a finding — cut it
- "Consider adding error handling" is banned — specify WHICH error, WHERE, and WHAT happens if unhandled
- If a category has no findings, omit it from the report — do NOT pad with "looks good"
- Do NOT re-litigate design decisions from 03-design-discussion.md — those are settled

SEVERITY:
- Critical: Would cause a bug, security vulnerability, data loss, or failure to meet a requirement
- Important: Would cause performance issues, maintenance burden, fragility, or missing edge case coverage
- Advisory: Would improve quality but absence won't cause failures

OUTPUT FORMAT (use this exactly):

## Plan Review Verdict: APPROVED | APPROVED WITH CONDITIONS | REVISE

Verdict criteria:
- APPROVED: no Critical or Important findings
- APPROVED WITH CONDITIONS: Important findings only, implementable with noted fixes
- REVISE: Critical findings that require plan changes before implementation

## Requirements Traceability
- [x] Requirement 1 → Task X.Y
- [ ] Requirement 3 → MISSING — no task covers <specific gap>

## Critical Issues
### [CATEGORY] Task X.Y: <short title>
**What:** <specific problem>
**Impact:** <what breaks, what's vulnerable, what data is lost>
**Fix:** <concrete action — add step, modify task, add test case>

## Important Issues
### [CATEGORY] Task X.Y: <short title>
**What:** <specific problem>
**Impact:** <concrete consequence>
**Fix:** <concrete action>

## Advisory
- [CATEGORY] Task X.Y: <observation> — <suggested improvement>

## Strengths
- <specific positive observation with task reference>
```

**When no artifact directory (bare file path, no siblings):**

Use a simplified prompt that omits the artifact-reading directive and the Requirements
Traceability section. Replace the opening with:

```
You are performing an adversarial review of an implementation plan. Assume the plan
has gaps until proven otherwise. Your job is to find problems, not confirm the plan
is good.

Read the plan file at <PLAN_FILE>. Then read any codebase source files referenced in
the plan. Do not assume file paths, function signatures, or module structures mentioned
in the plan are correct — verify them against the actual codebase.

Review the plan against EACH of these categories...
```

(Remainder of the prompt is the same — categories, calibration, severity, output format —
but omit Requirements Traceability from the output format since there's no ticket to
trace against.)

### Step 3: Dispatch via opencode-consult

Spawn a `general-purpose` subagent with this prompt:

```
Verify the opencode server is running by executing:
  ~/.claude/skills/opencode-consult/check-server.sh
If it fails, stop and report that the opencode server is not available.

Then run:
  opencode run --attach "$SERVER_URL" -m github-copilot/gpt-5.4 "<ADVERSARIAL_PROMPT>" 2>/dev/null

Return the result using this format:

## OpenCode Adversarial Plan Review
**Model:** github-copilot/gpt-5.4
**Plan:** <plan file path>
**Artifacts reviewed:** <list of artifact files read, or "plan only">

### Review Result
<OpenCode's full review output — preserve the verdict, findings, and format exactly>
```

### Step 4: Present result

Display the subagent's structured result to the user. Add guidance based on verdict:

- **APPROVED:** "OpenCode adversarial review found no blocking issues."
- **APPROVED WITH CONDITIONS:** "OpenCode found Important issues. Review findings above
  and decide whether to address before implementation."
- **REVISE:** "OpenCode found Critical issues. Review findings above and consider
  updating the plan."

The user decides what to do with the findings. This skill does not modify anything.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Reviewing a plan without the ticket/artifacts | Half the value is cross-referencing. Always use topic slug when artifacts exist. |
| Prompt too long for large plans | If the plan exceeds ~50 tasks, split the review — send one phase at a time and ask opencode to review each independently. |
| Treating opencode findings as authoritative | OpenCode is advisory. Findings may be wrong — the reviewing model lacks Claude's context. Always verify findings before acting. |
| Re-running without addressing prior findings | If a previous review found issues, fix them first. Don't hope a second run will disagree. |

## Failure Handling

| Failure | What to Tell the User |
|---------|----------------------|
| `check-server.sh` exits 1 | "OpenCode server is not running on port 4096. Start it before retrying." |
| `opencode run` returns empty output | "OpenCode returned no response. The server may have timed out or the model may be unavailable. Check the opencode web UI at http://localhost:4096 for session details." |
| `opencode run` returns truncated/garbled output | "OpenCode response appears incomplete. The prompt may be too long for the model's context. Try reviewing a smaller scope (single phase or subset of tasks)." |
| Subagent fails for any other reason | "Consultation failed. Check that opencode is configured correctly (`opencode auth list`) and the server is healthy." |
````

- [ ] **Step 3: Verify skill is discoverable**

Run: `claude --print-skills 2>/dev/null | grep adversary-plan`

If `--print-skills` is not available, verify by checking the file exists:

Run: `ls -la .claude/skills/adversary-plan-review-opencode/SKILL.md`

Expected: file exists with non-zero size.

- [ ] **Step 4: Commit**

```bash
git add .claude/skills/adversary-plan-review-opencode/SKILL.md
git commit -m "feat: add adversary-plan-review-opencode skill

Adversarial plan review routed through opencode for independent model
perspective. Supports topic slug (reads full artifact chain) and bare
file path. 13 review categories with strict calibration rules."
```

---

### Task 3: GREEN — Create adversary-code-review-opencode skill

**Files:**
- Create: `.claude/skills/adversary-code-review-opencode/SKILL.md`

- [ ] **Step 1: Create the skill directory**

Run: `mkdir -p .claude/skills/adversary-code-review-opencode`

- [ ] **Step 2: Write SKILL.md**

Create `.claude/skills/adversary-code-review-opencode/SKILL.md` with the following content.
This skill must address every baseline weakness documented in Task 1, Step 5.

````markdown
---
name: adversary-code-review-opencode
description: >
  Use when code changes need independent adversarial review from a different model.
  Accepts git SHA range and optional spec file for compliance checking.
---

# Adversary Code Review (OpenCode)

Route code changes to opencode for adversarial review. The reviewing model assumes
the implementation has defects and independently verifies correctness — not trusting
claims, commit messages, or test results at face value.

**This skill composes on `opencode-consult`.** It constructs the adversarial prompt;
opencode-consult handles the transport.

## When to Use

- After implementation — before merging or creating a PR
- When you want a genuinely independent code review (different model, different biases)
- When Claude already reviewed the code and you want a second opinion
- With `--spec` when you need to verify implementation matches a spec or plan

**When NOT to use:**
- For quick reviews — use `quick-review` or `pr-review` (faster, no server dependency)
- For plan review — use `adversary-plan-review-opencode` instead
- When the diff is very large (>500 changed lines) — split into logical chunks first

## Arguments

Parse `$ARGUMENTS` for:
- **SHA range** — in `base..head` format (e.g., `main..HEAD`, `abc123..def456`)
- **`--spec` flag** — optional path to a plan or spec file for compliance checking

If no arguments provided, default to `main..HEAD`.

Examples:
- `/adversary-code-review-opencode` → reviews `main..HEAD`
- `/adversary-code-review-opencode abc123..def456` → reviews that range
- `/adversary-code-review-opencode main..HEAD --spec docs/plan.md` → reviews with spec compliance

## Process

### Step 1: Resolve git range

Parse the SHA range from arguments. Default to `main..HEAD` if not provided.

Verify the range is valid:

```bash
git rev-parse --verify "$(echo "$RANGE" | cut -d. -f1)" >/dev/null 2>&1
git rev-parse --verify "$(echo "$RANGE" | cut -d. -f3-)" >/dev/null 2>&1
```

If invalid: "Invalid git range `$RANGE`. Use format `base..head`." **Stop.**

Get the diff stat to understand scope:

```bash
git diff --stat "$BASE".."$HEAD"
```

If empty diff: "No changes found in range `$RANGE`." **Stop.**

### Step 2: Resolve spec file (if provided)

If `--spec` was passed, verify the file exists. If not: "Spec file not found at
`$SPEC_FILE`." **Stop.**

### Step 3: Build the prompt

Construct the adversarial review prompt. The prompt has two forms depending on whether
a spec file was provided.

**Base prompt (always included):**

```
You are performing an adversarial code review. The implementation may be incomplete,
incorrect, or subtly broken. Do not assume it works — verify independently.

DO NOT:
- Trust commit messages about what changed
- Assume tests prove correctness
- Accept that passing CI means the code is correct
- Take any claim at face value

DO:
- Run: git diff <BASE>..<HEAD>
- Read the actual changed files in full (not just the diff hunks — surrounding context matters)
- Trace logic paths through the changes
- Check edge cases the author likely didn't consider
- Verify error handling covers real failure modes

Review the changes against EACH of these categories. For each finding, include a
file:line reference and explain why it matters:

| Category | What to Challenge |
|----------|-------------------|
| Correctness | Logic errors, race conditions, off-by-one, null/empty handling, state machine gaps, incorrect assumptions |
| Security | Injection vectors, auth/authz gaps, secret exposure, input validation, OWASP Top 10 relevance |
| Performance | N+1 queries, unbounded iterations, hot path allocations, missing caching, missing pagination |
| Architecture | Separation of concerns, coupling, cohesion, design pattern misuse, abstraction level |
| Error Handling | Missing error paths, swallowed errors, incorrect error propagation, missing retries, unclear failure modes |
| Testing | Missing tests for new logic, tests that test mocks not behavior, missing edge case coverage, flaky test patterns |

CALIBRATION RULES:
- Every finding MUST include a file:line reference and explain why it matters
- Only report findings with high confidence — if you're uncertain, say so explicitly rather than inflating certainty
- Vague findings ("improve error handling") are banned — specify WHICH error, WHERE, and WHAT happens
- If a category has no findings, omit it — do NOT pad with "looks good"

SEVERITY:
- Critical: Bugs, security vulnerabilities, data loss, broken functionality
- Important: Performance problems, missing error handling, architectural issues, test gaps
- Advisory: Style improvements, optimization opportunities, documentation gaps

OUTPUT FORMAT (use this exactly):

## Code Review Verdict: APPROVE | REQUEST CHANGES

Verdict criteria:
- APPROVE: no Critical or Important findings
- REQUEST CHANGES: Critical or Important findings exist

## Files Reviewed
<list each changed file>

## Critical Issues
### <file:line> — <short title>
**What:** <specific problem>
**Impact:** <what breaks>
**Fix:** <concrete action>

## Important Issues
### <file:line> — <short title>
**What:** <specific problem>
**Impact:** <concrete consequence>
**Fix:** <concrete action>

## Advisory
- <file:line>: <observation> — <suggested improvement>

## Strengths
- <specific positive observation with file reference>
```

**Spec compliance addendum (appended when `--spec` is provided):**

```
ADDITIONAL REVIEW DIMENSION — Spec Compliance:

Read the spec/plan file at <SPEC_FILE>. Then verify:
- Every requirement in the spec has a corresponding implementation in the diff
- No requirements were silently dropped or only partially implemented
- No extra features were added beyond what the spec calls for
- The implementation approach matches what the spec describes

Add a "Spec Compliance" section to your output between "Files Reviewed" and "Critical Issues":

## Spec Compliance
- [x] Requirement 1 — implemented in <file:line>
- [ ] Requirement 2 — MISSING or PARTIAL: <what's missing>
- [!] Extra: <file:line> — implements something not in the spec
```

### Step 4: Dispatch via opencode-consult

Spawn a `general-purpose` subagent with this prompt:

```
Verify the opencode server is running by executing:
  ~/.claude/skills/opencode-consult/check-server.sh
If it fails, stop and report that the opencode server is not available.

Then run:
  opencode run --attach "$SERVER_URL" -m github-copilot/gpt-5.4 "<ADVERSARIAL_PROMPT>" 2>/dev/null

Return the result using this format:

## OpenCode Adversarial Code Review
**Model:** github-copilot/gpt-5.4
**Range:** <base>..<head>
**Spec:** <spec file path, or "none">

### Review Result
<OpenCode's full review output — preserve the verdict, findings, and format exactly>
```

### Step 5: Present result

Display the subagent's structured result to the user. Add guidance based on verdict:

- **APPROVE:** "OpenCode adversarial review found no blocking issues."
- **REQUEST CHANGES:** "OpenCode found issues that should be addressed. Review
  findings above and decide how to proceed."

The user decides what to do with the findings. This skill does not modify anything.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Reviewing a huge diff (500+ lines) | Split into logical chunks. Large diffs overwhelm the reviewing model — findings get shallow. |
| Omitting `--spec` when a spec exists | Without spec compliance, the review only checks code quality — not whether you built the right thing. |
| Treating opencode findings as authoritative | OpenCode is advisory. Findings may be wrong — the reviewing model lacks Claude's full context. Verify before acting. |
| Using `main..HEAD` when main is stale | Make sure main is up to date (`git fetch origin main`) or the diff will include changes already merged. |

## Failure Handling

| Failure | What to Tell the User |
|---------|----------------------|
| `check-server.sh` exits 1 | "OpenCode server is not running on port 4096. Start it before retrying." |
| `git diff --stat` returns empty | "No changes found in range. Verify the SHA range is correct." |
| `opencode run` returns empty output | "OpenCode returned no response. The server may have timed out. Check http://localhost:4096 for session details." |
| `opencode run` returns truncated output | "OpenCode response appears incomplete. The diff may be too large. Try a smaller SHA range or split the review." |
| Subagent fails for any other reason | "Consultation failed. Check that opencode is configured correctly (`opencode auth list`) and the server is healthy." |
````

- [ ] **Step 3: Verify skill is discoverable**

Run: `ls -la .claude/skills/adversary-code-review-opencode/SKILL.md`

Expected: file exists with non-zero size.

- [ ] **Step 4: Commit**

```bash
git add .claude/skills/adversary-code-review-opencode/SKILL.md
git commit -m "feat: add adversary-code-review-opencode skill

Adversarial code review routed through opencode for independent model
perspective. Supports git SHA range (defaults to main..HEAD) and
optional --spec flag for compliance checking. Distrust-by-default
framing with 6 review categories and strict calibration."
```

---

### Task 4: GREEN — Verify skills improve over baseline

**Files:** None (testing only)

**Prerequisite:** Tasks 2 and 3 complete. opencode server running.

Re-run the same scenarios from Task 1, but this time use the skills. Compare against
baseline to verify the skills actually improve review quality.

- [ ] **Step 1: Test plan review WITH skill**

Invoke `/adversary-plan-review-opencode` against the same plan used in Task 1 baseline.
Use the topic slug or file path as appropriate.

- [ ] **Step 2: Compare plan review against baseline**

Verify improvements over Task 1 baseline:
- [ ] Findings reference specific task/step numbers (not vague)
- [ ] Multiple review categories covered (not just surface-level)
- [ ] Structured output format followed (verdict + severity sections)
- [ ] Calibration evident (no "consider adding error handling" fluff)
- [ ] Artifacts cross-referenced (if topic slug was used)

If any baseline weakness persists, note it for Task 5.

- [ ] **Step 3: Test code review WITH skill**

Invoke `/adversary-code-review-opencode` against the same git range used in Task 1
baseline.

- [ ] **Step 4: Compare code review against baseline**

Verify improvements over Task 1 baseline:
- [ ] Findings include file:line references (not vague)
- [ ] Multiple review categories covered (correctness, security, performance, etc.)
- [ ] Structured output format followed (verdict + severity sections)
- [ ] DO/DON'T distrust framing resulted in independent verification
- [ ] Calibration evident (high-confidence findings only)

If any baseline weakness persists, note it for Task 5.

- [ ] **Step 5: Commit test results**

No code changes needed. Verify clean state:

Run: `git status`

Expected: clean working tree.

Run: `git log --oneline -5`

Expected: commits for both skills visible.

---

### Task 5: REFACTOR — Close loopholes (if needed)

**Files:**
- Modify: `.claude/skills/adversary-plan-review-opencode/SKILL.md` (if issues found)
- Modify: `.claude/skills/adversary-code-review-opencode/SKILL.md` (if issues found)

Only execute this task if Task 4 found weaknesses that persisted from baseline.

- [ ] **Step 1: Review findings from Task 4**

List every baseline weakness that was NOT fixed by the skill.

- [ ] **Step 2: Identify prompt gaps**

For each persistent weakness, identify which part of the prompt failed:
- Missing category? Add it.
- Calibration not working? Strengthen the rule.
- Output format not followed? Add explicit "you MUST use this exact format" language.
- Artifacts not read? Make the directive more forceful.

- [ ] **Step 3: Update skill files**

Edit the affected SKILL.md file(s) to close the loopholes found.

- [ ] **Step 4: Re-test**

Re-run the affected review (plan or code) and verify the loophole is closed.

- [ ] **Step 5: Commit fixes**

```bash
git add .claude/skills/adversary-*/SKILL.md
git commit -m "refactor: close loopholes found during skill testing

Address persistent weaknesses from RED-GREEN-REFACTOR testing cycle."
```
