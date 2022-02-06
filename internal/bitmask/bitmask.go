package bitmask

import (
	"math/big"
	"sync"
)

type Value struct {
	i *big.Int
}

var (
	bigZero = big.NewInt(0)
	bigOne  = big.NewInt(1)
	bigTwo  = big.NewInt(2)

	valueCurrent = Value{i: big.NewInt(0)}
	valuesByName = make(map[string]Value)
	mutex        sync.RWMutex
)

func Names() (names []string) {
	defer mutex.RUnlock()
	mutex.RLock()

	names = make([]string, len(valuesByName))

	var i int
	for n := range valuesByName {
		names[i] = n
		i++
	}

	return
}

func Values() (values []Value) {
	defer mutex.RUnlock()
	mutex.RLock()

	values = make([]Value, len(valuesByName))

	var i int
	var found bool
	for n := range valuesByName {
		values[i], found = valueFromName(n, false)
		if !found {
			panic(ErrValueNotFound)
		}

		i++
	}

	return
}

func ValueFromName(name string) (value Value) {
	value, found := valueFromName(name, false)
	if found {
		return
	}

	value, found = valueFromName(name, true)
	if !found {
		panic(ErrValueNotFound)
	}

	return
}

func ValueFromNames(names []string) (value Value) {
	nameValues := make([]Value, len(names))

	for i, n := range names {
		nameValues[i] = ValueFromName(n)
	}

	i := new(big.Int)

	for _, v := range nameValues {
		i.Or(i, v.i)
	}

	value.i = i

	return
}

func valueFromName(name string, allowWrites bool) (value Value, found bool) {
	if allowWrites {
		defer mutex.Unlock()
		mutex.Lock()
	} else {
		defer mutex.RUnlock()
		mutex.RLock()
	}

	value, found = valuesByName[name]
	if !found {
		if !allowWrites {
			return
		}

		clonedValueCurrent := valueCurrent.Clone()

		valueCurrent.i.Add(valueCurrent.i, bigOne)

		valuesByName[name] = clonedValueCurrent

		value = clonedValueCurrent
	}

	shiftedValue := Value{i: new(big.Int)}
	shiftedValue.i.Exp(bigTwo, value.i, nil)

	value, found = shiftedValue, true

	return
}

func (v Value) Clone() (v2 Value) {
	i2 := new(big.Int)
	i2.Set(v.i)

	v2 = Value{
		i: i2,
	}

	return
}

func (v Value) Or(v2 Value) {
	v.i.Or(v.i, v2.i)
}

func (v Value) BigInt() *big.Int {
	return v.i
}

func (v Value) IsMatch(v2 Value) bool {
	vI, v2I := v.i, v2.i

	b := new(big.Int)
	b.And(vI, v2I)

	return b.Cmp(bigZero) == 1
}
