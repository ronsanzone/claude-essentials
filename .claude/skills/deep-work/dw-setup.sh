#!/usr/bin/env bash
# dw-setup.sh — Resolve the standard deep-work skill setup variables.
#
# Usage: dw-setup.sh <topic-slug>
#
# On success (exit 0), prints three KEY=VALUE lines to stdout:
#   REPO=<derived from `git remote get-url origin`, falls back to basename of pwd>
#   TOPIC_SLUG=<the slug passed in>
#   ARTIFACT_DIR=<$HOME/notes/context-engineering/<repo>/<slug>, mkdir -p applied>
#
# On missing slug (exit 2), prints "MISSING_SLUG" to stderr.

slug="${1:-}"
if [ -z "$slug" ]; then
    echo "MISSING_SLUG" >&2
    exit 2
fi

repo=$(basename "$(git remote get-url origin 2>/dev/null | sed 's/.git$//')" 2>/dev/null)
if [ -z "$repo" ]; then
    repo=$(basename "$(pwd)")
fi
if [ -z "$repo" ]; then
    echo "MISSING_REPO" >&2
    exit 3
fi
artifact_dir="$HOME/notes/context-engineering/$repo/$slug"

mkdir -p "$artifact_dir"

printf 'REPO=%s\n' "$repo"
printf 'TOPIC_SLUG=%s\n' "$slug"
printf 'ARTIFACT_DIR=%s\n' "$artifact_dir"
