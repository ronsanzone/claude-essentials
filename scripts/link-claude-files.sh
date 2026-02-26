#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: scripts/link-claude-files.sh [--force] [--dry-run] [--source-dir DIR] [--target-dir DIR]

Creates symlinks for each top-level entry in source .claude into target ~/.claude.
Does NOT symlink the source .claude directory itself.

Options:
  --force             Replace existing non-symlink targets.
  --dry-run           Print planned actions without changing files.
  --source-dir DIR    Source .claude directory (default: <repo>/.claude).
  --target-dir DIR    Target directory (default: ~/.claude).
  -h, --help          Show this help.
EOF
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEFAULT_SOURCE_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)/.claude"
SOURCE_DIR="${DEFAULT_SOURCE_DIR}"
TARGET_DIR="${HOME}/.claude"
FORCE=0
DRY_RUN=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --force)
      FORCE=1
      shift
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --source-dir)
      SOURCE_DIR="${2:-}"
      shift 2
      ;;
    --target-dir)
      TARGET_DIR="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ ! -d "$SOURCE_DIR" ]]; then
  echo "Source directory not found: $SOURCE_DIR" >&2
  exit 1
fi

run() {
  if [[ "$DRY_RUN" -eq 1 ]]; then
    echo "[dry-run] $*"
  else
    "$@"
  fi
}

echo "Source: $SOURCE_DIR"
echo "Target: $TARGET_DIR"
run mkdir -p "$TARGET_DIR"

shopt -s dotglob nullglob
for source_path in "$SOURCE_DIR"/*; do
  name="$(basename "$source_path")"
  target_path="${TARGET_DIR}/${name}"

  if [[ -L "$target_path" ]]; then
    current_link="$(readlink "$target_path" || true)"
    if [[ "$current_link" == "$source_path" ]]; then
      echo "ok: ${name} (already linked)"
      continue
    fi
    echo "relink: ${name}"
    run rm -f "$target_path"
    run ln -s "$source_path" "$target_path"
    continue
  fi

  if [[ -e "$target_path" ]]; then
    if [[ "$FORCE" -eq 1 ]]; then
      echo "replace: ${name}"
      run rm -rf "$target_path"
      run ln -s "$source_path" "$target_path"
    else
      echo "skip: ${name} (target exists; use --force to replace)"
    fi
    continue
  fi

  echo "link: ${name}"
  run ln -s "$source_path" "$target_path"
done
