package events

import (
	"sort"

	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

func (m *Manager) runnableUnorderedHandlerData(value *bitmask.Value) (handlers handlerData) {
	values := m.bitmaskGenerator.Values()

	defer m.dataMutex.RUnlock()
	m.dataMutex.RLock()

	handlers = make(handlerData, 0, len(m.data))

	for _, datum := range m.data {
		func() {
			defer datum.mainMutex.Unlock()
			datum.mainMutex.Lock()

			if value.Contains(datum.bitmaskValue) {
				for _, v := range values {
					if datum.bitmaskValue.Contains(v) {
						handlers = append(handlers, datum)
						break
					}
				}
			}
		}()
	}

	return
}

func (m *Manager) runnableUnorderedAllHandlerData() (handlers handlerData) {
	defer m.dataMutex.RUnlock()
	m.dataMutex.RLock()

	handlers = make(handlerData, 0, len(m.data))

	for _, datum := range m.data {
		handlers = append(handlers, datum)
	}

	return
}

func (m *Manager) runnableHandlerData(value *bitmask.Value) (handlers handlerData) {
	if value != nil {
		handlers = m.runnableUnorderedHandlerData(value)
	} else {
		handlers = m.runnableUnorderedAllHandlerData()
	}

	sort.Sort(handlers)

	return
}

func (m *Manager) runDatumPending(datum *handlerDatum) {
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

		var donePendingCount int

		func() {
			defer datum.doneMutex.Unlock()
			datum.doneMutex.Lock()

			if datum.donePendingCount > 0 {
				donePendingCount = datum.donePendingCount
				datum.donePendingCount = 0
			}
		}()

		for i := 0; i < donePendingCount; i++ {
			datum.doneChan <- struct{}{}
		}

		if end || handler == nil {
			break
		}

		handler()
	}
}

func (m *Manager) runDatum(datum *handlerDatum) {
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
		go m.runDatumPending(datum)
	} else {
		go func() {
			handler()
			m.runDatumPending(datum)
		}()
	}
}

func (m *Manager) runHandlers(handlers handlerData) {
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
			m.runDatum(datum)
		}
	}
}

func (m *Manager) Run(names ...string) {
	if m == nil {
		m = defaultManager
	}

	value := m.bitmaskGenerator.ValueFromNames(names...)
	handlers := m.runnableHandlerData(value)

	m.runHandlers(handlers)
}

func Run(names ...string) {
	defaultManager.Run(names...)
}

func (m *Manager) RunAll() {
	if m == nil {
		m = defaultManager
	}

	handlers := m.runnableHandlerData(nil)

	m.runHandlers(handlers)
}

func RunAll() {
	defaultManager.RunAll()
}
