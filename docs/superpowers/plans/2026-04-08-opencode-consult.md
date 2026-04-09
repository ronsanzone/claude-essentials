# OpenCode Consult Skill — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create a base building-block skill that lets Claude consult opencode (with full project read access) for a second opinion, using server mode for observability.

**Architecture:** Two files — a shell script that manages the opencode server lifecycle with read-only permissions, and a SKILL.md that documents the consultation pattern (subagent invocation, prompt framing, session management, output structure).

**Tech Stack:** Claude Code skills (markdown + YAML frontmatter), Bash, opencode CLI (`opencode web`, `opencode run`)

**Spec:** `docs/superpowers/specs/2026-04-08-opencode-consult-design.md`

---

## File Structure

| File | Action | Responsibility |
|------|--------|----------------|
| `.claude/skills/opencode-consult/ensure-server.sh` | Create | Check if opencode server is running on target port; start with read-only permissions if not; return server URL |
| `.claude/skills/opencode-consult/SKILL.md` | Create | Skill document: consultation pattern, subagent invocation, prompt framing, session management, configuration |

---

### Task 1: Create ensure-server.sh

**Files:**
- Create: `.claude/skills/opencode-consult/ensure-server.sh`

- [ ] **Step 1: Create the script**

```bash
#!/bin/bash
# ensure-server.sh — Start opencode web server if not already running.
# Usage: ensure-server.sh [port] [project-dir]
# Returns: The server URL on stdout.

PORT="${1:-4242}"
PROJECT_DIR="${2:-$(pwd)}"
URL="http://localhost:$PORT"

if curl -s --max-time 2 "$URL" >/dev/null 2>&1; then
  echo "$URL"
else
  export OPENCODE_PERMISSION='{"*":"allow","edit":"deny","external_directory":"deny"}'
  cd "$PROJECT_DIR" || exit 1
  opencode web --port "$PORT" &>/dev/null &
  sleep 2
  if curl -s --max-time 2 "$URL" >/dev/null 2>&1; then
    echo "$URL"
  else
    echo "ERROR: opencode server failed to start on port $PORT" >&2
    exit 1
  fi
fi
```

- [ ] **Step 2: Make it executable**

Run: `chmod +x .claude/skills/opencode-consult/ensure-server.sh`

- [ ] **Step 3: Verify the script parses correctly**

Run: `bash -n .claude/skills/opencode-consult/ensure-server.sh`
Expected: No output (clean parse, no syntax errors)

- [ ] **Step 4: Commit**

```bash
git add .claude/skills/opencode-consult/ensure-server.sh
git commit -m "feat: add ensure-server.sh for opencode server lifecycle"
```

---

### Task 2: Create SKILL.md

**Files:**
- Create: `.claude/skills/opencode-consult/SKILL.md`

- [ ] **Step 1: Write the complete skill document**

```markdown
---
name: opencode-consult
description: >
  Consult opencode for a second opinion on design, review, or analysis.
  Building block skill — invoke when another skill or the user requests
  consultation. Uses server mode for observability (attach via TUI or browser).
---

# OpenCode Consult

## Overview

Consult opencode for a second opinion during any task. Claude remains the decision-maker;
opencode provides an independent perspective using a different model.

Results are captured via subagent to protect the main conversation context. The opencode
server runs in persistent mode so you can observe consultations via TUI or web browser.

**This skill is a building block.** It does not decide when to consult — the calling
context (user, pipeline skill, or Claude's judgment) makes that call.

## Prerequisites

- `opencode` CLI installed and on `$PATH`
- At least one provider configured (`opencode auth list` to verify)
- `curl` available (used by server health check)

## Server Lifecycle

Before the first consultation, ensure the opencode server is running:

```bash
~/.claude/skills/opencode-consult/ensure-server.sh [port] [project-dir]
```

| Argument | Default | Description |
|----------|---------|-------------|
| `port` | `4242` | Port for the opencode web server |
| `project-dir` | `$(pwd)` | Project root — opencode's working directory |

The script checks if the server is already running (curl health check). If not, it starts
`opencode web` with read-only permissions and returns the server URL.

**Permissions applied at startup:**
- `read`, `grep`, `glob`, `list`, `bash` → **allow** (full codebase exploration, build/test commands)
- `edit` (covers `write`, `apply_patch`, `multiedit`) → **deny** (read-only consultation)
- `external_directory` → **deny** (stay within project boundary)

The server persists across consultations. Do not restart it between calls.

## Consultation Pattern

### Step 1: Spawn a subagent

Claude spawns a `general-purpose` subagent with three things:
1. Instructions to ensure the server is running
2. The consultation prompt (see Prompt Framing below)
3. Instructions to return the structured output format

### Step 2: Subagent runs the consultation

```bash
# Ensure server is running (idempotent)
SERVER_URL=$(~/.claude/skills/opencode-consult/ensure-server.sh 4242 /path/to/project)

# Run the consultation
opencode run --attach "$SERVER_URL" -m github-copilot/gpt-5.4 "Your prompt here" 2>/dev/null
```

The subagent can read project files (via its own Read/Grep/Glob tools) to include
relevant code context in the prompt string.

### Step 3: Subagent returns structured result

The subagent formats and returns:

```
## OpenCode Consultation Result
**Model:** <model used>
**Session:** <session ID, if available>
**Question:** <what was asked>

### Response
<OpenCode's response, summarized if lengthy>

### Key Points
- <bullet summary of actionable insights>
```

### Step 4: Claude synthesizes

Claude in the main context receives the summary and decides how to use it.
The consultation is advisory — Claude makes the final call.

## Prompt Framing

Frame consultation questions with context so the consulted model understands its role:

```
You are being consulted as a second opinion. Here's the context:

**Task:** <what we're working on>
**Question:** <specific thing we want your perspective on>
**Constraints:** <relevant constraints you should know>

Provide your analysis. Be direct and specific.
```

## Session Management

| Mode | How | When to Use |
|------|-----|-------------|
| New session | _(default, no extra flags)_ | Default. Clean context per consultation. |
| Continue session | Add `--continue` or `-s <session-id>` | Multiple related consultations in the same phase. OpenCode retains prior context. |
| Fork session | Add `--fork` with `--continue` or `-s` | Branch from a prior session to explore an alternative without losing the original. |

The calling context decides which mode to use. Default to new sessions unless
there's a specific reason to continue or fork.

## Configuration

| Setting | Default | Override |
|---------|---------|---------|
| Port | `4242` | First argument to `ensure-server.sh` |
| Model | `github-copilot/gpt-5.4` | `-m` flag on `opencode run` |
| Project dir | Current working directory | Second argument to `ensure-server.sh` |

## Observability

Because the skill uses `opencode web` (server mode), you can observe consultations:

- **Browser:** Open `http://localhost:4242` to see the web UI with all sessions
- **TUI:** Run `opencode attach http://localhost:4242` from any terminal
- **Session history:** All consultations persist as sessions, reviewable after the fact

## Example Invocations

**Simple one-shot consultation (subagent prompt):**

> Ensure the opencode server is running by executing `~/.claude/skills/opencode-consult/ensure-server.sh 4242 /path/to/project`.
> Then run: `opencode run --attach http://localhost:4242 -m github-copilot/gpt-5.4 "We're designing an error handling strategy for a CLI tool. The current approach uses panic/recover. Should we switch to explicit error returns? The codebase is Go, ~15k LOC." 2>/dev/null`
> Return the result using the structured output format documented in the opencode-consult skill.

**Continuing a prior session:**

> Run: `opencode run --attach http://localhost:4242 -m github-copilot/gpt-5.4 --continue "Based on your earlier analysis, what specific functions should we refactor first?" 2>/dev/null`

**Consultation with code context (subagent reads file first):**

> Read `src/api/handler.go` lines 45-90, then consult opencode:
> Run: `opencode run --attach http://localhost:4242 -m github-copilot/gpt-5.4 "Review this HTTP handler for error handling gaps: <paste code>" 2>/dev/null`
> Return the result using the structured output format.
```

- [ ] **Step 2: Commit**

```bash
git add .claude/skills/opencode-consult/SKILL.md
git commit -m "feat: add opencode-consult skill for second-opinion consultations"
```

---

### Task 3: Smoke test the skill

**Files:**
- None modified — validation only

- [ ] **Step 1: Verify skill is discoverable**

Run: `ls -la .claude/skills/opencode-consult/`
Expected: Two files — `SKILL.md` and `ensure-server.sh` (with execute bit)

- [ ] **Step 2: Verify frontmatter parses**

Run: `head -6 .claude/skills/opencode-consult/SKILL.md`
Expected: Valid YAML frontmatter with `name: opencode-consult` and `description`

- [ ] **Step 3: Verify server script syntax**

Run: `bash -n .claude/skills/opencode-consult/ensure-server.sh && echo "OK"`
Expected: `OK`

- [ ] **Step 4: Test server startup (if opencode is available)**

Run: `~/.claude/skills/opencode-consult/ensure-server.sh 4242`
Expected: `http://localhost:4242` printed to stdout. Verify with `curl -s http://localhost:4242 | head -5` that the web UI responds.

If opencode is not installed in this environment, skip this step — the user confirmed it's installed on their machine.

- [ ] **Step 5: Stop test server (cleanup)**

Run: `pkill -f "opencode web --port 4242" 2>/dev/null; echo "Cleaned up"`
Expected: `Cleaned up`
