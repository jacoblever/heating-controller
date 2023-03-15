package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/update-temperature/", UpdateTemperatureHandler)
	http.HandleFunc("/update-thermostat/", UpdateThermostatHandler)
	http.HandleFunc("/boiler-state/", BoilerStateHandler)
	http.HandleFunc("/smart-switch-alive/", SmartSwitchAliveHandler)

	fmt.Println("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

var currentTemperatureFilePath string = "./current-temperature.txt"
var currentThermostatThresholdFilePath string = "./current-thermostat-threshold.txt"
var smartSwitchLastAliveFilePath string = "./smart-switch-last-alive.txt"

type UpdateTemperatureResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs                int
	ThermostatThresholdCelsius float64
}

func UpdateTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	temperature := r.URL.Query().Get("temperature")

	writeToFile(currentTemperatureFilePath, temperature)

	response := UpdateTemperatureResponse{
		PollDelayMs:                1000,
		ThermostatThresholdCelsius: getThermostat(),
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

func UpdateThermostatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	threshold := r.URL.Query().Get("threshold")
	if threshold != "" {
		if _, err := strconv.ParseFloat(threshold, 32); err == nil {
			writeToFile(currentThermostatThresholdFilePath, threshold)
		}
	}

	response := UpdateThermostatResponse{
		PollDelayMs:                1000,
		BoilerState:                getBoilerState(),
		SmartSwitchOn:              getSmartSwitchStatus(),
		TemperatureCelsius:         getTemperature(),
		ThermostatThresholdCelsius: getThermostat(),
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

func BoilerStateHandler(w http.ResponseWriter, r *http.Request) {
	boilerState := getBoilerState()

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

func SmartSwitchAliveHandler(w http.ResponseWriter, r *http.Request) {
	writeToFile(smartSwitchLastAliveFilePath, time.Now().Format(time.RFC3339))

	response := SmartSwitchAliveResponse{
		PollDelayMs: 1000,
	}
	writeJSON(w, response)
}

func getBoilerState() string {
	currentTemperature := getTemperature()
	thermostatThreshold := getThermostat()
	smartSwitchOn := getSmartSwitchStatus()

	boilerState := "off"
	if currentTemperature < thermostatThreshold && smartSwitchOn {
		boilerState = "on"
	}
	return boilerState
}

func getSmartSwitchStatus() bool {
	timeValue, err := readFile(smartSwitchLastAliveFilePath)
	if err != nil {
		log.Fatal(err)
	}

	currentTime := time.Now()
	lastAliveTime, err := time.Parse(time.RFC3339, timeValue)
	smartSwitchOn := currentTime.Sub(lastAliveTime) < 3*time.Second
	return smartSwitchOn
}

func getThermostat() float64 {
	thermostatThresholdValue, err := readFile(currentThermostatThresholdFilePath)
	if err != nil {
		return 24
	}

	thermostatThreshold, err := strconv.ParseFloat(thermostatThresholdValue, 64)
	if err != nil {
		return 24
	}
	return thermostatThreshold
}

func getTemperature() float64 {
	temperatureValue, err := readFile(currentTemperatureFilePath)
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
