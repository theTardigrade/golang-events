package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitmaskValueFromNames(t *testing.T) {
	defer testMutex.Unlock()
	testMutex.Lock()

	bitmaskGenerator.Clear()

	value := bitmaskValueFromNames("test", nil)

	assert.Equal(t, "1", value.String())

	value = bitmaskValueFromNames("", []string{"test"})

	assert.Equal(t, "1", value.String())

	value = bitmaskValueFromNames("test", []string{"test2"})

	assert.Equal(t, "11", value.String())

	value = bitmaskValueFromNames("test2", []string{""})

	assert.Equal(t, "10", value.String())

	value = bitmaskValueFromNames("", []string{""})

	assert.Equal(t, "100", value.String())
}
