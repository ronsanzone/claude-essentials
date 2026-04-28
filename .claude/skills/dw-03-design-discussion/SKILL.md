---
name: dw-03-design-discussion
description: "Use when deep-work Phase 2 research is complete. Combines research findings with the original task to explore design options, evaluate tradeoffs, and make decisions interactively."
---

# Phase 3: Design Discussion

Combine objective research findings with the original prompt to identify design
decisions, enumerate options, and evaluate tradeoffs. Research is locked in —
the prompt safely re-enters the pipeline here.

**Announce at start:** "Starting deep-work Phase 3: Design Discussion."

## Setup

1. Run `~/.claude/skills/deep-work/dw-setup.sh "$ARGUMENTS"` and parse stdout for `REPO`, `TOPIC_SLUG`, `ARTIFACT_DIR`.
   - If the script exits 2 (`MISSING_SLUG` on stderr), ask user via AskUserQuestion for the topic slug, then re-run with the slug.

## Pre-flight Validation

- `02-research.md` exists → if not: "Research not found. Complete Phases 1-2 first." **Stop.**
- `00-ticket.md` exists → if not: "No ticket found. Run `/dw-research-questions` first." **Stop.**

## Process

### Step 1: Load context
1. Read `02-research.md` completely
2. Read `00-ticket.md` completely
3. Summarize: "The user wants to [goal from ticket]. Research found [key findings]."

### Step 2: Synthesize context into design sections

Before identifying decisions, distill research into structured sections:

**Summary of Changes Requested**
- Restate the ticket goal concisely — what are we doing and why?

**Current State**
- From research findings, extract the parts of the system that are relevant to the
  requested changes. Include file:line references. Focus on what exists today that
  the changes will touch, integrate with, or depend on.
- Pull from the research's "Structured Summary > System State" section if available,
  then filter to only what's relevant to the ticket.

**Desired End State**
- Combine the ticket goal with current state to describe concretely what the system
  looks like after the changes. Be specific — name the files, APIs, behaviors.

**What We're Not Doing**
- Identify explicit scope boundaries. What might a reader assume is included but isn't?
  What adjacent changes are we deferring?

**Patterns to Follow**
- From research findings (especially Pattern Discovery and Code Tracing answers),
  extract concrete codebase patterns that implementation should mirror.
- Each pattern MUST include file:line references to exemplar code.
- Pull from the research's "Structured Summary > Patterns Found" section if available.
- In brownfield codebases, there may be multiple patterns — choose which to follow
  and note why.

### Step 3: Identify design questions
Based on the gap between "current state" and "desired end state," identify every design question that needs resolution. Common types include placement (where new code lives), pattern selection (which existing pattern to mirror), integration (how to connect with existing code), API/interface shape, edge-case handling, and testability.

### Step 4: Build options
For EACH question, create 2-4 options. Every option MUST:
- Cite a specific research finding (e.g., "Research Q2 found that...")
- Include concrete pros and cons
- Be grounded in what exists in the codebase

Include your **recommendation** with rationale for each question.

**FORBIDDEN:** Options that ignore research findings or require uninvestigated changes.

### Step 5: Surface risks
Compile: constraints from research, INCOMPLETE research gaps, out-of-scope items.

### Step 5b: Targeted exploration (conditional)

Review the design questions from Step 3-4 and the risks from Step 5. Identify any
questions where:
- Research findings are INCOMPLETE and the gap affects a design decision
- A design question requires understanding code that research didn't cover
- You'd be guessing at behavior, interfaces, or constraints without reading the code

**If no gaps:** Skip to Step 6.

**If gaps exist:**
1. For each gap, formulate a specific, bounded lookup — a file to read, a pattern
   to grep, or a function signature to check. No open-ended exploration.
2. Execute the lookups using Read/Grep/Glob (or dispatch a codebase-locator agent
   for broader searches). Cap at **5 lookups total** — this is targeted, not a
   second research phase.
3. For each finding, create additional design questions in a separate
   `## Exploration-Driven Design Questions` section (numbered EDQ-1, EDQ-2, etc.).
   These follow the same format as DQ questions — options, pros/cons, recommendation.
   Each EDQ MUST cite what was found in the exploration and why it wasn't covered by
   research.

**FORBIDDEN:** Re-running research. This is surgical gap-filling, not Phase 2 redux.

### Step 6: Write draft design artifact
Write `03-design-discussion.md` to the artifact directory with ALL sections
populated and design questions marked as OPEN:

```yaml
---
phase: design-discussion
date: <today>
topic: <topic-slug>
repo: <repo>
git_sha: <HEAD>
input_artifacts: [00-ticket.md, 02-research.md]
decisions_count: <N>
exploration_decisions_count: <N or 0>
open_questions: <N total across DQ + EDQ>
status: draft
---
```

```markdown
## Summary of Changes Requested
<distilled from ticket — concise statement of what and why>

## Current State
<relevant system state from research, with file:line refs>
<only what the changes will touch, integrate with, or depend on>

## Desired End State
<concrete description of the system after changes>
<name files, APIs, behaviors>

## What We're Not Doing
<explicit scope boundaries and deferred work>

## Patterns to Follow
- **<pattern name>** — `file:line` — <brief description of when/how to use>
- ...

## Design Questions

### DQ-1: <title>
**Context:** <relevant research findings>

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| A | ... | <citing research> | ... |
| B | ... | <citing research> | ... |

**Recommendation:** <option and rationale>
**Decision:** OPEN

### DQ-2: <title>
...

## Exploration-Driven Design Questions
<only present if Step 5b found gaps — omit section entirely if not needed>

### EDQ-1: <title>
**Gap:** <what was missing from research and why it matters>
**Found:** <what targeted exploration revealed, with file:line refs>

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| A | ... | <citing exploration findings> | ... |
| B | ... | <citing exploration findings> | ... |

**Recommendation:** <option and rationale>
**Decision:** OPEN

## Constraints Discovered
<from research findings>

## Risks from Incomplete Research
<INCOMPLETE questions and their implications>
<exclude any gaps that were resolved by Step 5b exploration>
```

### Step 7: Present and resolve questions

Present a summary of the design document to the user, then ask via AskUserQuestion:

> "Design document written to `03-design-discussion.md` with N open questions
> (M from research, K from targeted exploration).
> How would you like to resolve them?
>
> 1. **Batch** — Answer all questions in one response, referenced by ID (e.g. 'DQ-1: A, DQ-3: B')
> 2. **Accept recommendations** — Use my recommendations for all open questions"

In **Batch mode**, parse the user's response and resolve each question accordingly. In **Accept recommendations mode**, resolve all OPEN questions using the stated recommendations.

For each resolved question, update the artifact:
- `**Decision:** <chosen option>`
- `**Rationale:** <from user's response or recommendation>`
- `**Implementation implication:** <one-liner about what this means for implementation>`

### Step 8: Finalize artifact
After all questions are resolved:
1. Update all `OPEN` decisions to show the chosen option
2. Update frontmatter: `open_questions: 0`, `status: complete`
3. Write the finalized `03-design-discussion.md`

## Completion

1. Present design decisions summary
2. Update `.state.json` with `current_phase: 3, completed_phases: [1, 2, 3]`
3. Instruct: "Run `/dw-04-outline <topic-slug>` in a **fresh conversation** to continue."
