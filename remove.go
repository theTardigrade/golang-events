package events

import "reflect"

// RemoveOptions is used by the Remove method on a Manager
// and the package-level Remove function.
// The Name and Names fields are used together to uniquely
// identify an event.
// The Handler field contains the function that will be
// removed.
type RemoveOptions struct {
	Name    string
	Names   []string
	Handler HandlerFunc
}

// Remove stops a handler function from being called
// when a named event is run.
func (m *Manager) Remove(options RemoveOptions) {
	if options.Handler == nil {
		return
	}

	initManagerMethod(&m)

	bitmaskValue := m.bitmaskValueFromNames(options.Name, options.Names)

	defer m.dataMutex.Unlock()
	m.dataMutex.Lock()

	{
		p1 := reflect.ValueOf(options.Handler).Pointer()

		for i, datum := range m.inner.data {
			if p2 := reflect.ValueOf(datum.handler).Pointer(); p1 == p2 {
				func() {
					defer datum.mainMutex.Unlock()
					datum.mainMutex.Lock()

					datum.bitmaskValue.Uncombine(bitmaskValue)

					if datum.bitmaskValue.IsEmpty() {
						lastIndex := len(m.inner.data) - 1

						if i != lastIndex {
							m.inner.data[lastIndex], m.inner.data[i] = m.inner.data[i], m.inner.data[lastIndex]
						}

						m.inner.data = m.inner.data[:lastIndex]
					}
				}()

				return
			}
		}
	}
}

// Remove calls the Remove method on the default manager.
func Remove(options RemoveOptions) {
	defaultManager.Remove(options)
}
