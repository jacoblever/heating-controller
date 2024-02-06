package logging_test

import (
	"testing"

	"github.com/jacoblever/heating-controller/brain/brain/logging"
	"github.com/jacoblever/heating-controller/brain/common"
	"github.com/stretchr/testify/assert"
)

var logKey = "per_day_test_dir"

func TestPerDayLogger(t *testing.T) {
	clock := common.FakeClock{}
	logger := logging.PerDayLogger{
		Key:      logKey,
		Clock:    &clock,
		Settings: logging.Settings{DaysToKeepFor: 3, SynchronousDeletion: true},
	}
	defer logger.DeleteAllLogs()

	t.Run("logs on day 1 work", func(t *testing.T) {
		clock.SetTimeRFC3339(t, "2024-01-20T00:00:10Z")
		logger.Log("Day 1 Log 1")
		clock.SetTimeRFC3339(t, "2024-01-20T00:00:12Z")
		logger.Log("Day 1 Log 2")

		logs, err := logger.GetLogs()
		assert.Nil(t, err)
		assert.Equal(t, []string{
			"2024-01-20T00:00:10Z,Day 1 Log 1",
			"2024-01-20T00:00:12Z,Day 1 Log 2",
		}, logs)
	})

	t.Run("logs on day 2 work", func(t *testing.T) {
		clock.SetTimeRFC3339(t, "2024-01-21T00:00:11Z")
		logger.Log("Day 2 Log 1")

		logs, err := logger.GetLogs()
		assert.Nil(t, err)
		assert.Equal(t, []string{
			"2024-01-20T00:00:10Z,Day 1 Log 1",
			"2024-01-20T00:00:12Z,Day 1 Log 2",
			"2024-01-21T00:00:11Z,Day 2 Log 1",
		}, logs)
	})

	t.Run("logs on day 3 work", func(t *testing.T) {
		clock.SetTimeRFC3339(t, "2024-01-22T00:00:11Z")
		logger.Log("Day 3 Log 1")
		clock.SetTimeRFC3339(t, "2024-01-22T00:00:13Z")
		logger.Log("Day 3 Log 2")

		logs, err := logger.GetLogs()
		assert.Nil(t, err)
		assert.Equal(t, []string{
			"2024-01-20T00:00:10Z,Day 1 Log 1",
			"2024-01-20T00:00:12Z,Day 1 Log 2",
			"2024-01-21T00:00:11Z,Day 2 Log 1",
			"2024-01-22T00:00:11Z,Day 3 Log 1",
			"2024-01-22T00:00:13Z,Day 3 Log 2",
		}, logs)
	})

	t.Run("logs on day 4 work and remove day 1", func(t *testing.T) {
		clock.SetTimeRFC3339(t, "2024-01-23T00:00:08Z")
		logger.Log("Day 4 Log 1")

		logs, err := logger.GetLogs()
		assert.Nil(t, err)
		assert.Equal(t, []string{
			"2024-01-21T00:00:11Z,Day 2 Log 1",
			"2024-01-22T00:00:11Z,Day 3 Log 1",
			"2024-01-22T00:00:13Z,Day 3 Log 2",
			"2024-01-23T00:00:08Z,Day 4 Log 1",
		}, logs)
	})
}
