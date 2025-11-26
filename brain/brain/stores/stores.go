package stores

import (
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/timeseries"
)

var tenMinutes time.Duration = 10 * time.Minute
var oneHour time.Duration = time.Hour
var defaultThermostatThreshold float64 = 18

type Stores struct {
	Temperature        timeseries.ValueStore[float64]
	Temperature1       timeseries.ValueStore[float64]
	Temperature2       timeseries.ValueStore[float64]
	OutsideTemperature timeseries.ValueStore[float64]
	Thermostat         timeseries.ValueStore[float64]
	SmartSwitch        timeseries.ValueStore[timeseries.OnOff]
	BoilerState        timeseries.ValueStore[timeseries.OnOff]
	BoilerMode         timeseries.ValueStore[string]
}

func MakeStores(c clock.Clock, config Config) Stores {
	return Stores{
		Temperature:        timeseries.MakeFloatValueStore(c, config.TemperatureLogFilePath, false, &tenMinutes, 0),
		Temperature1:       timeseries.MakeFloatValueStore(c, config.TemperatureLog1FilePath, false, &tenMinutes, 0),
		Temperature2:       timeseries.MakeFloatValueStore(c, config.TemperatureLog1FilePath, false, &tenMinutes, 0),
		OutsideTemperature: timeseries.MakeFloatValueStore(c, config.OutsideTemperatureLogFilePath, false, &oneHour, 0),

		Thermostat:  timeseries.MakeFloatValueStore(c, config.ThermostatThresholdLogFilePath, true, nil, defaultThermostatThreshold),
		SmartSwitch: timeseries.MakeOnOffValueStore(c, config.SmartSwitchStateLogFilePath, true, nil),
		BoilerState: timeseries.MakeOnOffValueStore(c, config.BoilerStateLogFilePath, true, nil),
		BoilerMode:  timeseries.MakeStringValueStore(c, config.BoilerModeLogFilePath, true, nil),
	}
}
