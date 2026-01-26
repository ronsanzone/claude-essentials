# Claude Essentials

A collection of Claude Code customizations including global instructions, agents, skills, commands, and utility scripts for enhanced AI-assisted development.

## What's Included

### Core Configuration
- **CLAUDE.md** - Global development instructions and workflow guidelines
- **settings.json** - Model configuration and plugin settings
- **settings.local.json.example** - Template for local permissions

### Rules
- **go-standards.md** - Go development standards with automated enforcement

### Skills (3)
- **tmux-stalker** - Read content from tmux panes for context gathering
- **tmux-stalker-summarized** - Summarize tmux content efficiently
- **log-analysis** - Direct log analysis tools for Splunk JSON exports

### Agents (5)
- **codebase-analyzer** - Analyze implementation details with file:line references
- **codebase-locator** - Find files and components by feature/topic
- **codebase-pattern-finder** - Find similar implementations and patterns
- **splunk-analyzer** - Analyze Splunk JSON logs for patterns and errors
- **web-search-researcher** - Deep web research for technical topics

### Commands (2)
- **crux** - Detailed code review for local files
- **crux_gh** - Code review for GitHub pull requests

### Scripts
- **log_analysis_lib.py** - Python library for log analysis
- **example_commands.md** - Usage examples for log analysis

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
│   ├── CLAUDE.md              # Global instructions
│   ├── settings.json          # Model/plugin config
│   ├── settings.local.json    # Permissions (gitignored)
│   ├── settings.local.json.example
│   ├── rules/
│   │   └── go-standards.md
│   ├── skills/
│   │   ├── tmux-stalker/
│   │   │   └── SKILL.md
│   │   ├── tmux-stalker-summarized/
│   │   │   └── skill.md
│   │   └── log-analysis.md
│   ├── agents/
│   │   ├── codebase-analyzer.md
│   │   ├── codebase-locator.md
│   │   ├── codebase-pattern-finder.md
│   │   ├── splunk-analyzer.md
│   │   └── web-search-researcher.md
│   └── commands/
│       ├── crux.md
│       └── crux_gh.md
└── scripts/
    ├── log_analysis_lib.py
    ├── example_commands.md
    └── test_log_analysis.py
```

## Customization

### Adding New Rules
Create a markdown file in `.claude/rules/` with frontmatter specifying file globs:
```yaml
---
globs:
  - "**/*.py"
---
# Python Standards
...
```

### Adding New Skills
Create a markdown file in `.claude/skills/` with frontmatter:
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
- tmux (for tmux-stalker skills)
- gh CLI (for GitHub-related commands)

## License

MIT
