package timeseries

import (
	"fmt"
	"log"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/clock"
)

func MakeOnOffValueStore(clock clock.Clock, filePath string, persistOnChangeOnly bool, onlyEvery *time.Duration) ValueStore[OnOff] {
	return &OnOffStore{
		store: Store{
			clock: clock,

			filePath:            filePath,
			persistOnChangeOnly: persistOnChangeOnly,
			onlyEvery:           onlyEvery,
		},
	}
}

type OnOffStore struct {
	store Store
}

func (s OnOffStore) Store(value OnOff) error {
	return s.store.Store(value.OnOffString())
}

func (s OnOffStore) StoreLazy(valueGetter func() (OnOff, error)) error {
	return s.store.StoreLazy(func() (string, error) {
		value, err := valueGetter()
		if err != nil {
			return "", err
		}

		return value.OnOffString(), nil
	})
}

func (s OnOffStore) GetLatestValue() (Value[OnOff], error) {
	stringValue, err := s.store.GetLatestValue()
	if err != nil {
		return Value[OnOff]{isEmpty: true}, err
	}

	return s.convertToBoolValue(stringValue)
}

func (s OnOffStore) GetLatestValueOrDefault() OnOff {
	value, err := s.GetLatestValue()
	if err != nil {
		return Off
	}

	return value.value
}

func (OnOffStore) convertToBoolValue(stringValue Value[string]) (Value[OnOff], error) {
	onOffValue := OnOff(stringValue.value == "on")

	if !onOffValue && stringValue.value != "off" {
		return Value[OnOff]{isEmpty: true}, fmt.Errorf("unknown on/off value '%s'", stringValue.value)
	}

	return Value[OnOff]{
		timestamp: stringValue.timestamp,
		value:     onOffValue,
	}, nil
}

func (s OnOffStore) GetAll(lastXdaysOnly *int) ([]Value[OnOff], error) {
	vals, err := s.store.GetAll(lastXdaysOnly)
	if err != nil {
		return nil, err
	}

	results := []Value[OnOff]{}
	for _, v := range vals {
		boolValue, err := s.convertToBoolValue(v)
		if err != nil {
			log.Printf("failed to parse value in '%s' in %s", v.value, s.store.filePath)
			return nil, err
		}

		results = append(results, boolValue)
	}

	return results, nil
}
