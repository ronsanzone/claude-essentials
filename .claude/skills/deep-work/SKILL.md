---
name: deep-work
description: "Use for complex engineering tasks where premature solutioning is dangerous. Runs a 6-phase pipeline: research-questions, research, design-discussion, outline, plan, implement — with bias isolation between prompt and research."
---

# Deep Work Pipeline

A 6-phase context engineering workflow that separates research from solutioning
through a bias firewall. The original prompt is "burned" after Phase 1 — the
research phase sees only objective questions, never the original intent.

**Announce at start:** "Starting deep-work pipeline. This will guide you through
6 phases with checkpoints between each."

## Pipeline Overview

```
Phase 1: Research Questions — decompose prompt into objective questions
Phase 2: Research — answer questions by investigating the codebase (NO prompt)
Phase 3: Design Discussion — combine research with prompt to explore options
Phase 4: Structure Outline — map decisions to concrete file changes
Phase 5: Plan — detailed implementation plan with exact code patterns
Phase 6: Implement — execute the plan
```

## Setup

1. Determine repo name:
   ```bash
   basename $(git remote get-url origin 2>/dev/null | sed 's/.git$//') 2>/dev/null || basename $(pwd)
   ```

2. Determine topic slug from `$ARGUMENTS`:
   - If argument is a file path, read the file and use its title
   - If argument is text, slugify it (lowercase, hyphens, no special chars)
   - Ask user to confirm the topic slug

3. Create artifact directory:
   ```bash
   mkdir -p ~/notes/context-engineering/<repo>/<topic-slug>
   ```

4. Write `00-ticket.md`:
   ```markdown
   ---
   phase: ticket
   date: <today>
   topic: <topic-slug>
   repo: <repo>
   git_sha: <HEAD>
   status: complete
   ---

   ## Ticket

   <user's prompt or file contents>
   ```

## Resume Support

If invoked with `--resume <topic-slug>`:
1. Read `.state.json` from the artifact directory
2. Load the last completed phase
3. Skip to the next phase
4. Present: "Resuming <topic> at Phase <N>. Last completed: Phase <N-1>."

## Phase Execution

For each phase:
1. Read the phase prompt from `phases/<phase-name>.md` (relative to this skill's directory)
2. Follow the instructions in the phase prompt
3. Write the output artifact to the artifact directory
4. Present a summary of the artifact to the user
5. Checkpoint via AskUserQuestion:
   - **Continue** — proceed to next phase
   - **Revise** — re-run current phase with user's feedback
   - **Stop** — update .state.json, exit

### Phase 1 → 2 Handoff (CRITICAL — Bias Firewall)
After Phase 1 writes `01-research-questions.md`:
- Present the questions to the user
- Instruct: "Copy the Research Questions section below and paste it when prompted.
  Edit or remove any questions you don't need. The research phase will ONLY see
  what you paste — it will not read the original prompt."
- Wait for user to paste questions before starting Phase 2

### Phase 2 → 3 Handoff (Prompt Re-introduction)
After Phase 2 writes `02-research.md`:
- Present research summary
- Note: "The original prompt will be re-introduced in Phase 3 alongside these
  research findings. The research is now locked in objectively."

## State Tracking

After each phase completes, update `.state.json`:

```json
{
  "topic": "<topic-slug>",
  "repo": "<repo>",
  "current_phase": <N>,
  "completed_phases": [1, 2, ...],
  "last_updated": "<ISO timestamp>"
}
```

## Individual Phase Invocation

Phases can also be run standalone via commands:
- `/dw-research-questions` — Phase 1
- `/dw-research` — Phase 2
- `/dw-design-discussion` — Phase 3
- `/dw-outline` — Phase 4
- `/dw-plan` — Phase 5
- `/dw-implement` — Phase 6

When run standalone, each command will ask for its required inputs.
