package events

import (
	"reflect"
	"sort"
	"sync"

	"github.com/theTardigrade/golang-events/internal/bitmask"
)

var (
	data      handlerData
	dataMutex sync.RWMutex
)

func Add(handler HandlerFunc, order HandlerOrder, names ...string) {
	value := bitmask.ValueFromNames(names)

	defer dataMutex.Unlock()
	dataMutex.Lock()

	// just update event bitmask value if handler function is already found
	{
		p1 := reflect.ValueOf(handler).Pointer()

		for _, datum := range data {
			if datum.order == order {
				p2 := reflect.ValueOf(datum.handler).Pointer()

				if p1 == p2 {
					datum.value.Or(value)
					return
				}
			}
		}
	}

	datum := handlerDatum{
		value:   value,
		order:   order,
		handler: handler,
	}

	data = append(data, &datum)
}

func AddUnordered(handler HandlerFunc, names ...string) {
	Add(handler, 0, names...)
}

func runnableHandlerData(value bitmask.Value) (handlers handlerData) {
	values := bitmask.Values()

	defer dataMutex.RUnlock()
	dataMutex.RLock()

	dataLen := len(data)
	handlers = make(handlerData, 0, dataLen)

	for i := 0; i < dataLen; i++ {
		if datum := data[i]; datum != nil {
			if value.IsMatch(datum.value) {
				for _, v := range values {
					if datum.value.IsMatch(v) {
						handlers = append(handlers, datum)
						break
					}
				}
			}
		}
	}

	return
}

func runDatum(datum *handlerDatum) {
	var handler HandlerFunc

	func() {
		defer datum.mutex.Unlock()
		datum.mutex.Lock()

		if !datum.isRunning {
			datum.isRunning = true
			handler = datum.handler
		} else {
			datum.isRunPending = true
		}
	}()

	if handler == nil {
		return
	}

	handler()

	for {
		var done bool

		func() {
			defer datum.mutex.Unlock()
			datum.mutex.Lock()

			if datum.isRunPending {
				datum.isRunPending = false
				handler = datum.handler
			} else {
				datum.isRunning = false
				done = true
			}
		}()

		if done || handler == nil {
			break
		}

		handler()
	}
}

func Run(names ...string) {
	value := bitmask.ValueFromNames(names)
	handlers := runnableHandlerData(value)

	sort.Sort(handlers)

	for _, datum := range handlers {
		go runDatum(datum)
	}
}
