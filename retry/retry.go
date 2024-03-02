package retry

import (
	"errors"
	"reflect"
	"sync/atomic"
	"time"
)

// Retry
type Retry struct {
	delay       time.Duration
	attempts    *atomic.Int32
	maxAttempts int32
	onErrors    []error
}

// New Retry  instance with given delay, max attempts and a list of errors to watch for to trigger a retry. If no error is matched panic
func New(maxAttempts int32, delay time.Duration, watchErrors []error) *Retry {
	return &Retry{
		delay:       delay,
		attempts:    &atomic.Int32{},
		maxAttempts: maxAttempts,
		onErrors:    watchErrors,
	}
}

// Retry execution of fn until success, max Attempts or panic
func (r *Retry) Do(fn func() error) error {
	for {
		// watched errors triggers a new execution iteration
		err := fn()
		if err == nil {
			// execution completed successfully
			return nil
		}

		if !r.isWatchedError(err) {
			panic(err)
		}

		if r.attempts.Load() == r.maxAttempts {
			return errors.New("maximum attempts reached")
		}

		// increment attempt and delay till next iteration
		r.attempts.Add(1)
		time.Sleep(r.delay)
	}
}

// check if the given error is in the list of watched errors. Used to catch mostly transient errors for retry scenario
func (r *Retry) isWatchedError(err error) bool {
	for _, error := range r.onErrors {
		if reflect.TypeOf(error) == reflect.TypeOf(&err) {
			//one of the errors we watching for
			return true
		}
	}
	return false
}
