package events

import (
	"reflect"
)

type AddOptions struct {
	Name               string
	Names              []string
	Handler            HandlerFunc
	Order              HandlerOrder
	ShouldWaitTillDone bool
}

func (m *Manager) Add(options AddOptions) {
	if m == nil {
		m = defaultManager
	}

	if options.Handler == nil {
		panic(ErrHandlerNil)
	}

	bitmaskValue := m.bitmaskValueFromNames(options.Name, options.Names)

	defer m.dataMutex.Unlock()
	m.dataMutex.Lock()

	// just update event bitmask value if handler function is already found
	{
		p1 := reflect.ValueOf(options.Handler).Pointer()

		for _, datum := range m.data {
			if datum.order == options.Order && datum.shouldWaitTillDone == options.ShouldWaitTillDone {
				p2 := reflect.ValueOf(datum.handler).Pointer()

				if p1 == p2 {
					func() {
						defer datum.mainMutex.Unlock()
						datum.mainMutex.Lock()

						datum.bitmaskValue.Combine(bitmaskValue)
					}()

					return
				}
			}
		}
	}

	datum := handlerDatum{
		bitmaskValue:       bitmaskValue,
		order:              options.Order,
		handler:            options.Handler,
		shouldWaitTillDone: options.ShouldWaitTillDone,
		doneChan:           make(chan struct{}),
	}

	m.data = append(m.data, &datum)
}

func Add(options AddOptions) {
	defaultManager.Add(options)
}
