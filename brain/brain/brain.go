package brain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/fileio"
	"github.com/jacoblever/heating-controller/brain/brain/logging"
)

var debounceBuffer = 0.5
var defaultThermostatThreshold float64 = 18
var boilerSwitchStepCount = 250

var DefaultConfig Config = Config{
	CurrentTemperatureFilePath:         "./current-temperature.txt",
	TemperatureLogFilePath:             "./temperature-log.txt",
	CurrentThermostatThresholdFilePath: "./current-thermostat-threshold.txt",
	SmartSwitchLastAliveFilePath:       "./smart-switch-last-alive.txt",
	BoilerStateFilePath:                "./boiler-state.txt",
	BoilerStateLogFilePath:             "./boiler-state-log.txt",
}

type handlers struct {
	config             Config
	clock              clock.Clock
	boilerCommandQueue []string
	logger             logging.Logger
}

func CreateRouter(config Config, c clock.Clock, logger logging.Logger) *http.ServeMux {
	router := http.NewServeMux()
	if c == nil {
		c = clock.CreateClock()
	}
	handlers := handlers{config: config, clock: c, boilerCommandQueue: make([]string, 0), logger: logger}
	router.HandleFunc("/update-temperature/", handlers.UpdateTemperatureHandler)
	router.HandleFunc("/update-thermostat/", handlers.UpdateThermostatHandler)
	router.HandleFunc("/boiler-state/", handlers.BoilerStateHandler)
	router.HandleFunc("/smart-switch-alive/", handlers.SmartSwitchAliveHandler)
	router.HandleFunc("/turn-boiler/", handlers.TurnBoilerHandler)
	router.HandleFunc("/graph-data/", handlers.GraphDataHandler)
	return router
}

type Config struct {
	CurrentTemperatureFilePath         string
	TemperatureLogFilePath             string
	CurrentThermostatThresholdFilePath string
	SmartSwitchLastAliveFilePath       string
	BoilerStateFilePath                string
	BoilerStateLogFilePath             string
}

func (c Config) AllFilePaths() []string {
	return []string{
		c.CurrentTemperatureFilePath,
		c.TemperatureLogFilePath,
		c.CurrentThermostatThresholdFilePath,
		c.SmartSwitchLastAliveFilePath,
		c.BoilerStateFilePath,
		c.BoilerStateLogFilePath,
	}
}

type UpdateTemperatureResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs                int
	ThermostatThresholdCelsius float64
}

func (h handlers) UpdateTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	temperature := r.URL.Query().Get("temperature")

	fileio.WriteToFile(h.config.CurrentTemperatureFilePath, temperature)

	lastTemp, err := fileio.ReadLastLine(h.config.TemperatureLogFilePath)
	if err != nil {
		log.Printf("error reading: %s", err)
		lastTemp = ""
	}
	lastTimeStr := strings.Split(lastTemp, ",")[0]
	lastTime, err := time.Parse(time.RFC3339, lastTimeStr)
	if err != nil {
		log.Printf("error parsing time: %s", err)
		lastTime = h.clock.Now().Add(-24 * time.Hour)
	}

	if lastTime.Add(10 * time.Minute).Before(h.clock.Now()) {
		fileio.AppendLineToFile(h.config.TemperatureLogFilePath, strings.Join([]string{h.clock.Now().Format(time.RFC3339), temperature}, ","))
	}

	response := UpdateTemperatureResponse{
		PollDelayMs:                1000,
		ThermostatThresholdCelsius: h.getThermostat(),
	}
	writeJSON(w, response)
}

type UpdateThermostatResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs                int
	BoilerState                string
	SmartSwitchOn              bool
	TemperatureCelsius         float64
	ThermostatThresholdCelsius float64
}

func (h handlers) UpdateThermostatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	threshold := r.URL.Query().Get("threshold")
	if threshold != "" {
		if _, err := strconv.ParseFloat(threshold, 32); err == nil {
			fileio.WriteToFile(h.config.CurrentThermostatThresholdFilePath, threshold)
			h.logger.Log(fmt.Sprintf("Thermostat set to %s", threshold))
		}
	}

	boilerState := h.getBoilerState(false)

	response := UpdateThermostatResponse{
		PollDelayMs:                1000,
		BoilerState:                boilerState,
		SmartSwitchOn:              h.getSmartSwitchStatus(),
		TemperatureCelsius:         h.getTemperature(),
		ThermostatThresholdCelsius: h.getThermostat(),
	}
	writeJSON(w, response)
}

type BoilerStateResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs   int
	BoilerState   string
	MotorSpeedRPM int
	StepsToTurn   int
	Command       string
}

func (h *handlers) BoilerStateHandler(w http.ResponseWriter, r *http.Request) {
	boilerState := h.getBoilerState(true)

	response := BoilerStateResponse{
		PollDelayMs:   1000,
		BoilerState:   boilerState,
		MotorSpeedRPM: 4,
		StepsToTurn:   boilerSwitchStepCount,
		Command:       h.getNextBoilerCommand(),
	}
	writeJSON(w, response)
}

func (h *handlers) getNextBoilerCommand() string {
	if len(h.boilerCommandQueue) == 0 {
		return ""
	}

	command := h.boilerCommandQueue[0]
	h.boilerCommandQueue = h.boilerCommandQueue[1:]
	return command
}

type SmartSwitchAliveResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs int
}

func (h handlers) SmartSwitchAliveHandler(w http.ResponseWriter, r *http.Request) {
	fileio.WriteToFile(h.config.SmartSwitchLastAliveFilePath, h.clock.Now().Format(time.RFC3339))

	response := SmartSwitchAliveResponse{
		PollDelayMs: 1000,
	}
	writeJSON(w, response)
}

func (h *handlers) TurnBoilerHandler(w http.ResponseWriter, r *http.Request) {
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

type TimePoint struct {
	Time  int64
	Value float64
}

type GraphDatResponse struct {
	Temperature []TimePoint
	BoilerState []TimePoint
}

func (h *handlers) GraphDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	temperatureData := h.getData(h.config.TemperatureLogFilePath, func(v string) (float64, error) {
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

	response := GraphDatResponse{Temperature: temperatureData, BoilerState: boilerStateData}

	writeJSON(w, response)
}

func (h *handlers) getData(filePath string, parseValue func(v string) (float64, error)) []TimePoint {
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
			output = append(output, TimePoint{
				Time:  t.UnixMilli(),
				Value: value,
			})
		}
	}
	return output
}

func (h handlers) getBoilerState(logChange bool) string {
	currentTemperature := h.getTemperature()
	thermostatThreshold := h.getThermostat()
	smartSwitchOn := h.getSmartSwitchStatus()

	currentBoilerState, err := fileio.ReadFile(h.config.BoilerStateFilePath)
	if err != nil {
		h.logger.Logf("error reading boiler state: %s", err)
	}
	if currentBoilerState == "" {
		currentBoilerState = "off"
	}

	boilerState := "off"
	if smartSwitchOn {
		if currentBoilerState == "off" && currentTemperature < thermostatThreshold-debounceBuffer {
			boilerState = "on"
		}
		if currentBoilerState == "on" && currentTemperature < thermostatThreshold {
			boilerState = "on"
		}
	}

	err = fileio.WriteToFile(h.config.BoilerStateFilePath, boilerState)
	if err != nil {
		h.logger.Logf("error reading boiler state: %s", err)
	}

	if logChange {
		if boilerState != currentBoilerState {
			err := fileio.AppendLineToFile(h.config.BoilerStateLogFilePath, strings.Join([]string{h.clock.Now().Format(time.RFC3339), boilerState}, ","))
			if err != nil {
				h.logger.Logf("error reading boiler state: %s", err)
			}
		}
	}

	return boilerState
}

func (h handlers) getSmartSwitchStatus() bool {
	currentTime := h.clock.Now()

	timeValue, err := fileio.ReadFile(h.config.SmartSwitchLastAliveFilePath)
	if err != nil {
		timeValue = currentTime.Add(-100 * time.Hour).Format(time.RFC3339)
	}

	lastAliveTime, err := time.Parse(time.RFC3339, timeValue)
	smartSwitchOn := currentTime.Sub(lastAliveTime) < 3*time.Second
	return smartSwitchOn
}

func (h handlers) getThermostat() float64 {
	thermostatThresholdValue, err := fileio.ReadFile(h.config.CurrentThermostatThresholdFilePath)
	if err != nil {
		return defaultThermostatThreshold
	}

	thermostatThreshold, err := strconv.ParseFloat(thermostatThresholdValue, 64)
	if err != nil {
		return defaultThermostatThreshold
	}
	return thermostatThreshold
}

func (h handlers) getTemperature() float64 {
	temperatureValue, err := fileio.ReadFile(h.config.CurrentTemperatureFilePath)
	if err != nil {
		log.Fatal(err)
	}

	currentTemperature, err := strconv.ParseFloat(temperatureValue, 64)
	if err != nil {
		log.Fatal(err)
	}
	return currentTemperature
}

func writeJSON(w http.ResponseWriter, response any) {
	jData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err2 := w.Write([]byte("{\"error\": \"marshal error\"}"))
		if err2 != nil {
			log.Fatal(err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jData)
	if err != nil {
		log.Fatal(err)
	}
}
