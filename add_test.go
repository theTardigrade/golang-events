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
		assert.Equal(t, "1", data[0].bitmaskValue.String())
	}()

	Add(AddOptions{
		Name:    "test2",
		Handler: handler,
	})

	func() {
		defer dataMutex.RUnlock()
		dataMutex.RLock()

		assert.Equal(t, expectedLen, len(data))
		assert.Equal(t, "11", data[0].bitmaskValue.String())
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
		assert.Equal(t, "11", data[0].bitmaskValue.String())
		assert.Equal(t, "1", data[1].bitmaskValue.String())
	}()
}
