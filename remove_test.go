package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	m := NewManager()

	handler := func() {}

	m.Add(AddOptions{
		Name:    "test",
		Handler: handler,
	})

	assert.Equal(t, 1, len(m.inner.data))

	m.Remove(RemoveOptions{
		Name:    "test",
		Handler: handler,
	})

	assert.Equal(t, 0, len(m.inner.data))

	m.Add(AddOptions{
		Names:   []string{"test1", "test2"},
		Handler: handler,
	})

	assert.Equal(t, 1, len(m.inner.data))

	m.Remove(RemoveOptions{
		Name:    "test1",
		Handler: handler,
	})

	assert.Equal(t, 1, len(m.inner.data))

	m.Remove(RemoveOptions{
		Name:    "test2",
		Handler: handler,
	})

	assert.Equal(t, 0, len(m.inner.data))
}
