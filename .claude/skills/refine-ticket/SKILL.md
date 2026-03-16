---
name: refine-ticket
description: "Use when starting work on a ticket that needs refinement before entering the deep-work pipeline. Reads from Jira, pasted text, or a file, then interactively refines it into a structured ticket.md."
---

# Refine Ticket

Interactively refine a ticket into a structured `ticket.md` file. Walks through 4 fixed categories asking targeted questions one at a time. Outputs a clean, consistent ticket ready for the deep-work pipeline.

**Announce at start:** "Starting ticket refinement."

## Hard Rules

- **One question at a time.** Never batch multiple questions. Use `AskUserQuestion` for every question.
- **Never ask directly in response text.** Always use `AskUserQuestion`.
- **Never auto-fill or guess content.** The user provides all information. You structure it.
- **Fixed categories only.** Do not add, skip, or reorder categories.
- **No technical investigation.** Do not suggest reproduction steps, environment details, investigation starting points, or code analysis. That is the pipeline's job.
- **Always write a file.** Output is `ticket.md` in the current working directory. Never output the ticket as conversation text only.

## Setup

1. Parse `$ARGUMENTS`:

   **Jira key** — If `$ARGUMENTS` matches a Jira issue key pattern (letters, hyphen, digits like `PROJ-123`):
   ```bash
   jira issue view $ARGUMENTS --raw
   ```
   Extract: summary, description, acceptance criteria, labels, priority.

   Fetch linked issues and comments for Gathered Context:
   ```bash
   jira issue view $ARGUMENTS --raw | jq -r '.fields.issuelinks[]? | select(.inwardIssue) | .inwardIssue.key + ": " + .inwardIssue.fields.summary'
   jira issue view $ARGUMENTS --raw | jq -r '.fields.issuelinks[]? | select(.outwardIssue) | .outwardIssue.key + ": " + .outwardIssue.fields.summary'
   jira issue view $ARGUMENTS --raw | jq -r '.fields.comment.comments[]? | "**" + .author.displayName + "** (" + .created[:10] + "): " + .body'
   ```
   Set `<source>` to the Jira key.

   **File path** — If `$ARGUMENTS` is a readable file path, read it. Set `<source>` to `"file"`.

   **Paste mode** — If `$ARGUMENTS` is empty or doesn't match above patterns, ask via `AskUserQuestion`:
   - Question: "Please paste the ticket content."
   - Options: "Ready to paste" (user provides content via Other/free-text)

   Set `<source>` to `"pasted"`.

2. Store raw content as `<raw-ticket>`.

3. Display:
   > **Raw ticket:**
   > <raw-ticket content>

## Process

Walk through each category **in order**. For each:
1. Show what the ticket currently says about it (extract from `<raw-ticket>`, or "Not specified in ticket.")
2. Ask **one** question via `AskUserQuestion`
3. Store the response as refined content for that category

### Category 1: Problem Statement

> **Problem Statement — here's what the ticket says:**
> <extracted or "Not specified in ticket.">

`AskUserQuestion`:
- Question: "What problem is this solving, and who is affected?"
- Options:
  - "What's shown above is sufficient"
  - "I'll provide more detail" (user provides via Other/free-text)

If "sufficient", use existing content. Otherwise use the user's response (combined with existing if relevant).

### Category 2: Scope & Boundaries

> **Scope & Boundaries — here's what the ticket says:**
> <extracted or "Not specified in ticket.">

`AskUserQuestion`:
- Question: "What is explicitly in scope and out of scope for this work?"
- Options:
  - "What's shown above is sufficient"
  - "I'll define the boundaries"

### Category 3: Acceptance Criteria

> **Acceptance Criteria — here's what the ticket says:**
> <extracted or "Not specified in ticket.">

`AskUserQuestion`:
- Question: "What are the concrete, testable conditions for this work to be considered done?"
- Options:
  - "What's shown above is sufficient"
  - "I'll specify criteria"

### Category 4: Gathered Context

For Jira input, pre-populate with linked issues and comments from Setup. For other input, start empty.

> **Gathered Context — here's what we have so far:**
> <linked issues + comments, or "No additional context found.">

`AskUserQuestion`:
- Question: "Are there related tickets, CS reports, error logs, docs, or Slack threads that should be included?"
- Options:
  - "What's shown above is sufficient"
  - "I have additional context to add"

## Assembly

After all 4 categories:

1. Assemble the refined ticket using the **exact** Output format below — including YAML frontmatter. Do not invent a different format.
2. Display the complete assembled ticket to the user
3. Ask via `AskUserQuestion`:
   - Question: "Does this refined ticket look good?"
   - Options:
     - "Looks good, save it"
     - "Needs edits" (user provides corrections via Other/free-text, apply them, re-present)

## Output

Write `ticket.md` to the current working directory using the Write tool.

```markdown
---
source: <source>
refined_date: <today YYYY-MM-DD>
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
<if any>

### Comments
<if any>

### Additional Context
<if any>

## Original Ticket
<raw-ticket>
```

Omit empty subsections under Gathered Context rather than leaving them blank.

## Completion

Print: "Refined ticket saved to `ticket.md`."
