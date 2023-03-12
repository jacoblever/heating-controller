package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	// Register handler for default route
	http.HandleFunc("/update-temperature/", UpdateTemperatureHandler)
	http.HandleFunc("/boiler-state/", BoilerStateHandler)
	http.HandleFunc("/smart-switch-alive/", SmartSwitchAliveHandler)

	// Start server listening
	fmt.Println("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

var currentTemperatureFilePath string = "./current-temperature.txt"
var smartSwitchLastAliveFilePath string = "./smart-switch-last-alive.txt"

func UpdateTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	temperature := r.URL.Query().Get("temperature")

	filePath := currentTemperatureFilePath

	writeToFile(filePath, temperature)

	fmt.Fprintf(w, "Got temp: %s", temperature)
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

	fmt.Fprintf(w, "%s", boilerState)
}

func readFile(filePath string) (string, error) {
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(buffer), err
}

func SmartSwitchAliveHandler(w http.ResponseWriter, r *http.Request) {
	writeToFile(smartSwitchLastAliveFilePath, time.Now().Format(time.RFC3339))
}

func writeToFile(filePath string, value string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}

	n, err := f.WriteString(value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("wrote %d bytes\n", n)
	f.Sync()
}
