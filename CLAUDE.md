<!-- Last updated: 2025-01-25 -->
# Development Standards

## Environment
- **Languages**: Go, TypeScript, Python, Java
- **Platform**: macOS / zsh
- **Editor**: NeoVim + Claude Code

## Verification Commands
All checks must pass before completion. No exceptions. Use repository specific verifications commands always. If you can't find any, ask the user for them.

## Code Philosophy
- **Simple over clever** - Minimal changes for the task at hand
- **Delete old code** - No migration layers or versioned functions
- **Clarity over abstraction** - Three similar lines beats a premature abstraction

## Workflow
- Use plan mode (`/plan`) for non-trivial changes
- Research the codebase before modifying unfamiliar areas
- Run verification after each logical change

## When Stuck
Ask: "I see approaches [A] vs [B]. Which do you prefer?"

## Progress Updates
Keep them concise:
```
Implemented X (tests passing)
Found issue with Y - investigating
```

## Language-Specific Rules
Auto-loaded from `.claude/rules/` based on file type:
- Go: `go-standards.md`
- TypeScript: `typescript-standards.md` (add when needed)
- Python: `python-standards.md` (add when needed)
