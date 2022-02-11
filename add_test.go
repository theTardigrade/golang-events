package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	var expectedLen int

	func() {
		defer dataMutex.RUnlock()
		dataMutex.RLock()

		assert.Equal(t, expectedLen, len(data))
	}()

	handler := func() {}

	Add(AddOptions{
		Name:    "test",
		Handler: handler,
	})

	func() {
		expectedLen++

		defer dataMutex.RUnlock()
		dataMutex.RLock()

		assert.Equal(t, expectedLen, len(data))
	}()

	Add(AddOptions{
		Name:    "test",
		Handler: handler,
	})

	func() {
		defer dataMutex.RUnlock()
		dataMutex.RLock()

		assert.Equal(t, expectedLen, len(data))
	}()

	handler = func() {}

	Add(AddOptions{
		Name:    "test",
		Handler: handler,
	})

	func() {
		expectedLen++

		defer dataMutex.RUnlock()
		dataMutex.RLock()

		assert.Equal(t, expectedLen, len(data))
	}()
}
