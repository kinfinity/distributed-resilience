package timeout

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock function
type MockFunction struct {
	mock.Mock
}

func (m *MockFunction) Execute() error {
	args := m.Called()
	return args.Error(0)
}

func simulateTimeOut(completionChan chan<- bool) {
	time.Sleep(10 * time.Second)
	completionChan <- true
}

// TestTimeOutWait
func TestTimeOutWait(t *testing.T) {
	//
	compChan := make(chan bool)
	go func() { simulateTimeOut(compChan) }()

	t.Run("TimeOut", func(t *testing.T) {
		time_delay := time.Duration(8 * time.Second)
		mockFn := new(MockFunction)
		mockFn.On("Fallback").Return(nil)

		timeout := NewTimeOut(time_delay)
		response := timeout.Watch(compChan)

		assert.Empty(t, response.result)
	})

	t.Run("Completion", func(t *testing.T) {
		time_delay := time.Duration(12 * time.Second)
		mockFn := new(MockFunction)
		mockFn.On("Fallback").Return(nil)

		timeout := NewTimeOut(time_delay)
		response := timeout.Watch(compChan)

		assert.NotEmpty(t, response.result)
	})

}

// TestTimeOutWait
func TestTimeOutWaitFallback(t *testing.T) {
	//
	time_delay := time.Duration(4 * time.Second)
	compChan := make(chan bool)
	go func() { simulateTimeOut(compChan) }()

	t.Run("FallbackExecutes", func(t *testing.T) {
		mockFn := new(MockFunction)
		mockFn.On("Execute").Return(nil).Times(int(1))

		timeout := NewTimeOutWithFallback(time_delay, mockFn.Execute)
		response := timeout.Watch(compChan)

		assert.Empty(t, response.result)
		mockFn.AssertNumberOfCalls(t, "Execute", int(1))
	})

}
