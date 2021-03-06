package meteolt

import (
	"fmt"
	"math"
	"time"

	"github.com/Pawka/esp32-eink-smart-display/entities"
	"github.com/Pawka/esp32-eink-smart-display/gateway"
	"github.com/Pawka/esp32-eink-smart-display/lib"
)

type service struct {
	c Gateway
	t lib.Clock
}

// New creates meteolt service.
func New() gateway.WeatherInterface {
	return &service{
		c: NewGateway(),
		t: lib.NewClock(),
	}
}

type dailyForecast map[string]*dayForecast

type dayForecast struct {
	NightTemp     float64
	DayTemp       float64
	ConditionCode string
}

const (
	dateLayout = "2006-01-02"
	// lastNightHour defines the hour when the night period is over. Required to
	// separate night/day hour ranges.
	lastNightHour = 6
	// tempNone indicates when temperature is not set.
	tempNone = -999
)

func (s *service) Forecast(place string) (*entities.ForecastResponse, error) {
	w, err := s.c.Get(place)
	if err != nil {
		return nil, fmt.Errorf("querying client: %w", err)
	}

	dailyFc := s.getDailyForecast(w.ForecastTimestamps)
	w.ForecastTimestamps = s.dropPastTimestamps(w.ForecastTimestamps)

	todayFormat := w.ForecastTimestamps[0].ForecastTimeUTC.ToTime().Format(dateLayout)
	today := mapFromForecast(entities.Today, w.ForecastTimestamps[0])
	today.NightTemperature = int(dailyFc[todayFormat].NightTemp)
	today.DayTemperature = int(dailyFc[todayFormat].DayTemp)

	const furtherDaysCount = 3
	days := s.getFurtherDays(dailyFc, w.ForecastTimestamps, furtherDaysCount)
	f := []entities.Forecast{today}
	f = append(f, days...)

	resp := &entities.ForecastResponse{
		Place:    w.Place.Name,
		Forecast: f,
	}
	return resp, nil
}

func (s *service) getFurtherDays(df dailyForecast, ft []Forecast, days int) []entities.Forecast {
	i := 0
	t := ft[0].ForecastTimeUTC.ToTime().Add(time.Hour * 24)
	result := make([]entities.Forecast, days)
	for ; i < days; i++ {
		dayFormat := t.Format(dateLayout)
		day := mapFromDayForecast(df[dayFormat])
		day.Day = entities.ForecastDay(t.Weekday().String()[:3])
		result[i] = day
		t = t.Add(time.Hour * 24)
	}

	return result
}

func (s *service) getDailyForecast(timestamps []Forecast) dailyForecast {
	var cd string
	result := make(dailyForecast)
	for _, v := range timestamps {
		currentDate := v.ForecastTimeUTC.ToTime().Format(dateLayout)
		if cd != currentDate {
			cd = currentDate
			result[cd] = &dayForecast{
				NightTemp: tempNone,
				DayTemp:   tempNone,
			}
		}

		if v.ForecastTimeUTC.ToTime().Hour() <= lastNightHour {
			if result[cd].NightTemp == tempNone || result[cd].NightTemp > v.AirTemperature {
				result[cd].NightTemp = v.AirTemperature
			}
		} else {
			if result[cd].DayTemp == tempNone || result[cd].DayTemp < v.AirTemperature {
				result[cd].DayTemp = v.AirTemperature
			}

			condition := result[cd].ConditionCode
			if condition == "" || conditionOrder[condition] < conditionOrder[v.ConditionCode] {
				result[cd].ConditionCode = v.ConditionCode
			}
		}
	}

	return result
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

func mapFromForecast(day entities.ForecastDay, data Forecast) entities.Forecast {
	return entities.Forecast{
		Day:            day,
		AirTemperature: data.AirTemperature,
		ConditionCode:  data.ConditionCode,
		WindDirection:  data.WindDirection,
		WindGust:       data.WindGust,
		WindSpeed:      data.WindSpeed,
		Icon:           string(data.Icon),
	}
}

func mapFromDayForecast(forecast *dayForecast) entities.Forecast {
	return entities.Forecast{
		NightTemperature: int(math.Round(forecast.NightTemp)),
		DayTemperature:   int(math.Round(forecast.DayTemp)),
		Icon:             string(getIcon(forecast.ConditionCode)),
	}
}
