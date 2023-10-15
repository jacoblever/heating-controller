package brain

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/clock"
)

var debounceBuffer = 0.5
var defaultThermostatThreshold float64 = 18

var DefaultConfig Config = Config{
	CurrentTemperatureFilePath:         "./current-temperature.txt",
	TemperatureLogFilePath:             "./temperature-log.txt",
	CurrentThermostatThresholdFilePath: "./current-thermostat-threshold.txt",
	SmartSwitchLastAliveFilePath:       "./smart-switch-last-alive.txt",
	BoilerStateFilePath:                "./boiler-state.txt",
}

type handlers struct {
	config Config
	clock  clock.Clock
}

func CreateRouter(config Config, c clock.Clock) *http.ServeMux {
	router := http.NewServeMux()
	if c == nil {
		c = clock.CreateClock()
	}
	handlers := handlers{config: config, clock: c}
	router.HandleFunc("/update-temperature/", handlers.UpdateTemperatureHandler)
	router.HandleFunc("/update-thermostat/", handlers.UpdateThermostatHandler)
	router.HandleFunc("/boiler-state/", handlers.BoilerStateHandler)
	router.HandleFunc("/smart-switch-alive/", handlers.SmartSwitchAliveHandler)
	return router
}

type Config struct {
	CurrentTemperatureFilePath         string
	TemperatureLogFilePath             string
	CurrentThermostatThresholdFilePath string
	SmartSwitchLastAliveFilePath       string
	BoilerStateFilePath                string
}

func (c Config) AllFilePaths() []string {
	return []string{
		c.CurrentTemperatureFilePath,
		c.TemperatureLogFilePath,
		c.CurrentThermostatThresholdFilePath,
		c.SmartSwitchLastAliveFilePath,
		c.BoilerStateFilePath,
	}
}

type UpdateTemperatureResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs                int
	ThermostatThresholdCelsius float64
}

func (h handlers) UpdateTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	temperature := r.URL.Query().Get("temperature")

	writeToFile(h.config.CurrentTemperatureFilePath, temperature)

	lastTemp, err := readLastLine(h.config.TemperatureLogFilePath)
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
		appendLineToFile(h.config.TemperatureLogFilePath, strings.Join([]string{h.clock.Now().Format(time.RFC3339), temperature}, ","))
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
			writeToFile(h.config.CurrentThermostatThresholdFilePath, threshold)
		}
	}

	response := UpdateThermostatResponse{
		PollDelayMs:                1000,
		BoilerState:                h.getBoilerState(),
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
}

func (h handlers) BoilerStateHandler(w http.ResponseWriter, r *http.Request) {
	boilerState := h.getBoilerState()

	response := BoilerStateResponse{
		PollDelayMs:   1000,
		BoilerState:   boilerState,
		MotorSpeedRPM: 4,
		StepsToTurn:   300,
	}
	writeJSON(w, response)
}

type SmartSwitchAliveResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs int
}

func (h handlers) SmartSwitchAliveHandler(w http.ResponseWriter, r *http.Request) {
	writeToFile(h.config.SmartSwitchLastAliveFilePath, h.clock.Now().Format(time.RFC3339))

	response := SmartSwitchAliveResponse{
		PollDelayMs: 1000,
	}
	writeJSON(w, response)
}

func (h handlers) getBoilerState() string {
	currentTemperature := h.getTemperature()
	thermostatThreshold := h.getThermostat()
	smartSwitchOn := h.getSmartSwitchStatus()

	currentBoilerState, _ := readFile(h.config.BoilerStateFilePath)
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

	writeToFile(h.config.BoilerStateFilePath, boilerState)
	return boilerState
}

func (h handlers) getSmartSwitchStatus() bool {
	currentTime := h.clock.Now()

	timeValue, err := readFile(h.config.SmartSwitchLastAliveFilePath)
	if err != nil {
		timeValue = currentTime.Add(-100 * time.Hour).Format(time.RFC3339)
	}

	lastAliveTime, err := time.Parse(time.RFC3339, timeValue)
	smartSwitchOn := currentTime.Sub(lastAliveTime) < 3*time.Second
	return smartSwitchOn
}

func (h handlers) getThermostat() float64 {
	thermostatThresholdValue, err := readFile(h.config.CurrentThermostatThresholdFilePath)
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
	temperatureValue, err := readFile(h.config.CurrentTemperatureFilePath)
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
