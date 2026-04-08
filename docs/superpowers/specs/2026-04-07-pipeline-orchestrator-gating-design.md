# Deep-Work Pipeline: Configurable Gating & Thin Lead Orchestrator

**Date:** 2026-04-07
**Status:** Approved
**Scope:** Orchestrator skill (`deep-work-pipeline/SKILL.md`) + phase config (`deep-work-pipeline/phase-config.md`)

## Problem

The deep-work pipeline has strong intellectual architecture (6 phases with deliberate context isolation, bias firewall, artifact-based state transfer) but poor ergonomics:

1. **Conversation switching overhead** — running 6+ separate conversations with manual copy-paste between them is high friction
2. **No gate flexibility** — every phase requires human review even when the pipeline's judgment is trustworthy
3. **Orchestrator token waste** — the team lead accumulates teammate summaries in its context, growing linearly with pipeline length
4. **Phase 3 proxying is awkward** — the lead acts as an intermediary for design question resolution, adding friction vs. direct interaction

## Design

### Gate Mode System

Replace the hardcoded `human` gate with a configurable `--mode` argument supporting four presets:

```
/deep-work-pipeline <topic-slug> [--mode full|research-gate|design-gate|auto]
```

Default: `full` (preserves current behavior).

#### Gate Mode Matrix

Each mode automates a strictly larger prefix of phases. The gate name indicates where automation **stops** — everything after is human-gated.

| Phase | Skill | full | research-gate | design-gate | auto |
|-------|-------|------|---------------|-------------|------|
| 1 | dw-01-research-questions | human | auto | auto | auto |
| 2 | dw-02-research | human | human | auto | auto |
| 3 | dw-03-design-discussion | human | human | human | accept-recs |
| 4 | dw-04-outline | human | human | human | auto |
| 5 | dw-05-plan | human | human | human | auto |
| 6 | dw-06b-implement-subagents | human | human | human | human |

#### Gate Types

- **human** — Present artifact summary + path to user. Offer: Approve / Revise / Take over / Abort.
- **auto** — Check `.state.json` updated after teammate completes. Advance immediately. Stop on `STATUS: error`.
- **accept-recs** — Send "Accept all recommendations" to the Phase 3 teammate for design question resolution, then auto-advance.

#### Mode Descriptions

Modes form a progression of increasing automation: `full` < `research-gate` < `design-gate` < `auto`.

| Mode | Automates | Human from | Use case |
|------|-----------|------------|----------|
| `full` | Nothing | Phase 1 | Unfamiliar codebase, high-stakes changes. Review everything. |
| `research-gate` | Phase 1 | Phase 2 | Trust question generation. Review research before proceeding manually through design, outline, plan, and implementation. |
| `design-gate` | Phases 1-2 | Phase 3 | Trust research. Review design decisions and proceed manually through outline, plan, and implementation. |
| `auto` | Phases 1-5 | Phase 6 | Re-runs, well-understood features. Auto-advance everything (accept-recs for P3) except implementation. |

### Thin Lead Architecture

The orchestrator lead becomes a thin dispatcher. It never reads artifacts or accumulates phase content.

#### Auto-gated phases

```
Lead spawns teammate -> teammate works -> agent returns -> lead checks .state.json -> lead spawns next phase
```

The lead does not process the teammate's output. It confirms `.state.json` advanced, then moves on. Context growth is minimal: just the agent spawn/return cycle.

#### Human-gated phases

```
Lead spawns teammate -> teammate works -> agent returns with final message -> lead presents artifact path + teammate summary bullets to user -> user decides -> lead acts on decision
```

The lead receives the teammate's final message (Agent tool return value) and forwards the summary bullets and artifact path. It does not interpret or accumulate artifact content.

### Phase 3 Design Question Flow

Phase 3 is the only phase requiring mid-flight user interaction (design question resolution). Handled per mode:

#### Interactive modes (`full`, `research-gate`, `design-gate`)

1. Teammate writes draft artifact with OPEN design questions
2. Teammate sends `STATUS: needs-input` with question summary (titles + options + recommendations)
3. Lead presents questions to user via batch mode proxy
4. User provides answers (e.g., "DQ-1: A, DQ-3: B") or says "accept all"
5. Lead relays answers to teammate via SendMessage
6. Teammate finalizes artifact, reports `STATUS: complete`
7. Normal gate proceeds

#### Auto mode (`auto`)

1. Teammate prompt includes directive: "For Step 7 (question resolution), choose 'Accept recommendations' mode."
2. Teammate runs to completion autonomously, writes finalized artifact
3. Lead checks `.state.json`, advances

### Manual Takeover Protocol

At any human gate, the user can choose **Take over** to work directly with the teammate session. This applies to all phases, not just Phase 3.

1. User selects "Take over" at a gate
2. Lead provides the teammate session identifier and waits
3. User works directly with the teammate (editing artifacts, answering questions, iterating)
4. User tells the lead "done" / "continue" / "phase complete"
5. Lead checks `.state.json` to confirm the phase completed, then advances to the next phase

This handles complex scenarios where the batch proxy is insufficient — the user can always drop down to direct interaction and return control to the lead when ready.

### Error Handling

#### Auto-gated phase errors

| Condition | Action |
|-----------|--------|
| Teammate returns `STATUS: complete`, `.state.json` updated | Advance to next phase |
| Teammate returns `STATUS: error` | Stop auto-advance. Present error to user. Offer: Retry / Take over / Abort |
| Teammate returns `STATUS: needs-input` on non-P3 phase | Unexpected. Present to user. |
| Teammate returns `STATUS: complete`, `.state.json` NOT updated | Warn user: "Teammate reported complete but state not updated." Offer: Continue anyway / Investigate / Abort |

#### Resumability

No changes needed. Existing `.state.json` resume check handles mid-pipeline failures:
- On restart, lead reads `.state.json`, reports completed phases, offers to resume from next incomplete phase
- The `gate_mode` field is stored in `.state.json` at pipeline start so resume uses the same mode

### Model Selection

All teammate phases use `opus`. The current config incorrectly specified `sonnet` for Phase 6 — the orchestrator teammate needs strong reasoning for task coordination, review interpretation, and deviation handling. Implementation *subagents* dispatched internally by Phase 6 use `sonnet` as specified in their own prompt templates.

| Phases | Model | Rationale |
|--------|-------|-----------|
| 1-6 | opus | All orchestrator teammates need strong reasoning. Phase 6 subagents use sonnet internally. |

## Files Changed

| File | Change |
|------|--------|
| `deep-work-pipeline/SKILL.md` | Rewrite orchestrator with gate mode system, thin lead, manual takeover, error handling |
| `deep-work-pipeline/phase-config.md` | Replace single Gate column with per-mode gate matrix. Fix Phase 6 model to opus. |

## Out of Scope

- Individual phase skill changes (research scoping, plan size, outline/plan consolidation)
- `.state.json` schema standardization
- Cross-feature linking between pipeline runs
- Bias firewall fix (prompt in `01-research-questions.md`)
- New gate modes beyond the four defined here

## Risks

- **Auto-accept design recommendations may produce suboptimal designs** — Mitigated by: user chooses mode per-run, can always switch to `design-gate` for complex features. The pipeline's design question recommendations are grounded in research findings.
- **Manual takeover complexity** — The lead must handle the user returning control at any point. Mitigated by: lead simply checks `.state.json` rather than trying to understand what the user did.
- **`gate_mode` in `.state.json`** — If user resumes a pipeline in a different mode than it started, phase completion may not make sense. Mitigated by: on resume, show the original `gate_mode` and ask if user wants to continue with it or switch.
