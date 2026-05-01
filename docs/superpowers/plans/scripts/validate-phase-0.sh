#!/usr/bin/env bash
# Phase 0 validation: skill dirs exist and reviewer prompts match dw-06 sources.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

[ -d .claude/skills/write-plan ]      || { echo "FAIL: write-plan dir missing"; exit 1; }
[ -d .claude/skills/implement-plan ]  || { echo "FAIL: implement-plan dir missing"; exit 1; }

for f in implementer-prompt.md spec-reviewer-prompt.md code-quality-reviewer-prompt.md; do
    src=".claude/skills/dw-06-implement/$f"
    dst=".claude/skills/implement-plan/$f"
    [ -f "$dst" ] || { echo "FAIL: $dst missing"; exit 1; }
    diff -q "$src" "$dst" >/dev/null || { echo "FAIL: $dst differs from $src"; exit 1; }
done

echo "PASS: Phase 0 invariants hold."
