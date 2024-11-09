package endpoints

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jacoblever/heating-controller/brain/brain/timeseries"
)

type TemperatureResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs int
}

func (h Handlers) TemperatureHandler(w http.ResponseWriter, r *http.Request) {
	temperature := r.URL.Query().Get("temperature")
	id := r.URL.Query().Get("id")

	temperatureStore, err := h.getTemperatureStore(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	float, err := strconv.ParseFloat(temperature, 64)
	if err != nil {
		writeErrorWithStatus(w, fmt.Errorf("invalid temperature %s: %s", temperature, err.Error()), http.StatusBadRequest)
		return
	}

	err = temperatureStore.Store(float)
	if err != nil {
		writeError(w, err)
		return
	}

	response := TemperatureResponse{
		PollDelayMs: 1000,
	}
	writeJSON(w, response)
}

func (h Handlers) getTemperatureStore(id string) (timeseries.ValueStore[float64], error) {
	switch id {
	case "1":
		return h.stores.Temperature1, nil
	case "2":
		return h.stores.Temperature2, nil
	default:
		return nil, fmt.Errorf("getTemperatureLogFilePath: unknown device id %s", id)
	}
}
