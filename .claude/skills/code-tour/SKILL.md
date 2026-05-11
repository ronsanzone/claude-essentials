---
name: code-tour
description: Use when you need to deeply understand a PR, branch, or system/feature end-to-end. Triggers on "explain this PR", "walk me through", "how does X work", "tour of this change", "help me understand this system", "ramp me up on".
---

# Code Tour

Deep end-to-end guided tour of a PR, branch, or system feature. Dispatches 5 parallel research agents to trace flows across the codebase, then compiles findings and writes a self-contained HTML document — using the `frontend-design` skill for styling — with SVG diagrams, syntax-highlighted code, and collapsible sections.

## Usage

```
/code-tour <PR-URL-or-number>          # Tour centered on a PR
/code-tour                             # Tour of current branch diff vs main
/code-tour <topic>                     # Tour of a system/feature (no diff)
```

## Scoping Principle

All research and output is scoped to the specific subsystem under review. Do NOT explain the broader application, build system, or unrelated subsystems. Assume the reader works in this codebase daily.

## Implementation

### Step 1: Parse Input & Determine Mode

Determine mode from arguments:
- **PR mode:** Argument is numeric or contains `github.com`
- **Branch mode:** No arguments and current branch is not main/master
- **Topic mode:** Free-text argument that isn't a PR reference

### Step 2: Gather Raw Material

**PR mode:**
```bash
gh pr view -R <owner/repo> <number>
gh pr diff -R <owner/repo> <number>
```

**Branch mode:**
```bash
git diff main...HEAD
git log main..HEAD --oneline
```

**Topic mode:** The argument text is the scope description. No diff to fetch.

### Step 3: Analyze Scope

Read the diff or topic description and identify:
- Key files and packages touched or relevant
- The subsystem name (used to scope all research agents)
- The central abstraction or entry point
- **Reader profile**: ramp-up tour (surface-level prior knowledge) vs. PR review tour (deep prior knowledge). Ramp-up tours need a Lifecycle/Workflow section (see Step 5); PR review tours can usually skip it.

### Step 4: Dispatch 5 Parallel Research Agents

Launch all 5 agents in a single message using `codebase-analyzer` subagent type. Each agent gets:
- The diff or topic description for context
- The key files/packages from Step 3
- This instruction: "Scope your research to this subsystem. Do not explain the broader application. Include file:line citations on every claim. Report findings concisely."

| Agent | Research Dimension | Key Questions |
|-------|--------------------|---------------|
| **Framework** | Base class, interface, or framework the change extends | What contract does it implement? What are the abstractions and extension points? What must subclasses provide? |
| **Upstream** | Who calls into the changed code | What events, actions, or callers trigger this code? What's the entry point? |
| **Downstream** | Who consumes the output | Where does the data flow next? What decisions or behaviors does it affect? |
| **Model/Persistence** | Data structures, DAOs, schemas | How is state stored and queried? What are the key fields and their semantics? |
| **The Change** (PR/branch) or **System Overview** (topic) | The diff itself, or the feature's boundaries | What specifically changed and why? Or: what are the components and boundaries of this system? |

For ramp-up tours where the reader is new to the subsystem, also have at least one agent surface recent activity themed by feature area (`git log --since` grouped by ticket/theme) — this powers the "What's Actively Being Built" section.

### Step 5: Compile Tour

Wait for all 5 agents to complete, then compile findings using the section structure below. The narrative must build incrementally — each section should be understandable given only the sections before it.

#### Required sections (always include)

**1. What This System Does** — 1–2 paragraphs of plain English describing the specific subsystem, not the application. No jargon without explanation. Someone unfamiliar with this subsystem should understand its purpose after reading this. End with a short list of "things that make this different / worth tracing carefully" if the research surfaced any.

**2. The Big Picture** — Inline SVG diagram showing the full end-to-end flow from trigger to effect. This is the "map" that everything else refers back to. Label key decision points and data transformations.

**3. Lifecycle / Workflow** — *Required for ramp-up tours; optional for PR review tours.* A conceptual walkthrough of how the workflow unfolds across stages, with a stage-by-stage diagram. Make the **shared infrastructure vs. specialization split** explicit (e.g., "stages 1–4 are cloud-agnostic; 5–6 are where Alibaba code lives"). Then a short paragraph per stage saying what happens and which class owns it. End with a "why this matters when reading X code" note that anchors the abstraction back to where the reader will encounter it. This section gives the mental model that everything else hangs off — without it, jumping from a system diagram into base classes is too steep.

**4. Step-by-Step Tour** — Numbered collapsible sections walking through the flow in execution order. Each step includes:
- What happens and why (2–4 paragraphs max)
- Syntax-highlighted code snippets (10–30 lines, with `file:line` references)
- Additional SVG diagrams where they aid understanding
- Connection back to the big picture diagram

After each conceptual section, end with a brief "anchor" — one sentence mapping the abstraction to where the reader will encounter it in concrete code. Small but very high signal.

#### Supplementary sections (choose based on research findings)

Promote these three when the research surfaces material — they're frequently the highest-value additions to a tour, not just nice-to-haves:

- **What's Actively Being Built** — themed `git log --since="..."` summary grouped by feature area. Turns the tour into a current-state document, not just a static reference. Especially valuable for ramp-up tours.
- **What's Not Yet Wired / Gaps** — open TODOs, unsupported features, in-flight epics, embedded `// TODO(CLOUDP-...)` comments.
- **Sibling/Comparable Implementations Cheat-Sheet** — when the system is one of N parallel implementations (cloud providers, payment processors, auth flows), a side-by-side terminology + behavior table is high signal.

Other useful sections (include only when the research surfaces something interesting):
- **What The Change Adds/Modifies** — how the diff maps onto the system (PR/branch mode)
- **Key Design Decisions** — why the system works this way, not just how
- **The Data Model** — when persistence is complex enough to warrant its own section
- **Error Handling & Recovery** — when the system has interesting failure modes
- **Testing Strategy** — when the test approach is non-obvious

These are not a checklist. Use judgment — if research surfaced something that doesn't fit these categories, add a section for it.

### Step 6: Style and Write Output

Write the compiled tour as a **self-contained HTML file** (all CSS/SVG inline, no external dependencies — fonts via Google Fonts CDN are fine) to `~/notes/04_Research/<topic-slug>-tour.html`.

Derive the slug from:
- PR mode: ticket ID or PR title (e.g., `CLOUDP-398944-alibaba-capacity-denylist`)
- Branch mode: branch name (e.g., `feature-lcm-update-api`)
- Topic mode: slugified topic (e.g., `capacity-denylist-system`)

**Invoke the `frontend-design` skill before writing the HTML, with this brief:**

> Long-form technical reading document (~20–30 min, ~3000 words). Reader will consume sequentially in a focused session. Prose-heavy with mono-font code blocks and inline SVG diagrams.
>
> **Aesthetic constraint:** prefer light editorial themes (warm parchment, off-white, soft cream) over dark or stark white. Prioritize readability for long sessions. Distinctive typography (variable serif display + clean sans body + characterful mono) — avoid generic Inter/Arial defaults.
>
> **Required UX:**
> - Sticky TOC sidebar with scroll-spy (active link tracks current section)
> - Section numbers as marginalia in the gutter (not inline with the heading)
> - Reading-progress bar at the top
> - Constrained prose column (~700–740px) with diagrams/tables allowed to break wider
> - Hero with masthead-style meta-row (compiled date, branch, reading time, section count)
> - Pull quotes for the 1–2 stickiest insights
> - Collapsible `<details>` styled as elegant chapter sections (no buttony chrome)
> - SVG diagrams in cards with italic figcaptions
> - Syntax-highlighted code blocks (manual `<span>`-class highlighting)
>
> Write the final HTML directly — informed by both the compiled content and frontend-design's guidance. One pass.

After writing, display a brief summary in conversation with the file path, and open the file in the browser with `open <path>` so the user sees the result immediately.

### Step 7: Iteration is Expected

Tours typically need a second pass — the user may request content additions (e.g. a new high-level section), aesthetic adjustments, or a section reordering. Treat this as a normal part of the flow, not a failure mode. When renumbering sections, edit in **reverse order** (last → first) so each search string stays unique mid-edit.
