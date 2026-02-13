# Session Retrospective Skill — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a single-file skill that analyzes end-of-session process efficiency across 5 categories, scores each 1–5, and writes a structured markdown report to `~/notes/retrospectives/`.

**Architecture:** Single SKILL.md file with embedded analysis heuristics, scoring rubric, report template, and Bash data collection commands. Invoked at end-of-session via `/session-retrospective`. No subagents, no external dependencies.

**Tech Stack:** Markdown skill file, Bash for metrics collection, Write tool for report output.

**Design doc:** `docs/plans/2026-02-12-session-retrospective-design.md`

---

### Task 1: Create skill directory and SKILL.md skeleton

**Files:**
- Create: `.claude/skills/session-retrospective/SKILL.md`

**Step 1: Create the directory**

```bash
mkdir -p /Users/ron.sanzone/code/claude-essentials/.claude/skills/session-retrospective
```

**Step 2: Write the SKILL.md skeleton with frontmatter and top-level structure**

Create `.claude/skills/session-retrospective/SKILL.md` with:

```markdown
---
name: session-retrospective
description: End-of-session process analysis. Scores context engineering, tool usage, sub-agents, cost, and provides improvement insights. Writes report to ~/notes/retrospectives/.
---

# Session Retrospective

Analyze the current session's process efficiency and write a scored report.

**Announce at start:** "Starting session retrospective..."

## Setup

1. Derive repo name:
   `basename $(git remote get-url origin 2>/dev/null | sed 's/.git$//') 2>/dev/null || basename $(pwd)`
2. Set report directory: `~/notes/retrospectives/<repo>/`
3. Set report file: `~/notes/retrospectives/<repo>/YYYY-MM-DD-HHMMSS.md` (use current timestamp)

## Step 1: Collect Supplementary Metrics

## Step 2: Scan Available Skills

## Step 3: Analyze Session

## Step 4: Write Report

## Step 5: Print Inline Summary
```

**Step 3: Verify the file exists**

```bash
ls -la /Users/ron.sanzone/code/claude-essentials/.claude/skills/session-retrospective/SKILL.md
```

Expected: file exists with correct path.

**Step 4: Commit**

```bash
git add .claude/skills/session-retrospective/SKILL.md
git commit -m "feat: add session-retrospective skill skeleton"
```

---

### Task 2: Write the data collection section (Step 1 & Step 2)

**Files:**
- Modify: `.claude/skills/session-retrospective/SKILL.md`

**Step 1: Replace the `## Step 1: Collect Supplementary Metrics` section**

Replace with the full data collection instructions:

```markdown
## Step 1: Collect Supplementary Metrics

Run these Bash commands in parallel to gather external data:

**Git activity during this session:**
```bash
git log --oneline --since="2 hours ago"
```

**Session log files (if available):**
```bash
ls -la ~/.claude/projects/ 2>/dev/null | head -20
```

**Cost data (if available):**
Check for any cost/token reporting the CLI exposes. If not available, note "cost data unavailable — analyze from context only."

Record the results for use in the analysis sections below. If a command returns no data, skip that metric in the report.
```

**Step 2: Replace the `## Step 2: Scan Available Skills` section**

Replace with:

```markdown
## Step 2: Scan Available Skills

List all skills the user has installed. This is used in the Tool Usage analysis to detect skills that were available but not invoked during the session.

```bash
ls ~/.claude/skills/ 2>/dev/null
ls ~/.claude/plugins/cache/*/skills/ 2>/dev/null
```

Compare this list against the skills actually invoked during the session (visible in conversation context as Skill tool calls).
```

**Step 3: Verify formatting**

Read the file and confirm the sections are well-formed markdown with proper fencing.

**Step 4: Commit**

```bash
git add .claude/skills/session-retrospective/SKILL.md
git commit -m "feat: add data collection steps to session-retrospective"
```

---

### Task 3: Write the analysis framework (Step 3)

This is the core of the skill — the analysis instructions with embedded heuristics.

**Files:**
- Modify: `.claude/skills/session-retrospective/SKILL.md`

**Step 1: Replace the `## Step 3: Analyze Session` section**

Replace with the full analysis framework. This section instructs the model how to systematically analyze the session across 5 categories:

```markdown
## Step 3: Analyze Session

Analyze the full conversation context across these 5 categories. For each category, assign a score using the rubric below and provide specific examples from the session.

### Scoring Rubric

| Score | Meaning |
|-------|---------|
| 5 | Excellent — near-optimal choices, minimal waste |
| 4 | Good — mostly effective with minor improvements possible |
| 3 | Adequate — functional but notable inefficiencies |
| 2 | Below average — significant missed opportunities |
| 1 | Poor — major process failures or waste |

**Grade mapping**: 4–5 = 🟢, 3 = 🟡, 1–2 = 🔴

---

### Category 1: Context Engineering

Evaluate how efficiently the session managed context (the model's working memory).

**Signals of good context engineering:**
- Subagents used for discrete research/exploration tasks
- Read tool used with offset/limit for large files
- Minimal re-reading of the same files
- Skills loaded only when relevant to the task
- Parallel tool calls where possible

**Signals of poor context engineering:**
- Multiple Read calls on the same file without offset/limit
- Large file reads where only a small portion was used
- Grep/Glob used directly when an Explore agent would have preserved context
- Skills loaded that weren't relevant to the task at hand
- Redundant information gathered across multiple tool calls
- TaskOutput called on background agents (loads full transcript into context)

Score this category and cite specific examples from the session.

---

### Category 2: Tool Usage

Evaluate whether the right tools were used for each task.

**Signals of good tool usage:**
- Read tool for file reading (not Bash cat/head/tail)
- Edit tool for file editing (not Bash sed/awk)
- Explore agent for open-ended codebase questions
- Grep/Glob for targeted searches with known patterns
- Skills invoked when applicable
- Parallel tool calls for independent operations

**Signals of poor tool usage:**
- Bash used for file reading (cat, head, tail) instead of Read
- Bash used for file editing (sed, awk) instead of Edit
- Direct Grep/Glob for open-ended exploration instead of Explore agent
- Available skills not invoked when they applied (compare against Step 2 scan)
- Sequential tool calls that could have been parallelized
- Tools used that returned no useful information

Score this category. Include a tool inventory table listing each tool used, approximate count, and whether usage was effective.

---

### Category 3: Sub-agent Work

Evaluate how effectively the session used subagents/subtasks to preserve context.

**Signals of good sub-agent usage:**
- Discrete research/exploration delegated to Explore agents
- Complex multi-step tasks delegated to general-purpose agents
- Background agents used for long-running operations (builds, tests)
- Subagent results used without re-doing the same work
- Appropriate agent types chosen (Explore for search, general-purpose for implementation)

**Signals of poor sub-agent usage:**
- Long exploration sequences done in main context (could have been delegated)
- TaskOutput called on background agents (defeats context-saving purpose)
- Subagents launched but results not used effectively
- Excessive steering required (unclear prompts to subagents)
- Missed opportunities where delegation would have saved context

If no subagents were used, evaluate whether there were missed opportunities for delegation. Score accordingly.

---

### Category 4: Cost Efficiency

Evaluate token and cost efficiency of the session.

**Signals of good cost efficiency:**
- Haiku model used for simple/quick tasks (via Task tool model parameter)
- Minimal redundant tool calls
- Efficient file reads (targeted, not full files when unnecessary)
- Good prompt engineering (clear, concise instructions to subagents)

**Signals of poor cost efficiency:**
- All tasks run on Opus when some could use Haiku
- Duplicate information gathering (same data fetched multiple times)
- Large file reads where only small portions were needed
- Verbose or unclear subagent prompts leading to extra iterations
- Unnecessary tool calls that returned no useful information

Include any cost/token metrics gathered in Step 1. If unavailable, estimate relative efficiency based on tool call patterns.

---

### Category 5: Actionable Insights

This is not scored. Instead, synthesize findings from the other 4 categories into concrete recommendations:

**Prompt improvements:** How could the initial user prompt have been structured to get better results faster?

**Skill improvements:** Were any skills used that could be enhanced? What specific changes would help?

**Process improvements:** What workflow changes would improve future sessions? Be specific — not "use more subagents" but "delegate codebase exploration to Explore agents before starting implementation."
```

**Step 2: Verify the analysis section is complete**

Read the file and confirm all 5 categories are present with heuristics and scoring instructions.

**Step 3: Commit**

```bash
git add .claude/skills/session-retrospective/SKILL.md
git commit -m "feat: add analysis framework with heuristics to session-retrospective"
```

---

### Task 4: Write the report output section (Step 4 & Step 5)

**Files:**
- Modify: `.claude/skills/session-retrospective/SKILL.md`

**Step 1: Replace the `## Step 4: Write Report` section**

Replace with:

```markdown
## Step 4: Write Report

Create the report directory and write the full report:

```bash
mkdir -p ~/notes/retrospectives/<repo>/
```

Write the report to `~/notes/retrospectives/<repo>/YYYY-MM-DD-HHMMSS.md` using the Write tool with this structure:

```markdown
# Session Retrospective — <repo> — <YYYY-MM-DD HH:MM>

## Session Summary
- **Goal**: [Infer from the initial user prompt and session context]
- **Duration**: [Estimate from timestamps if available, otherwise "unknown"]
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
[Specific examples from the session]

### What could improve
[Specific examples with recommendations]

### Key metrics
- Estimated context utilization: [high/medium/low]
- Duplication detected: [yes/no, with specifics]
- Unused context loaded: [yes/no, with specifics]

## 2. Tool Usage (X/5)

### What went well
[Specific tool choices that were efficient]

### What could improve
[Tools that should have been used differently]

### Tool inventory

| Tool | Uses | Effective? | Notes |
|------|------|-----------|-------|
| [tool] | [count] | ✅/⚠️/❌ | [brief note] |

## 3. Sub-agent Work (X/5)

### What went well
[Effective delegations]

### What could improve
[Missed opportunities for delegation]

### Agent inventory

| Agent | Task | Result | Context saved? |
|-------|------|--------|---------------|
| [type] | [task] | ✅/❌ | [yes/no + estimate] |

## 4. Cost Efficiency (X/5)

### Metrics
- Total session cost: [if available, otherwise "unavailable"]
- Subagent cost: [if available, otherwise "unavailable"]
- Token usage: [if available, otherwise "unavailable"]

### What could improve
[Specific cost optimizations]

## 5. Actionable Insights

### Prompt improvements
[How the initial prompt could have been better]

### Skill improvements
[Specific skills that could be enhanced, with suggestions]

### Process improvements
[Workflow changes for future sessions]

## Top 3 Takeaways
1. [Most impactful improvement]
2. [Second most impactful]
3. [Third most impactful]
```

Populate every section with specific findings from the analysis in Step 3. Do not leave any section as a placeholder — if a category has no findings, state "No issues identified" with a brief explanation of why.
```

**Step 2: Replace the `## Step 5: Print Inline Summary` section**

Replace with:

```markdown
## Step 5: Print Inline Summary

After writing the report, print a brief summary to the conversation:

1. The scorecard table (copied from the report)
2. The Top 3 Takeaways
3. The report file path so the user knows where to find the full report

Keep the inline summary concise — the full details are in the saved report.
```

**Step 3: Verify the complete file**

Read the entire SKILL.md and verify:
- Frontmatter is valid YAML
- All 5 steps are present and complete
- No placeholder text remains
- Markdown formatting is correct (proper fencing, headers, tables)

**Step 4: Commit**

```bash
git add .claude/skills/session-retrospective/SKILL.md
git commit -m "feat: add report output and summary to session-retrospective"
```

---

### Task 5: Manual validation — invoke the skill

**Files:**
- Read: `.claude/skills/session-retrospective/SKILL.md` (verify final state)

**Step 1: Read the complete SKILL.md one final time**

Verify the full file is well-formed and all sections flow logically.

**Step 2: Verify skill is discoverable**

```bash
ls /Users/ron.sanzone/code/claude-essentials/.claude/skills/session-retrospective/
```

Expected: `SKILL.md`

**Step 3: Note for user**

The skill can be tested by invoking `/session-retrospective` at the end of any future session. The current session IS a valid test — the implementer could invoke it after completing all tasks to generate a retrospective of the implementation session itself.

**Step 4: Final commit (if any cleanup was needed)**

```bash
git add .claude/skills/session-retrospective/SKILL.md
git commit -m "feat: finalize session-retrospective skill"
```
