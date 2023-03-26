package end2end

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jacoblever/heating-controller/brain/brain"
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

func (th Thermometer) UpdateTemperature(t *testing.T, ctx Context, temperature float64) (rawResponse *httptest.ResponseRecorder, responseModel map[string]string) {
	path := fmt.Sprintf("/update-temperature/?temperature=%f", temperature)
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	return SendTestRequestJSON(ctx.BrainRouter, request)
}

type Boiler struct {
}

func (b Boiler) BoilerState(t *testing.T, ctx Context) (rawResponse *httptest.ResponseRecorder, responseModel map[string]string) {
	path := "/boiler-state/"
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	requestJSON, model := SendTestRequestJSON(ctx.BrainRouter, request)
	return requestJSON, model
}

type SmartSwitchAdapter struct {
}

func (s SmartSwitchAdapter) SmartSwitchAlive(t *testing.T, ctx Context) (rawResponse *httptest.ResponseRecorder, responseModel map[string]string) {
	path := "/smart-switch-alive/"
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	return SendTestRequestJSON(ctx.BrainRouter, request)
}

type Dashboard struct {
}

func (d Dashboard) Poll(t *testing.T, ctx Context) (rawResponse *httptest.ResponseRecorder, responseModel map[string]string) {
	path := "/update-thermostat/"
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	return SendTestRequestJSON(ctx.BrainRouter, request)
}

func (d Dashboard) UpdateThermostat(t *testing.T, ctx Context, threshold float64) (rawResponse *httptest.ResponseRecorder, responseModel map[string]string) {
	path := fmt.Sprintf("/update-thermostat/?threshold=%f", threshold)
	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("failed to create reqeust: %s", err)
	}
	return SendTestRequestJSON(ctx.BrainRouter, request)
}

func CreateTestContext(t *testing.T) Context {
	router := brain.CreateRouter(brain.DefaultConfig)
	ctx := Context{
		Context:     context.Background(),
		BrainRouter: router,
	}

	t.Cleanup(func() {
		for _, f := range brain.DefaultConfig.AllFilePaths() {
			err := os.Remove(f)
			t.Logf("clean up error: %s", err)
		}
	})
	return ctx
}

func SendTestRequestJSON(router *http.ServeMux, req *http.Request) (rawResponse *httptest.ResponseRecorder, responseModel map[string]string) {
	rawResponse = httptest.NewRecorder()
	router.ServeHTTP(rawResponse, req)
	responseModel = nil
	_ = json.Unmarshal(rawResponse.Body.Bytes(), &responseModel)
	return rawResponse, responseModel
}
