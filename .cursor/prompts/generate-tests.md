---
description: Generate tests for Go code
---

Generate comprehensive tests for this Go code:

## Requirements

1. **Use table-driven tests** with clear test case names
2. **Cover these scenarios:**
   - Happy path (valid input)
   - Edge cases (empty, nil, zero values)
   - Error cases (invalid input, failures)
   - Boundary conditions

3. **Test structure:**
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

4. **Use testify for assertions:**
   - `assert.Equal(t, expected, actual)`
   - `assert.NoError(t, err)`
   - `require.NotNil(t, result)`

5. **Mock dependencies** using interfaces

6. **Add `t.Parallel()`** for independent tests

