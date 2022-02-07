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

func addValueFromNames(nameOne string, nameMany []string) (value bitmask.Value) {
	if len(nameMany) > 0 {
		if nameOne == "" {
			value = bitmask.ValueFromNames(nameMany)
		} else {
			var names []string

			names = append(names, nameMany...)
			names = append(names, nameOne)

			value = bitmask.ValueFromNames(names)
		}
	} else {
		value = bitmask.ValueFromName(nameOne)
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
