package timeseries

import "time"

type ValueStore[T any] interface {
	Store(value T) error
	GetLatestValue() (Value[T], error)
	GetLatestValueOrDefault() T
	GetAll(lastXdaysOnly *int) ([]Value[T], error)
}

type Value[T any] struct {
	timestamp time.Time
	value     T
	isEmpty   bool
}

func (v *Value[T]) Timestamp() time.Time {
	return v.timestamp
}

func (v *Value[T]) Value() T {
	return v.value
}

type OnOff bool

const On OnOff = true
const Off OnOff = false

func (b OnOff) OnOffString() string {
	if b {
		return "on"
	}
	return "off"
}

func (b OnOff) Bool() bool {
	return bool(b)
}
