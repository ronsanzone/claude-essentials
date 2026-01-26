# Installation Guide

This guide covers installing Claude Essentials as your global Claude Code configuration.

## Prerequisites

- **Claude Code CLI** - The official Anthropic CLI for Claude
- **Python 3.8+** - For log analysis scripts
- **tmux** (optional) - For tmux-stalker skills
- **gh CLI** (optional) - For GitHub-related commands

## Installation Steps

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/claude-essentials.git ~/code/claude-essentials
```

Or if you already have it elsewhere, adjust paths accordingly.

### 2. Backup Existing Configuration

If you have an existing `~/.claude` directory:

```bash
# Check if it exists
ls -la ~/.claude

# Backup if present
[ -d ~/.claude ] && mv ~/.claude ~/.claude.backup.$(date +%Y%m%d)
```

### 3. Create Symlink

Link the `.claude` directory from this repo to your home directory:

```bash
ln -s ~/code/claude-essentials/.claude ~/.claude
```

### 4. Configure Local Settings

The `settings.local.json` file contains permissions and should be customized:

```bash
# Copy the example file
cp ~/.claude/settings.local.json.example ~/.claude/settings.local.json

# Edit to customize permissions
# Add domains you want to allow for WebFetch
# Add bash commands you want to pre-approve
```

Example permissions you might add:
```json
{
  "permissions": {
    "allow": [
      "Bash(go test:*)",
      "Bash(npm:*)",
      "WebFetch(domain:docs.example.com)",
      "mcp__gopls__go_workspace"
    ],
    "deny": []
  }
}
```

### 5. Verify Installation

Start Claude Code and verify:

```bash
claude

# Inside Claude Code, check:
# 1. Model should match settings.json
# 2. Run /skills to see available skills
# 3. Run /crux to test command loading
```

## Verification Checklist

- [ ] `~/.claude` symlink points to `~/code/claude-essentials/.claude`
- [ ] `ls -la ~/.claude` shows the correct symlink target
- [ ] Claude Code starts without errors
- [ ] Skills appear in `/skills` output
- [ ] Commands work (try `/crux` with a file argument)
- [ ] Log analysis scripts work: `python3 ~/code/claude-essentials/scripts/log_analysis_lib.py --help`

## Updating

Since this is a git repository, updates are simple:

```bash
cd ~/code/claude-essentials
git pull
```

Your local `settings.local.json` won't be affected as it's gitignored.

## Uninstallation

To remove Claude Essentials:

```bash
# Remove the symlink
rm ~/.claude

# Restore backup if you had one
[ -d ~/.claude.backup.* ] && mv ~/.claude.backup.* ~/.claude
```

## Troubleshooting

### Symlink Not Working

If Claude Code isn't picking up the configuration:

```bash
# Verify symlink is correct
ls -la ~/.claude

# Should show something like:
# .claude -> /Users/yourname/code/claude-essentials/.claude

# If it's a regular directory instead, remove and recreate
rm -rf ~/.claude  # Be careful!
ln -s ~/code/claude-essentials/.claude ~/.claude
```

### Permissions Issues

If scripts aren't executable:

```bash
chmod +x ~/code/claude-essentials/scripts/*.py
```

### Python Dependencies

The log analysis scripts use only standard library modules. No additional packages needed.

### Skills Not Appearing

Skills require specific frontmatter format. Verify files have:
```yaml
---
name: skill-name
description: Description here
---
```

## File Locations Summary

| Purpose | Location |
|---------|----------|
| Global config | `~/.claude` (symlink) |
| Actual files | `~/code/claude-essentials/.claude` |
| Scripts | `~/code/claude-essentials/scripts` |
| Local permissions | `~/.claude/settings.local.json` |

## Multiple Configurations

If you need different configurations for different contexts, you can:

1. Create branches in this repo for different setups
2. Use environment-specific `settings.local.json` files
3. Swap symlinks between different config directories

## Support

For issues with:
- **Claude Code itself**: https://github.com/anthropics/claude-code/issues
- **This configuration**: Open an issue in this repository
