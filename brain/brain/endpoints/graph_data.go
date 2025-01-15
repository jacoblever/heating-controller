package endpoints

import (
	"net/http"
	"strconv"

	"github.com/jacoblever/heating-controller/brain/brain/timeseries"
)

type TimePoint struct {
	Time  int64
	Value float64
}

type GraphDatResponse struct {
	Temperature        []TimePoint
	Temperature1       []TimePoint
	Temperature2       []TimePoint
	Thermostat         []TimePoint
	SmartSwitchState   []TimePoint
	BoilerState        []TimePoint
	OutsideTemperature []TimePoint
}

func (h *Handlers) GraphDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	days, err := strconv.Atoi(r.URL.Query().Get("days"))
	if err != nil {
		writeErrorWithStatus(w, err, http.StatusBadRequest)
		return
	}

	temperatureData, err := h.getFloatData(h.stores.Temperature, days)
	if err != nil {
		writeError(w, err)
		return
	}

	temperature1Data, err := h.getFloatData(h.stores.Temperature1, days)
	if err != nil {
		writeError(w, err)
		return
	}

	temperature2Data, err := h.getFloatData(h.stores.Temperature2, days)
	if err != nil {
		writeError(w, err)
		return
	}

	thermostatData, err := h.getFloatData(h.stores.Thermostat, days)
	if err != nil {
		writeError(w, err)
		return
	}

	thermostatData = append(thermostatData, TimePoint{
		Time:  h.clock.Now().UnixMilli(),
		Value: h.stores.Thermostat.GetLatestValueOrDefault(),
	})

	smartSwitchData, err := h.getOnOffData(h.stores.SmartSwitch, days)
	if err != nil {
		writeError(w, err)
		return
	}

	boilerStateData, err := h.getOnOffData(h.stores.BoilerState, days)
	if err != nil {
		writeError(w, err)
		return
	}

	outsideTemperatureData, err := h.getFloatData(h.stores.OutsideTemperature, days)
	if err != nil {
		writeError(w, err)
		return
	}

	response := GraphDatResponse{
		Temperature:        temperatureData,
		Temperature1:       temperature1Data, // yellow
		Temperature2:       temperature2Data, // orange
		Thermostat:         thermostatData,
		SmartSwitchState:   smartSwitchData,
		BoilerState:        boilerStateData,
		OutsideTemperature: outsideTemperatureData,
	}

	writeJSON(w, response)
}

func (h *Handlers) getFloatData(store timeseries.ValueStore[float64], lastXdaysOnly int) ([]TimePoint, error) {
	values, err := store.GetAll(&lastXdaysOnly)
	if err != nil {
		return nil, err
	}

	timePoints := []TimePoint{}
	for _, v := range values {
		timePoints = append(timePoints, TimePoint{
			Time:  v.Timestamp().UnixMilli(),
			Value: v.Value(),
		})
	}

	return timePoints, nil
}

func (h *Handlers) getOnOffData(store timeseries.ValueStore[timeseries.OnOff], lastXdaysOnly int) ([]TimePoint, error) {
	values, err := store.GetAll(&lastXdaysOnly)
	if err != nil {
		return nil, err
	}

	timePoints := []TimePoint{}
	for _, v := range values {
		timePoints = append(timePoints, TimePoint{
			Time:  v.Timestamp().UnixMilli(),
			Value: v.Value().OneOrZero(),
		})
	}

	if timePoints[len(timePoints)-1].Value == timeseries.On.OneOrZero() {
		timePoints = append(timePoints, TimePoint{
			Time:  h.clock.Now().UnixMilli(),
			Value: timeseries.Off.OneOrZero(),
		})
	}

	return timePoints, nil
}
