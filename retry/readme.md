# Retry Pattern

The retry package provides functionality to execute a function with retry logic. This pattern is particularly useful for handling transient errors by retrying the execution of a function until it succeeds, reaches the maximum number of attempts, or encounters an unrecoverable error.

## BackOff

The BackOff type allows for configuring the delay between retries. It supports different strategies such as regular, exponential, with or without jitter.

## Usage

```
package main

import (
    "errors"
    "fmt"
    "time"

    "github.com/kinfinity/distributed-resilience/retry"
)

func main() {
    // Retry Configuration
    maxAttempts := int32(3)
    delay := 100 * time.Millisecond
    watchErrors := []error{CustomError{"custom error"}, errors.New("temporary error")} // whiteList to trigger retry

    // Retry
    retryPolicy := NewWithBackOff(maxAttempts, delay, watchErrors, NewBackOff(100*time.Millisecond, true, true))


    // Function to execute with retry logic
    myFunction := func() error {
        // - Simulate an operation that may fail
        // network call | database operation
        fmt.Println("Executing function...")
        return errors.New("temporary error")
    }

    // Execute function with retry
    err := retryPolicy.Do(myFunction)
    if err != nil {
        fmt.Println("Function failed after maximum attempts:", err)
    } else {
        fmt.Println("Function executed successfully.")
    }
}
```

# **References**

- https://learn.microsoft.com/en-us/azure/architecture/patterns/retry
