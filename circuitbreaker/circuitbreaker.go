package circuitbreaker

// Models Circuit Breaker and Handles lifecycle up to re connection
import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/kinfinity/distributed-resilience/retry"

)

// Circuit Breaker
type CircuitBreaker struct {
	recoveryTime   time.Duration
	recovered      bool
	maxFailures    int32
	failCount      *atomic.Int32
	successCount   *atomic.Int32 // what if a certain number of concurrent successful executions are required?
	maxSuccess     int32         // ?
	currentState   breakerState
	lastChangeTime int64
	ctx            context.Context // Set Life Cycle Expiry Timer ?
	completion     chan bool
	lastError      error
	fallback       func(error) // Default FallBack  Functionality - Kicks in when state switches to Open
	retryPin       *retry.Retry
}

// BreakerOptions
type BreakerOptions struct {
	RecoveryTime time.Duration
	MaxFailures  int32
	MaxSuccess   int32
	ctx          context.Context
}

// State
type breakerState uint8

const (
	closed   breakerState = iota // Forward Requests
	halfOpen                     // Error detected w no ful fail over yet
	open                         // Requests fail immediately - Not forwarded
)

// New Breaker Initialized in close state
func NewCircuitBreaker(bo BreakerOptions) *CircuitBreaker {
	return &CircuitBreaker{
		recoveryTime: bo.RecoveryTime,
		maxFailures:  bo.MaxFailures,
		maxSuccess:   bo.MaxSuccess,
		failCount:    &atomic.Int32{},
		successCount: &atomic.Int32{},
		currentState: closed,
		ctx:          bo.ctx, // Need context deadline for lifecycle  of circuit breaker?
		fallback: func(err error) {
			//
		},
		completion: make(chan bool, 1), // channel to watch for completion of circuit Execution
		recovered:  false,              // Flag indicating whether the Circuit has been recovered or not
	}
}

// Execute and ensure the  function is executed within a circuit breaker w context
func (cb *CircuitBreaker) Do(f func() error) error {
	cb.failCount.Store(0)
	cb.successCount.Store(0)

	select {
	case <-cb.ctx.Done():
		// Lifecycle time is up
		return errors.New("circuit breaker expired")
	default:
		// Continue with the circuit breaker logic
	}

	// kick out
	if cb.currentState == closed && cb.failCount.Load() > cb.maxFailures {
		return errors.New("circuit open")
	}

	// completion
	if cb.currentState == closed && cb.successCount.Load() >= cb.maxSuccess {
		cb.completion <- true
		return nil
	}

	switch cb.currentState {
	case closed:
		err := f()
		if err != nil {
			// Execution of f  failed mark as failure
			cb.recordFailure()
			cb.lastError = err
			// check if max retries has been hit & retry
			if cb.failCount.Load() == cb.maxFailures {
				// Move to half open
				cb.switchState()
			}
			cb.Do(f) // final Objective is to get maxSuccess on successCount ? base scenario 1
		}
		// Execution of f  succeeded
		cb.recordSuccess()
		// No need for circuit breaker Lifecyle continuation
		// cleanup
		return nil
	case halfOpen:
		err := f()
		if err == nil {
			// half-open, allowed one call to succeed
			cb.recordSuccess()
			// if it is back to normal operation with normal consecutive successes then move to close state
			if cb.successCount.Load() == cb.maxSuccess {
				cb.recovered = true
				cb.switchState()
			}
		} else {
			// still in half-open, treat as a failure
			cb.recordFailure()
			cb.retryPin = retry.NewWithBackOff(
				int32(2),
				cb.recoveryTime,
				[]error{cb.lastError},
				retry.NewBackOff(100*time.Millisecond, true, true),
			)
			cb.switchState()
		}
		cb.lastError = err
		return err
	case open:
		// evaluate and decide whether or not to switch back to closed state
		cb.fallback(cb.lastError)
	default:
		return errors.New("unexpected state")
	}

	return cb.Do(f)
}

// Change the state of the circuit breaker based on current configuration
func (cb *CircuitBreaker) switchState() {
	if !isState(cb.currentState) {
		panic("unexpected state")
	}
	log.Println("Changing Circuit Breaker state from", cb.currentState)
	switch cb.currentState {
	case open:
		cb.currentState = halfOpen
	case halfOpen:
		if cb.recovered {
			cb.currentState = closed
		} else {
			cb.currentState = open
		}
	default:
		cb.currentState++
	}
	cb.lastChangeTime = time.Now().UnixNano() / int64(time.Millisecond)
	log.Println("Changing Circuit Breaker state to", cb.currentState)
}

func isState(state breakerState) bool {
	return state < open
}

// Record a success - resets the fail count
func (cb *CircuitBreaker) recordSuccess() {
	cb.successCount.Add(1)
	cb.failCount.Store(0)
}

// Record a failure - increments the fail count, reset successCount
func (cb *CircuitBreaker) recordFailure() {
	cb.failCount.Add(1)
	cb.successCount.Store(0)
}

// Checks whether enough time has passed since the last state change
// func (cb *CircuitBreaker) isStale() bool {
// 	now := time.Now().UnixNano()
// 	delta := float64(now-cb.lastChangeTime) / 1e9
// 	return delta >= float64(cb.recoveryTime)
// }

// Duration in Context for Lifecyle Time
// post execute / free resources
