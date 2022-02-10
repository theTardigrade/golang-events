package events

import (
	"sync"

	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

type (
	HandlerOrder int
	HandlerFunc  func()
)

type handlerDatum struct {
	// constant
	order              HandlerOrder
	handler            HandlerFunc
	shouldWaitTillDone bool

	// mutable
	bitmaskValue *bitmask.Value
	isRunning    bool
	isRunPending bool
	mainMutex    sync.Mutex

	// constant
	doneChan chan struct{}

	// mutable
	donePendingCount int
	doneMutex        sync.Mutex
}

type handlerData []*handlerDatum

func (d handlerData) Len() int           { return len(d) }
func (d handlerData) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d handlerData) Less(i, j int) bool { return (d[j].order - d[i].order) > 0 }
