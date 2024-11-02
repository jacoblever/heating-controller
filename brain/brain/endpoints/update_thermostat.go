package endpoints

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jacoblever/heating-controller/brain/brain/fileio"
)

type UpdateThermostatResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs                int
	StateOfBoiler              string
	CalculatedBoilerState      string
	SmartSwitchOn              bool
	TemperatureCelsius         float64
	ThermostatThresholdCelsius float64
}

func (h Handlers) UpdateThermostatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	threshold := r.URL.Query().Get("threshold")
	if threshold != "" {
		if _, err := strconv.ParseFloat(threshold, 32); err == nil {
			fileio.WriteToFile(h.config.CurrentThermostatThresholdFilePath, threshold)
			h.loggers.SlackLogger.Log(fmt.Sprintf("Thermostat set to %s", threshold))
		}
	}

	boilerState := h.boiler.GetBoilerState(false)

	response := UpdateThermostatResponse{
		PollDelayMs:                1000,
		StateOfBoiler:              boilerState.StateOfBoiler,
		CalculatedBoilerState:      boilerState.CalculatedBoilerState,
		SmartSwitchOn:              h.boiler.GetSmartSwitchStatus(),
		TemperatureCelsius:         h.boiler.GetTemperature(),
		ThermostatThresholdCelsius: h.boiler.GetThermostat(),
	}
	writeJSON(w, response)
}
