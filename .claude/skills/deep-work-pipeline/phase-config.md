# Deep Work Pipeline — Phase Configuration

## Phases

| Phase | Skill | Artifact | Model | Interaction | Firewall | Gate |
|-------|-------|----------|-------|-------------|----------|------|
| 1 | dw-01-research-questions | 00-ticket.md, 01-research-questions.md | opus | none | no | human |
| 2 | dw-02-research | 02-research.md | opus | none | yes | human |
| 3 | dw-03-design-discussion | 03-design-discussion.md | opus | batch-qa | no | human |
| 4 | dw-04-outline | 04-structure-outline.md | opus | none | no | human |
| 5 | dw-05-plan | 05-plan.md | opus | none | no | human |
| 6 | dw-06b-implement-subagents | 06-completion.md | sonnet | none | no | human |

## Artifact Dependencies

| Phase | Reads | Writes |
|-------|-------|--------|
| 1 | — | 00-ticket.md, 01-research-questions.md |
| 2 | 01-research-questions.md (questions section only) | 02-research.md |
| 3 | 00-ticket.md, 02-research.md | 03-design-discussion.md |
| 4 | 00-ticket.md, 02-research.md, 03-design-discussion.md | 04-structure-outline.md |
| 5 | all prior artifacts | 05-plan.md |
| 6 | 05-plan.md | 06-completion.md |

## Gate Modes

- **human**: Present artifact summary to user, ask approve/revise/abort
- **auto** (future): Run quality checks, auto-advance if passing

## Firewall Constraint (Phase 2)

Phase 2 MUST NOT receive the original prompt or read 00-ticket.md.
The Phase 2 skill handles extraction internally via `extract-research-questions.sh` —
the pipeline orchestrator does not need to embed questions in the teammate prompt.
This ensures research objectivity before the prompt re-enters in Phase 3.
