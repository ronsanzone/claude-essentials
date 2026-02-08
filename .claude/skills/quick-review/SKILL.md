---
name: quick-review
description: Use when you need fast, interactive code review of a PR with detailed explanations. Triggers on exploratory reviews, learning-focused feedback, or when you want to discuss findings before any action.
---

# Quick Review

Single-pass expert code review for fast, interactive feedback with detailed explanations of every issue found.

## Overview

A direct, conversational code review that explains the reasoning behind every finding. No multi-agent orchestration - just thorough expert analysis with interactive clarification.

**When to use this vs pr-review:**
- **quick-review**: Fast feedback, learning/discussion, exploring issues interactively
- **pr-review**: Formal review, high-confidence filtering, ready-to-post findings

## Usage

```
/quick-review <owner/repo> <pr-number>
```

Examples:
- `/quick-review anthropics/claude-code 123`
- `/quick-review myorg/myrepo 456`

## Review Priorities (in order)

1. **Correctness** - Does it work as specified?
2. **Security** - Are there vulnerabilities?
3. **Performance** - Are there obvious bottlenecks?
4. **Maintainability** - Will future developers understand it?
5. **Style** - Does it follow conventions?

## Analysis Framework

Apply these in priority order:

| Priority | Focus | What to Look For |
|----------|-------|------------------|
| 0 | **Correctness** | Bugs, spec violations, logic errors, edge cases |
| 1 | **Security** | Injection, auth flaws, data exposure, insecure defaults |
| 2 | **Performance** | Algorithmic complexity, memory, queries, bottlenecks |
| 3 | **Code Quality** | Readability, naming, structure, organization |
| 4 | **Best Practices** | Language idioms, design patterns, standards |
| 5 | **Error Handling** | Exceptions, validation, failure recovery |
| 6 | **Testing** | Testability, missing test cases, test strategies |

## Implementation

### Step 1: Get the PR

```bash
gh pr view -R <owner/repo> <pr-number>
gh pr diff -R <owner/repo> <pr-number>
```

### Step 2: Check for CLAUDE.md (improvement over original)

Look for project-specific standards:
```bash
# Check root and modified directories for CLAUDE.md
gh pr view -R <owner/repo> <pr-number> --json files -q '.files[].path' | \
  xargs -I {} dirname {} | sort -u | \
  while read dir; do echo "$dir/CLAUDE.md"; done
```

If CLAUDE.md exists, factor its guidance into the review.

### Step 3: Review the Diff

Apply the Analysis Framework systematically. For each issue found:

1. **Identify** the problem clearly
2. **Explain** why it's a problem (consequences if not fixed)
3. **Suggest** a concrete fix with code example
4. **Categorize** by severity

### Step 4: Ask Clarifying Questions

If context is needed to provide better feedback, ask specific questions:
- "Is this intentionally handling X this way, or should it also cover Y?"
- "What's the expected behavior when Z occurs?"
- "Does this need to maintain backwards compatibility with...?"

## Output Structure

```markdown
## Code Review: <PR title>

### Executive Summary
<2-3 sentences on overall quality and key concerns>

### Critical Issues
Issues requiring immediate attention - security vulnerabilities, major bugs.

1. **<Issue title>**
   - **File:** `path/to/file.go:123`
   - **Problem:** <Clear description>
   - **Why it matters:** <Consequences>
   - **Suggested fix:**
   ```go
   // concrete code example
   ```

### Significant Improvements
Important but non-critical enhancements.

1. **<Issue title>**
   - **File:** `path/to/file.go:456`
   - **Problem:** <Description>
   - **Suggestion:** <How to improve>

### Minor Suggestions
Style, readability, optimization opportunities.

- `file.go:78` - Consider renaming `x` to `userCount` for clarity
- `file.go:92` - This loop could use `range` instead of index access

### Positive Observations
Well-implemented patterns worth noting.

- Good use of error wrapping in `handler.go`
- Clean separation of concerns in the service layer
- Comprehensive input validation

### Questions
<Any clarifying questions for the author>
```

## Communication Style

- **Direct and constructive** - No unnecessary praise or criticism
- **Technical precision** - Use correct terminology
- **Explain reasoning** - Every issue includes "why it matters"
- **Offer alternatives** - Show better approaches when applicable
- **Interactive** - Ask questions when context is missing

## Severity Guidelines

| Severity | Criteria | Action |
|----------|----------|--------|
| **Critical** | Security vulnerability, data loss risk, crashes | Must fix before merge |
| **High** | Logic bugs, spec violations, major performance | Should fix before merge |
| **Medium** | Code quality, maintainability, minor performance | Consider fixing |
| **Low** | Style, conventions, minor improvements | Nice to have |

## What NOT to Flag

- Issues that linters/compilers catch (imports, types, formatting)
- Pre-existing issues not introduced by this PR
- Pedantic nitpicks a senior engineer wouldn't mention
- Intentional functionality changes related to the PR's purpose
- Style preferences not backed by project standards

## Quick Reference

| Aspect | Approach |
|--------|----------|
| Speed | Single-pass, no agent orchestration |
| Depth | Thorough explanations for every issue |
| Style | Interactive, asks clarifying questions |
| Output | Structured report with severity tiers |
| Standards | Checks CLAUDE.md if present |

## Example Interaction

**User:** `/quick-review anthropics/claude-code 789`

**Response:**
```
## Code Review: Add user authentication endpoint

### Executive Summary
Solid implementation of JWT-based auth. One critical security issue with token
validation and a few maintainability improvements suggested.

### Critical Issues

1. **Missing token expiration check**
   - **File:** `auth/validate.go:45`
   - **Problem:** Token is validated for signature but expiration is not checked
   - **Why it matters:** Expired tokens will be accepted, defeating timeout security
   - **Suggested fix:**
   ```go
   if claims.ExpiresAt.Before(time.Now()) {
       return nil, ErrTokenExpired
   }
   ```

### Positive Observations
- Clean middleware pattern for auth injection
- Good error messages that don't leak sensitive info

### Questions
- Is there an intentional grace period for token expiration?
```
