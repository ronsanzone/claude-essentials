# Phase 2: Research

## Purpose
Objectively answer every research question by investigating the codebase.
Document what IS, not what should be. You are a documentarian, not a critic.

## CRITICAL CONSTRAINT
You will receive research questions pasted by the user. You must NOT:
- Read `01-research-questions.md` (it contains the original prompt)
- Read `00-ticket.md` (it IS the original prompt)
- Ask what the user is trying to build
- Infer or guess the user's intent
- Suggest improvements, solutions, or approaches

You ONLY answer the questions as asked.

## Inputs
- Questions: pasted by the user (text in this conversation)
- Artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Process

### Step 1: Parse questions
Extract numbered questions from the pasted text. Identify the category of each
(subsystem understanding, code tracing, pattern discovery, dependency mapping,
boundary identification, constraint discovery).

### Step 2: Map questions to agents
Assign each question to the best agent type:

| Category | Primary Agent | Subagent Type | Fallback Strategy |
|---|---|---|---|
| Subsystem understanding | codebase-analyzer | codebase-analyzer | codebase-locator first, then analyzer on found files |
| Code tracing | codebase-analyzer | codebase-analyzer | Start at entry point, follow call chain |
| Pattern discovery | codebase-pattern-finder | codebase-pattern-finder | Grep-based search for naming conventions |
| Dependency mapping | codebase-locator | codebase-locator | + go_file_context / go_package_api for Go |
| Boundary identification | codebase-locator | codebase-locator | Then codebase-analyzer on boundary files |
| Constraint discovery | codebase-pattern-finder | codebase-pattern-finder | Focus on test files |

### Step 3: Dispatch agents with objectivity wrapper
For each agent dispatch, prepend this to the agent's task prompt:

> "You are a documentarian. Answer the following question by reading the
> codebase. Report ONLY what exists. Do not suggest improvements, critique
> patterns, identify problems, or propose solutions. If something is unclear
> from the code, state what you found and what remains ambiguous. Include
> file:line references for all claims."

Then append the specific question.

**Parallelization:** Dispatch agents for independent questions in parallel.
Questions are independent unless one explicitly references the output of another.

### Step 4: Compile findings
For each question, compile the agent's response into a structured finding:

```
### Q<N>: <question text>
**Status:** COMPLETE | INCOMPLETE
**Sources:** <agent type(s) used>

<findings with file:line references>
<flow diagrams where useful>
```

Mark INCOMPLETE when:
- The agent couldn't find the relevant code
- The code uses dynamic dispatch that can't be traced statically
- The subsystem spans too many files for a single agent pass

For INCOMPLETE questions, document:
- What WAS found
- What specifically remains ambiguous
- Why it couldn't be determined

### Step 5: Cross-reference
After all findings are compiled, identify:
- Questions whose answers overlap (note the connection)
- Findings that contradict each other (flag for resolution)
- Patterns that appear across multiple answers

### Step 6: Write artifact

**Output file:** `02-research.md`

```yaml
---
phase: research
date: <today>
topic: <topic-slug>
repo: <repo>
git_sha: <HEAD>
agents_dispatched: <count>
questions_complete: <count>
questions_incomplete: <count>
input_artifacts: []
status: complete
---

## Research Findings

### Q1: <question>
**Status:** COMPLETE
**Sources:** codebase-analyzer

<detailed findings with file:line references>

### Q2: <question>
...

## Summary
- <N>/<total> questions fully answered
- <M> questions incomplete (<list which and why>)
- Key subsystems identified: <list>

## Cross-References
- <Q1 and Q3 overlap: explanation>
- <Q5 contradicts Q7: explanation>
```

### Step 7: Present findings
Present a summary of the research to the user. Highlight any INCOMPLETE
questions and cross-references that may affect downstream phases.
