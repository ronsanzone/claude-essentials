
# Workflows
These are commands and patterns to use in our interactions. They should always be prioritized:

## Problem-Solving Together
When you're stuck or confused:
1. **Stop** - Don't spiral into complex solutions
2. **Delegate** - Consider spawning agents for parallel investigation
3. **Step back** - Re-read the requirements
4. **Simplify** - The simple solution is usually correct
5. **Ask** - Use `AskUserQuestion` "I see two approaches: [A] vs [B]. Which do you prefer?"

## Available Specialized Agents
`codebase-analyzer`: A specialized agent for analyzing large amounts of code and returning summarized results. Used for context management while investigating large code bases.
`codebase-locator`: An agent to find code. Used to return code locations without the main context window needing to read hundreds of files.
`codebase-pattern-finder`: An agent to find patterns across files. 
`splunk-analyzer`: Analyze large splunk log downloads. Used to keep the main context window clean while investigating issues.
`web-search-researcher`: Search and summarize topics across the web.
