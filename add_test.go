package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	m := NewManager()

	assert.Equal(t, 0, len(m.inner.data))

	handler := func() {}

	m.Add(AddOptions{
		Name:    "test",
		Handler: handler,
	})

	assert.Equal(t, 1, len(m.inner.data))
	assert.Equal(t, "1", m.inner.data[0].bitmaskValue.String())
	assert.Equal(t, false, m.inner.data[0].shouldWaitTillDone)
	assert.Equal(t, HandlerOrder(0), m.inner.data[0].order)
	assert.Equal(t, false, m.inner.data[0].isRunning)
	assert.Equal(t, false, m.inner.data[0].isRunPending)
	assert.Equal(t, 0, m.inner.data[0].donePendingCount)

	m.Add(AddOptions{
		Name:    "test2",
		Handler: handler,
	})

	assert.Equal(t, 1, len(m.inner.data))
	assert.Equal(t, "11", m.inner.data[0].bitmaskValue.String())

	handler = func() {}

	m.Add(AddOptions{
		Name:    "test",
		Handler: handler,
	})

	assert.Equal(t, 2, len(m.inner.data))
	assert.Equal(t, "11", m.inner.data[0].bitmaskValue.String())
	assert.Equal(t, "1", m.inner.data[1].bitmaskValue.String())

	m.Add(AddOptions{
		Name:    "test3",
		Handler: handler,
	})

	assert.Equal(t, 2, len(m.inner.data))
	assert.Equal(t, "11", m.inner.data[0].bitmaskValue.String())
	assert.Equal(t, "101", m.inner.data[1].bitmaskValue.String())
}
