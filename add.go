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

func addValueFromNames(options AddOptions) (value bitmask.Value) {
	if len(options.Names) > 0 {
		if options.Name == "" {
			value = bitmask.ValueFromNames(options.Names)
		} else {
			var names []string

			names = append(names, options.Names...)
			names = append(names, options.Name)

			value = bitmask.ValueFromNames(names)
		}
	} else {
		value = bitmask.ValueFromName(options.Name)
	}

	return
}

func Add(options AddOptions) {
	value := addValueFromNames(options)

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
