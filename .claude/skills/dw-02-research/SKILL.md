---
name: dw-02-research
description: "Use when you have research questions from deep-work Phase 1. Objectively investigates the codebase to answer pasted questions without access to the original task description."
---

# Phase 2: Research

Objectively answer every research question by investigating the codebase.
Document what IS, not what should be. You are a documentarian, not a critic.

**Announce at start:** "Starting deep-work Phase 2: Research."

## BIAS FIREWALL — CRITICAL CONSTRAINTS

You MUST NOT:
- Read, open, or cat `01-research-questions.md` or `00-ticket.md` directly — ever
- Ask what the user is trying to build
- Infer or guess the user's intent
- Suggest improvements, solutions, or approaches

You MUST:
- Obtain research questions ONLY via the `extract-research-questions.sh` script (see Pre-flight Step 4)
- Treat the extracted questions as your complete and only input — do not seek additional context
- Answer ONLY the questions as written

## Setup

1. Parse `$ARGUMENTS` as `<topic-slug>`
   - If empty, ask user for topic-slug via AskUserQuestion
2. Derive repo name:
   ```bash
   basename $(git remote get-url origin 2>/dev/null | sed 's/.git$//') 2>/dev/null || basename $(pwd)
   ```
3. Set artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`
4. Locate skill script directory:
   ```bash
   SKILL_DIR="$(dirname "$(readlink -f "${BASH_SOURCE[0]:-$0}" 2>/dev/null || echo "$HOME/.claude/skills/dw-02-research")")"
   ```

## Pre-flight Validation

1. Verify artifact directory exists → if not: "No artifact directory found. Run `/dw-research-questions <slug>` first in a separate conversation." **Stop.**
2. Verify `00-ticket.md` exists in directory (confirms Phase 1 ran) → if not: "Phase 1 hasn't completed. Run `/dw-research-questions <slug>` first." **Stop.**
3. **Do NOT read `00-ticket.md`** — only check existence via bash `test -f`.
4. **Extract research questions** — run the co-located extraction script:
   ```bash
   ~/.claude/skills/dw-02-research/extract-research-questions.sh <repo> <topic-slug>
   ```
   This script outputs ONLY the `## Research Questions` section. It never exposes the original prompt.
   - If it exits non-zero → display its stderr message and **Stop.**

## Input

Use the output of `extract-research-questions.sh` from Pre-flight Step 4 as your input.

If the user pastes additional or edited questions, use those **instead** — user-provided questions always take precedence over the script output.

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

### Step 6: Build structured summary
Synthesize findings into tagged categories. This summary is based ONLY on what
the questions asked about — you have no task context, so do not filter by
"relevance." Summarize everything you investigated.

**System State (as investigated)**
- Distill the current state of the systems/components that the questions covered.
  Include file:line references. Focus on factual descriptions of what exists.

**Patterns Found**
- Extract concrete code patterns discovered while answering questions.
  Format: `<pattern name> — file:line — <brief description>`
- Include: naming conventions, architectural patterns, integration patterns,
  test patterns — anything structural that was observed.

**Constraints & Invariants**
- Document any constraints, invariants, or rules enforced by the investigated
  code. Include: validation rules, type constraints, test assertions, config
  requirements.

### Step 7: Write artifact
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
input_artifacts: [01-research-questions.md (questions section only, via extract script)]
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

## Structured Summary

### System State (as investigated)
- <factual description of current system state for investigated components>
- <include file:line references>

### Patterns Found
- <pattern name> — `file:line` — <brief description>
- ...

### Constraints & Invariants
- <constraint> — <source/file:line>
- ...
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
