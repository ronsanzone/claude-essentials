---
name: evergreen
description: Use when debugging CI failures from MongoDB Evergreen, fetching patch/PR status, or getting test failure details for local debugging
---

# Evergreen CI Skill

Fetch patch status, test failures, and logs from MongoDB's Evergreen CI system.

## Prerequisites

- `evergreen` CLI installed (`which evergreen`)
- Config at `~/.evergreen.yml` with `api_key` and `user`
- `jq` installed for JSON parsing

## Commands

### `/evergreen patch <patch_id|pr_url>`

Get patch overview and identify failed tasks.

### `/evergreen failures <patch_id>`

Get all test failures with verbatim errors and commands.

### `/evergreen debug <task_id>`

Map a single task's failures to local code paths.

## Critical: Getting Task IDs

**The CLI alone cannot list task IDs.** You must use the REST API:

```bash
# Read credentials from config
API_KEY=$(grep api_key ~/.evergreen.yml | awk '{print $2}')
API_USER=$(grep "^user:" ~/.evergreen.yml | awk '{print $2}')
API_HOST="https://evergreen.mongodb.com/api"

# Get builds (contains task IDs)
curl -s -H "Api-User: $API_USER" -H "Api-Key: $API_KEY" \
  "$API_HOST/rest/v2/versions/<patch_id>/builds" | jq '.'
```

Each build contains a `tasks` array with full task IDs.

## Implementation Reference

### Identifier Detection

| Input | Pattern | Action |
|-------|---------|--------|
| Patch ID | 24-char hex (`^[a-f0-9]{24}$`) | Use directly |
| PR URL | `github.com/.*/pull/(\d+)` | Search `list-patches` for matching `github_patch_data.pr_number` |

### `/evergreen patch` Implementation

```bash
# 1. Normalize to patch_id
if [[ "$ID" =~ ^[a-f0-9]{24}$ ]]; then
    PATCH_ID="$ID"
elif [[ "$ID" =~ github.com/.*/pull/([0-9]+) ]]; then
    PR_NUM="${BASH_REMATCH[1]}"
    PATCH_ID=$(evergreen list-patches -n 20 --json | jq -r \
        ".[] | select(.github_patch_data.pr_number == $PR_NUM) | .patch_id" | head -1)
fi

# 2. Get patch info
PATCH=$(evergreen list-patches -i "$PATCH_ID" --json)
echo "$PATCH" | jq '{patch_id, status, author, description}'

# 3. Get builds with task IDs (REST API required)
BUILDS=$(curl -s -H "Api-User: $API_USER" -H "Api-Key: $API_KEY" \
  "$API_HOST/rest/v2/versions/$PATCH_ID/builds")

# 4. For each task, get status
echo "$BUILDS" | jq -r '.[].tasks[]' | while read TASK_ID; do
    TASK=$(curl -s -H "Api-User: $API_USER" -H "Api-Key: $API_KEY" \
      "$API_HOST/rest/v2/tasks/$TASK_ID")
    echo "$TASK" | jq '{display_name, status}'
done
```

### `/evergreen failures` Implementation

```bash
# Get failed tasks
BUILDS=$(curl -s -H "Api-User: $API_USER" -H "Api-Key: $API_KEY" \
  "$API_HOST/rest/v2/versions/$PATCH_ID/builds")

echo "$BUILDS" | jq -r '.[].tasks[]' | while read TASK_ID; do
    TASK=$(curl -s -H "Api-User: $API_USER" -H "Api-Key: $API_KEY" \
      "$API_HOST/rest/v2/tasks/$TASK_ID")

    STATUS=$(echo "$TASK" | jq -r '.status')
    if [[ "$STATUS" == "failed" ]]; then
        # Get test results
        TESTS=$(curl -s -H "Api-User: $API_USER" -H "Api-Key: $API_KEY" \
          "$API_HOST/rest/v2/tasks/$TASK_ID/tests")

        # Show failed tests
        echo "$TESTS" | jq '.[] | select(.status == "fail")'

        # Get task logs
        evergreen task build TaskLogs --task_id "$TASK_ID" --type task_log --tail_limit 200
    fi
done
```

### `/evergreen debug` Implementation

```bash
# Get test failures with file references
TESTS=$(curl -s -H "Api-User: $API_USER" -H "Api-Key: $API_KEY" \
  "$API_HOST/rest/v2/tasks/$TASK_ID/tests")

# Extract failed tests
echo "$TESTS" | jq -r '.[] | select(.status == "fail") | .test_file' | while read TEST_FILE; do
    # Map to local file
    LOCAL_FILE=$(find . -name "$(basename $TEST_FILE)" -type f | head -1)
    echo "→ Local file: $LOCAL_FILE"
done

# Get logs for context
evergreen task build TestLogs --task_id "$TASK_ID" --tail_limit 100
```

## Output Format Contract

For composability with other skills:

| Element | Format |
|---------|--------|
| Failed task IDs | Under `Failed Task IDs:` header |
| Verbatim errors | In `Error Output (verbatim):` block |
| Local file paths | Prefixed with `→ Test:` or `→ Source:` |
| Commands | In `Command:` block |

## Error Handling

| Condition | Response |
|-----------|----------|
| `evergreen` not in PATH | "Install evergreen CLI: `evergreen get-update`" |
| `~/.evergreen.yml` missing | "Run: `evergreen client setup`" |
| Patch not found | "Patch ID not found. Verify ID or check `evergreen list-patches`" |
| PR has no patch | "No Evergreen patch found for PR #X" |
| All tasks passed | "No failures - all X tasks passed ✓" |

## Quick Reference

```bash
# Setup: Read credentials once
API_KEY=$(grep api_key ~/.evergreen.yml | awk '{print $2}')
API_USER=$(grep "^user:" ~/.evergreen.yml | awk '{print $2}')

# Get your recent patches
evergreen list-patches -n 5 --json

# Get builds for a patch (contains task IDs)
curl -s -H "Api-User: $API_USER" -H "Api-Key: $API_KEY" \
  "https://evergreen.mongodb.com/api/rest/v2/versions/<patch_id>/builds"

# Get task status
curl -s -H "Api-User: $API_USER" -H "Api-Key: $API_KEY" \
  "https://evergreen.mongodb.com/api/rest/v2/tasks/<task_id>"

# Get test results for a task
curl -s -H "Api-User: $API_USER" -H "Api-Key: $API_KEY" \
  "https://evergreen.mongodb.com/api/rest/v2/tasks/<task_id>/tests"

# Get task logs via CLI
evergreen task build TaskLogs --task_id "<task_id>" --type task_log
evergreen task build TestLogs --task_id "<task_id>"
```
