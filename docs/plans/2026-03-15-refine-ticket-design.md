# Refine Ticket Skill — Design Document

**Date:** 2026-03-15
**Status:** Approved

## Problem

Deep-work pipeline quality degrades when input tickets are vague, missing acceptance criteria, lacking scope boundaries, or missing scattered context (CS tickets, related Jira issues, error logs, docs). Phase 1 (research questions) writes `00-ticket.md` from raw input — garbage in, garbage out.

A pre-pipeline refinement step draws out human knowledge that can't be derived from code, consolidates scattered references, and produces a structured ticket ready for the deep-work pipeline.

## Skill Overview

- **Name:** `refine-ticket`
- **Location:** `.claude/skills/refine-ticket/SKILL.md`
- **Invocation:** `/refine-ticket <JIRA-KEY>` or `/refine-ticket` (paste/file mode)
- **Output:** `ticket.md` in the current working directory
- **Approach:** Category Walkthrough — walk through fixed categories, ask one question per gap

## Input Modes

1. **Jira key** — `$ARGUMENTS` matches pattern like `PROJ-123`. Fetch via `jira issue view --raw` and extract: summary, description, acceptance criteria, labels, priority, linked issues, comments.
2. **Paste mode** — If `$ARGUMENTS` is empty or doesn't match a Jira key/file path, prompt the user to paste ticket content.
3. **File path** — If `$ARGUMENTS` is a readable file path, read it as raw ticket content.

## Refinement Categories

| # | Category | Purpose |
|---|----------|---------|
| 1 | Problem Statement | Why does this work need to happen? Who's affected? Business/customer impact. |
| 2 | Scope & Boundaries | What's in vs. out. Prevents scope creep in research phase. |
| 3 | Acceptance Criteria | Concrete, testable conditions for "done". |
| 4 | Gathered Context | Consolidate scattered references — CS tickets, related Jira issues, error logs, docs, Slack threads — into a single summarized section. |

For Jira input, linked issues and comments are auto-fetched to pre-populate "Gathered Context".

## Skill Flow

```
/refine-ticket PROJ-123
        │
        ▼
┌─ Ingest ────────────────────┐
│ Jira key → fetch via jira   │
│ Paste → prompt user         │
│ File path → read file       │
└──────────┬──────────────────┘
           │
           ▼
┌─ Display raw ticket ────────┐
│ Show what we're working with│
└──────────┬──────────────────┘
           │
     ┌─────┴─────┐
     │  For each  │
     │  category  │
     └─────┬─────┘
           │
           ▼
┌─ Show current state ────────┐
│ "Here's what the ticket     │
│  says about [category]:"    │
│  → existing content or      │
│    "not specified"           │
└──────────┬──────────────────┘
           │
           ▼
┌─ Ask one question ──────────┐
│ Targeted question to fill   │
│ the gap. User responds or   │
│ confirms existing is fine.  │
└──────────┬──────────────────┘
           │
           ▼
     (next category)
           │
           ▼
┌─ Assemble & present ────────┐
│ Show complete refined ticket│
│ User confirms or adjusts    │
└──────────┬──────────────────┘
           │
           ▼
┌─ Write ticket.md ───────────┐
│ Save to working directory   │
└─────────────────────────────┘
```

## Output Format

```markdown
---
source: PROJ-123          # or "pasted" / "file"
refined_date: 2026-03-15
---

## Problem Statement
<refined content>

## Scope & Boundaries
<refined content>

## Acceptance Criteria
- [ ] <criterion 1>
- [ ] <criterion 2>

## Gathered Context
### Related Tickets
- PROJ-100: <summary>
- CS-456: <summary>

### Error References
<summarized errors/logs if any>

### Additional Context
<docs, slack threads, other references>

## Original Ticket
<raw ticket content preserved for reference>
```

## Key Design Decisions

1. **Fixed categories over dynamic analysis** — Predictable, consistent output. The user decides what's sufficient, not Claude.
2. **Interactive Q&A over auto-enrichment** — The goal is drawing out human/domain knowledge that isn't in the codebase.
3. **Standalone output (no pipeline coupling)** — Writes `ticket.md` to CWD. User decides when/how to feed it into deep-work.
4. **Categories focus on human knowledge** — Technical context, constraints, and edge cases are left to the research phase. This skill captures business context, scope decisions, and scattered references.
5. **Preserves original ticket** — Raw content kept at the bottom for reference and auditability.

## Non-Goals

- Does not replace Phase 1 (research questions) or any pipeline phase
- Does not perform codebase analysis — that's the pipeline's job
- Does not auto-start the deep-work pipeline
