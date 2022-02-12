package events

import "reflect"

type RemoveOptions struct {
	Name    string
	Names   []string
	Handler HandlerFunc
}

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
