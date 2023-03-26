package end2end_test

import (
	"net/http"
	"testing"

	"github.com/jacoblever/heating-controller/brain/end2end"
	"github.com/stretchr/testify/assert"
)

func TestBasicFlow(t *testing.T) {
	ctx := end2end.CreateTestContext(t)
	home := end2end.CreateHome()

	t.Run("the current temperature is set", func(t *testing.T) {
		response, _ := home.Thermometer.UpdateTemperature(t, ctx, 17.3)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("the boiler asks for its state", func(t *testing.T) {
		response, jsonResponse := home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, jsonResponse["BoilerState"], "off")
	})

	t.Run("the smart switch turns on", func(t *testing.T) {
		var response, _ = home.SmartSwitchAdapter.SmartSwitchAlive(t, ctx)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("the boiler asks for its state", func(t *testing.T) {
		response, jsonResponse := home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, jsonResponse["BoilerState"], "on")
	})

	t.Run("the current temperature is set hotter", func(t *testing.T) {
		response, _ := home.Thermometer.UpdateTemperature(t, ctx, 20.4)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("the boiler asks for its state", func(t *testing.T) {
		response, jsonResponse := home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, jsonResponse["BoilerState"], "off")
	})

	t.Run("the thermostat is set hotter", func(t *testing.T) {
		response, _ := home.Dashboard.UpdateThermostat(t, ctx, 22)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("the boiler asks for its state", func(t *testing.T) {
		response, jsonResponse := home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, jsonResponse["BoilerState"], "on")
	})
}
