---
name: pr-description
description: Use when creating or populating a PR description from git changes. Triggers on requests to write PR summaries, fill out PR templates, or prepare PR descriptions before opening a pull request.
---

# PR Description Generator

Generates a concise, reviewer-focused PR description from git changes and optional context documents.

## Usage

```
/pr-description [context-file-paths...]
```

Examples:
- `/pr-description` - Generate from current changes only
- `/pr-description docs/plans/feature-design.md` - Include design doc as context
- `/pr-description research.md plan.md` - Multiple context files

## Critical Rules

**DO NOT:**
- Modify anything below the main sections (merge checklist, "Carefully Review", etc.)
- Check or uncheck any checklist items - those are for reviewers
- Auto-discover context files - ONLY use files explicitly provided as arguments
- Exceed length limits (see below)

**DO:**
- Remove "Open Questions" section entirely if no unresolved questions exist
- Remove "Performance" section entirely if changes are not performance-sensitive
- Keep Summary under 200 words with 2-4 bullet points max

## Implementation

### Step 1: Find PR Template

Search these locations in order:
```
.github/PULL_REQUEST_TEMPLATE.md
.github/pull_request_template.md
.github/PULL_REQUEST_TEMPLATE/default.md
docs/pull_request_template.md
PULL_REQUEST_TEMPLATE.md
```

**If not found:** Use AskUserQuestion to request the template.

### Step 2: Gather Git Context

```bash
# Changes on this branch vs main
git diff origin/$(git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@')...HEAD

# Commit messages
git log origin/$(git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@')..HEAD --oneline

# Branch name (often has ticket number)
git branch --show-current
```

### Step 3: Read Context Files (if provided)

Only read files explicitly passed as arguments. Do NOT search for context files.

### Step 4: Fill Template

| Section | Rules |
|---------|-------|
| **Ticket** | Extract from branch name or commits, format as link |
| **Summary** | 2-4 bullets, 150-200 words max. Lead with what/why. Include reviewer hot topics. |
| **Open Questions** | Include only if questions exist. Otherwise DELETE section. |
| **Testing** | How work was tested. Be specific. |
| **Performance** | Include only if performance-relevant. Otherwise DELETE section. |
| **Everything below** | DO NOT TOUCH. Leave exactly as-is in template. |

### Step 5: Output

Present filled template for user review. Do NOT create PR automatically.

## Length Limits (Strict)

**Total editable content: Under one page (~400 words)**

| Section | Max |
|---------|-----|
| Summary | 200 words |
| Testing | 100 words |
| Open Questions | 50 words (if included) |
| Performance | 100 words (if included) |

## Hot Topics to Include in Summary

When these patterns appear in changes, mention them for reviewers:
- External API calls (latency, error handling)
- Database/schema changes (migration concerns)
- Auth/permission changes (security)
- Concurrency (goroutines, locks, channels)
- Caching (invalidation strategy)
- Config changes (deployment coordination)

## Quick Reference

| Step | Action |
|------|--------|
| 1 | Find template or ask for it |
| 2 | Get git diff, commits, branch |
| 3 | Read ONLY user-provided context files |
| 4 | Fill sections, DELETE unused optional sections |
| 5 | Output for review (don't create PR) |
