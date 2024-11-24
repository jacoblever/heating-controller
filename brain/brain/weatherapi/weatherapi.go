package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type WeatherData struct {
	// Location struct {
	// 	Name           string  `json:"name"`
	// 	Region         string  `json:"region"`
	// 	Country        string  `json:"country"`
	// 	Lat            float64 `json:"lat"`
	// 	Lon            float64 `json:"lon"`
	// 	TzID           string  `json:"tz_id"`
	// 	LocaltimeEpoch int64   `json:"localtime_epoch"`
	// 	Localtime      string  `json:"localtime"`
	// } `json:"location"`
	Current struct {
		// LastUpdatedEpoch int64   `json:"last_updated_epoch"`
		// LastUpdated      string  `json:"last_updated"`
		TempC float64 `json:"temp_c"`
		// TempF            float64 `json:"temp_f"`
		// Condition        struct {
		// 	Text string `json:"text"`
		// 	Icon string `json:"icon"`
		// 	Code int    `json:"code"`
		// } `json:"condition"`
		// IsDay      int     `json:"is_day"`
		// WindMph    float64 `json:"wind_mph"`
		// WindKph    float64 `json:"wind_kph"`
		// WindDegree int     `json:"wind_degree"`
		// WindDir    string  `json:"wind_dir"`
		// PressureMb float64 `json:"pressure_mb"`
		// PressureIn float64 `json:"pressure_in"`
		// PrecipMm   float64 `json:"precip_mm"`
		// PrecipIn   float64 `json:"precip_in"`
		// Humidity   int     `json:"humidity"`
		// Cloud      int     `json:"cloud"`
		// FeelslikeC float64 `json:"feelslike_c"`
		// FeelslikeF float64 `json:"feelslike_f"`
		// VisKm      float64 `json:"vis_km"`
		// VisMiles   float64 `json:"vis_miles"`
		// Uv         int     `json:"uv"`
		// GustMph    float64 `json:"gust_mph"`
		// GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

type API interface {
	GetWeather(q string) (WeatherData, error)
}

func Make() API {
	return &api{
		apiKey: os.Getenv("WEATHER_API_KEY"),
	}
}

type api struct {
	apiKey string
}

func (a *api) GetWeather(q string) (WeatherData, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", a.apiKey, q)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherData{}, err
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return WeatherData{}, err
	}

	return weatherData, nil
}
