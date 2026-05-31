# Claude Essentials

A collection of Claude Code customizations including global instructions, skills, and utility scripts for enhanced AI-assisted development.

## Skills Overview

### Deep-Work Pipeline Orchestrator

Runs the full deep-work pipeline (Phases 1-6) in a single session using agent teams with configurable review gates between each phase.

> **Note:** The individual deep-work phase skills (`dw-01` through `dw-06`) and the RPI skills (`rpi-research`, `rpi-plan`, `rpi-implement`) have moved to the [context-engineering-workflows](https://github.com/ronsanzone/context-engineering-workflows) repo. This repo retains only the single-session orchestrator (`/deep-work-pipeline`) that coordinates the phases.

---

### Code Review & PR Skills

| Skill | Command | Purpose |
|-------|---------|---------|
| **quick-review** | `/quick-review` | Single-pass expert review with severity-ranked findings (critical → minor) |
| **local-code-review** | `/local-code-review` | Local code review of changes in the working tree |
| **submit-pr** | `/submit-pr` | Full PR submission workflow — creates draft PRs or pushes updates to existing ones |

---

### Workflow Skills

| Skill | Command | Purpose |
|-------|---------|---------|
| **refine-ticket** | `/refine-ticket` | Interactively refine a Jira ticket, pasted text, or file into a structured `ticket.md` |
| **investigate-and-fix** | `/investigate-and-fix <ticket>` | Single-session alternative to the full pipeline — investigate, research, propose, plan, and implement for well-scoped bug fixes or small features |
| **session-retrospective** | `/session-retrospective` | Analyze session process efficiency — scores context engineering, tool usage, sub-agent work, and cost efficiency (1-5) |

---

### Reporting & Presentation

| Skill | Command | Purpose |
|-------|---------|---------|
| **code-tour** | `/code-tour` | Generate an interactive HTML tour of a codebase or feature area |
| **html-report** | `/html-report` | Render compiled content as a self-contained HTML technical report with TOC, scroll-spy, and collapsible sections |

---

## Quick Start

1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/claude-essentials.git ~/code/claude-essentials
   ```

2. Create symlink to enable globally:
   ```bash
   # Backup existing config if present
   [ -d ~/.claude ] && mv ~/.claude ~/.claude.backup

   # Create symlink
   ln -s ~/code/claude-essentials/.claude ~/.claude
   ```

3. Copy and customize settings:
   ```bash
   cp ~/.claude/settings.local.json.example ~/.claude/settings.local.json
   # Edit to add your permissions
   ```

4. Verify installation:
   ```bash
   claude  # Start Claude Code
   # Check that settings load correctly
   ```

See [INSTALL.md](INSTALL.md) for detailed installation instructions.

## Directory Structure

```
claude-essentials/
├── README.md
├── INSTALL.md
├── .claude/
│   ├── CLAUDE.md                        # Global instructions
│   ├── settings.json                    # Model/plugin config
│   ├── settings.local.json.example
│   ├── docs/
│   │   └── software-design-philosophy.md
│   ├── skills/
│   │   ├── code-tour/                   # Interactive HTML codebase tours
│   │   ├── deep-work-pipeline/          # Single-session pipeline orchestrator
│   │   ├── html-report/                 # Self-contained HTML reports
│   │   ├── investigate-and-fix/         # Single-session bug fix workflow
│   │   ├── local-code-review/           # Local working tree review
│   │   ├── quick-review/                # Fast severity-ranked review
│   │   ├── refine-ticket/               # Pre-pipeline ticket refinement
│   │   ├── session-retrospective/       # Session efficiency analysis
│   │   └── submit-pr/                   # PR creation/update workflow
│   └── agents/
└── scripts/
    ├── log_analysis_lib.py
    ├── example_commands.md
    └── test_log_analysis.py
```

## Related Repositories

- **[context-engineering-workflows](https://github.com/ronsanzone/context-engineering-workflows)** — Individual deep-work phase skills (dw-01 through dw-06) and RPI workflow skills (rpi-research, rpi-plan, rpi-implement)

## Customization

### Adding New Skills
Create a directory in `.claude/skills/` with a markdown file containing frontmatter:
```yaml
---
name: my-skill
description: What this skill does
---
```

### Adding New Agents
Create a markdown file in `.claude/agents/` with frontmatter:
```yaml
---
name: my-agent
description: What this agent does
tools: Read, Grep, Glob, LS
model: sonnet
---
```

## Requirements

- Claude Code CLI
- Python 3.8+ (for log analysis scripts)
- gh CLI (for PR review skills)

## License

MIT
