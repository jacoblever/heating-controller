package brain

import (
	"net/http"

	"github.com/jacoblever/heating-controller/brain/brain/boiler"
	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/endpoints"
	"github.com/jacoblever/heating-controller/brain/brain/logging"
	"github.com/jacoblever/heating-controller/brain/brain/stores"
)

func CreateRouter(config stores.Config, c clock.Clock, loggers logging.Loggers) *http.ServeMux {
	router := http.NewServeMux()
	if c == nil {
		c = clock.CreateClock()
	}

	stores := stores.MakeStores(c, config)
	boiler := boiler.MakeBoiler(config, c, loggers, stores)
	handlers := endpoints.MakeHandlers(config, c, loggers, stores, boiler)

	router.HandleFunc("/update-temperature/", handlers.UpdateTemperatureHandler)
	router.HandleFunc("/temperature/", handlers.TemperatureHandler)
	router.HandleFunc("/update-thermostat/", handlers.UpdateThermostatHandler)
	router.HandleFunc("/boiler-state/", handlers.BoilerStateHandler)
	router.HandleFunc("/smart-switch-alive/", handlers.SmartSwitchAliveHandler)
	router.HandleFunc("/turn-boiler/", handlers.TurnBoilerHandler)
	router.HandleFunc("/graph-data/", handlers.GraphDataHandler)
	router.HandleFunc("/logs/", handlers.LogsHandler)

	loggers.NewPerDayLogger("boiler", logging.Settings{DaysToKeepFor: 14})
	loggers.NewPerDayLogger("brain", logging.Settings{DaysToKeepFor: 14})
	return router
}
