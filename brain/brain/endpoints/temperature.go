package endpoints

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/timeseries"
)

type TemperatureResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs int
}

func (h Handlers) TemperatureHandler(w http.ResponseWriter, r *http.Request) {
	temperature := r.URL.Query().Get("temperature")
	id := r.URL.Query().Get("id")

	filePath, err := h.getTemperatureLogFilePath(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	every := 10 * time.Minute
	timeseries.Append(filePath, h.clock, temperature, &every)

	response := TemperatureResponse{
		PollDelayMs: 1000,
	}
	writeJSON(w, response)
}

func (h Handlers) getTemperatureLogFilePath(id string) (string, error) {
	switch id {
	case "1":
		return h.config.TemperatureLog1FilePath, nil
	case "2":
		return h.config.TemperatureLog2FilePath, nil
	default:
		return "", fmt.Errorf("getTemperatureLogFilePath: unknown device id %s", id)
	}
}
