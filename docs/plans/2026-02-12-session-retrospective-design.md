# Session Retrospective Skill — Design

## Overview

A single-file skill invoked at end-of-session that analyzes process efficiency across 5 categories, scores each 1–5, and writes a structured markdown report to disk.

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Architecture | Single-pass inline | Model already has session context; no need to extract/pass it |
| Data sources | Context + CLI metrics + log files | Maximize insight from all available data |
| Output | Saved markdown report | Historical record for trend tracking |
| Timing | End-of-session only | Full session analysis, simpler scope |
| Scoring | 1–5 numerical + qualitative | Enables trend tracking across sessions |
| Scope | Process-focused | Context engineering, tools, agents, cost — not outcome quality |

## Skill Structure

```
.claude/skills/session-retrospective/
  SKILL.md          # Main skill prompt (single file)
```

**Frontmatter:**
```yaml
---
name: session-retrospective
description: End-of-session process analysis. Scores context engineering, tool usage, sub-agents, cost, and provides improvement insights. Writes report to ~/notes/retrospectives/.
---
```

## Invocation Flow

```
Invoke /session-retrospective
  → Announce start
  → Collect CLI/log metrics via Bash
  → Scan available skills list (detect missed opportunities)
  → Analyze session context across 5 categories using embedded heuristics
  → Score each category 1–5 with qualitative narrative
  → Write report to disk
  → Print inline summary: scorecard table + top 3 takeaways
```

## Report Location

```
~/notes/retrospectives/<repo-name>/YYYY-MM-DD-HHMMSS.md
```

Timestamped to the second to avoid collisions from multiple sessions per day.

## Report Format

```markdown
# Session Retrospective — <repo> — <date>

## Session Summary
- **Goal**: [What the session set out to accomplish]
- **Duration**: [Approximate based on timestamps if available]
- **Outcome**: [Brief factual summary of what was accomplished]

## Scorecard

| Category | Score | Grade |
|----------|-------|-------|
| Context Engineering | X/5 | 🟢/🟡/🔴 |
| Tool Usage | X/5 | 🟢/🟡/🔴 |
| Sub-agent Work | X/5 | 🟢/🟡/🔴 |
| Cost Efficiency | X/5 | 🟢/🟡/🔴 |
| **Overall** | **X.X/5** | **🟢/🟡/🔴** |

## 1. Context Engineering (X/5)
### What went well
### What could improve
### Key metrics

## 2. Tool Usage (X/5)
### What went well
### What could improve
### Tool inventory table

## 3. Sub-agent Work (X/5)
### What went well
### What could improve
### Agent inventory table

## 4. Cost Efficiency (X/5)
### Metrics
### What could improve

## 5. Actionable Insights
### Prompt improvements
### Skill improvements
### Process improvements

## Top 3 Takeaways
```

## Scoring Rubric

| Score | Meaning |
|-------|---------|
| 5 | Excellent — near-optimal choices, minimal waste |
| 4 | Good — mostly effective with minor improvements possible |
| 3 | Adequate — functional but notable inefficiencies |
| 2 | Below average — significant missed opportunities |
| 1 | Poor — major process failures or waste |

**Grade mapping**: 4–5 = 🟢, 3 = 🟡, 1–2 = 🔴

## Data Collection

### From conversation context (no tool calls)
- Tool calls and results (Read, Grep, Glob, Edit, Write, Bash, Task, Skill, etc.)
- Skill invocations and which skills were loaded
- Subagent launches (Task tool with subagent_type, run_in_background, etc.)
- AskUserQuestion interactions
- Error messages, retries, and course corrections
- Initial user prompt and mid-session steering

### From Bash commands
- Cost data from Claude CLI (if available)
- Session log files from ~/.claude/
- Git activity during session
- Available skills list (for missed opportunity detection)

## Analysis Heuristics

### Context Engineering
- Multiple Read calls on the same file → duplication
- Large file reads without offset/limit → inefficiency
- Grep/Glob when Explore agent would have been better → wrong tool
- Skills loaded that weren't relevant → wasted context

### Tool Usage
- Bash for file reading instead of Read tool → anti-pattern
- Direct Grep/Glob instead of Explore agent for open-ended searches → inefficiency
- Available skills not invoked when applicable → missed opportunity

### Sub-agent Work
- Sequential tool calls that could have been parallelized → missed parallelization
- Long exploration in main context vs. delegating to Explore → context waste
- TaskOutput called on background agents → anti-pattern (per CLAUDE.md)

### Cost Efficiency
- Haiku-eligible tasks run on Opus → cost flag
- Large file reads where only a portion was needed → waste
- Duplicate information gathering → waste
