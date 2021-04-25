package meteolt

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pawka/esp32-eink-smart-display/entities"
	"github.com/Pawka/esp32-eink-smart-display/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockServer(t *testing.T, fixture string) func() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(fixture)
		require.NoError(t, err)
		w.Write(data)
	}))

	oldURLFunc := forecastURL
	forecastURL = func(place string) string {
		return ts.URL
	}

	return func() {
		defer ts.Close()
		defer func() { forecastURL = oldURLFunc }()
	}
}

func TestForecast(t *testing.T) {
	const forecastFixture string = "testdata/long-term-full.json"
	defer mockServer(t, forecastFixture)()

	wantForecast := []entities.Forecast{
		{
			Day:              entities.Today,
			AirTemperature:   4.7,
			NightTemperature: -999,
			DayTemperature:   5.7,
			WindSpeed:        5,
			WindGust:         12,
			WindDirection:    304,
			ConditionCode:    "light-rain",
			Icon:             "Q",
		},
		{
			Day:              "Sun",
			NightTemperature: -0.1,
			DayTemperature:   5.4,
			Icon:             "R",
		},
		{
			Day:              "Mon",
			NightTemperature: -0.1,
			DayTemperature:   6.2,
			Icon:             "V",
		},
		{
			Day:              "Tue",
			NightTemperature: -0.8,
			DayTemperature:   8.4,
			Icon:             "R",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clock := mocks.NewMockClock(ctrl)
	c := service{
		c: NewGateway(),
		t: clock,
	}

	clock.EXPECT().UTC().Times(1).DoAndReturn(func() time.Time {
		currentTime, err := time.Parse(ctLayout, "2021-04-24 15:00:00")
		require.NoError(t, err)
		return currentTime
	})
	want := &entities.ForecastResponse{
		Place:    "Vilnius",
		Forecast: wantForecast,
	}

	res, err := c.Forecast("vilnius")
	assert.NoError(t, err)
	assert.Equal(t, want, res)
}
