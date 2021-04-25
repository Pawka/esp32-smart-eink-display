package meteolt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	forecastsURL string = "https://api.meteo.lt/v1/places/%s/forecasts/long-term"
	ctLayout            = "2006-01-02 15:04:05"
	iconNotFound rune   = ')'
)

// conditionToMeteocon is weather conditions map to meteocon icons font.
// URL: alessioatzeni.com/meteocons/
var conditionToMeteocon = map[string]rune{
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

// conditionOrder is weather conditions order to decide which condition should
// represend the day. Higher value means higher importance.
var conditionOrder = map[string]int{
	"na":               0,
	"clear":            1,
	"scattered-clouds": 2,
	"isolated-clouds":  3,
	"overcast":         4,
	"fog":              5,
	"light-rain":       6,
	"light-snow":       7,
	"moderate-rain":    8,
	"moderate-snow":    9,
	"sleet":            10,
	"heavy-rain":       11,
	"heavy-snow":       12,
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

// Gateway defines an interface to fetch weather forecast information from
// meteo.lt API.
type Gateway interface {
	Get(place string) (*Weather, error)
}

type client struct{}

// NewGateway creates a new gateway.
func NewGateway() Gateway {
	return &client{}
}

var forecastURL = func(place string) string {
	return fmt.Sprintf(forecastsURL, place)
}

func (c *client) Get(place string) (*Weather, error) {
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

	for i, f := range response.ForecastTimestamps {
		response.ForecastTimestamps[i].Icon = getIcon(f.ConditionCode)
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
