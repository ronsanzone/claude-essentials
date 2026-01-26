---
globs:
  - "**/*.go"
  - "go.mod"
  - "go.sum"
---
# Go Development Standards

These rules apply when working with Go code.

## FORBIDDEN - NEVER DO THESE

- **NO `interface{}`** or **`any{}`** - use concrete types!
- **NO `time.Sleep()`** or busy waits - use channels for synchronization!
- **NO** keeping old and new code together
- **NO** migration functions or compatibility layers
- **NO** versioned function names (`processV2`, `handleNew`)
- **NO** custom error struct hierarchies
- **NO** TODOs in final code

> **AUTOMATED ENFORCEMENT**: The smart-lint hook will BLOCK commits that violate these rules.
> When you see `FORBIDDEN PATTERN`, you MUST fix it immediately!

## Required Standards

- **Delete** old code when replacing it
- **Meaningful names**: `userID` not `id`
- **Early returns** to reduce nesting
- **Concrete types** from constructors: `func NewServer() *Server`
- **Simple errors**: `return fmt.Errorf("context: %w", err)`
- **Table-driven tests** for complex logic
- **Channels for synchronization**: Use channels to signal readiness, not sleep
- **Select for timeouts**: Use `select` with timeout channels, not sleep loops

## Error Handling Pattern

```go
// Good - simple error wrapping
if err != nil {
    return fmt.Errorf("failed to connect: %w", err)
}

// Bad - custom error types for simple cases
type ConnectionError struct {
    Underlying error
    Timestamp  time.Time
}
```

## Concurrency Pattern

```go
// Good - channel for synchronization
ready := make(chan struct{})
go func() {
    // setup work
    close(ready)
}()
<-ready

// Bad - sleep-based synchronization
go func() {
    // setup work
}()
time.Sleep(100 * time.Millisecond)
```

## Testing

- Use table-driven tests for functions with multiple cases
- Name test cases descriptively
- Use `t.Helper()` in test helper functions
- Prefer `testify/assert` or `testify/require` for assertions
