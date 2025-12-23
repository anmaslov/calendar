---
description: Refactor Go code
---

Refactor this Go code with focus on:

## Clean Code Principles

1. **Single Responsibility**
   - Each function does one thing
   - Each package has one purpose

2. **DRY (Don't Repeat Yourself)**
   - Extract common logic into functions
   - Use interfaces for shared behavior

3. **Readability**
   - Clear variable and function names
   - Small functions (< 30 lines)
   - Minimal nesting (< 3 levels)

## Go-Specific Improvements

1. **Error Handling**
   - Wrap errors with context
   - Handle errors at appropriate level
   - Use custom error types when needed

2. **Interfaces**
   - Define interfaces where they're used
   - Keep interfaces small (1-3 methods)
   - Use interface composition

3. **Concurrency**
   - Use channels for communication
   - Avoid shared state when possible
   - Use context for cancellation

## Output

Provide refactored code with:
- Explanation of changes
- Before/after comparison for significant changes
- Any potential breaking changes noted

