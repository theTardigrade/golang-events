package events

import (
	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

func (m *Manager) bitmaskValueFromNames(nameOne string, nameMany []string) (value *bitmask.Value) {
	nameSet := make(map[string]struct{})

	for _, n := range nameMany {
		if n != "" {
			nameSet[n] = struct{}{}
		}
	}

	if nameOne != "" {
		nameSet[nameOne] = struct{}{}
	}

	switch len(nameSet) {
	case 0:
		value = m.bitmaskGenerator.ValueFromName("")
	case 1:
		for n := range nameSet {
			value = m.bitmaskGenerator.ValueFromName(n)
		}
	default:
		nameAll := make([]string, len(nameSet))

		var i int
		for n := range nameSet {
			nameAll[i] = n
			i++
		}

		value = m.bitmaskGenerator.ValueFromNames(nameAll...)
	}

	return
}
