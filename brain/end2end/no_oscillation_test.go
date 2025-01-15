package end2end_test

import (
	"testing"

	"github.com/jacoblever/heating-controller/brain/end2end"
	"github.com/stretchr/testify/assert"
)

func TestNoOscillation(t *testing.T) {
	ctx := end2end.CreateTestContext(t)
	home := end2end.CreateHome()

	_, _ = home.SmartSwitchAdapter.SmartSwitchAlive(t, ctx)
	_, _ = home.Thermometer.UpdateTemperature(t, ctx, 17.3)
	_, _ = home.Dashboard.UpdateThermostat(t, ctx, 20.2)

	t.Run("when the boiler asks for its state it should be on", func(t *testing.T) {
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, "on", home.Boiler.State)
	})

	t.Run("when the house gets close to the target temperature the boiler should be on", func(t *testing.T) {
		_, _ = home.Thermometer.UpdateTemperature(t, ctx, 20.1)
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, "on", home.Boiler.State)
	})

	t.Run("when the house gets hotter than the target temperature the boiler should be off", func(t *testing.T) {
		_, _ = home.Thermometer.UpdateTemperature(t, ctx, 20.3)
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, "off", home.Boiler.State)
	})

	t.Run("when the house gets slightly cooler than the target temperature the boiler should stay off", func(t *testing.T) {
		_, _ = home.Thermometer.UpdateTemperature(t, ctx, 20.1)
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, "off", home.Boiler.State)
	})

	t.Run("when the house gets close to 0.25 below the target temperature the boiler should stay off", func(t *testing.T) {
		_, _ = home.Thermometer.UpdateTemperature(t, ctx, 20.0)
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, "off", home.Boiler.State)
	})

	t.Run("when the house gets cooler than 0.5 below the target temperature the boiler should be on", func(t *testing.T) {
		_, _ = home.Thermometer.UpdateTemperature(t, ctx, 19.6)
		_, _ = home.Boiler.BoilerState(t, ctx)
		assert.Equal(t, "on", home.Boiler.State)
	})
}
