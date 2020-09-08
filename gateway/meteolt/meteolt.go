package meteolt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Pawka/esp32-eink-smart-display/lib"
)

const (
	forecastsURL string = "https://api.meteo.lt/v1/places/%s/forecasts/long-term"
	ctLayout            = "2006-01-02 15:04:05"
	iconNotFound rune   = ')'
)

// conditionToMeteocon is weather conditions map to meteocon icons font.
// URL: alessioatzeni.com/meteocons/
var conditionToMeteocon map[string]rune = map[string]rune{
	"clear":            'B',
	"isolated-clouds":  'H',
	"scattered-clouds": 'H',
	"overcast":         'N',
	"light-rain":       'Q',
	"moderate-rain":    'R',
	"heavy-rain":       'R',
	"sleet":            'V',
	"light-snow":       'U',
	"moderate-snow":    'U',
	"heavy-snow":       'W',
	"fog":              'L',
	"na":               ')',
}

type Weather struct {
	Place              Place      `json:"place"`
	ForecastTimestamps []Forecast `json:"forecastTimestamps"`
}

type Place struct {
	Name string `json:"name"`
}

type Forecast struct {
	ForecastTimeUTC Time    `json:"forecastTimeUtc"`
	AirTemperature  float32 `json:"airTemperature"`
	WindSpeed       int     `json:"WindSpeed"`
	WindGust        int     `json:"windGust"`
	WindDirection   int     `json:"windDirection"`
	ConditionCode   string  `json:"conditionCode"`
	Icon            rune    `json:"icon"`
}

type Time time.Time

// UnmarshalJSON Parses the json string in the custom format
func (ct *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(ctLayout, s)
	*ct = Time(nt)
	return
}

// MarshalJSON writes a quoted string in the custom format
func (ct Time) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

// String returns the time in the custom format
func (ct *Time) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(ctLayout))
}

// ToTime returns value converted to time.Time struct.
func (ct *Time) ToTime() time.Time {
	return time.Time(*ct)
}

type Service interface {
	Forecast(place string) (*Weather, error)
}

type service struct {
	c *client
	t lib.Clock
}

func New() Service {
	return &service{
		c: newClient(),
		t: lib.NewClock(),
	}
}

func (s *service) Forecast(place string) (*Weather, error) {
	w, err := s.c.Forecast(place)
	if err != nil {
		return nil, fmt.Errorf("querying client: %w", err)
	}

	w.ForecastTimestamps = s.dropPastTimestamps(w.ForecastTimestamps)
	for i, f := range w.ForecastTimestamps {
		w.ForecastTimestamps[i].Icon = getIcon(f.ConditionCode)
	}
	return w, nil
}

type client struct{}

func newClient() *client {
	return &client{}
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

var forecastURL = func(place string) string {
	return fmt.Sprintf(forecastsURL, place)
}

func (c *client) Forecast(place string) (*Weather, error) {
	resp, err := http.Get(forecastURL(place))
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

func getIcon(code string) rune {
	icon, exists := conditionToMeteocon[code]
	if exists == false {
		return iconNotFound
	}
	return icon
}
