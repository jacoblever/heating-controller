package timeseries

import (
	"log"
	"strconv"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/clock"
)

func MakeFloatValueStore(clock clock.Clock, filePath string, persistOnChangeOnly bool, onlyEvery *time.Duration, defaultValue float64) ValueStore[float64] {
	return &FloatStore{
		store: Store{
			clock: clock,

			filePath:            filePath,
			persistOnChangeOnly: persistOnChangeOnly,
			onlyEvery:           onlyEvery,
		},
		defaultValue: defaultValue,
	}
}

type FloatStore struct {
	store        Store
	defaultValue float64
}

func (s FloatStore) Store(value float64) error {
	return s.store.Store(strconv.FormatFloat(value, 'f', -1, 64))
}

func (s FloatStore) GetLatestValue() (Value[float64], error) {
	stringValue, err := s.store.GetLatestValue()
	if err != nil {
		return Value[float64]{isEmpty: true}, err
	}

	return s.convertToFloatValue(stringValue)
}

func (s FloatStore) GetLatestValueOrDefault() float64 {
	value, err := s.GetLatestValue()
	if err != nil {
		return s.defaultValue
	}

	return value.value
}

func (FloatStore) convertToFloatValue(stringValue Value[string]) (Value[float64], error) {
	floatValue, err := strconv.ParseFloat(stringValue.Value(), 64)
	if err != nil {
		return Value[float64]{isEmpty: true}, err
	}

	return Value[float64]{
		timestamp: stringValue.timestamp,
		value:     floatValue,
	}, nil
}

func (s FloatStore) GetAll(lastXdaysOnly *int) ([]Value[float64], error) {
	vals, err := s.store.GetAll(lastXdaysOnly)
	if err != nil {
		return nil, err
	}

	results := []Value[float64]{}
	for _, v := range vals {
		floatValue, err := s.convertToFloatValue(v)
		if err != nil {
			log.Printf("failed to parse value in '%s' in %s", v.value, s.store.filePath)
			return nil, err
		}

		results = append(results, floatValue)
	}

	return results, nil
}
