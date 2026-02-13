---
name: dw-02-research
description: "Use when you have research questions from deep-work Phase 1. Objectively investigates the codebase to answer pasted questions without access to the original task description."
---

# Phase 2: Research

Objectively answer every research question by investigating the codebase.
Document what IS, not what should be. You are a documentarian, not a critic.

**Announce at start:** "Starting deep-work Phase 2: Research."

## BIAS FIREWALL — CRITICAL CONSTRAINTS

You will receive research questions pasted by the user. You MUST NOT:
- Read `01-research-questions.md` or `00-ticket.md` from the artifact directory
- Ask what the user is trying to build
- Infer or guess the user's intent
- Suggest improvements, solutions, or approaches

You ONLY answer the questions as asked.

## Setup

1. Parse `$ARGUMENTS` as `<topic-slug>`
   - If empty, ask user for topic-slug via AskUserQuestion
2. Derive repo name:
   ```bash
   basename $(git remote get-url origin 2>/dev/null | sed 's/.git$//') 2>/dev/null || basename $(pwd)
   ```
3. Set artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Pre-flight Validation

- Verify artifact directory exists → if not: "No artifact directory found. Run `/dw-research-questions <slug>` first in a separate conversation." **Stop.**
- Verify `00-ticket.md` exists in directory (confirms Phase 1 ran) → if not: "Phase 1 hasn't completed. Run `/dw-research-questions <slug>` first." **Stop.**
- **Do NOT read `00-ticket.md`** — only check existence via bash `test -f`.

## Input

If research questions were not included in `$ARGUMENTS`, ask:
"Paste the research questions from Phase 1 (everything below '## Research Questions')."

## Process

### Step 1: Parse questions
Extract numbered questions from pasted text. Identify the category of each.

### Step 2: Map questions to agents

| Category | Agent Type |
|----------|-----------|
| Subsystem understanding | codebase-analyzer |
| Code tracing | codebase-analyzer |
| Pattern discovery | codebase-pattern-finder |
| Dependency mapping | codebase-locator |
| Boundary identification | codebase-locator → codebase-analyzer |
| Constraint discovery | codebase-pattern-finder |

### Step 3: Dispatch agents
For each agent, prepend this objectivity wrapper to the task prompt:

> "You are a documentarian. Answer the following question by reading the
> codebase. Report ONLY what exists. Do not suggest improvements, critique
> patterns, or propose solutions. Include file:line references for all claims."

Dispatch independent questions in parallel.

### Step 4: Compile findings
For each question:
```
### Q<N>: <question text>
**Status:** COMPLETE | INCOMPLETE
**Sources:** <agent type(s) used>

<findings with file:line references>
```

Mark INCOMPLETE when: code can't be found, uses dynamic dispatch, or spans too
many files. For INCOMPLETE, document what WAS found and what remains ambiguous.

### Step 5: Cross-reference
Identify overlapping answers, contradictions, and cross-cutting patterns.

### Step 6: Write artifact
Write `02-research.md` to the artifact directory:
```yaml
---
phase: research
date: <today>
topic: <topic-slug>
repo: <repo>
git_sha: <HEAD>
agents_dispatched: <count>
questions_complete: <count>
questions_incomplete: <count>
input_artifacts: []
status: complete
---

## Research Findings

### Q1: <question>
**Status:** COMPLETE
**Sources:** codebase-analyzer

<detailed findings with file:line references>

...

## Summary
- <N>/<total> questions fully answered
- <M> questions incomplete (<list which and why>)

## Cross-References
- <overlaps, contradictions, patterns>
```

## Completion

1. Present findings summary, highlighting INCOMPLETE questions
2. Update `.state.json`:
   ```json
   {
     "topic": "<topic-slug>",
     "repo": "<repo>",
     "current_phase": 2,
     "completed_phases": [1, 2],
     "last_updated": "<ISO timestamp>"
   }
   ```
3. Instruct: "Research is locked in. Run `/dw-03-design-discussion <topic-slug>`
   in a **fresh conversation** to continue. The original prompt will be
   re-introduced alongside these findings."
