package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	func() {
		defer dataMutex.RUnlock()
		dataMutex.RLock()

		assert.Equal(t, 0, len(data))
	}()

	handler := func() {}

	Add(AddOptions{
		Name:    "test",
		Handler: handler,
	})

	func() {
		defer dataMutex.RUnlock()
		dataMutex.RLock()

		assert.Equal(t, 1, len(data))
		assert.Equal(t, "1", data[0].bitmaskValue.String())
		assert.Equal(t, false, data[0].shouldWaitTillDone)
		assert.Equal(t, HandlerOrder(0), data[0].order)
		assert.Equal(t, false, data[0].isRunning)
		assert.Equal(t, false, data[0].isRunPending)
		assert.Equal(t, 0, data[0].donePendingCount)
	}()

	Add(AddOptions{
		Name:    "test2",
		Handler: handler,
	})

	func() {
		defer dataMutex.RUnlock()
		dataMutex.RLock()

		assert.Equal(t, 1, len(data))
		assert.Equal(t, "11", data[0].bitmaskValue.String())
	}()

	handler = func() {}

	Add(AddOptions{
		Name:    "test",
		Handler: handler,
	})

	func() {
		defer dataMutex.RUnlock()
		dataMutex.RLock()

		assert.Equal(t, 2, len(data))
		assert.Equal(t, "11", data[0].bitmaskValue.String())
		assert.Equal(t, "1", data[1].bitmaskValue.String())
	}()
}
