package events

import (
	"sort"

	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

func runnableUnorderedHandlerData(value bitmask.Value) (handlers handlerData) {
	values := bitmask.Values()

	defer dataMutex.RUnlock()
	dataMutex.RLock()

	dataLen := len(data)
	handlers = make(handlerData, 0, dataLen)

	for i := 0; i < dataLen; i++ {
		if datum := data[i]; datum != nil {
			if value.Contains(datum.value) {
				for _, v := range values {
					if datum.value.Contains(v) {
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

func runnableUnorderedAllHandlerData() (handlers handlerData) {
	defer dataMutex.RUnlock()
	dataMutex.RLock()

	dataLen := len(data)
	handlers = make(handlerData, 0, dataLen)

	for i := 0; i < dataLen; i++ {
		if datum := data[i]; datum != nil {
			handlers = append(handlers, datum)
		}
	}

	return
}

func runnableAllHandlerData() (handlers handlerData) {
	handlers = runnableUnorderedAllHandlerData()

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

func runHandlers(handlers handlerData) {
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

func Run(names ...string) {
	value := bitmask.ValueFromNames(names)
	handlers := runnableHandlerData(value)

	runHandlers(handlers)
}

func RunAll() {
	handlers := runnableAllHandlerData()

	runHandlers(handlers)
}
