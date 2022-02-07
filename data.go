package events

import (
	"sync"
)

var (
	data      handlerData
	dataMutex sync.RWMutex
)
