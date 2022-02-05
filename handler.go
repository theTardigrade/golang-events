package events

import (
	"sync"

	"github.com/theTardigrade/golang-events/internal/bitmask"
)

type (
	HandlerOrder uint
	HandlerFunc  func()
)

type handlerDatum struct {
	value        bitmask.Value
	order        HandlerOrder
	handler      HandlerFunc
	isRunning    bool
	isRunPending bool
	mutex        sync.Mutex
}

type handlerData []*handlerDatum

func (d handlerData) Len() int           { return len(d) }
func (d handlerData) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d handlerData) Less(i, j int) bool { return (d[j].order - d[i].order) > 0 }
