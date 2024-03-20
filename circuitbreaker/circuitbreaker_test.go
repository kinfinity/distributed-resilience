package circuitbreaker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCircuitBreakerDo tests the Do method of CircuitBreaker
func TestCircuitBreakerDo(t *testing.T) {
	maxFailures := int32(3)
	maxSuccess := int32(1)
	recoveryTime := 100 * time.Millisecond

	// Mock context with cancelation after a certain duration
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	options := BreakerOptions{
		RecoveryTime: recoveryTime,
		MaxFailures:  maxFailures,
		MaxSuccess:   maxSuccess,
		ctx:          ctx,
	}

	cb := NewCircuitBreaker(options)

	t.Run("Success", func(t *testing.T) {
		err := cb.Do(func() error {
			return nil
		})

		// Wait for transition
		time.Sleep(2 * time.Second)
		assert.NoError(t, err)
	})

}

// Test state switching mechanism
func TestCircuitBreakerSwitchState(t *testing.T) {
	recoveryTime := 100 * time.Millisecond

	// Mock context with cancellation after a certain duration
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	options := BreakerOptions{
		RecoveryTime: recoveryTime,
		MaxFailures:  3,
		MaxSuccess:   1,
		ctx:          ctx,
	}

	cb := NewCircuitBreaker(options)

	t.Run("SwitchFromClosedToHalfOpen", func(t *testing.T) {
		cb.currentState = closed
		cb.failCount.Store(3)
		cb.switchState()
		assert.Equal(t, halfOpen, cb.currentState)
	})

	t.Run("SwitchFromHalfOpenToClosed", func(t *testing.T) {
		cb.currentState = halfOpen
		cb.successCount.Store(1)
		cb.recovered = true
		cb.switchState()
		assert.Equal(t, closed, cb.currentState)
	})

	t.Run("SwitchFromHalfOpenToOpen", func(t *testing.T) {
		cb.currentState = halfOpen
		cb.failCount.Store(3)
		cb.recovered = false
		cb.switchState()
		assert.Equal(t, open, cb.currentState)
	})

	t.Run("UnexpectedState", func(t *testing.T) {
		cb.currentState = breakerState(100) // Invalid state
		assert.Panics(t, func() { cb.switchState() })
	})
}
