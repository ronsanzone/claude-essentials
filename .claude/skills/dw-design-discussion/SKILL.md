---
name: dw-design-discussion
description: "Use when deep-work Phase 2 research is complete. Combines research findings with the original task to explore design options, evaluate tradeoffs, and make decisions interactively."
---

# Phase 3: Design Discussion

Combine objective research findings with the original prompt to identify design
decisions, enumerate options, and evaluate tradeoffs. Research is locked in —
the prompt safely re-enters the pipeline here.

**Announce at start:** "Starting deep-work Phase 3: Design Discussion."

## Setup

1. Parse `$ARGUMENTS` as `<topic-slug>`
   - If empty, ask user via AskUserQuestion
2. Derive repo: `basename $(git remote get-url origin 2>/dev/null | sed 's/.git$//') 2>/dev/null || basename $(pwd)`
3. Set artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Pre-flight Validation

- `02-research.md` exists → if not: "Research not found. Complete Phases 1-2 first." **Stop.**
- `00-ticket.md` exists → if not: "No ticket found. Run `/dw-research-questions` first." **Stop.**

## Process

### Step 1: Load context
1. Read `02-research.md` completely
2. Read `00-ticket.md` completely
3. Summarize: "The user wants to [goal from ticket]. Research found [key findings]."

### Step 2: Identify design decisions
Based on the gap between "what the user wants" and "what exists," identify every
design decision. Common types:
- Where should new code live? (module/package/directory)
- What pattern should it follow? (based on existing patterns from research)
- How should it integrate with existing code? (based on boundaries from research)
- What should the API/interface look like?
- How should edge cases be handled?
- What should be tested and how?

### Step 3: Build options
For EACH decision, create 2-4 options. Every option MUST:
- Cite a specific research finding (e.g., "Research Q2 found that...")
- Include concrete pros and cons
- Be grounded in what exists in the codebase

**FORBIDDEN:** Options that ignore research findings or require uninvestigated changes.

### Step 4: Present decisions interactively
Present ONE AT A TIME via AskUserQuestion:
- Show context (relevant research findings)
- Present options table with pros/cons
- Include your recommendation with rationale
- Record user's choice as CHOSEN, others as REJECTED

### Step 5: Surface risks
Compile: constraints from research, INCOMPLETE research gaps, out-of-scope items.

### Step 6: Write artifact
Write `03-design-discussion.md` to the artifact directory:
```yaml
---
phase: design-discussion
date: <today>
topic: <topic-slug>
repo: <repo>
git_sha: <HEAD>
input_artifacts: [00-ticket.md, 02-research.md]
decisions_count: <N>
status: complete
---

## Goal
<restated from ticket, grounded in research findings>

## Design Decisions

### Decision 1: <title>
**Context:** <relevant research findings>

| Option | Description | Pros | Cons | Verdict |
|--------|-------------|------|------|---------|
| A | ... | <citing research> | ... | **CHOSEN** |
| B | ... | <citing research> | ... | REJECTED |

**Rationale:** <why chosen option was selected>

## Constraints Discovered
## Risks from Incomplete Research
## Out of Scope
```

## Completion

1. Present design decisions summary
2. Update `.state.json` with `current_phase: 3, completed_phases: [1, 2, 3]`
3. Instruct: "Run `/dw-outline <topic-slug>` in a **fresh conversation** to continue."
