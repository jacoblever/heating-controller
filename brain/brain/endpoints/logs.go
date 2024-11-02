package endpoints

import (
	"fmt"
	"net/http"
)

type LogLine struct {
	Time    int64
	Message string
}

type LogsResponse struct {
	Boiler []LogLine
	Brain  []LogLine
}

func (h *Handlers) LogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	response := LogsResponse{
		Boiler: h.getLogs("boiler"),
		Brain:  h.getLogs("brain"),
	}

	writeJSON(w, response)
}

func (h *Handlers) getLogs(key string) []LogLine {
	lines, err := h.loggers.Get(key).GetLogs()
	if err != nil {
		return []LogLine{{Message: fmt.Sprintf("error getting logger: %s", err.Error())}}
	}

	var logLines []LogLine = []LogLine{}
	for _, line := range lines {
		logLines = append(logLines, LogLine{Message: line})
	}
	return logLines
}
