---
description: Code review for Go
---

Review this Go code for:

1. **Error Handling**
   - Are all errors handled properly?
   - Are errors wrapped with context?
   - Are there any ignored errors?

2. **Concurrency**
   - Are there potential race conditions?
   - Is proper synchronization used?
   - Are goroutines properly managed?

3. **Resource Management**
   - Are all resources closed (files, connections, etc.)?
   - Is `defer` used correctly?
   - Are there potential memory leaks?

4. **Go Idioms**
   - Does the code follow Go conventions?
   - Are interfaces used appropriately?
   - Is the code readable and maintainable?

5. **Testing**
   - Is the code testable?
   - What tests should be added?
   - Are edge cases covered?

6. **Security**
   - Is input validated?
   - Are there SQL injection risks?
   - Is sensitive data handled properly?

