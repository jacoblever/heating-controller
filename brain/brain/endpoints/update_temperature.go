package endpoints

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jacoblever/heating-controller/brain/brain/fileio"
)

type UpdateTemperatureResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs                int
	ThermostatThresholdCelsius float64
}

func (h Handlers) UpdateTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	temperature := r.URL.Query().Get("temperature")
	float, err := strconv.ParseFloat(temperature, 64)
	if err != nil {
		writeErrorWithStatus(w, fmt.Errorf("invalid temperature %s: %s", temperature, err.Error()), http.StatusBadRequest)
		return
	}

	fileio.WriteToFile(h.config.CurrentTemperatureFilePath, temperature)

	err = h.stores.Temperature.Store(float)
	if err != nil {
		writeError(w, err)
		return
	}

	response := UpdateTemperatureResponse{
		PollDelayMs:                1000,
		ThermostatThresholdCelsius: h.stores.Thermostat.GetLatestValueOrDefault(),
	}
	writeJSON(w, response)
}
