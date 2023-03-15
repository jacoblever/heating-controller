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
	http.HandleFunc("/boiler-state/", BoilerStateHandler)
	http.HandleFunc("/smart-switch-alive/", SmartSwitchAliveHandler)

	fmt.Println("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

var currentTemperatureFilePath string = "./current-temperature.txt"
var smartSwitchLastAliveFilePath string = "./smart-switch-last-alive.txt"

type UpdateTemperatureResponse struct {
	// PollDelayMs is the number of milliseconds the Arduino should wait before making another request
	PollDelayMs                int
	ThermostatThresholdCelsius float32
}

func UpdateTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	temperature := r.URL.Query().Get("temperature")

	filePath := currentTemperatureFilePath

	writeToFile(filePath, temperature)

	response := UpdateTemperatureResponse{
		PollDelayMs:                1000,
		ThermostatThresholdCelsius: 24,
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
	filePath := currentTemperatureFilePath
	temperatureValue, err := readFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	currentTemperature, err := strconv.ParseFloat(temperatureValue, 32)
	if err != nil {
		log.Fatal(err)
	}

	timeValue, err := readFile(smartSwitchLastAliveFilePath)
	if err != nil {
		log.Fatal(err)
	}

	currentTime := time.Now()
	lastAliveTime, err := time.Parse(time.RFC3339, timeValue)

	boilerState := "off"
	if currentTemperature < 24 && (currentTime.Sub(lastAliveTime) < 3*time.Second) {
		boilerState = "on"
	}

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
