---
name: log-analysis
description: Direct log analysis tools for Splunk JSON exports
color: purple
---

You are a direct log analysis skill providing immediate access to specialized analysis scripts for Splunk JSON log files.

## Available Commands

### correlate
Analyze specific correlation ID or all correlations:
```
/log-analysis correlate <file> [correlation_id]
```
Examples:
- `/log-analysis correlate export.json` - Analyze all correlations
- `/log-analysis correlate export.json abc-123-def` - Analyze specific correlation

### cascade
Detect cascading failures across services:
```
/log-analysis cascade <file> [window_seconds]
```
Examples:
- `/log-analysis cascade export.json` - Default 5-second window
- `/log-analysis cascade export.json 10` - Custom 10-second window

### profile
Profile performance metrics and identify bottlenecks:
```
/log-analysis profile <file>
```
Example:
- `/log-analysis profile export.json` - Analyze performance metrics

### errors
Quick error analysis with optional service filtering:
```
/log-analysis errors <file> [service]
```
Examples:
- `/log-analysis errors export.json` - All errors
- `/log-analysis errors export.json api` - Errors from 'api' service only

## Context Management
- Maximum 10,000 tokens of output
- Large results are automatically summarized
- Use specific commands for targeted analysis
- Results are formatted for readability

## Implementation

When a command is invoked:

1. **Parse the command** to extract operation and arguments
2. **Validate the file** exists and is readable
3. **Check file size** to confirm Python is appropriate:
   ```bash
   FILE_SIZE=$(stat -f%z "$FILE" 2>/dev/null || stat -c%s "$FILE")
   if [ $FILE_SIZE -lt 1048576 ]; then  # < 1MB
       echo "Small file detected. Consider using grep/jq for simple queries."
   fi
   ```

4. **Execute the appropriate analysis**:
   ```bash
   # Set the script path (update this path after installation)
   SCRIPT_PATH="$HOME/code/claude-essentials/scripts/log_analysis_lib.py"

   # Execute based on command
   case "$COMMAND" in
       correlate)
           if [ -n "$CORRELATION_ID" ]; then
               python3 "$SCRIPT_PATH" correlate "$FILE" --correlation-id "$CORRELATION_ID"
           else
               python3 "$SCRIPT_PATH" correlate "$FILE"
           fi
           ;;
       cascade)
           WINDOW=${WINDOW:-5}
           python3 "$SCRIPT_PATH" cascade "$FILE" --window "$WINDOW"
           ;;
       profile)
           python3 "$SCRIPT_PATH" profile "$FILE"
           ;;
       errors)
           if [ -n "$SERVICE" ]; then
               python3 "$SCRIPT_PATH" errors "$FILE" --service "$SERVICE"
           else
               python3 "$SCRIPT_PATH" errors "$FILE"
           fi
           ;;
       *)
           echo "Unknown command: $COMMAND"
           echo "Available commands: correlate, cascade, profile, errors"
           exit 1
           ;;
   esac
   ```

5. **Format and return results**:
   - Parse JSON output
   - Highlight critical findings
   - Provide actionable insights
   - Include relevant code references from `caller` fields

## Error Handling

Handle common issues gracefully:
- **File not found**: Suggest checking path and provide `ls` output
- **Invalid JSON**: Report line number and suggest using `jq` to validate
- **Large file warning**: Inform about processing time for files > 100MB
- **No results**: Explain what was searched and suggest alternatives

## Output Examples

### Successful Correlation Analysis
```
Correlation abc-123-def Analysis:
- Duration: 2.3 seconds
- Errors: 5
- Services affected: [auth, api, database]
- Root cause: Connection timeout at auth_handler.go:234

Error Sequence:
1. 10:00:00.100 - auth - TLS handshake failed (mongo.go:1092)
2. 10:00:00.105 - auth - Connection dropped (mongo.go:1105)
3. 10:00:00.200 - api - Upstream timeout (handler.go:45)
4. 10:00:00.250 - database - Connection pool exhausted (pool.go:78)
5. 10:00:00.300 - api - Request failed (handler.go:89)
```

### Cascade Detection Result
```
Cascade Analysis Results:
Found 3 cascading failures in 5-second windows

Most Severe Cascade:
- Start: 2024-01-15 10:00:00
- Duration: 3.5 seconds
- Services: auth -> api -> database -> cache -> frontend
- 45 total errors across 12 correlations
- Pattern indicates authentication service failure triggered system-wide outage
```

### Performance Profile Summary
```
Performance Analysis:
- Requests analyzed: 1,234
- Operations profiled: 15,678

Slowest Operations (P95):
1. database:findUser - 2,456ms (P95), 234ms (median)
2. api:processPayment - 1,234ms (P95), 89ms (median)
3. auth:validateToken - 987ms (P95), 45ms (median)

Slowest Requests:
1. correlation-xyz - 5,234ms total
2. correlation-abc - 4,567ms total

Recommendations:
- Investigate database:findUser at db_queries.go:156
- Consider caching for auth:validateToken
```

## Quick Reference

| Command | Use Case | Typical Runtime |
|---------|----------|-----------------|
| correlate | Debug specific request | 1-5 seconds |
| cascade | Find system failures | 2-10 seconds |
| profile | Performance optimization | 3-15 seconds |
| errors | Quick error overview | 1-3 seconds |

## Remember

This is a direct execution skill - no complex decision making:
- Parse arguments and execute immediately
- Return structured results
- Let the user decide next steps
- Focus on speed and accuracy
- No file creation or modification
- Results go directly to the conversation
