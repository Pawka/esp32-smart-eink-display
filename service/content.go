package service

import (
	"fmt"
	"time"
)

const WeatherLocation string = "vilnius"

type Content struct {
	TS      time.Time        `json:"ts"`
	Weather ForecastResponse `json:"weather"`
}

func NewContent() (Content, error) {
	w := GetWeather()
	wr, err := w.Forecast(WeatherLocation)
	if err != nil {
		return Content{}, fmt.Errorf("creating content: %v", err)
	}
	return Content{
		TS:      time.Now(),
		Weather: *wr,
	}, nil
}
