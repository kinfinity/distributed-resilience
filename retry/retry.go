package retry

import (
	"errors"
	"math/rand"
	"reflect"
	"sync/atomic"
	"time"
)

// BackOff
type BackOff struct {
	delay       time.Duration
	exponential bool
	jitter      bool
	randomness  *rand.Rand
}

// BackOff
func NewBackOff(Delay time.Duration, exponential bool, jitter bool) *BackOff {
	return &BackOff{
		delay:       Delay,
		exponential: exponential,
		jitter:      jitter,
		randomness:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Triggers time.sleep() based of current number of retries & BackOff configuration - regular | jitter | exponential | exponential jitter
func (b *BackOff) Delay(attemptIteration int32) {
	if !b.jitter && !b.exponential {
		time.Sleep(time.Duration(attemptIteration) * b.delay)
	}
	if b.jitter && b.exponential {
		sleepValue := uint64(b.delay.Seconds())
		sleepTime := (time.Duration(attemptIteration) * time.Duration(sleepValue)) + time.Duration(float64(b.randomness.Int63())/(1<<63))
		time.Sleep(sleepTime)
	}
	if b.jitter {
		//
		// scale down the 63-bit integer value to a floating-point number between 0 and 1.0.
		// multiply by sleep duration in seconds to get the actual delay in milliseconds.
		sleepTime := (time.Duration(attemptIteration) * b.delay) + time.Duration(float64(b.randomness.Int63())/(1<<63)) // 1 << 63 - bitwise left shift operation (2^63)
		time.Sleep(sleepTime)
	}
	if b.exponential {
		sleepValue := uint64(b.delay.Seconds())
		time.Sleep(time.Duration(attemptIteration) * time.Duration(atomic.AddUint64(&sleepValue, 1)))
	}
}

// Retry
type Retry struct {
	attempts    *atomic.Int32
	maxAttempts int32
	onErrors    []error
	backoff     *BackOff
}

// New Retry  instance with given delay, max attempts and a list of errors to watch for to trigger a retry. If no error is matched panic
func NewWithBackOff(maxAttempts int32, delay time.Duration, watchErrors []error, backOff *BackOff) *Retry {
	return &Retry{
		attempts:    &atomic.Int32{},
		maxAttempts: maxAttempts,
		onErrors:    watchErrors,
		backoff:     backOff,
	}
}

func New(maxAttempts int32, delay time.Duration, watchErrors []error) *Retry {
	return &Retry{
		attempts:    &atomic.Int32{},
		maxAttempts: maxAttempts,
		onErrors:    watchErrors,
	}
}

// Retry execution of fn until success, max Attempts or panic
func (r *Retry) Do(fn func() error) error {
	//Always start at count 0
	r.attempts.Store(0)
	for {
		// watched errors triggers a new execution iteration
		err := fn()
		// increment attempt
		r.attempts.Add(1)
		if err == nil {
			// execution completed successfully
			return nil
		}

		if !IsErrorWatchedByRetry(err, r) {
			panic(err)
		}

		if r.attempts.Load() == r.maxAttempts {
			return errors.New("maximum attempts reached")
		}

		//  delay till next iteration
		if r.backoff != nil {
			r.backoff.Delay(r.attempts.Load())
		}

	}
}

// check if the given error is in the list of watched errors.
// Used to catch mostly transient errors for retry scenario and
// returns true when the error matches one of the provided errors by Type and Messages
func IsErrorWatchedByRetry(err error, r *Retry) bool {
	for _, watchedError := range r.onErrors {
		// match by type and messages
		if reflect.TypeOf(err) == reflect.TypeOf(watchedError) && err.Error() == watchedError.Error() {
			// One of the errors we are watching for
			return true
		}
	}
	return false
}
