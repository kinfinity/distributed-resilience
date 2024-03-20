package timeout

import (
	"context"
	"fmt"
	"time"
)

// Timeout
type TimeOut struct {
	duration time.Duration
	fallback func() error
}

// Returns nil if the operation times out
type TimeOutResult struct {
	result bool
}

// TimeOut
func NewTimeOut(Delay time.Duration) *TimeOut {
	return &TimeOut{
		duration: Delay,
	}
}

// TimeOut with options via
// fallback() error
func NewTimeOutWithFallback(Delay time.Duration, options ...interface{}) *TimeOut {

	timeout := &TimeOut{
		duration: Delay,
	}

	for _, option := range options {
		switch opt := option.(type) {
		case func() error:
			timeout.fallback = opt
		default:
			panic(fmt.Sprintf("Unknown option type: %T", opt))
		}
	}

	return timeout
}

// Creates a context with the duration lifetime
// func executes till the end| Context Lifetime expires
func (to *TimeOut) Watch(executionCompletionChan chan bool) *TimeOutResult {
	// Build Context Lifetime
	ctx, cancel := context.WithTimeout(context.Background(), to.duration)
	defer cancel()

	for {
		select {
		case result := <-executionCompletionChan:
			return &TimeOutResult{
				result: result,
			}
		case <-ctx.Done():
			if to.fallback != nil {
				to.fallback() // fallback on timeout
			}
			// Timeout
			return &TimeOutResult{}
		default:
		}
		// random delay before we check conditions again
		// give timeout grace
		time.Sleep(to.duration / 4)
	}

}
