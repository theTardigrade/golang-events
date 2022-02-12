package events

import (
	"sync"

	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

type Manager struct {
	bitmaskGenerator *bitmask.Generator
	data             handlerData
	dataMutex        sync.RWMutex
}

var (
	defaultManager = NewManager()
)

func NewManager() *Manager {
	return &Manager{
		bitmaskGenerator: bitmask.NewGenerator(),
		data:             make(handlerData, 0, 1024),
	}
}
