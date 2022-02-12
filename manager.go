package events

import (
	"sync"

	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

// Manager contains data and methods to handle all
// event management.
// You can, however, call the public package-level
// functions, if you only want to use a default manager.
type Manager struct {
	inner       managerInner
	innerInited bool
	dataMutex   sync.RWMutex
}

type managerInner struct {
	bitmaskGenerator *bitmask.Generator
	data             handlerData
}

var (
	defaultManager = NewManager()
)

// NewManager returns a pointer to an initialized
// Manager struct.
func NewManager() (m *Manager) {
	m = &Manager{}

	m.initInner()

	return
}

func (m *Manager) initInner() {
	defer m.dataMutex.Unlock()
	m.dataMutex.Lock()

	if m.innerInited {
		return
	}

	m.inner.bitmaskGenerator = bitmask.NewGenerator()
	m.inner.data = make(handlerData, 0, 1024)

	m.innerInited = true
}

func (m *Manager) checkInner() {
	var innerInited bool

	func() {
		defer m.dataMutex.RUnlock()
		m.dataMutex.RLock()

		innerInited = m.innerInited
	}()

	if !innerInited {
		m.initInner()
	}
}

func initManagerMethod(m **Manager) {
	if mValue := *m; mValue == nil {
		*m = defaultManager
	} else if mValue != defaultManager {
		mValue.checkInner()
	}
}
