package endpoints

import (
	"net/http"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/fileio"
	"github.com/jacoblever/heating-controller/brain/brain/timeseries"
)

type UpdateTemperatureResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs                int
	ThermostatThresholdCelsius float64
}

func (h Handlers) UpdateTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	temperature := r.URL.Query().Get("temperature")

	fileio.WriteToFile(h.config.CurrentTemperatureFilePath, temperature)

	every := 10 * time.Minute
	timeseries.Append(h.config.TemperatureLogFilePath, h.clock, temperature, &every)

	response := UpdateTemperatureResponse{
		PollDelayMs:                1000,
		ThermostatThresholdCelsius: h.boiler.GetThermostat(),
	}
	writeJSON(w, response)
}
