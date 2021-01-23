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
	Date    string           `json:"date"`
	TS      int32            `json:"ts"`
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
		Date:    t.Format(layoutISO),
		TS:      int32(t.Unix()),
		Weather: *wr,
	}, nil
}
