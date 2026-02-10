# Phase 1: Research Questions

## Purpose
Decompose the user's prompt into objective, investigative questions that can be
answered by reading the codebase. These questions must NOT assume any particular
solution.

## Inputs
- User's prompt: from `00-ticket.md` in the artifact directory, or pasted directly
- Artifact directory: `~/notes/context-engineering/<repo>/<topic-slug>/`

## Process

### Step 1: Read the prompt
Read the user's prompt from 00-ticket.md or accept it as pasted text.

### Step 2: Targeted codebase scan
Gather lightweight structural context (NOT deep implementation details):
- List root directory structure
- Read CLAUDE.md files for project context and conventions
- Dispatch a codebase-locator agent to find areas relevant to what the prompt
  mentions. The locator prompt should be: "Find files and directories related to:
  <extract key nouns/systems from prompt>. Return locations grouped by purpose."

### Step 3: Generate research questions
Generate 5-15 questions. EVERY question must be:
- **Objective** — answerable by reading code, not by making design decisions
- **Specific** — references concrete subsystems, not abstract concepts
- **Grounded** — uses real module/file names from the codebase scan

Distribute across these categories:

| Category | Question Pattern | Example |
|----------|-----------------|---------|
| Subsystem Understanding | "How does [component] work?" | "How does the auth middleware chain process requests?" |
| Code Tracing | "What is the [data/request] flow from [A] to [B]?" | "What is the request lifecycle from HTTP handler to DB write?" |
| Pattern Discovery | "What patterns exist for [action]?" | "What patterns exist for adding new API endpoints?" |
| Dependency Mapping | "What does [module] depend on / what depends on [module]?" | "What does the handlers package import?" |
| Boundary Identification | "Where do [A] and [B] integrate?" | "Where do the HTTP and storage layers connect?" |
| Constraint Discovery | "What invariants does [system] enforce?" | "What does the test suite enforce for handler responses?" |

**FORBIDDEN question patterns:**
- "How should we..." — this is solutioning
- "What's the best way to..." — this is evaluation
- "Would it be better to..." — this is comparison
- "Can we..." — this is feasibility for a specific solution

### Step 4: Present and write artifact

Present the questions to the user grouped by category. Then write the artifact:

**Output file:** `01-research-questions.md`

```yaml
---
phase: research-questions
date: <today>
topic: <topic-slug>
repo: <repo>
git_sha: <HEAD>
status: complete
---

## Original Prompt
<the user's full prompt — stored for traceability, NOT passed to Phase 2>

## Research Questions

### Subsystem Understanding
1. <question>
2. <question>

### Code Tracing
3. <question>

### Pattern Discovery
4. <question>
5. <question>

### Dependency Mapping
6. <question>

### Boundary Identification
7. <question>

### Constraint Discovery
8. <question>
```

### Step 5: Handoff
After writing, tell the user:
"Here are the research questions. Review and edit as needed. When ready, **copy
the Research Questions section** (everything below '## Research Questions') and
paste it to start the research phase. The research agent will ONLY see what you
paste."
