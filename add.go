package events

import (
	"reflect"
)

// AddOptions is used by the Add method on a Manager
// and the package-level Add function.
// The Name and Names fields are used together to uniquely
// identify an event.
// The Handler field contains the function that will be
// added.
// The Order field is used to determine whether a handler
// will be called before or after another: lower values
// get called first.
// The ShouldWaitTillDone field determines whether the Run
// method (or package-level Run function) will wait for
// the handler to finish its work before returning.
type AddOptions struct {
	Name               string
	Names              []string
	Handler            HandlerFunc
	Order              HandlerOrder
	ShouldWaitTillDone bool
}

// Add sets up a handler function that will be called
// when a named event is run.
func (m *Manager) Add(options AddOptions) {
	if options.Handler == nil {
		return
	}

	initManagerMethod(&m)

	bitmaskValue := m.bitmaskValueFromNames(options.Name, options.Names)

	defer m.dataMutex.Unlock()
	m.dataMutex.Lock()

	// just update event bitmask value if handler function is already found
	{
		p1 := reflect.ValueOf(options.Handler).Pointer()

		for _, datum := range m.inner.data {
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

	m.inner.data = append(m.inner.data, &datum)
}

// Add calls the Add method on the default manager.
func Add(options AddOptions) {
	defaultManager.Add(options)
}
