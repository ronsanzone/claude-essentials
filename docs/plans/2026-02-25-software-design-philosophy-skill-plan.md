# Software Design Philosophy Skill — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a skill that evaluates code and designs against "A Philosophy of Software Design" principles, surfacing only the 2-3 most relevant findings contextually.

**Architecture:** Single markdown skill file with YAML frontmatter, embedded reference material, and two-mode (review/design) instructions. No scripts, no dependencies.

**Tech Stack:** Markdown skill file following project conventions.

---

### Task 1: Create the skill file

**Files:**
- Create: `.claude/skills/software-design-philosophy.md`

**Step 1: Write the skill file**

Create `.claude/skills/software-design-philosophy.md` with the following exact content:

```markdown
---
name: software-design-philosophy
description: Use when reviewing code, designing features, or evaluating architecture decisions. Surfaces the most relevant design principles and red flags from "A Philosophy of Software Design" contextually.
---

# Software Design Philosophy

Evaluate code or designs against the principles from "A Philosophy of Software Design." Surface only the 2-3 most relevant principles and red flags — no checklist, no boilerplate.

## Mode Detection

Determine mode from the input:

- **Review mode**: You have code, a diff, file paths, or a PR to evaluate.
- **Design mode**: You have a feature description, design doc, or architecture proposal to evaluate.

## Review Mode

1. Read/analyze the code or diff provided
2. Evaluate against all principles and red flags below
3. Select the 2-3 most relevant findings (violations or exemplary adherence worth noting)
4. Output each finding using the format below

## Design Mode

1. Read/analyze the design or feature description provided
2. Evaluate against all principles below
3. Select the 2-3 most relevant principles that should guide this design
4. Output each finding using the format below, with "Action" describing what to watch for during implementation

## Output Format

For each finding:

**[Principle or Red Flag Name]** — [one-line summary of relevance]
- **Where:** [file:line, function name, or design component]
- **Why:** [1-2 sentences explaining the issue or opportunity]
- **Action:** [concrete suggestion or implementation guidance]

No preamble. No summary. Just the findings.

---

## Reference: Design Principles

1. Complexity is incremental: you have to sweat the small stuff.
2. Working code isn't enough.
3. Make continual small investments to improve system design.
4. Modules should be deep.
5. Interfaces should be designed to make the most common usage as simple as possible.
6. It's more important for a module to have a simple interface than a simple implementation.
7. General-purpose modules are deeper.
8. Separate general-purpose and special-purpose code.
9. Different layers should have different abstractions.
10. Pull complexity downward.
11. Define errors (and special cases) out of existence.
12. Design it twice.
13. Comments should describe things that are not obvious from the code.
14. Software should be designed for ease of reading, not ease of writing.
15. The increments of software development should be abstractions, not features.

## Reference: Red Flags

- **Shallow Module**: the interface for a class or method isn't much simpler than its implementation.
- **Information Leakage**: a design decision is reflected in multiple modules.
- **Temporal Decomposition**: the code structure is based on the order in which operations are executed, not on information hiding.
- **Overexposure**: An API forces callers to be aware of rarely used features in order to use commonly used features.
- **Pass-Through Method**: a method does almost nothing except pass its arguments to another method with a similar signature.
- **Repetition**: a nontrivial piece of code is repeated over and over.
- **Special-General Mixture**: special-purpose code is not cleanly separated from general-purpose code.
- **Conjoined Methods**: two methods have so many dependencies that it's hard to understand the implementation of one without understanding the implementation of the other.
- **Comment Repeats Code**: all of the information in a comment is immediately obvious from the code next to the comment.
- **Implementation Documentation Contaminates Interface**: an interface comment describes implementation details not needed by users of the thing being documented.
- **Vague Name**: the name of a variable or method is so imprecise that it doesn't convey much useful information.
```

**Step 2: Verify the skill file**

Run: `cat .claude/skills/software-design-philosophy.md | head -5`
Expected: YAML frontmatter with `name: software-design-philosophy`

**Step 3: Commit**

```bash
git add .claude/skills/software-design-philosophy.md
git commit -m "feat: add software-design-philosophy skill"
```

---

### Task 2: Verify skill is discoverable

**Step 1: Check skill appears in skills listing**

Run: `ls .claude/skills/`
Expected: `software-design-philosophy.md` appears in the listing alongside other skills.

**Step 2: Verify frontmatter is valid**

Run: `head -4 .claude/skills/software-design-philosophy.md`
Expected:
```
---
name: software-design-philosophy
description: Use when reviewing code, designing features, or evaluating architecture decisions...
---
```

---

### Task 3: Manual smoke test

**Step 1: Test review mode**

Invoke `/software-design-philosophy` and provide a code snippet or file path. Verify it:
- Detects review mode
- Surfaces 2-3 relevant findings
- Uses the correct output format
- Cites specific code locations

**Step 2: Test design mode**

Invoke `/software-design-philosophy` with a feature description. Verify it:
- Detects design mode
- Surfaces 2-3 relevant principles
- Uses the correct output format
- Provides implementation guidance in the Action field
