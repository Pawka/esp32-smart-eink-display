package meteolt

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	const forecastFixture string = "testdata/long-term.json"
	defer mockServer(t, forecastFixture)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := NewGateway()

	forecastTime1, err := time.Parse(ctLayout, "2020-08-16 11:00:00")
	require.NoError(t, err)
	forecastTime2, err := time.Parse(ctLayout, "2020-08-16 12:00:00")
	require.NoError(t, err)

	want := &Weather{
		Place: Place{
			Name: "Vilnius",
		},
		ForecastTimestamps: []Forecast{
			{
				ForecastTimeUTC: Time(forecastTime1),
				AirTemperature:  28.5,
				WindSpeed:       2,
				WindGust:        6,
				WindDirection:   271,
				ConditionCode:   "clear",
				Icon:            'B',
			},
			{
				ForecastTimeUTC: Time(forecastTime2),
				AirTemperature:  29,
				WindSpeed:       3,
				WindGust:        7,
				WindDirection:   259,
				ConditionCode:   "clear",
				Icon:            'B',
			},
		},
	}

	res, err := c.Get("vilnius")
	assert.NoError(t, err)
	assert.Equal(t, want, res)
}
