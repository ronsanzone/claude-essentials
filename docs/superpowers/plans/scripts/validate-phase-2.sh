#!/usr/bin/env bash
# Phase 2 validation: implement-plan SKILL.md exists with correct frontmatter and required sections.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

f=.claude/skills/implement-plan/SKILL.md
[ -f "$f" ] || { echo "FAIL: $f missing"; exit 1; }

# Frontmatter
head -5 "$f" | grep -q "^name: implement-plan$" || { echo "FAIL: frontmatter name not 'implement-plan'"; exit 1; }
head -5 "$f" | grep -Eq "^description: .+"      || { echo "FAIL: frontmatter description missing or empty"; exit 1; }

# Required sections
for heading in \
    "^## Setup" \
    "^## Pre-flight" \
    "^## Tooling" \
    "^## Model Selection" \
    "^## The Process" \
    "^## Session Review" \
    "^## Resume" \
    "^## Red [Ff]lags"; do
    grep -Eq "$heading" "$f" || { echo "FAIL: missing section matching /$heading/"; exit 1; }
done

# References to required infrastructure
grep -q "implementer-prompt.md"             "$f" || { echo "FAIL: must reference ./implementer-prompt.md";             exit 1; }
grep -q "spec-reviewer-prompt.md"           "$f" || { echo "FAIL: must reference ./spec-reviewer-prompt.md";           exit 1; }
grep -q "code-quality-reviewer-prompt.md"   "$f" || { echo "FAIL: must reference ./code-quality-reviewer-prompt.md";   exit 1; }
grep -q "/quick-review"                     "$f" || { echo "FAIL: must reference /quick-review for session review";    exit 1; }
grep -q "plan.md"                           "$f" || { echo "FAIL: must reference plan.md as input artifact";           exit 1; }
grep -q "Research Context"                  "$f" || { echo "FAIL: must reference Research Context as implementer ref"; exit 1; }
grep -q "Task Completion"                   "$f" || { echo "FAIL: must update Task Completion table per task";         exit 1; }

echo "PASS: Phase 2 invariants hold."
