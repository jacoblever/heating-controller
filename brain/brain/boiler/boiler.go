package boiler

import (
	"log"
	"strconv"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/boiler/commandqueue"
	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/fileio"
	"github.com/jacoblever/heating-controller/brain/brain/logging"
	"github.com/jacoblever/heating-controller/brain/brain/stores"
	"github.com/jacoblever/heating-controller/brain/brain/timeseries"
)

var debounceBuffer = 0.25
var smartSwitchOffTimeout = 60 * time.Second

type Boiler struct {
	config  stores.Config
	clock   clock.Clock
	loggers logging.Loggers
	stores  stores.Stores

	BoilerCommandQueue commandqueue.CommandQueue
}

func MakeBoiler(config stores.Config, clock clock.Clock, loggers logging.Loggers, stores stores.Stores) Boiler {
	return Boiler{
		config:  config,
		clock:   clock,
		loggers: loggers,
		stores:  stores,

		BoilerCommandQueue: commandqueue.Make(),
	}
}

type BoilerState struct {
	StateOfBoiler         string
	CalculatedBoilerState string
}

func (b Boiler) GetBoilerState(logChange bool) BoilerState {
	currentTemperature := b.GetTemperature()
	thermostatThreshold := b.stores.Thermostat.GetLatestValueOrDefault()
	smartSwitchOn := b.GetSmartSwitchStatus()
	mode := b.stores.BoilerMode.GetLatestValueOrDefault()
	if mode == "" {
		mode = "auto"
	}

	currentBoilerStateRecord, err := b.stores.BoilerState.GetLatestValue()
	if err != nil {
		b.loggers.Get("brain").Logf("error reading current boiler state log: %s", err)
		currentBoilerStateRecord = timeseries.Value[timeseries.OnOff]{}
	}

	currentBoilerState := currentBoilerStateRecord.Value()

	boilerState := timeseries.Off
	if mode == "on" {
		boilerState = timeseries.On
	} else if mode == "off" {
		boilerState = timeseries.Off
	} else {
		if smartSwitchOn {
			if currentBoilerState == timeseries.Off && currentTemperature < thermostatThreshold-debounceBuffer {
				boilerState = timeseries.On
			}
			if currentBoilerState == timeseries.On && currentTemperature < thermostatThreshold {
				boilerState = timeseries.On
			}
		}
	}

	if logChange {
		if boilerState != currentBoilerState {
			err := b.stores.BoilerState.Store(boilerState)
			if err != nil {
				b.loggers.Get("brain").Logf("error writing boiler state: %s", err)
			}

			b.loggers.Get("brain").Logf(
				"boiler state changing from %s to %s (smartSwitchOn: %t, currentTemperature: %.2f, thermostatThreshold: %.2f)",
				currentBoilerState, boilerState, smartSwitchOn, currentTemperature, thermostatThreshold,
			)

			currentBoilerState = boilerState
		}
	}

	return BoilerState{
		StateOfBoiler:         currentBoilerState.OnOffString(),
		CalculatedBoilerState: boilerState.OnOffString(),
	}
}

func (b Boiler) GetSmartSwitchStatus() timeseries.OnOff {
	currentTime := b.clock.Now()

	timeValue, err := fileio.ReadFile(b.config.SmartSwitchLastAliveFilePath)
	if err != nil {
		timeValue = currentTime.Add(-100 * time.Hour).Format(time.RFC3339)
	}

	lastAliveTime, err := time.Parse(time.RFC3339, timeValue)
	smartSwitchOn := timeseries.OnOff(currentTime.Sub(lastAliveTime) < smartSwitchOffTimeout)
	err = b.stores.SmartSwitch.Store(smartSwitchOn)
	if err != nil {
		b.loggers.Get("brain").Logf("failed to write smart switch state '%b': %s", smartSwitchOn, err.Error())
	}
	return smartSwitchOn
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
