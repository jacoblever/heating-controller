package end2end

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/jacoblever/heating-controller/brain/brain"
	"github.com/stretchr/testify/assert"
)

type Context struct {
	context.Context

	BrainRouter *http.ServeMux
}

type Home struct {
	Boiler             Boiler
	Thermometer        Thermometer
	SmartSwitchAdapter SmartSwitchAdapter
	Dashboard          Dashboard
}

func CreateHome() Home {
	return Home{
		Boiler:             Boiler{},
		Thermometer:        Thermometer{},
		SmartSwitchAdapter: SmartSwitchAdapter{},
		Dashboard:          Dashboard{},
	}
}

type Thermometer struct {
}

func (th Thermometer) UpdateTemperature(t *testing.T, ctx Context, temperature float64) (rawResponse *httptest.ResponseRecorder, responseModel map[string]interface{}) {
	path := fmt.Sprintf("/update-temperature/?temperature=%f", temperature)
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	return SendTestRequestJSON(ctx.BrainRouter, request)
}

type Boiler struct {
	State    string
	Position int
}

func (b *Boiler) BoilerState(t *testing.T, ctx Context) (rawResponse *httptest.ResponseRecorder, responseModel map[string]interface{}) {
	path := "/boiler-state/"
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	response, model := SendTestRequestJSON(ctx.BrainRouter, request)
	assert.Equal(t, http.StatusOK, response.Code)

	newBoilerState := model["BoilerState"].(string)
	stepsToTurn := model["StepsToTurn"].(float64)
	command := model["Command"].(string)

	if newBoilerState == "on" {
		if b.State == "off" {
			b.Position = b.Position + int(stepsToTurn)
		}
	} else if newBoilerState == "off" {
		if b.State == "on" {
			b.Position = b.Position - int(stepsToTurn)
		}
	} else {
		t.Errorf("[Boiler] unknown boiler state: %s", newBoilerState)
	}
	b.State = newBoilerState

	if command != "" {
		steps, err := strconv.Atoi(command)
		if err != nil {
			t.Errorf("could not convert command (%s) to int: %s", command, err)
		}

		b.Position = b.Position + steps
	}

	return response, model
}

type SmartSwitchAdapter struct {
}

func (s SmartSwitchAdapter) SmartSwitchAlive(t *testing.T, ctx Context) (rawResponse *httptest.ResponseRecorder, responseModel map[string]interface{}) {
	path := "/smart-switch-alive/"
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	return SendTestRequestJSON(ctx.BrainRouter, request)
}

type Dashboard struct {
}

func (d Dashboard) Poll(t *testing.T, ctx Context) (rawResponse *httptest.ResponseRecorder, responseModel map[string]interface{}) {
	path := "/update-thermostat/"
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	return SendTestRequestJSON(ctx.BrainRouter, request)
}

func (d Dashboard) UpdateThermostat(t *testing.T, ctx Context, threshold float64) (rawResponse *httptest.ResponseRecorder, responseModel map[string]interface{}) {
	path := fmt.Sprintf("/update-thermostat/?threshold=%f", threshold)
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	return SendTestRequestJSON(ctx.BrainRouter, request)
}

func (d Dashboard) GetGraphData(t *testing.T, ctx Context) (rawResponse *httptest.ResponseRecorder, responseModel map[string]interface{}) {
	path := "/graph-data/"
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	return SendTestRequestJSON(ctx.BrainRouter, request)
}

func (d Dashboard) TurnBoiler(t *testing.T, ctx Context, command string) (rawResponse *httptest.ResponseRecorder, responseModel map[string]interface{}) {
	path := fmt.Sprintf("/turn-boiler/?command=%s", command)
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	response, model := SendTestRequestJSON(ctx.BrainRouter, request)
	assert.Equal(t, http.StatusOK, response.Code)
	return response, model
}

func CreateTestContext(t *testing.T) Context {
	router := brain.CreateRouter(brain.DefaultConfig)
	ctx := Context{
		Context:     context.Background(),
		BrainRouter: router,
	}

	t.Cleanup(func() {
		for _, f := range brain.DefaultConfig.AllFilePaths() {
			_ = os.Remove(f)
		}
	})
	return ctx
}

func SendTestRequestJSON(router *http.ServeMux, req *http.Request) (rawResponse *httptest.ResponseRecorder, responseModel map[string]interface{}) {
	rawResponse = httptest.NewRecorder()
	router.ServeHTTP(rawResponse, req)
	responseModel = nil
	_ = json.Unmarshal(rawResponse.Body.Bytes(), &responseModel)
	return rawResponse, responseModel
}
