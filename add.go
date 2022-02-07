package events

import (
	"reflect"

	bitmask "github.com/theTardigrade/golang-infiniteBitmask"
)

type AddOptions struct {
	Name               string
	Names              []string
	Handler            HandlerFunc
	Order              HandlerOrder
	ShouldWaitTillDone bool
}

func addValueFromNames(nameOne string, nameMany []string) (value *bitmask.Value) {
	if nameManyLen := len(nameMany); nameManyLen > 0 {
		var nameAll []string

		if nameOne == "" {
			nameAll = nameMany
		} else {
			nameAll = make([]string, nameManyLen+1)

			for i, n := range nameMany {
				nameAll[i] = n
			}

			nameAll[nameManyLen] = nameOne
		}

		value = bitmaskGenerator.ValueFromNames(nameAll...)
	} else {
		value = bitmaskGenerator.ValueFromName(nameOne)
	}

	return
}

func Add(options AddOptions) {
	value := addValueFromNames(options.Name, options.Names)

	defer dataMutex.Unlock()
	dataMutex.Lock()

	// just update event bitmask value if handler function is already found
	{
		p1 := reflect.ValueOf(options.Handler).Pointer()

		for _, datum := range data {
			if datum.order == options.Order && datum.shouldWaitTillDone == options.ShouldWaitTillDone {
				p2 := reflect.ValueOf(datum.handler).Pointer()

				if p1 == p2 {
					datum.value.Combine(value)
					return
				}
			}
		}
	}

	datum := handlerDatum{
		value:              value,
		order:              options.Order,
		handler:            options.Handler,
		shouldWaitTillDone: options.ShouldWaitTillDone,
	}

	data = append(data, &datum)
}
