package events

import (
	"sync"

	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

// HandlerOrder is used in the AddOptions struct
// to determine whether a handler will be called
// before or after another: lower values get
// called first.
type HandlerOrder int

// HandlerFunc is used as the type of the function
// that will be called when a relevant event is run.
type HandlerFunc func()

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
