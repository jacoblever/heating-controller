package boiler

import (
	"log"
	"strconv"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/boiler/commandqueue"
	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/fileio"
	"github.com/jacoblever/heating-controller/brain/brain/logging"
	"github.com/jacoblever/heating-controller/brain/brain/timeseries"
)

var debounceBuffer = 0.5
var defaultThermostatThreshold float64 = 18

type Boiler struct {
	config  Config
	clock   clock.Clock
	loggers logging.Loggers

	BoilerCommandQueue commandqueue.CommandQueue
}

func MakeBoiler(config Config, clock clock.Clock, loggers logging.Loggers) Boiler {
	return Boiler{
		config:  config,
		clock:   clock,
		loggers: loggers,

		BoilerCommandQueue: commandqueue.Make(),
	}
}

type BoilerState struct {
	StateOfBoiler         string
	CalculatedBoilerState string
}

func (b Boiler) GetBoilerState(logChange bool) BoilerState {
	currentTemperature := b.GetTemperature()
	thermostatThreshold := b.GetThermostat()
	smartSwitchOn := b.GetSmartSwitchStatus()

	currentBoilerStateRecord, err := timeseries.ReadLastRecord(b.config.BoilerStateLogFilePath)
	if err != nil {
		b.loggers.Get("brain").Logf("error reading current boiler state log: %s", err)
		currentBoilerStateRecord = []string{""}
	}

	currentBoilerState := currentBoilerStateRecord[0]

	if currentBoilerState != "on" && currentBoilerState != "off" {
		currentBoilerState = "off"
	}

	boilerState := "off"
	if smartSwitchOn {
		if currentBoilerState == "off" && currentTemperature < thermostatThreshold-debounceBuffer {
			boilerState = "on"
		}
		if currentBoilerState == "on" && currentTemperature < thermostatThreshold {
			boilerState = "on"
		}
	}

	if logChange {
		if boilerState != currentBoilerState {
			err := timeseries.Append(b.config.BoilerStateLogFilePath, b.clock, boilerState, nil)
			if err != nil {
				b.loggers.Get("brain").Logf("error reading boiler state: %s", err)
			}

			b.loggers.Get("brain").Logf(
				"boiler state changing from %s to %s (smartSwitchOn: %t, currentTemperature: %.2f, thermostatThreshold: %.2f)",
				currentBoilerState, boilerState, smartSwitchOn, currentTemperature, thermostatThreshold,
			)

			currentBoilerState = boilerState
		}
	}

	return BoilerState{
		StateOfBoiler:         currentBoilerState,
		CalculatedBoilerState: boilerState,
	}
}

func (b Boiler) GetSmartSwitchStatus() bool {
	currentTime := b.clock.Now()

	timeValue, err := fileio.ReadFile(b.config.SmartSwitchLastAliveFilePath)
	if err != nil {
		timeValue = currentTime.Add(-100 * time.Hour).Format(time.RFC3339)
	}

	lastAliveTime, err := time.Parse(time.RFC3339, timeValue)
	smartSwitchOn := currentTime.Sub(lastAliveTime) < 6*time.Second
	return smartSwitchOn
}

func (b Boiler) GetThermostat() float64 {
	thermostatThresholdValue, err := fileio.ReadFile(b.config.CurrentThermostatThresholdFilePath)
	if err != nil {
		return defaultThermostatThreshold
	}

	thermostatThreshold, err := strconv.ParseFloat(thermostatThresholdValue, 64)
	if err != nil {
		return defaultThermostatThreshold
	}
	return thermostatThreshold
}

func (b Boiler) GetTemperature() float64 {
	temperatureValue, err := fileio.ReadFile(b.config.CurrentTemperatureFilePath)
	if err != nil {
		log.Fatal(err)
	}

	currentTemperature, err := strconv.ParseFloat(temperatureValue, 64)
	if err != nil {
		log.Fatal(err)
	}
	return currentTemperature
}
