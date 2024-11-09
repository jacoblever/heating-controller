package timeseries_test

import (
	"os"
	"testing"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/timeseries"
	"github.com/jacoblever/heating-controller/brain/common"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	t.Run("when the file does not exist the latest value is blank", func(t *testing.T) {
		store := timeseries.MakeStringValueStore(clock.CreateClock(), "./no-file.txt", false, nil)
		val, err := store.GetLatestValue()

		assert.NoError(t, err)
		assert.Equal(t, "", val.Value())
	})

	t.Run("when the file does not exist, values can be stored", func(t *testing.T) {
		mockClock := common.FakeClock{}
		store := timeseries.MakeStringValueStore(&mockClock, "./file.txt", false, nil)
		t.Cleanup(func() { os.Remove("./file.txt") })

		mockClock.SetTimeRFC3339(t, "2024-01-02T03:04:05Z")
		err := store.Store("first value")
		assert.NoError(t, err)

		mockClock.SetTimeRFC3339(t, "2024-01-03T03:04:05Z")
		err = store.Store("second value")
		assert.NoError(t, err)

		t.Run("the latest value can be read", func(t *testing.T) {
			val, err := store.GetLatestValue()
			assert.NoError(t, err)
			assert.Equal(t, "second value", val.Value())
		})

		t.Run("all values can be read", func(t *testing.T) {
			vals, err := store.GetAll(nil)
			assert.NoError(t, err)
			assert.Equal(t, 2, len(vals))
			assert.Equal(t, "first value", vals[0].Value())
			assert.Equal(t, "2024-01-02T03:04:05Z", vals[0].Timestamp().Format(time.RFC3339))
			assert.Equal(t, "second value", vals[1].Value())
			assert.Equal(t, "2024-01-03T03:04:05Z", vals[1].Timestamp().Format(time.RFC3339))
		})
	})

	t.Run("only every works", func(t *testing.T) {
		mockClock := common.FakeClock{}
		tenMinutes := 10 * time.Minute
		store := timeseries.MakeStringValueStore(&mockClock, "./file.txt", false, &tenMinutes)
		t.Cleanup(func() { os.Remove("./file.txt") })

		mockClock.SetTimeRFC3339(t, "2024-01-01T10:01:00Z")
		err := store.Store("first")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first"}, getValues(t, store))

		mockClock.Advance(t, "7m")
		err = store.Store("not stored")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first"}, getValues(t, store))

		mockClock.Advance(t, "4m")
		err = store.Store("third")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first", "third"}, getValues(t, store))
	})

	t.Run("only persist on change works", func(t *testing.T) {
		store := timeseries.MakeStringValueStore(clock.CreateClock(), "./file.txt", true, nil)
		t.Cleanup(func() { os.Remove("./file.txt") })

		err := store.Store("first")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first"}, getValues(t, store))

		err = store.Store("first")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first"}, getValues(t, store))

		err = store.Store("third")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first", "third"}, getValues(t, store))

		err = store.Store("first")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first", "third", "first"}, getValues(t, store))
	})

	t.Run("only on change and only every work together", func(t *testing.T) {
		mockClock := common.FakeClock{}
		fiveMinutes := 5 * time.Minute
		store := timeseries.MakeStringValueStore(&mockClock, "./file.txt", true, &fiveMinutes)
		t.Cleanup(func() { os.Remove("./file.txt") })

		mockClock.SetTimeRFC3339(t, "2024-01-01T10:01:00Z")
		err := store.Store("first")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first"}, getValues(t, store))

		mockClock.Advance(t, "2m")
		err = store.Store("changed")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first"}, getValues(t, store))

		mockClock.Advance(t, "4m")
		err = store.Store("first")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first"}, getValues(t, store))

		err = store.Store("changed")
		assert.NoError(t, err)
		assert.Equal(t, []string{"first", "changed"}, getValues(t, store))
	})

	t.Run("GetAll limit days", func(t *testing.T) {
		mockClock := common.FakeClock{}
		store := timeseries.MakeStringValueStore(&mockClock, "./file.txt", false, nil)
		t.Cleanup(func() { os.Remove("./file.txt") })

		mockClock.SetTimeRFC3339(t, "2024-01-01T10:01:00Z")
		err := store.Store("first")
		assert.NoError(t, err)

		mockClock.Advance(t, "25h")
		err = store.Store("second")
		assert.NoError(t, err)

		mockClock.Advance(t, "4m")
		err = store.Store("third")
		assert.NoError(t, err)

		oneDay := 1
		vals, err := store.GetAll(&oneDay)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(vals))
		assert.Equal(t, "second", vals[0].Value())
		assert.Equal(t, "third", vals[1].Value())
	})
}

func getValues(t *testing.T, store timeseries.ValueStore[string]) []string {
	vals, err := store.GetAll(nil)
	assert.NoError(t, err)

	values := []string{}
	for _, v := range vals {
		values = append(values, v.Value())
	}
	return values
}
