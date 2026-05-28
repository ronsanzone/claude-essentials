# deprecated-skills/

Skills moved out of `.claude/skills/` so the harness no longer auto-loads them. Kept in-tree for git history and reference.

## Why

`write-plan` and `implement-plan` were a 2-phase compression of the deep-work pipeline. They worked well, but they inlined research into the plan, had no internal gates during plan creation, and had no end-of-implementation audit pass.

Replaced by the `rpi-*` skill family, which is closer to humanlayer's `research_codebase` / `create_plan` / `implement_plan` pattern while preserving our dw-05 plan format and fresh-subagent-per-task execution loop.

## Successor map

| Deprecated | Successor(s) |
|------------|--------------|
| `write-plan` | `/rpi-research` (if standalone research artifact is wanted) → `/rpi-plan` (gated plan creation, consumes research.md or runs inline research) |
| `implement-plan` | `/rpi-implement` (same per-task subagent + 2-stage review, plus a baked-in end-of-impl audit) |

The `rpi-implement` prompt files (`implementer-prompt.md`, `spec-reviewer-prompt.md`, `code-quality-reviewer-prompt.md`) are independent copies of what lived here.

## If you need these back

Move the directory back under `.claude/skills/`:

```sh
git mv deprecated-skills/write-plan .claude/skills/write-plan
git mv deprecated-skills/implement-plan .claude/skills/implement-plan
```

The skills are unchanged from when they were active — they reference `~/.claude/skills/deep-work/dw-setup.sh` which still exists, so they will still function if restored.
