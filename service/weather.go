package service

import (
	"fmt"
	"time"

	"github.com/Pawka/esp32-eink-smart-display/gateway/meteolt"
)

type ForecastDay string

var (
	today    ForecastDay = "Today"
	tomorrow ForecastDay = "Tomorrow"
	monday   ForecastDay = "Monday"
	tuesday  ForecastDay = "Tuesday"
)

type ForecastResponse struct {
	Place    string     `json:"place"`
	Forecast []Forecast `json:"forecast"`
}

type Forecast struct {
	Day            ForecastDay `json:"day"`
	AirTemperature float32     `json:"temp"`
	WindSpeed      int         `json:"wind"`
	WindDirection  int         `json:"direction"`
	ConditionCode  string      `json:"condition"`
}

type Weather interface {
	Forecast(place string) (*ForecastResponse, error)
}

var weatherService Weather

// GetWeather initializes a new Weather service if it was not created yet.
// Returns the instance of service.
func GetWeather() Weather {
	if weatherService == nil {
		weatherService = &weather{
			Client: meteolt.New(),
		}
	}
	return weatherService
}

type weather struct {
	Client meteolt.Client
	ts     time.Time
	last   *meteolt.Weather
}

const requestCacheTTL = time.Second * 120

func (w *weather) Forecast(place string) (*ForecastResponse, error) {
	var weather *meteolt.Weather
	var err error
	now := time.Now()

	fmt.Printf("ts: %#v", w.ts)
	if w.ts.Add(requestCacheTTL).After(now) {
		weather = w.last
	} else {
		if weather, err = w.Client.Forecast(place); err != nil {
			return nil, fmt.Errorf("weather service forecast: %v", err)
		}
		w.last = weather
		w.ts = now
	}

	resp := &ForecastResponse{
		Place:    w.last.Place.Name,
		Forecast: mapFromMeteoltResponse(weather.ForecastTimestamps),
	}
	return resp, nil
}

func mapFromMeteoltResponse(data []meteolt.Forecast) []Forecast {
	today := Forecast{
		Day:            today,
		AirTemperature: data[0].AirTemperature,
		ConditionCode:  data[0].ConditionCode,
		WindDirection:  data[0].WindDirection,
		WindSpeed:      data[0].WindSpeed,
	}

	// TODO: implement mapping for remaining days
	return []Forecast{today}
}
