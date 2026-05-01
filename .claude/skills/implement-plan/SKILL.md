---
name: implement-plan
description: "Use when you have a plan written by /write-plan, or any dw-05-format plan at ~/notes/context-engineering/<repo>/<slug>/plan.md, and want to execute it. Use this rather than executing-plans when subagents are available."
---

# /implement-plan

Execute a plan written by `/write-plan` (or any plan in dw-05 format living at `~/notes/context-engineering/<repo>/<slug>/plan.md`) by dispatching a fresh subagent per task, with two-stage review (spec compliance → code quality) after each.

**Core principle:** Fresh subagent per task + two-stage review = high quality, fast iteration.

**Announce at start:** "Starting /implement-plan."

## Setup

1. Run `~/.claude/skills/deep-work/dw-setup.sh "$ARGUMENTS"` and parse stdout for `REPO`, `TOPIC_SLUG`, `ARTIFACT_DIR`.
   - If the script exits 2 (`MISSING_SLUG` on stderr), use `AskUserQuestion` to ask the user for a topic slug, then re-run.

## Pre-flight Validation

- `<ARTIFACT_DIR>/plan.md` exists → if not: "Plan not found at `<path>`. Run `/write-plan <slug>` first." **Stop.**

## Tooling

Use the agent's native task tools (in Claude Code: `TaskCreate`, `TaskUpdate`, `TaskList`). Manage task dependencies via `TaskUpdate`'s `dependency` field.

## Model Selection

All Task tool dispatches (implementer, spec reviewer, code quality reviewer, session quick-review) use `model: "sonnet"`.

## Plan Structure Expectations

The plan should be in dw-05 format with clear headers for phases and tasks. The subagent-driven approach relies on dispatching a new subagent for as small a scope as possible to preserve context. Ideally each task in the plan is its own subagent loop. Avoid sending whole phases or multiple tasks to a single subagent.

The plan's `## Research Context` section is the implementer's only reference document for codebase context — there are no separate `00-ticket.md` / `02-research.md` artifacts in the lighter `/write-plan` flow. When dispatching the implementer subagent, include the relevant Research Context excerpts (typically `### Files in scope` and `### Patterns to follow`) alongside the task text.
