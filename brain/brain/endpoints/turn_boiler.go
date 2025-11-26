package endpoints

import (
	"fmt"
	"net/http"
)

func (h *Handlers) TurnBoilerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	command := r.URL.Query().Get("command")
	if command != "" {
		switch command {
		case "turn-clockwise":
			h.boiler.BoilerCommandQueue.Add(fmt.Sprintf("%d", boilerSwitchStepCountOn))
		case "turn-anticlockwise":
			h.boiler.BoilerCommandQueue.Add(fmt.Sprintf("-%d", boilerSwitchStepCountOff))
		default:
			h.boiler.BoilerCommandQueue.Add(command)
		}
	}

	writeJSON(w, struct{}{})
}
