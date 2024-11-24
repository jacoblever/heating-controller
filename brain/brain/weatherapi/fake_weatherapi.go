package weatherapi

import (
	"encoding/json"
	"fmt"
)

func MakeFakeAPI() API {
	return &fakeApi{
		mockTemp: 18.6,
	}
}

type fakeApi struct {
	mockTemp float64
}

func (a *fakeApi) GetWeather(q string) (WeatherData, error) {
	body := fmt.Sprintf("{\"current\": {\"temp_c\": \"%f\"}}", a.mockTemp)

	var weatherData WeatherData
	err := json.Unmarshal([]byte(body), &weatherData)
	if err != nil {
		return WeatherData{}, err
	}

	return weatherData, nil
}
