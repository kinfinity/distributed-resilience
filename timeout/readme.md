# Timeout Pattern in Go

The Timeout pattern is used to limit the execution time of a function or operation. It ensures that the operation completes within a specified time duration, and if it exceeds that duration, it either returns a default value or invokes a fallback function.

## Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/kinfinity/distributed-resilience/timeout"
)

func main() {
    // Create a Timeout instance with a fallback function
    timeout := timeout.NewTimeOutWithFallback(5 * time.Second, func() error {
        // Custom fallback logic here
        return nil
    })

    // Watch for timeout
    result := timeout.Watch(executionCompletionChan)
    if result.result {
    // Operation completed within the timeout duration
    } else {
    // Operation timed out
    }
}
```

# **References**

- [codecentric resilience-design-patterns](https://www.codecentric.de/wissens-hub/blog/resilience-design-patterns-retry-fallback-timeout-circuit-breaker)
