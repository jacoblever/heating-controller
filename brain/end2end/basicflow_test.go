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

	t.Run("when the current temperature is set", func(t *testing.T) {
		response, _ := home.Thermometer.UpdateTemperature(t, ctx, 17.3)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("when the boiler asks for its state it should be off", func(t *testing.T) {
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, "off", home.Boiler.State)
	})

	t.Run("the boiler does a manual turn", func(t *testing.T) {
		assert.Equal(t, 0, home.Boiler.Position)
		home.Dashboard.TurnBoiler(t, ctx, "100")
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, 100, home.Boiler.Position)
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, 100, home.Boiler.Position)
		home.Dashboard.TurnBoiler(t, ctx, "-25")
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, 75, home.Boiler.Position)
		home.Dashboard.TurnBoiler(t, ctx, "1")
		home.Dashboard.TurnBoiler(t, ctx, "-2")
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, 76, home.Boiler.Position)
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, 74, home.Boiler.Position)
	})

	t.Run("when the smart switch turns on", func(t *testing.T) {
		var response, _ = home.SmartSwitchAdapter.SmartSwitchAlive(t, ctx)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("when the boiler asks for its state it should be on", func(t *testing.T) {
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, "on", home.Boiler.State)
	})

	t.Run("when the current temperature is set hotter", func(t *testing.T) {
		response, _ := home.Thermometer.UpdateTemperature(t, ctx, 20.4)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("when the boiler asks for its state it should be off", func(t *testing.T) {
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, "off", home.Boiler.State)
	})

	t.Run("when the thermostat is set hotter", func(t *testing.T) {
		response, _ := home.Dashboard.UpdateThermostat(t, ctx, 22)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("when the boiler asks for its state it should be on", func(t *testing.T) {
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, "on", home.Boiler.State)
	})

	home.Dashboard.GetGraphData(t, ctx)
}
