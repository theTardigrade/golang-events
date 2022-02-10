package events

import (
	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

var (
	bitmaskGenerator = bitmask.NewGenerator()
)

func bitmaskValueFromNames(nameOne string, nameMany []string) (value *bitmask.Value) {
	nameSet := make(map[string]struct{})

	for _, n := range nameMany {
		if n != "" {
			nameSet[n] = struct{}{}
		}
	}

	if nameOne != "" {
		nameSet[nameOne] = struct{}{}
	}

	if len(nameSet) == 0 {
		nameSet[""] = struct{}{}
	}

	if len(nameSet) == 1 {
		for n := range nameSet {
			value = bitmaskGenerator.ValueFromName(n)
		}
	} else {
		nameAll := make([]string, len(nameSet))

		var i int
		for n := range nameSet {
			nameAll[i] = n
			i++
		}

		value = bitmaskGenerator.ValueFromNames(nameAll...)
	}

	return
}
