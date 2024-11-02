package endpoints

import (
	"net/http"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/fileio"
)

type SmartSwitchAliveResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs int
}

func (h Handlers) SmartSwitchAliveHandler(w http.ResponseWriter, r *http.Request) {
	fileio.WriteToFile(h.config.SmartSwitchLastAliveFilePath, h.clock.Now().Format(time.RFC3339))

	response := SmartSwitchAliveResponse{
		PollDelayMs: 1000,
	}
	writeJSON(w, response)
}
