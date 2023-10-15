package end2end_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jacoblever/heating-controller/brain/end2end"
	"github.com/stretchr/testify/assert"
)

func TestTempLog(t *testing.T) {
	ctx := end2end.CreateTestContext(t)
	home := end2end.CreateHome()

	ctx.Clock.SetTimeRFC3339(t, "2020-03-26T19:00:00+01:00")

	t.Run("the current temperature is sent", func(t *testing.T) {
		response, _ := home.Thermometer.UpdateTemperature(t, ctx, 17.3)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("the temperature is recorded in the log", func(t *testing.T) {
		c := readFile(t, ctx.Config.TemperatureLogFilePath)
		assert.Equal(t, "\n2020-03-26T19:00:00+01:00,17.300000", c)
	})

	ctx.Clock.Advance(t, "1m")

	t.Run("a minute later, the current temperature is sent", func(t *testing.T) {
		response, _ := home.Thermometer.UpdateTemperature(t, ctx, 18.2)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("the temperature is not recorded", func(t *testing.T) {
		c := readFile(t, ctx.Config.TemperatureLogFilePath)
		assert.Equal(t, "\n2020-03-26T19:00:00+01:00,17.300000", c)
	})

	ctx.Clock.Advance(t, "10m")

	t.Run("10 more minutes later, the current temperature is sent", func(t *testing.T) {
		response, _ := home.Thermometer.UpdateTemperature(t, ctx, 18.2)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("the temperature is not recorded in the log", func(t *testing.T) {
		c := readFile(t, ctx.Config.TemperatureLogFilePath)
		assert.Equal(t, "\n2020-03-26T19:00:00+01:00,17.300000\n2020-03-26T19:11:00+01:00,18.200000", c)
	})
}

func readFile(t *testing.T, filePath string) string {
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("error reading file %s: %s\n", filePath, err)
		return ""
	}

	return string(buffer)
}
