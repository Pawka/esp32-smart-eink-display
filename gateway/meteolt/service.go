package meteolt

import (
	"fmt"
	"time"

	"github.com/Pawka/esp32-eink-smart-display/entities"
	"github.com/Pawka/esp32-eink-smart-display/gateway"
	"github.com/Pawka/esp32-eink-smart-display/lib"
)

type service struct {
	c Gateway
	t lib.Clock
}

func New() gateway.WeatherInterface {
	return &service{
		c: NewGateway(),
		t: lib.NewClock(),
	}
}

func (s *service) Forecast(place string) (*entities.ForecastResponse, error) {
	w, err := s.c.Get(place)
	if err != nil {
		return nil, fmt.Errorf("querying client: %w", err)
	}

	w.ForecastTimestamps = s.dropPastTimestamps(w.ForecastTimestamps)

	resp := &entities.ForecastResponse{
		Place:    w.Place.Name,
		Forecast: mapFromMeteoltResponse(w.ForecastTimestamps),
	}
	return resp, nil
}

func (s *service) dropPastTimestamps(timestamps []Forecast) []Forecast {
	startTime := s.t.UTC().Add(-time.Hour)
	for i := range timestamps {
		if startTime.Before(timestamps[i].ForecastTimeUTC.ToTime()) == true {
			return timestamps[i:]
		}
	}

	return []Forecast{}
}

func mapFromMeteoltResponse(data []Forecast) []entities.Forecast {
	today := entities.Forecast{
		Day:            entities.Today,
		AirTemperature: data[0].AirTemperature,
		ConditionCode:  data[0].ConditionCode,
		WindDirection:  data[0].WindDirection,
		WindGust:       data[0].WindGust,
		WindSpeed:      data[0].WindSpeed,
		Icon:           string(data[0].Icon),
	}

	// TODO: implement mapping for remaining days
	return []entities.Forecast{today}
}
