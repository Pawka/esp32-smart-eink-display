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
	Place    string
	Forecast []Forecast
}

type Forecast struct {
	Day            ForecastDay
	AirTemperature float32
	WindSpeed      int
	WindDirection  int
	ConditionCode  string
}

type Weather interface {
	Forecast(place string) (*ForecastResponse, error)
}

func NewWeather() Weather {
	return &weather{
		Client: meteolt.New(),
	}
}

type weather struct {
	Client meteolt.Client
	ts     time.Time
	last   *meteolt.Forecast
}

const requestCacheTTL = time.Second

func (w *weather) Forecast(place string) (*ForecastResponse, error) {
	// TODO: Add caching
	weather, err := w.Client.Forecast(place)
	if err != nil {
		return nil, fmt.Errorf("weather service forecast: %v", err)
	}

	resp := &ForecastResponse{
		Place:    weather.Place.Name,
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
