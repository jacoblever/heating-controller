package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jacoblever/heating-controller/brain/brain/boiler"
	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/logging"
	"github.com/jacoblever/heating-controller/brain/brain/stores"
)

type Handlers struct {
	config  stores.Config
	clock   clock.Clock
	loggers logging.Loggers
	stores  stores.Stores
	boiler  boiler.Boiler
}

func MakeHandlers(config stores.Config, clock clock.Clock, loggers logging.Loggers, stores stores.Stores, boiler boiler.Boiler) Handlers {
	return Handlers{
		config:  config,
		clock:   clock,
		loggers: loggers,
		stores:  stores,
		boiler:  boiler,
	}
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

type ErrorResponse struct {
	error string
}

func writeError(w http.ResponseWriter, err error) {
	writeErrorWithStatus(w, err, http.StatusInternalServerError)
}

func writeErrorWithStatus(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	response := ErrorResponse{error: err.Error()}
	writeJSON(w, response)
}
