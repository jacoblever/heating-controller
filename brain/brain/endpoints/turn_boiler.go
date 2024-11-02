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
			h.boilerCommandQueue = append(h.boilerCommandQueue, fmt.Sprintf("%d", boilerSwitchStepCount))
		case "turn-anticlockwise":
			h.boilerCommandQueue = append(h.boilerCommandQueue, fmt.Sprintf("-%d", boilerSwitchStepCount))
		default:
			h.boilerCommandQueue = append(h.boilerCommandQueue, command)
		}
	}

	writeJSON(w, struct{}{})
}
