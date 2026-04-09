# OpenCode Consult Skill — Design Spec

**Date:** 2026-04-08
**Status:** Draft
**Author:** Ron Sanzone + Claude

## Summary

A base building-block skill that lets Claude consult opencode for a second opinion via its server mode. Claude remains the decision-maker; opencode provides a second perspective using a configurable model (default: `github-copilot/gpt-5.4`). Results are captured via subagent to protect the main conversation context.

The skill is generic and unopinionated — it does not dictate when to consult. Calling contexts (users, pipeline skills, wrapper skills) decide when consultation adds value.

## Goals

1. Enable Claude to get a second opinion from opencode during any task
2. Use opencode's server mode for observability — the user can attach via TUI (`opencode attach`) or web browser to watch/interact with consultations
3. Capture consultation results via subagent to keep the main context clean
4. Provide a stable interface that future deep-work pipeline phases can wrap

## Non-Goals

- Proactive auto-consultation (no decision framework baked in)
- Pipeline-specific integration (comes later as wrapper skills)
- Model comparison workflows
- Automated decision-making based on opencode's response

## Architecture

### Three Layers

```
┌─────────────────────────────────────┐
│  Claude (main context)              │
│  - Decides to consult               │
│  - Formulates question              │
│  - Receives structured summary      │
│  - Synthesizes with own analysis    │
└──────────────┬──────────────────────┘
               │ spawns subagent
┌──────────────▼──────────────────────┐
│  Subagent                           │
│  - Ensures server is running        │
│  - Runs opencode run --attach       │
│  - Captures and structures output   │
│  - Returns summary to main context  │
└──────────────┬──────────────────────┘
               │ opencode run --attach
┌──────────────▼──────────────────────┐
│  OpenCode Server (persistent)       │
│  - Runs configured model             │
│  - Sessions persist for observation │
│  - User can attach via TUI/web      │
└─────────────────────────────────────┘
```

### Layer 1: Server Lifecycle

An `ensure-server.sh` helper script handles server lifecycle:

1. Check if `http://localhost:<port>` responds (simple curl with 2s timeout)
2. If running, return the URL
3. If not, `cd` to the project directory, set `OPENCODE_PERMISSION` for read-only access, start `opencode web --port <port>` in background, wait briefly, return the URL

Default port: `4242`, default project dir: `$(pwd)` (both overridable via arguments).

Using `opencode web` (not `serve`) so the user gets both HTTP API and a web interface for observation. The server starts in the project directory so opencode's file tools resolve paths correctly, and permissions are locked to read-only consultation (see **Permissions** section).

The server persists across consultations within a session. No complex health checks or PID management.

### Layer 2: Consultation Invocation

Each consultation follows this flow:

1. **Claude (main context)** formulates the consultation question
2. **Claude spawns a `general-purpose` subagent** with:
   - Instructions to run `ensure-server.sh` 
   - The consultation prompt
   - Instructions on output format
3. **Subagent runs:**
   ```bash
   ~/.claude/skills/opencode-consult/ensure-server.sh <port>
   opencode run --attach http://localhost:<port> -m <model> "<prompt>" 2>/dev/null
   ```
4. **Subagent structures** the response and returns it
5. **Claude (main context)** receives the structured summary

OpenCode itself has full read access to the project — its built-in `read`, `grep`, `glob`, `list`, and `bash` tools are all set to `allow` so the consulted model can explore the codebase without permission prompts. File modifications (`edit`) are set to `deny` since this is a read-only consultation. See the **Permissions** section below for the full configuration.

### Layer 3: Output Handling

The subagent returns a structured response to the main context:

```markdown
## OpenCode Consultation Result
**Model:** <model used>
**Session:** <session ID, if available>
**Question:** <what was asked>

### Response
<OpenCode's response, summarized if lengthy>

### Key Points
- <bullet summary of actionable insights>
```

## Consultation Protocol

### Invocation Patterns

The skill is invoked explicitly, not proactively:

1. **User requests:** "Get a second opinion from opencode on this"
2. **Skill references:** Another skill includes a step that invokes `opencode-consult`
3. **Claude's judgment:** When the calling context suggests consultation adds value

### Prompt Framing Template

The subagent frames questions to opencode with context:

```
You are being consulted as a second opinion. Here's the context:

**Task:** <what we're working on>
**Question:** <specific thing we want your perspective on>
**Constraints:** <relevant constraints you should know>

Provide your analysis. Be direct and specific.
```

### Session Management

Three modes, documented for the caller to choose:

| Mode | Flag | When to Use |
|------|------|-------------|
| New session | _(default)_ | Default. Clean context per consultation. |
| Continue session | `--continue` or `-s <id>` | Multiple related consultations in the same phase. OpenCode retains prior context. |
| Fork session | `--fork` | Branch from a prior session to explore an alternative. |

## File Structure

```
.claude/skills/opencode-consult/
├── SKILL.md              # Main skill document
└── ensure-server.sh      # Helper: check port, start server if needed
```

### SKILL.md Sections

1. **Overview** — consultation model, skill purpose
2. **Prerequisites** — opencode installed, provider configured
3. **Server Lifecycle** — run `ensure-server.sh`, reuse across session
4. **Consultation Pattern** — subagent invocation, prompt framing, output structure
5. **Session Management** — new vs. continue vs. fork
6. **Configuration** — default model, default port (overridable)
7. **Example Invocations** — concrete examples for common consultation scenarios

### ensure-server.sh

```bash
#!/bin/bash
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
  echo "$URL"
fi
```

## Permissions

OpenCode must run with read-only project access and no permission prompts (non-interactive). The `ensure-server.sh` script sets this via the `OPENCODE_PERMISSION` env var:

```json
{
  "*": "allow",
  "edit": "deny",
  "external_directory": "deny"
}
```

**What this gives the consulted model:**
- `read`, `grep`, `glob`, `list` — full codebase exploration, no prompts
- `bash` — can run build/test commands to understand the project (allowed by `*`)
- `edit` (covers `write`, `apply_patch`, `multiedit`) — **denied**. Consultation is read-only.
- `external_directory` — **denied**. Consulted model stays within the project boundary.

**Working directory:** The server must be started from the project root so opencode's file tools resolve paths correctly. The `ensure-server.sh` script accepts the project directory as an argument and `cd`s before starting.

**Bash scoping (optional future tightening):** If `bash: allow` is too permissive, it can be narrowed:

```json
{
  "bash": {
    "*": "ask",
    "git *": "allow",
    "make *": "allow",
    "go *": "allow",
    "npm *": "allow",
    "bazel *": "allow"
  }
}
```

For v1, `bash: allow` is acceptable since `edit: deny` prevents file mutation and the consultation is advisory.

## Configuration Defaults

| Setting | Default | Override |
|---------|---------|---------|
| Port | `4242` | Argument to `ensure-server.sh` or skill parameter |
| Model | `github-copilot/gpt-5.4` | `-m` flag in consultation command |
| Project dir | Current working directory | Argument to `ensure-server.sh` |

## Observability

Because the skill uses `opencode web` (server mode):

- **Browser:** Navigate to `http://localhost:4242` to see the web UI with all sessions
- **TUI:** Run `opencode attach http://localhost:4242` from any terminal to connect the interactive TUI
- **Session history:** All consultations are persisted as sessions, reviewable after the fact

This gives the user full visibility into what Claude asked and what opencode responded.

## Future: Pipeline Integration

The base skill is designed as a building block. Future wrapper skills can invoke it for specific pipeline phases:

- **dw-02/dw-03 (Research/Design):** Consult opencode for alternative design perspectives
- **dw-05b (Plan Review):** Use opencode as an adversarial reviewer of implementation plans
- **dw-06 checkpoints:** Consult opencode for code review at implementation review gates

Each wrapper would import the consultation pattern and add phase-specific prompt framing and session strategy.

## Risks and Mitigations

| Risk | Mitigation |
|------|-----------|
| OpenCode server fails to start | Subagent reports failure; Claude proceeds without consultation |
| OpenCode returns low-quality response | Claude is the decision-maker; it evaluates and can discard |
| Server port conflict | Port is configurable; document how to change |
| Cold-start delay on first consultation | `opencode web` starts once, reused. 2s initial wait is acceptable |
