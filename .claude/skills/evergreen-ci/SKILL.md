---
name: evergreen-cicd
description: Diagnose and resolve Evergreen CI/CD failures. Task priority ordering, MCP tool usage, compilation/test failure analysis, and deduplication strategies. Use when investigating CI failures, flaky tests, or build errors.
---

# Evergreen CI/CD Failure Resolution Guide

This guide provides comprehensive troubleshooting steps for Evergreen CI failures in the MongoDB Atlas (MMS) repository, with specific guidance for AI agents using MCP tools.

## MCP Tool Usage Strategy

**ALWAYS use MCP tools first** for Evergreen failure analysis:

1. Use `list_user_recent_patches_evergreen-mcp-server` to get build patch status
2. Use `get_patch_failed_jobs_evergreen-mcp-server` to analyze failures
3. Use `get_task_logs_evergreen-mcp-server` only when job name and general failure are insufficient

IMPORTANT: If job is failing on master please do not attempt to fix it.

## Task Pattern Quick Reference

| Task Pattern                     | Issue Type                | MCP Analysis                           | Primary Action                                     |
| -------------------------------- | ------------------------- | -------------------------------------- | -------------------------------------------------- |
| `COMPILE_*` / `COMPILE_CLIENT_*` | Compilation failures      | Use MCP, check logs if needed          | **PRIORITY 1**: Fix syntax/dependency errors first |
| `BLOCK_COMMIT_TASK`              | Code health/linting       | Use MCP, failure details sufficient    | Fix linting/formatting issues                      |
| `LINT_PULL_REQUEST_TITLE`        | PR title format           | Use MCP, no logs needed                | Fix PR title with Jira ticket prefix               |
| `UNIT_*` / `*_UNIT_TESTS`        | Unit test failures        | Use MCP, check logs for test details   | Run locally and fix broken functionality           |
| `INT_*` / `*_INT_TESTS`          | Integration test failures | Use MCP, check logs for service issues | Check service dependencies and config              |
| `E2E_*` (non-Cypress)            | End-to-end test failures  | Use MCP, check logs for test failures  | Investigate test logic and timing                  |
| `E2E_CYPRESS_*` / `CYPRESS_*`    | Cypress E2E tests         | **DO NOT FIX** - Known flaky           | Ask user if action needed                          |
| `OAS_*` / `FOAS_*`               | OpenAPI validation        | Use MCP, agent can run rebuilds        | **CROSS-CUTTING**: Run OpenAPI rebuild             |

## Failure Prioritization Rules

### 1. **COMPILATION FAILURES FIRST** (Highest Priority)

- **Tasks**: `COMPILE_BAZEL`, `COMPILE_CLIENT_BAZEL`
- **Strategy**: Fix compilation issues before investigating other failures
- **Rationale**: Compilation failures often cause cascading test failures
- **MCP Usage**: Use MCP tools, check logs only if error details insufficient

### 2. **Code Health Issues** (High Priority)

- **Tasks**: `BLOCK_COMMIT_TASK`, `LINT_PULL_REQUEST_TITLE`, `UNIT_JAVA_CHECKSTYLE`
- **Strategy**: Fix linting, formatting, and code quality issues
- **MCP Usage**: Use MCP tools, failure details usually sufficient

### 3. **Unit Test Failures** (Medium Priority)

- **Tasks**: `UNIT_*`, `*_UNIT_TESTS`, `UNIT_JS_*`, `UNIT_PYTHON_*`
- **Strategy**: Run tests locally, fix broken functionality
- **MCP Usage**: Use MCP tools, check logs for specific test failure details

### 4. **Integration Test Failures** (Medium Priority)

- **Tasks**: `INT_*`, `*_INT_TESTS`, `INT_JS_*`
- **Strategy**: Check service dependencies, database connections
- **MCP Usage**: Use MCP tools, check logs for service startup/connection issues

### 5. **End-to-End Test Failures** (Lowest Priority/Ignore)

- **Tasks**: `E2E_Local_*`, `E2E_NDS_*` (non-Cypress)
- **Strategy**: Investigate test logic, timing issues, environment setup
- **MCP Usage**: Use MCP tools, check logs for test execution details

## Compilation Failures (`COMPILE_*`, `COMPILE_CLIENT_*`)

**Issue**: Code fails to compile due to syntax errors or missing dependencies

**MCP Strategy**:

1. Use `get_patch_failed_jobs_evergreen-mcp-server` to get failure details
2. Check logs - Do not attempt to compile codebase localy as this might take time.
3. Focus on first compilation error - fixing it may resolve others

**Java Compilation** (`COMPILE_BAZEL`):

```bash
# Build specific component
bazel build //server/src/main/com/path/to/component

# Build all server code (only when really needed)
bazel build //server/src/main/...
```

**TypeScript/Frontend Compilation** (`COMPILE_CLIENT_BAZEL`):

```bash
# Build client assets
bazel build //client:assets

# Typecheck specific package
bazel test //client/packages/common:typecheck_test --test_output=all

# Global typecheck
pnpm run typecheck
```

## Unit Test Failures (`UNIT_*`, `*_UNIT_TESTS`)

**Issue**: Unit tests failing due to code changes
**Priority**: Medium - Fix after compilation issues

**MCP Strategy**:

1. Use `get_patch_failed_jobs_evergreen-mcp-server` to identify failed tests
2. Use `get_task_logs_evergreen-mcp-server` to get specific test failure details
3. Look for test method names and assertion failures in logs

**Task Patterns**:

- `UNIT_JS_*` - JavaScript/TypeScript unit tests
- `UNIT_PYTHON_*` - Python unit tests
- `UNIT_JAVA_*` - Java unit tests
- `UNIT_GO_*` - Go unit tests

**Java Unit Tests**:

```bash
# Run specific test class
bazel test //server/src/unit/com/path/to/test:TestClassName --test_output=streamed

# Run specific test method
bazel test //server/src/unit/com/path/to/test:TestClassName --test_filter=testMethodName$

# Run all unit tests in package
bazel test //server/src/unit/com/path/to/package:all
```

**TypeScript Unit Tests**:

```bash
# Run specific test file
pnpm test:unit client/packages/common/utils/envUtils.test.ts

# Run all tests in package
bazel test //client/packages/common:js_unit_test --test_output=all

# Watch mode for development
pnpm test:unit --watch
```

**Debugging Steps**:

1. Read test failure messages from MCP logs carefully
2. Check if code changes broke test assumptions
3. Update tests if behavior change is intentional
4. Fix code if tests reveal actual bugs

## Integration Test Failures (`INT_*`, `*_INT_TESTS`)

**Issue**: Integration tests failing due to service/database issues
**Priority**: Medium - Fix after compilation and unit test issues

**MCP Strategy**:

1. Use `get_patch_failed_jobs_evergreen-mcp-server` to identify failed integration tests
2. Use `get_task_logs_evergreen-mcp-server` to check for service startup failures
3. Look for database connection errors, timeout issues, and service dependencies

Important: Log file is very large. Please fetch log from the bottom and stop after identifying first issue.
Important: tests require running database locally (`./scripts/mongodb-start-standalone.bash -p 26000`)

**Task Patterns**:

- `INT_JS_*` - JavaScript integration tests
- `INT_JAVA_EXTERNAL_*` - External Java integration tests (AWS, Azure, GCP)
- `INT_TEST_*` - Specific integration test suites

**Java Integration Tests**:

```bash
# Run specific integration test
bazel test //server/src/test/com/path/to/test:TestClassName --test_output=streamed

# Run with debug output and no cache
bazel test //server/src/test/com/path/to/test:TestClassName --test_output=all --nocache_test_results

# Run with increased timeout
bazel test //server/src/test/com/path/to/test:TestClassName --test_timeout=3000
```

## End-to-End Test Failures

### Cypress Tests (`E2E_CYPRESS_*`, `CYPRESS_*`) - **DO NOT FIX - Leave to user**

**Issue**: Cypress E2E tests are inherently flaky and non-blocking
**Priority**: **EXCLUDED** - Do not attempt to fix

**MCP Strategy**:

1. Use `get_patch_failed_jobs_evergreen-mcp-server` to identify Cypress failures
2. **DO NOT** fetch logs or attempt fixes
3. **Always ask user** if action should be taken

**Important**: Cypress tests are known to be flaky due to:

- Browser timing issues
- Network latency
- UI rendering delays
- External service dependencies
- Selenium WebDriver instability

**Agent Response**:

```
"Cypress E2E tests have failed. These tests are known to be flaky and non-blocking.
Should I investigate these failures or focus on other issues?"
```

### Non-Cypress E2E Tests (`E2E_Local_*`, `E2E_NDS_*`) - **INVESTIGATE**

**Issue**: Backend E2E tests failing due to test logic or environment issues
**Priority**: Lower - Fix after compilation, unit, and integration tests

**MCP Strategy**:

1. Use `get_patch_failed_jobs_evergreen-mcp-server` to identify failed E2E tests
2. Use `get_task_logs_evergreen-mcp-server` to check for test execution details
3. Look for cucumber test failures, environment setup issues

**Task Patterns**:

- `E2E_Local_Core_*` - Core functionality E2E tests
- `E2E_Local_ATM_*` - Automation-related E2E tests
- `E2E_Local_Monitoring_*` - Monitoring E2E tests
- `E2E_NDS_*` - NDS (Atlas) E2E tests

## Code Health Failures (`BLOCK_COMMIT_TASK`, `LINT_*`)

**Issue**: Code quality checks failing
**Priority**: High - Fix early to prevent blocking commits

**MCP Strategy**:

1. Use `get_patch_failed_jobs_evergreen-mcp-server` to get failure details
2. Failure details usually sufficient - logs rarely needed
3. Focus on specific linting/formatting violations

**Task Patterns**:

- `BLOCK_COMMIT_TASK` - General code health blocking commits
- `LINT_PULL_REQUEST_TITLE` - PR title format validation
- `UNIT_JAVA_CHECKSTYLE` - Java code style validation

**Actions**:

```bash
# Run all code health checks
pnpm run code_health

# Run specific checks
pnpm run code_health:lint
pnpm run code_health:deadcode

# Fix formatting
pnpm run format
bazel run @aspect_rules_format//format

# Lint specific package
bazel test //client/packages/common:lint --test_output=all
```

## Cross-Cutting Concerns - **AGENTS CAN HANDLE**

### OpenAPI/FOAS Tasks (`OAS_*`, `FOAS_*`)

**Issue**: OpenAPI specification validation or generation failures
**Agent Action**: **Can run OpenAPI rebuilds and breaking change detection**

**Task Patterns**:

- `OAS_VALIDATE_*` - OpenAPI validation
- `FOAS_NON_BREAKING_*` - FOAS breaking change detection
- `OPENAPI_GENERATE_SPECS` - OpenAPI spec generation

## Exclusions - **DO NOT FIX**

### Known Flaky Tests

- **Cypress E2E tests** (`E2E_CYPRESS_*`, `CYPRESS_*`) - Always flaky, non-blocking
- **Tests with `foliage_check_task_only` tag** - Monitoring only, not blocking
- **Nightly/cron tests** - Run on schedule, not patch-blocking

### Infrastructure Issues

- **Database connection timeouts** in CI environment - Escalate to Evergreen team
- **Resource exhaustion** - Escalate to infrastructure team
- **External service outages** - Wait for service restoration

## Deduplication Strategy

### 1. **Fix Compilation First**

- Compilation failures often cause cascading test failures
- Fix `COMPILE_*` tasks before investigating test failures
- Re-evaluate other failures after compilation fixes

### 2. **Group Related Failures**

- Multiple unit test failures in same package → likely single root cause
- Integration test failures + unit test failures → check for shared dependencies
- E2E failures + integration failures → check for service configuration issues

### 3. **Wait for Master Fixes**

- If failures appear unrelated to patch changes, wait for fixes to land on master
- Use MCP tools to monitor patch status and auto-rebase when master is fixed
- Avoid fixing issues that are already being addressed upstream

## Emergency Procedures

### Getting Help

- **Task Ownership**: Check `assigned_to_jira_team_*` tags for responsible team
- **Foliage Integration**: Use team mapping for automated issue assignment
- **Escalation Path**: Infrastructure → Team Lead → Engineering Manager

---

# CI/CD Failure Investigation Workflow

Systematic workflow for investigating and fixing CI/CD test failures in Evergreen.

**Repository**: MongoDB Atlas (10gen/mms)
**Tools**: Evergreen MCP (`list_user_recent_patches`, `get_patch_failed_jobs`, `get_task_logs`, `get_task_test_results`), `codebase-retrieval`, Bazel, pnpm
**Critical**: NEVER investigate or commit without user approval, NEVER fix master branch failures

## Workflow Phases

### Phase 1: Retrieve Failing Tasks

Use `list_user_recent_patches_evergreen` (limit: 5-10), identify failed patch, get failed jobs with `get_patch_failed_jobs_evergreen`, extract task names/IDs, variants, test counts, errors.

### Phase 2: Analyze and Prioritize Failures

**CRITICAL**: Do NOT investigate yet. Only analyze and categorize.

**Categorize** (using task patterns above): PRIORITY 1: Compilation (`COMPILE_*`) → PRIORITY 2: Code health (`LINT_*`) → PRIORITY 3: Unit (`UNIT_*`) → PRIORITY 4: Integration (`INT_*`) → EXCLUDE: Cypress E2E (flaky, ask user) → LOWER: Non-Cypress E2E

**Determine**: Code relationship (compare modified files), flaky tests (no related changes), group related failures (same package).

**Analysis Format**:

```markdown
## Evergreen Failure Analysis - Patch #{patch_id}

### PRIORITY 1: Compilation Failures

- Task: COMPILE_BAZEL | Variant: ubuntu2204 | Status: Failed
  Error: "Cannot find symbol: UserService.validateUser()"
  Related: YES (modified UserService.java)

### PRIORITY 2: Code Health

- Task: BLOCK_COMMIT_TASK | Status: Failed
  Error: "Linting errors in 3 files"
  Related: YES (modified files need formatting)

### PRIORITY 3: Unit Test Failures

- Task: UNIT_JAVA_CORE | Failed Tests: 2
  Tests: UserServiceUnitTests.testValidation, AuthServiceUnitTests.testPermissions
  Related: YES (modified UserService.java, AuthService.java)

### EXCLUDED: Cypress E2E

- Task: E2E_CYPRESS_SMOKE | Status: Failed
  Note: Known flaky, non-blocking. Should I investigate?

RECOMMENDED ORDER: Fix compilation → code health → unit tests → (await user decision on Cypress)
```

### Phase 3: Get User Approval on Investigation Scope

Present analysis, recommend compilation/code health first, ask about Cypress E2E (known flaky). **Wait for explicit approval** before proceeding.

### Phase 4: Investigate Failures Using Evergreen Logs

For approved failures: Fetch logs (`get_task_logs_evergreen`, filter_errors=true, max_lines=500-1000), get test results (`get_task_test_results_evergreen`, failed_only=true), analyze errors, use `codebase-retrieval` for context.

**Important**: Integration test logs are large - fetch from bottom, stop after identifying issue.

### Phase 5: Run Tests Locally if Needed

If Evergreen logs insufficient, run specific failing tests locally:

**Java Tests**:

```bash
# Specific test method
bazel test //server/src/unit/com/path/to/package:TestClassName --test_filter=testMethodName$ --test_output=streamed

# Specific test class (unit tests)
bazel test //server/src/unit/com/path/to/package:TestClassName --test_output=streamed

# Specific test class (integration tests)
bazel test //server/src/test/com/path/to/package:TestClassName --test_output=streamed
```

**TypeScript Tests**:

```bash
# Specific test file (pnpm - recommended)
pnpm test:unit client/packages/common/utils/file.test.ts

# Specific test file (Bazel)
bazel test //client/packages/common:js_unit_test --test_filter="**/*/utils/file.test.ts" --test_output=all

# Package tests
bazel test //client/packages/common:js_unit_test --test_output=all
```

**Compilation**:

```bash
# Java compilation (specific component)
bazel build //server/src/main/com/path/to/component

# TypeScript typecheck (specific package)
bazel test //client/packages/common:typecheck_test --test_output=all

# TypeScript typecheck (all packages)
pnpm run typecheck
```

**Important**: Do NOT attempt to compile entire codebase locally - use Evergreen logs for compilation errors.

### Phase 6: Implement Fixes

Use `codebase-retrieval` for context, make minimal changes, update tests if behavior changed, fix bugs, handle downstream impacts (callers, types, interfaces).

### Phase 7: Verify Fixes Locally

**CRITICAL**: Confirm fixes work before committing. Run previously failing tests, related tests (no regressions), code quality checks (`pnpm run code_health`). DO NOT run `bazel run @aspect_rules_format//format` - formatting is handled automatically by precommit hooks. Mark tasks COMPLETE.

### Phase 8: Commit and Push

**CRITICAL**: NEVER commit/push without approval. Request approval with summary (fixes, files, test results), then commit with descriptive message (include ticket/patch), push to branch, confirm success.

## Key Principles

1. **NEVER** investigate or commit/push without user approval, **NEVER** fix master branch failures
2. **Always** use Evergreen MCP tools first, prioritize compilation, ask about Cypress (flaky)
3. **Follow** TDD (fail → fix → pass), make minimal focused changes
4. **Quality**: All tests pass locally, follow style guidelines, no unrelated changes

## Error Handling

**API**: Rate limiting (wait/retry), auth (verify MCP), patch not found (check ID)
**Investigation**: Insufficient logs (run locally), persistent failures (escalate), master failures (notify, don't fix)
**Local Tests**: Build failures (check deps), test environment (verify DB for integration), timeouts (increase flag)

## Quality Checklist

- [ ] Failures analyzed, user approved scope, prioritized correctly (compilation → code health → unit → integration)
- [ ] Cypress E2E consulted, Evergreen logs analyzed, local tests only if needed
- [ ] Fixes minimal/focused, failing tests now pass, related tests pass (no regressions)
- [ ] Code quality checks pass
- [ ] User approved commit/push

## Best Practices

**Before**: Use Evergreen MCP first, prioritize by type, group related, ask about Cypress
**During**: Use codebase-retrieval, minimal changes, TDD (fail → fix → pass), run tests frequently
**After**: Verify locally, related tests (no regressions), code quality checks
