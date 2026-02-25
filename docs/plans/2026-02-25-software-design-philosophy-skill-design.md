# Software Design Philosophy Skill — Design

## Summary

A standalone skill that evaluates code and designs against the principles from "A Philosophy of Software Design." Surfaces only the 2-3 most relevant principles and red flags contextually — no checklist fatigue.

## Skill Identity

- **Name**: `software-design-philosophy`
- **Location**: `.claude/skills/software-design-philosophy.md`
- **Format**: Single markdown file with YAML frontmatter
- **Trigger**: Code review, feature design, architecture evaluation

## Two Modes

### Review Mode
Activated when given code, a diff, or file paths.

1. Analyze against all 14 principles and 10 red flags
2. Surface only the 2-3 most relevant findings
3. For each finding: name the principle/red flag, cite the specific code location, explain why it applies, suggest a concrete improvement (if a violation)

### Design Mode
Activated when given a feature description, design doc, or architecture proposal.

1. Evaluate the design against the principles
2. Surface the 2-3 most relevant principles that should guide this design
3. For each: explain how it applies and what to watch for during implementation

## Output Format

Per finding:

```
**[Principle/Red Flag Name]** — [one-line summary of relevance]
- Where: [file:line or design component]
- Why: [1-2 sentences explaining the issue/opportunity]
- Action: [concrete suggestion]
```

No preamble, no summary section. Just the findings.

## Reference Material

The skill embeds the full principles and red flags from `software-design-philosophy.md` directly — self-contained, no runtime file reads.

## Decisions

- **Standalone** over integrated: composable, can be referenced by other skills later
- **Contextual** over checklist: surfaces only relevant items to preserve context budget
- **Embedded reference** over file reads: avoids runtime I/O, keeps skill self-contained
