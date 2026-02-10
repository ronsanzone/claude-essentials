# Phase 3: Design Discussion

## Purpose
Combine objective research findings with the original prompt to identify design
decisions, enumerate options, and evaluate tradeoffs. The prompt safely re-enters
the pipeline here — research is already locked in.

## Inputs
- `02-research.md` from the artifact directory
- `00-ticket.md` from the artifact directory (prompt re-introduced)
- Artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Process

### Step 1: Load context
1. Read `02-research.md` completely
2. Read `00-ticket.md` completely
3. Summarize: "The user wants to [goal from ticket]. Research found [key findings]."

### Step 2: Identify design decisions
Based on the gap between "what the user wants" and "what exists," identify every
design decision that needs to be made. Common decision types:
- Where should new code live? (module/package/directory)
- What pattern should it follow? (based on existing patterns from research)
- How should it integrate with existing code? (based on boundaries from research)
- What should the API/interface look like?
- How should edge cases be handled?
- What should be tested and how?

### Step 3: Build options for each decision
For EACH decision, create 2-4 options. Every option MUST:
- Cite a specific research finding (e.g., "Research Q2 found that...")
- Include concrete pros and cons
- Be grounded in what exists in the codebase

**FORBIDDEN:**
- Options that ignore research findings
- Options that require changes the research didn't investigate
- "We could..." without citing evidence from the research

### Step 4: Present decisions interactively
Present decisions ONE AT A TIME using AskUserQuestion:
- Show the context (which research findings are relevant)
- Present options as a table with pros/cons
- Include your recommendation with rationale
- Record the user's choice as CHOSEN, mark others as REJECTED

### Step 5: Surface risks
After all decisions, compile:
- Constraints discovered from research that affect the implementation
- INCOMPLETE research findings that create uncertainty
- Out of scope items (things research revealed we should NOT touch)

### Step 6: Write artifact

**Output file:** `03-design-discussion.md`

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

### Decision 1: <decision title>
**Context:** <which research findings are relevant>

| Option | Description | Pros | Cons | Verdict |
|--------|-------------|------|------|---------|
| A | <description> | <pros citing research> | <cons> | **CHOSEN** |
| B | <description> | <pros citing research> | <cons> | REJECTED |

**Rationale:** <why the chosen option was selected>

### Decision 2: ...

## Constraints Discovered
- <constraint from research with Q# reference>

## Risks from Incomplete Research
- <Q# was INCOMPLETE — this affects Decision N because...>

## Out of Scope
- <things we explicitly decided NOT to do>
```
