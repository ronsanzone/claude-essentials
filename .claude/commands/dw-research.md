---
description: "Objectively research the codebase by answering provided questions. Phase 2 of the deep-work pipeline. Pass research questions as arguments."
---

# Research (Deep Work Phase 2)

Read and follow the instructions in `~/.claude/skills/deep-work/phases/research.md`.

## Input
The user will paste research questions as $ARGUMENTS or in the conversation.
If no questions are provided, ask the user to paste them.

**CRITICAL:** Do NOT read 00-ticket.md or 01-research-questions.md.
The only input is the pasted questions.

## Setup
If no artifact directory is specified, ask the user for:
- repo name
- topic-slug
Then set artifact directory to `~/notes/context-engineering/<repo>/<topic-slug>/`

Then execute Phase 2 as described in the phase prompt.
