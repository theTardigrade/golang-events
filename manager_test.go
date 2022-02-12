package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	m := NewManager()

	assert.True(t, m.innerInited)
	assert.NotEqual(t, nil, m.inner.bitmaskGenerator)
	assert.NotEqual(t, nil, m.inner.data)
}
