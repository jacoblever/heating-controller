package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jacoblever/heating-controller/brain/brain/boiler"
	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/logging"
)

type Handlers struct {
	config  boiler.Config
	clock   clock.Clock
	loggers logging.Loggers
	boiler  boiler.Boiler
}

func MakeHandlers(config boiler.Config, clock clock.Clock, loggers logging.Loggers, boiler boiler.Boiler) Handlers {
	return Handlers{
		config:  config,
		clock:   clock,
		loggers: loggers,
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
