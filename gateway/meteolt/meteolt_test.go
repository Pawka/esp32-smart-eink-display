package meteolt

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestForecast(t *testing.T) {
	const forecastFixture string = "testdata/long-term.json"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(forecastFixture)
		if err != nil {
			t.Fatal(err)
		}
		w.Write(data)
	}))

	defer ts.Close()
	oldURLFunc := forecastURL
	forecastURL = func(place string) string {
		return ts.URL
	}
	defer func() { forecastURL = oldURLFunc }()

	c := New()
	res, err := c.Forecast("vilnius")
	t1, _ := time.Parse(ctLayout, "2020-08-16 11:00:00")
	t2, _ := time.Parse(ctLayout, "2020-08-16 12:00:00")

	want := &Weather{
		Place: Place{
			Name: "Vilnius",
		},
		ForecastTimestamps: []Forecast{
			{
				ForecastTimeUTC: Time(t1),
				AirTemperature:  28.5,
				WindSpeed:       2,
				WindGust:        6,
				WindDirection:   271,
				ConditionCode:   "clear",
			},
			{
				ForecastTimeUTC: Time(t2),
				AirTemperature:  29,
				WindSpeed:       3,
				WindGust:        7,
				WindDirection:   259,
				ConditionCode:   "clear",
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, res)
}
