package meteolt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
)

const forecastsURL string = "https://api.meteo.lt/v1/places/%s/forecasts/long-term"

type Weather struct {
	Place              Place      `json:"place"`
	ForecastTimestamps []Forecast `json:"forecastTimestamps"`
}

type Place struct {
	Name string `json:"name"`
}

type Forecast struct {
	ForecastTimeUTC string  `json:"forecastTimeUtc"`
	AirTemperature  float32 `json:"airTemperature"`
	WindSpeed       int     `json:"WindSpeed"`
	WindGust        int     `json:"windGust"`
	WindDirection   int     `json:"windDirection"`
	ConditionCode   string  `json:"conditionCode"`
}

type Client interface {
	Forecast(place string) (*Weather, error)
}

func New() Client {
	return &client{}
}

type client struct {
}

func (c *client) Forecast(place string) (*Weather, error) {
	url := fmt.Sprintf(forecastsURL, place)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("requesting meteo.lt: %v", err)
	}

	defer resp.Body.Close()
	forecast, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading forecasts response %v", err)
	}

	var response Weather
	if err := json.Unmarshal(forecast, &response); err != nil {
		return nil, fmt.Errorf("decoding forecasts response to JSON: %v", err)
	}

	return &response, nil
}
