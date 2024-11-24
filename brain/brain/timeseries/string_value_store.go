package timeseries

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/fileio"
)

func MakeStringValueStore(clock clock.Clock, filePath string, persistOnChangeOnly bool, onlyEvery *time.Duration) ValueStore[string] {
	return &Store{
		clock: clock,

		filePath:            filePath,
		persistOnChangeOnly: persistOnChangeOnly,
		onlyEvery:           onlyEvery,
	}
}

type Store struct {
	clock clock.Clock

	filePath            string
	persistOnChangeOnly bool
	onlyEvery           *time.Duration
}

func (s Store) Store(value string) error {
	return s.StoreLazy(func() (string, error) {
		return value, nil
	})
}

func (s Store) StoreLazy(valueGetter func() (string, error)) error {
	lastValue, err := s.GetLatestValue()
	if err != nil {
		return err
	}

	timestampToWrite := Value[string]{
		timestamp: s.clock.Now(),
	}

	if !s.shouldStore(lastValue, timestampToWrite) {
		return nil
	}

	value, err := valueGetter()
	if err != nil {
		return err
	}

	valueToWrite := Value[string]{
		timestamp: s.clock.Now(),
		value:     value,
	}

	if s.shouldStore(lastValue, valueToWrite) {
		timestamp := valueToWrite.timestamp.Format(time.RFC3339)
		line := strings.Join([]string{timestamp, valueToWrite.value}, ",")
		return fileio.AppendLineToFile(s.filePath, line)
	}
	return nil
}

func (s Store) GetLatestValue() (Value[string], error) {
	line, err := fileio.ReadLastLine(s.filePath)
	if os.IsNotExist(err) {
		return Value[string]{isEmpty: true}, nil
	}
	if err != nil {
		return Value[string]{}, err
	}

	if line == "" {
		return Value[string]{isEmpty: true}, nil
	}

	parts := strings.Split(line, ",")
	timestamp, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return Value[string]{}, fmt.Errorf("last record in %s has invalid time %s", s.filePath, parts[0])
	}
	return Value[string]{
		timestamp: timestamp,
		value:     strings.Join(parts[1:], ","),
	}, nil
}

func (s Store) GetLatestValueOrDefault() string {
	value, err := s.GetLatestValue()
	if err != nil {
		return ""
	}

	return value.value
}

func (s Store) shouldStore(latestValue Value[string], valueToWrite Value[string]) bool {
	if latestValue.isEmpty {
		return true
	}

	if s.onlyEvery != nil && latestValue.timestamp.Add(*s.onlyEvery).After(s.clock.Now()) {
		return false
	}

	if s.persistOnChangeOnly && latestValue.value == valueToWrite.value {
		return false
	}

	return true
}

func (s Store) GetAll(lastXdaysOnly *int) ([]Value[string], error) {
	data, err := fileio.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(data, "\n")
	var output []Value[string]

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) > 1 {
			t, err := time.Parse(time.RFC3339, parts[0])
			if err != nil {
				log.Printf("failed to parse timestamp in '%s' in %s", line, s.filePath)
				continue
			}
			if lastXdaysOnly != nil {
				days := time.Duration(*lastXdaysOnly)
				if t.Before(s.clock.Now().Add(-24 * time.Hour * days)) {
					// Ignore logs from more than a week ago
					continue
				}
			}
			output = append(output, Value[string]{
				timestamp: t,
				value:     parts[1],
			})
		}
	}
	return output, nil
}
