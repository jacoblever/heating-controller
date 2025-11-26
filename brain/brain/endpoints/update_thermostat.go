package endpoints

import (
	"fmt"
	"net/http"
	"strconv"
)

type UpdateThermostatResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs                int
	StateOfBoiler              string
	CalculatedBoilerState      string
	SmartSwitchOn              bool
	TemperatureCelsius         float64
	ThermostatThresholdCelsius float64
	BoilerMode                 string
}

func (h Handlers) UpdateThermostatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	threshold := r.URL.Query().Get("threshold")
	if threshold != "" {
		if float, err := strconv.ParseFloat(threshold, 64); err == nil {
			err := h.stores.Thermostat.Store(float)
			if err != nil {
				h.loggers.Get("brain").Logf("error writing thermostat value: %s", err.Error())
			}
			h.loggers.SlackLogger.Log(fmt.Sprintf("Thermostat set to %s", threshold))
		}
	}

	mode := r.URL.Query().Get("mode")
	if mode != "" {
		err := h.stores.BoilerMode.Store(mode)
		if err != nil {
			h.loggers.Get("brain").Logf("error writing boiler mode value: %s", err.Error())
		}
		h.loggers.SlackLogger.Log(fmt.Sprintf("Boiler mode set to %s", mode))
	}

	boilerState := h.boiler.GetBoilerState(false)

	boilerMode := h.stores.BoilerMode.GetLatestValueOrDefault()
	if boilerMode == "" {
		boilerMode = "auto"
	}

	response := UpdateThermostatResponse{
		PollDelayMs:                1000,
		StateOfBoiler:              boilerState.StateOfBoiler,
		CalculatedBoilerState:      boilerState.CalculatedBoilerState,
		SmartSwitchOn:              h.boiler.GetSmartSwitchStatus().Bool(),
		TemperatureCelsius:         h.boiler.GetTemperature(),
		ThermostatThresholdCelsius: h.stores.Thermostat.GetLatestValueOrDefault(),
		BoilerMode:                 boilerMode,
	}
	writeJSON(w, response)
}
