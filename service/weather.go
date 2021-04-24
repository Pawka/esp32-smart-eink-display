package service

import (
	"fmt"
	"time"

	"github.com/Pawka/esp32-eink-smart-display/entities"
	"github.com/Pawka/esp32-eink-smart-display/gateway"
	"github.com/Pawka/esp32-eink-smart-display/gateway/meteolt"
)

var weatherService gateway.WeatherInterface

// GetWeather initializes a new Weather service if it was not created yet.
// Returns the instance of service.
func GetWeather() gateway.WeatherInterface {
	if weatherService == nil {
		weatherService = &weather{
			Client: meteolt.New(),
		}
	}
	return weatherService
}

type weather struct {
	Client gateway.WeatherInterface
	ts     time.Time
	last   *entities.ForecastResponse
}

const requestCacheTTL = time.Second * 120

func (w *weather) Forecast(place string) (*entities.ForecastResponse, error) {
	var weather *entities.ForecastResponse
	var err error
	now := time.Now()

	if w.ts.Add(requestCacheTTL).After(now) {
		weather = w.last
	} else {
		if weather, err = w.Client.Forecast(place); err != nil {
			return nil, fmt.Errorf("weather service forecast: %v", err)
		}
		w.last = weather
		w.ts = now
	}

	resp := &entities.ForecastResponse{
		Place:    w.last.Place,
		Forecast: weather.Forecast,
	}
	return resp, nil
}
