package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock function to simulate a function with errors
type MockFunction struct {
	mock.Mock
}

func (m *MockFunction) Execute() error {
	args := m.Called()
	return args.Error(0)
}

// TestRetryDo tests the retry mechanism
func TestRetryDo(t *testing.T) {
	maxAttempts := int32(3)
	delay := 100 * time.Millisecond
	watchErrors := []error{errors.New("temporary error"), errors.New("maximum attempts reached")}

	retry := NewWithBackOff(maxAttempts, delay, watchErrors, NewBackOff(100*time.Millisecond, true, true))

	t.Run("Success", func(t *testing.T) {
		mockFn := new(MockFunction)
		mockFn.On("Execute").Return(nil)

		err := retry.Do(mockFn.Execute)

		assert.NoError(t, err)
		mockFn.AssertCalled(t, "Execute")
	})

	t.Run("MaxAttemptsReached", func(t *testing.T) {
		mockFn := new(MockFunction)
		mockFn.On("Execute").Return(errors.New("temporary error")).Times(int(maxAttempts))

		err := retry.Do(mockFn.Execute)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "maximum attempts reached")
		// mockFn.AssertExpectations(t)
		mockFn.AssertNumberOfCalls(t, "Execute", int(maxAttempts))
	})

	t.Run("PanicOnError", func(t *testing.T) {
		mockFn := new(MockFunction)
		mockFn.On("Execute").Return(errors.New("fatal error"))

		assert.PanicsWithError(t, "fatal error", func() {
			retry.Do(mockFn.Execute)
		})
		mockFn.AssertCalled(t, "Execute")
	})
}

// TestisErrorWatchedByRetry tests the isErrorWatchedByRetry function
func TestIsErrorWatchedByRetry(t *testing.T) {
	maxAttempts := int32(3)
	delay := 100 * time.Millisecond
	watchErrors := []error{CustomError{"custom error"}, errors.New("temporary error")}

	retry := NewWithBackOff(maxAttempts, delay, watchErrors, NewBackOff(100*time.Millisecond, true, true))

	t.Run("WatchedError", func(t *testing.T) {
		customError := CustomError{"custom error"}
		assert.True(t, IsErrorWatchedByRetry(customError, retry))
	})

	t.Run("UnwatchedError", func(t *testing.T) {
		err := errors.New("permanent error")
		assert.False(t, IsErrorWatchedByRetry(err, retry))
	})
}

type CustomError struct {
	Message string
}

func (ce CustomError) Error() string {
	return ce.Message
}
