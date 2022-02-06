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

type AddOptions struct {
	Name               string
	Names              []string
	Handler            HandlerFunc
	Order              HandlerOrder
	ShouldWaitTillDone bool
}

func addValue(options *AddOptions) (value bitmask.Value) {
	if len(options.Names) > 0 {
		if options.Name == "" {
			value = bitmask.ValueFromNames(options.Names)
		} else {
			var names []string

			names = append(names, options.Names...)
			names = append(names, options.Name)

			value = bitmask.ValueFromNames(names)
		}
	} else {
		value = bitmask.ValueFromName(options.Name)
	}

	return
}

func Add(options AddOptions) {
	value := addValue(&options)

	defer dataMutex.Unlock()
	dataMutex.Lock()

	// just update event bitmask value if handler function is already found
	{
		p1 := reflect.ValueOf(options.Handler).Pointer()

		for _, datum := range data {
			if datum.order == options.Order {
				p2 := reflect.ValueOf(datum.handler).Pointer()

				if p1 == p2 {
					datum.value.Or(value)
					return
				}
			}
		}
	}

	datum := handlerDatum{
		value:              value,
		order:              options.Order,
		handler:            options.Handler,
		shouldWaitTillDone: options.ShouldWaitTillDone,
	}

	data = append(data, &datum)
}

func runnableUnorderedHandlerData(value bitmask.Value) (handlers handlerData) {
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

func runnableHandlerData(value bitmask.Value) (handlers handlerData) {
	handlers = runnableUnorderedHandlerData(value)

	sort.Sort(handlers)

	return
}

func runDatumPending(datum *handlerDatum) {
	for {
		var handler HandlerFunc
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

func runDatum(datum *handlerDatum) {
	var handler HandlerFunc
	var shouldWaitTillDone bool

	func() {
		defer datum.mutex.Unlock()
		datum.mutex.Lock()

		if !datum.isRunning {
			datum.isRunning = true
			handler = datum.handler
		} else {
			datum.isRunPending = true
		}

		shouldWaitTillDone = datum.shouldWaitTillDone
	}()

	if handler == nil {
		return
	}

	if shouldWaitTillDone {
		handler()
		go runDatumPending(datum)
	} else {
		go func() {
			handler()
			runDatumPending(datum)
		}()
	}
}

func Run(names ...string) {
	value := bitmask.ValueFromNames(names)
	handlers := runnableHandlerData(value)

	for _, datum := range handlers {
		var done bool

		func(datum *handlerDatum) {
			defer datum.mutex.Unlock()
			datum.mutex.Lock()

			if datum.isRunning {
				datum.isRunPending = true
				done = true
			}
		}(datum)

		if !done {
			runDatum(datum)
		}
	}
}
