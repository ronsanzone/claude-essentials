#!/usr/bin/env bash
# Phase 1 validation: write-plan SKILL.md exists with correct frontmatter and required sections.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

f=.claude/skills/write-plan/SKILL.md
[ -f "$f" ] || { echo "FAIL: $f missing"; exit 1; }

# Frontmatter
head -5 "$f" | grep -q "^name: write-plan$"   || { echo "FAIL: frontmatter name not 'write-plan'"; exit 1; }
head -5 "$f" | grep -Eq "^description: .+"     || { echo "FAIL: frontmatter description missing or empty"; exit 1; }

# Required sections (must appear as Markdown headings)
for heading in \
    "^## Setup" \
    "^## Pre-flight" \
    "^## Process" \
    "^### Step 1: Parse input" \
    "^### Step 2: Light research" \
    "^### Step 3: Draft plan" \
    "^### Step 4: Write artifact" \
    "^## Completion"; do
    grep -Eq "$heading" "$f" || { echo "FAIL: missing section matching /$heading/"; exit 1; }
done

# References to required infrastructure
grep -q "dw-setup.sh"            "$f" || { echo "FAIL: must reference dw-setup.sh";        exit 1; }
grep -q "codebase-locator"       "$f" || { echo "FAIL: must reference codebase-locator";   exit 1; }
grep -q "codebase-analyzer"      "$f" || { echo "FAIL: must reference codebase-analyzer";  exit 1; }
grep -q "Research Context"       "$f" || { echo "FAIL: must instruct to add Research Context section to plan"; exit 1; }
grep -q "Phase Progress"         "$f" || { echo "FAIL: must instruct to write Phase Progress table"; exit 1; }
grep -q "Task Completion"        "$f" || { echo "FAIL: must instruct to write Task Completion table"; exit 1; }
grep -q "Deviation Log"          "$f" || { echo "FAIL: must instruct to write Deviation Log"; exit 1; }

echo "PASS: Phase 1 invariants hold."
