package service

import (
	"fmt"
	"time"
)

const (
	WeatherLocation string = "vilnius"
	layoutISO              = "2006-01-02 15:04:05"
)

type Content struct {
	TS      string           `json:"ts"`
	Weather ForecastResponse `json:"weather"`
}

func NewContent() (Content, error) {
	w := GetWeather()
	wr, err := w.Forecast(WeatherLocation)
	if err != nil {
		return Content{}, fmt.Errorf("creating content: %v", err)
	}
	t := time.Now()
	return Content{
		TS:      t.Format(layoutISO),
		Weather: *wr,
	}, nil
}
