package main

import (
	"fmt"
	"net/http"

	"github.com/jacoblever/heating-controller/brain/brain"
	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/logging"
	"github.com/jacoblever/heating-controller/brain/brain/stores"
	"github.com/jacoblever/heating-controller/brain/brain/weatherapi"
)

var port = 8080

func main() {
	clock := clock.CreateClock()
	slackLogger := logging.CreateSlackLogger()
	loggers := logging.InitLoggers(clock, slackLogger)
	weatherAPI := weatherapi.Make()

	router := brain.CreateRouter(stores.DefaultConfig, clock, loggers, weatherAPI)

	fmt.Println(fmt.Sprintf("Listening on port %d...", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		panic(err)
	}
}
