package meteolt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"net/http"
)

const (
	forecastsURL string = "https://api.meteo.lt/v1/places/%s/forecasts/long-term"
	ctLayout            = "2006-01-02 15:04:05"
)

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

type Client interface {
	Forecast(place string) (*Weather, error)
}

func New() Client {
	return &client{}
}

type client struct {
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
