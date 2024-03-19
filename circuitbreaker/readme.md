# Circuit Breaker Pattern

The Circuit Breaker Pattern is designed to enhance the resilience of a system by monitoring for failures and preventing further requests to a failing service. It consists of three states: **Closed**, **Open**, and **Half-Open**. The circuit breaker transitions between these states based on the observed failures

- **Closed State:** In this state, the circuit breaker allows requests to pass through. It monitors for failures, and if the failure rate exceeds a predefined threshold, it transitions to the **Open State**.
- **Open State:** In this state, the circuit breaker prevents requests from reaching the service, providing a fast-fail mechanism. After a predefined timeout, the circuit breaker transitions to the **Half-Open State** to test if the service has recovered.
- **Half-Open State:** In this state, the circuit breaker allows a limited number of requests to pass through. If these requests succeed, the circuit breaker transitions back to the **Closed State**; otherwise, it returns to the **Open State**.

## Execution

- Starts up in the closed state
- on failure move to half closed
- if it doesn't recovery move to open state
- after recovery duration move back to half closed
- if it recovers move to closed state

### Components

- **Circuit Breaker**: Represents the main circuit breaker object, responsible for managing the state transitions and handling requests based on the current state.
- **BreakerOptions**: Defines the options for configuring the circuit breaker, such as recovery time and maximum allowed failures.
- **breakerState**: Enumerates the possible states of the circuit breaker, including closed, half-open, and open.
- **Retry**: Utilizes the retry package for handling retry logic when transitioning from the half-open state back to the closed state.

## Usage

### Creating a Circuit Breaker

To create a new circuit breaker, use the `NewCircuitBreaker` function:

```go
package main

import (
    "context"
    "errors"
    "log"
    "time"

    "github.com/example/circuitbreaker"
)

func main() {
    options := circuitbreaker.BreakerOptions{
        RecoveryTime: time.Minute,
        MaxFailures:  5,
        // ctx:       Provide context for lifecycle time (optional)
    }

    cb := circuitbreaker.NewCircuitBreaker(options)

    err := cb.Do(func() error {
        // Simulate an operation that may fail
        log.Println("Executing function...")
        return errors.New("temporary error")
    })

    if err != nil {
        log.Println("Operation failed:", err)
    } else {
        log.Println("Operation succeeded.")
    }
}
```

# **References:**

- **Michael Nygard**  https://pragprog.com/titles/mnee2/release-it-second-edition/ Second Edition - Stability Patterns => Circuit Breaker
- https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
