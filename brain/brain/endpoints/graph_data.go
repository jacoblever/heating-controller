package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/fileio"
)

type TimePoint struct {
	Time  int64
	Value float64
}

type GraphDatResponse struct {
	Temperature  []TimePoint
	Temperature1 []TimePoint
	Temperature2 []TimePoint
	BoilerState  []TimePoint
}

func (h *Handlers) GraphDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	temperatureData := h.getData(h.config.TemperatureLogFilePath, func(v string) (float64, error) {
		return strconv.ParseFloat(v, 32)
	})
	temperature1Data := h.getData(h.config.TemperatureLog1FilePath, func(v string) (float64, error) {
		return strconv.ParseFloat(v, 32)
	})
	temperature2Data := h.getData(h.config.TemperatureLog2FilePath, func(v string) (float64, error) {
		return strconv.ParseFloat(v, 32)
	})
	boilerStateData := h.getData(h.config.BoilerStateLogFilePath, func(v string) (float64, error) {
		switch v {
		case "on":
			return 1, nil
		case "off":
			return 0, nil
		default:
			return 0, fmt.Errorf("boiler state %s invalid", v)
		}
	})

	response := GraphDatResponse{
		Temperature:  temperatureData,
		Temperature1: temperature1Data, // yellow
		Temperature2: temperature2Data, // orange
		BoilerState:  boilerStateData,
	}

	writeJSON(w, response)
}

func (h *Handlers) getData(filePath string, parseValue func(v string) (float64, error)) []TimePoint {
	data, _ := fileio.ReadFile(filePath)
	lines := strings.Split(data, "\n")
	var output []TimePoint

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) > 1 {
			t, err := time.Parse(time.RFC3339, parts[0])
			if err != nil {
				log.Printf("failed to parse timestamp in '%s' in %s", line, filePath)
				continue
			}
			value, err := parseValue(parts[1])
			if err != nil {
				log.Printf("failed to parse value in '%s' in %s", line, filePath)
				continue
			}
			if t.Before(time.Now().Add(-7 * 24 * time.Hour)) {
				// Ignore logs from more than a week ago
				continue
			}
			output = append(output, TimePoint{
				Time:  t.UnixMilli(),
				Value: value,
			})
		}
	}
	return output
}
