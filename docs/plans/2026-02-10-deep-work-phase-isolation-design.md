# Deep-Work Phase Isolation Design

**Date:** 2026-02-10
**Status:** Draft
**Branch:** feat/deep-work-pipeline

## Problem

The deep-work skill orchestrates all 6 phases in a single context window. This defeats the core design principle: each phase must run with fresh context to prevent bias contamination. The "bias firewall" between Phase 1 (research questions) and Phase 2 (research) is especially compromised — the model has already seen the original prompt and cannot truly "forget" it.

## Solution

Split the monolithic orchestrator into **6 independent skill commands** + **1 non-executing guide**. Each phase runs in a fresh conversation with its own context. Phases are connected only through filesystem artifacts and explicit user handoff.

## Architecture

### Pipeline

```
/deep-work              -> Non-executing guide (explains pipeline, shows progress)
/dw-research-questions  -> Phase 1: reads ticket, produces research questions
/dw-research            -> Phase 2: user pastes questions, produces research findings
/dw-design-discussion   -> Phase 3: reads ticket + research, produces design decisions
/dw-outline             -> Phase 4: reads research + design, produces structure outline
/dw-plan                -> Phase 5: reads research + outline, produces implementation plan
/dw-implement           -> Phase 6: reads plan, executes implementation
```

### Artifact Directory Convention

All phases use the same convention-based artifact directory:

```
~/notes/context-engineering/<repo>/<topic-slug>/
```

- `<repo>` is derived from `basename(git remote get-url origin)`
- `<topic-slug>` is passed as `$ARGUMENTS` to each command

### Artifact Files

| File | Written By | Read By |
|------|-----------|---------|
| `00-ticket.md` | Phase 1 | Phase 3 |
| `01-research-questions.md` | Phase 1 | None (user copies questions manually) |
| `02-research.md` | Phase 2 | Phases 3, 4, 5 |
| `03-design-discussion.md` | Phase 3 | Phase 4 |
| `04-structure-outline.md` | Phase 4 | Phase 5 |
| `05-plan.md` | Phase 5 | Phase 6 |
| `06-completion.md` | Phase 6 | None |
| `.state.json` | All phases | Guide |

### Bias Firewall

Phase 2 enforces true context isolation:
- Runs in a fresh conversation (no prior context)
- Does NOT read any files from the artifact directory
- User pastes research questions directly into the conversation
- Phase 2 only writes `02-research.md` — it never reads `00-ticket.md` or `01-research-questions.md`

### Pre-flight Validation

Each phase validates required artifacts before executing. If artifacts are missing, the phase warns and exits with instructions on which phase to run first.

| Phase | Required Artifacts | Error Message |
|-------|-------------------|---------------|
| 1 | None (creates setup) | — |
| 2 | `00-ticket.md` exists | "Run /dw-research-questions first" |
| 3 | `02-research.md` + `00-ticket.md` | "Run Phases 1-2 first" |
| 4 | `03-design-discussion.md` + `02-research.md` | "Run Phases 1-3 first" |
| 5 | `04-structure-outline.md` + `02-research.md` | "Run Phases 1-4 first" |
| 6 | `05-plan.md` | "Run Phases 1-5 first" |

### State Tracking

Each phase updates `.state.json` on completion:

```json
{
  "topic": "<topic-slug>",
  "repo": "<repo>",
  "current_phase": 3,
  "completed_phases": [1, 2, 3],
  "last_updated": "2026-02-10T14:30:00Z"
}
```

The guide (`/deep-work <topic-slug>`) reads this to show progress and which command to run next.

## File Structure

### New Structure

```
.claude/skills/
├── deep-work/
│   └── SKILL.md                <- Guide: explains pipeline, shows progress
├── dw-research-questions/
│   └── SKILL.md                <- Self-contained Phase 1
├── dw-research/
│   └── SKILL.md                <- Self-contained Phase 2
├── dw-design-discussion/
│   └── SKILL.md                <- Self-contained Phase 3
├── dw-outline/
│   └── SKILL.md                <- Self-contained Phase 4
├── dw-plan/
│   └── SKILL.md                <- Self-contained Phase 5
└── dw-implement/
    └── SKILL.md                <- Self-contained Phase 6
```

### Deleted

- `.claude/skills/deep-work/phases/` directory (content moves into individual SKILL.md files)
- `.claude/commands/dw-*.md` files (replaced by skills)

### Unchanged

- Artifact directory convention and file names
- `.state.json` format
- Phase prompt content (moved, not rewritten)

## Guide Skill (`/deep-work`)

The orchestrator becomes a non-executing reference:

1. Explains the 6-phase pipeline and bias firewall concept
2. If `$ARGUMENTS` is a topic-slug:
   - Reads `.state.json` from the artifact directory
   - Shows completed phases and current progress
   - Tells user which `/dw-*` command to run next in a fresh conversation
3. If no arguments: shows general pipeline documentation

## Skill Authoring Standards

All skills MUST be written using the `superpowers:writing-skills` methodology:

### TDD Cycle (RED-GREEN-REFACTOR)
1. **RED:** Run pressure scenarios WITHOUT the skill. Document baseline failures — what does the model do wrong without guidance?
2. **GREEN:** Write minimal skill addressing those specific failures. Re-run scenarios to verify compliance.
3. **REFACTOR:** Identify new rationalizations from testing. Add explicit counters. Re-test until bulletproof.

### Frontmatter Requirements
- Only `name` and `description` fields
- `description` MUST start with "Use when..." and describe triggering conditions only
- Max 1024 chars total frontmatter

### Token Budget
- Each phase skill should target <500 words (the phase prompt content is the bulk)
- Guide skill can be larger since it's reference documentation

### CSO (Claude Search Optimization)
- Description answers "Should I read this skill right now?"
- Include concrete triggers and symptoms
- No `@` syntax for force-loading other skills

## Implementation Order

Each skill should go through the full writing-skills TDD cycle individually:

1. `/deep-work` guide (simplest — no execution logic)
2. `/dw-research-questions` (Phase 1 — entry point, creates setup)
3. `/dw-research` (Phase 2 — bias firewall, most critical to get right)
4. `/dw-design-discussion` (Phase 3)
5. `/dw-outline` (Phase 4)
6. `/dw-plan` (Phase 5)
7. `/dw-implement` (Phase 6)

After all skills pass their tests: delete `phases/` directory and `commands/dw-*.md` files.

## Future Work

- **Plugin packaging:** Convert to `/deep-work:research-questions` naming via marketplace plugin structure
- **Automated orchestration:** Script or hook that launches fresh `claude` sessions per phase, passing only the topic-slug
- **Phase resumption:** If a phase was interrupted, detect partial artifacts and offer to resume
