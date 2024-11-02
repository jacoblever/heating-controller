package endpoints

import "net/http"

var boilerSwitchStepCount = 250

type BoilerStateResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs   int
	BoilerState   string
	MotorSpeedRPM int
	StepsToTurn   int
	Command       string
}

func (h *Handlers) BoilerStateHandler(w http.ResponseWriter, r *http.Request) {
	boilerState := h.boiler.GetBoilerState(true)
	err := r.ParseForm()
	if err != nil {
		h.loggers.SlackLogger.Log("[BoilerStateHandler] failed to parse form from request")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		logLines := r.Form["Log"]
		for _, line := range logLines {
			h.loggers.Get("boiler").Log(line)
		}
	}

	response := BoilerStateResponse{
		PollDelayMs:   1000,
		BoilerState:   boilerState.CalculatedBoilerState,
		MotorSpeedRPM: 4,
		StepsToTurn:   boilerSwitchStepCount,
		Command:       h.boiler.BoilerCommandQueue.Pop(),
	}
	writeJSON(w, response)
}
