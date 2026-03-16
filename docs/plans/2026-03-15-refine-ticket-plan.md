# Refine Ticket Skill Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build an interactive skill that refines Jira tickets (or pasted text) into structured `ticket.md` files ready for the deep-work pipeline.

**Architecture:** Single SKILL.md file following project conventions (frontmatter, announce, setup, process, completion). Uses AskUserQuestion for interactive Q&A. Invokes jira-cli via Bash for Jira ingestion. Writes output via Write tool.

**Tech Stack:** Markdown skill file, jira CLI, AskUserQuestion tool

---

### Task 1: Create skill directory and SKILL.md with frontmatter + input parsing

**Files:**
- Create: `.claude/skills/refine-ticket/SKILL.md`

**Step 1: Create the skill file with frontmatter, announcement, and setup section**

Write `.claude/skills/refine-ticket/SKILL.md` with this exact content:

````markdown
---
name: refine-ticket
description: "Use when starting work on a ticket that needs refinement before entering the deep-work pipeline. Reads from Jira, pasted text, or a file, then interactively refines it into a structured ticket.md."
---

# Refine Ticket

Interactively refine a ticket into a structured `ticket.md` file. Walks through fixed categories (Problem Statement, Scope & Boundaries, Acceptance Criteria, Gathered Context) asking targeted questions to fill gaps. Outputs a clean ticket ready for the deep-work pipeline.

**Announce at start:** "Starting ticket refinement."

## Setup

1. Parse `$ARGUMENTS`:
   - **Jira key** — If `$ARGUMENTS` matches a Jira issue key pattern (e.g., `PROJ-123`, `ABC-1`), fetch the ticket:
     ```bash
     jira issue view $ARGUMENTS --raw
     ```
     Extract and store: summary, description, acceptance criteria, labels, priority, linked issues, comments.
     Also fetch linked issues for Gathered Context pre-population:
     ```bash
     jira issue view $ARGUMENTS --raw | jq -r '.fields.issuelinks[]? | select(.inwardIssue) | .inwardIssue.key + ": " + .inwardIssue.fields.summary'
     jira issue view $ARGUMENTS --raw | jq -r '.fields.issuelinks[]? | select(.outwardIssue) | .outwardIssue.key + ": " + .outwardIssue.fields.summary'
     jira issue view $ARGUMENTS --raw | jq -r '.fields.comment.comments[]? | "**" + .author.displayName + "** (" + .created[:10] + "): " + .body'
     ```
     Set `<source>` to the Jira key (e.g., `PROJ-123`).

   - **File path** — If `$ARGUMENTS` is a readable file path, read it as the raw ticket content. Set `<source>` to `"file"`.

   - **Paste mode** — If `$ARGUMENTS` is empty or doesn't match a Jira key or file path, ask the user to paste the ticket content using `AskUserQuestion` with a single option "Ready to paste" and let them use the "Other" free-text input. Set `<source>` to `"pasted"`.

2. Store the raw ticket content as `<raw-ticket>` for later inclusion in the output.

3. Display the raw ticket to the user:
   > **Raw ticket:**
   > <raw-ticket content>

## Process

Walk through each category in order. For each category:
- Show what the ticket currently says about it (extract from `<raw-ticket>`, or state "Not specified in ticket.")
- Ask **one** targeted question using `AskUserQuestion` to fill the gap
- Store the user's response as the refined content for that category

### Category 1: Problem Statement

Show any existing problem/summary context from the ticket, then ask:

> **Problem Statement — here's what the ticket says:**
> <extracted content or "Not specified in ticket.">

Ask via `AskUserQuestion`:
- Question: "What problem is this solving, and who is affected?"
- Options:
  - "What's shown above is sufficient" — use the existing ticket content as-is
  - "I'll provide more detail" — user provides additional context via Other/free-text
- If the user selects "What's shown above is sufficient", use the existing content. Otherwise, combine the existing content with the user's additions.

### Category 2: Scope & Boundaries

Show any existing scope information from the ticket, then ask:

> **Scope & Boundaries — here's what the ticket says:**
> <extracted content or "Not specified in ticket.">

Ask via `AskUserQuestion`:
- Question: "What is explicitly in scope and out of scope for this work?"
- Options:
  - "What's shown above is sufficient"
  - "I'll define the boundaries"

### Category 3: Acceptance Criteria

Show any existing acceptance criteria from the ticket, then ask:

> **Acceptance Criteria — here's what the ticket says:**
> <extracted content or "Not specified in ticket.">

Ask via `AskUserQuestion`:
- Question: "What are the concrete, testable conditions for this work to be considered done?"
- Options:
  - "What's shown above is sufficient"
  - "I'll specify criteria"

### Category 4: Gathered Context

For Jira input, pre-populate this section with:
- Linked issues (fetched in Setup)
- Comments from the Jira ticket (fetched in Setup)

Show the pre-populated context (or "No additional context found." for non-Jira input), then ask:

> **Gathered Context — here's what we have so far:**
> <linked issues, comments, or "No additional context found.">

Ask via `AskUserQuestion`:
- Question: "Are there related tickets, CS reports, error logs, docs, or Slack threads that should be included?"
- Options:
  - "What's shown above is sufficient"
  - "I have additional context to add"

## Assembly

After all categories are complete:

1. Assemble the refined ticket in the output format (see below)
2. Display the complete refined ticket to the user for review
3. Ask via `AskUserQuestion`:
   - Question: "Does this refined ticket look good?"
   - Options:
     - "Looks good, save it" — proceed to write
     - "Needs edits" — user provides corrections, apply them, then re-present

## Output

Write `ticket.md` to the current working directory using the Write tool.

Format:

```markdown
---
source: <source>
refined_date: <today's date YYYY-MM-DD>
---

## Problem Statement
<refined content>

## Scope & Boundaries
<refined content>

## Acceptance Criteria
- [ ] <criterion 1>
- [ ] <criterion 2>
...

## Gathered Context
### Related Tickets
<linked issues if any, otherwise omit subsection>

### Comments
<jira comments if any, otherwise omit subsection>

### Additional Context
<user-provided context if any, otherwise omit subsection>

## Original Ticket
<raw-ticket content>
```

## Completion

Print: "Refined ticket saved to `ticket.md`."
````

**Step 2: Verify the file was created**

Run: `cat .claude/skills/refine-ticket/SKILL.md | head -5`
Expected: frontmatter lines with `name: refine-ticket`

**Step 3: Commit**

```bash
git add .claude/skills/refine-ticket/SKILL.md
git commit -m "feat: add refine-ticket skill for pre-pipeline ticket refinement"
```

---

### Task 2: Manual test — paste mode

**Step 1: Invoke the skill in paste mode**

Run: `/refine-ticket`

**Step 2: Verify behavior**

- Skill should announce "Starting ticket refinement."
- Skill should prompt to paste ticket content
- Walk through all 4 categories with AskUserQuestion
- Present assembled ticket for review
- Write `ticket.md` to CWD

**Step 3: Inspect output**

Read `ticket.md` and verify:
- Frontmatter has `source: pasted` and `refined_date`
- All 4 sections present
- Original ticket preserved at bottom

---

### Task 3: Manual test — Jira mode

**Step 1: Invoke with a real Jira key**

Run: `/refine-ticket <REAL-JIRA-KEY>` (use an actual ticket from your project)

**Step 2: Verify Jira-specific behavior**

- Ticket fetched via `jira issue view`
- Linked issues and comments pre-populated in Gathered Context
- Frontmatter `source` matches the Jira key

**Step 3: Commit test fix if needed**

If any adjustments are needed based on testing, fix and commit:
```bash
git add .claude/skills/refine-ticket/SKILL.md
git commit -m "fix: adjust refine-ticket skill based on manual testing"
```
