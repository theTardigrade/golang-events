package events

import (
	"sort"

	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

func runnableUnorderedHandlerData(value *bitmask.Value) (handlers handlerData) {
	values := bitmaskGenerator.Values()

	defer dataMutex.RUnlock()
	dataMutex.RLock()

	dataLen := len(data)
	handlers = make(handlerData, 0, dataLen)

	for i := 0; i < dataLen; i++ {
		if datum := data[i]; datum != nil {
			if value.Contains(datum.bitmaskValue) {
				for _, v := range values {
					if datum.bitmaskValue.Contains(v) {
						handlers = append(handlers, datum)
						break
					}
				}
			}
		}
	}

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

func runnableHandlerData(value *bitmask.Value) (handlers handlerData) {
	if value != nil {
		handlers = runnableUnorderedHandlerData(value)
	} else {
		handlers = runnableUnorderedAllHandlerData()
	}

	sort.Sort(handlers)

	return
}

func runDatumPending(datum *handlerDatum) {
	for {
		var handler HandlerFunc
		var end bool

		func() {
			defer datum.mainMutex.Unlock()
			datum.mainMutex.Lock()

			if datum.isRunPending {
				datum.isRunPending = false
				handler = datum.handler
			} else {
				datum.isRunning = false
				end = true
			}
		}()

		func() {
			defer datum.doneMutex.Unlock()
			datum.doneMutex.Lock()

			if datum.donePendingCount > 0 {
				for i := datum.donePendingCount; i > 0; i-- {
					datum.doneChan <- struct{}{}
				}

				datum.donePendingCount = 0
			}
		}()

		if end || handler == nil {
			break
		}

		handler()
	}
}

func runDatum(datum *handlerDatum) {
	var handler HandlerFunc
	var shouldWaitTillDone bool

	func() {
		defer datum.mainMutex.Unlock()
		datum.mainMutex.Lock()

		if !datum.isRunning {
			datum.isRunning = true
			handler = datum.handler
		} else {
			datum.isRunPending = true
		}

		shouldWaitTillDone = datum.shouldWaitTillDone
	}()

	if handler == nil {
		if shouldWaitTillDone {
			func() {
				defer datum.doneMutex.Unlock()
				datum.doneMutex.Lock()

				datum.donePendingCount++
			}()

			<-datum.doneChan
		}

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
			defer datum.mainMutex.Unlock()
			datum.mainMutex.Lock()

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
	value := bitmaskGenerator.ValueFromNames(names...)
	handlers := runnableHandlerData(value)

	runHandlers(handlers)
}

func RunAll() {
	handlers := runnableHandlerData(nil)

	runHandlers(handlers)
}
