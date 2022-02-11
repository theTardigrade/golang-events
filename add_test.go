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

	Add(AddOptions{
		Name: "test",
	})

	func() {
		defer dataMutex.RUnlock()
		dataMutex.RLock()

		assert.Equal(t, 1, len(data))
	}()
}
