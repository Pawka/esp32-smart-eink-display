package meteolt

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Pawka/esp32-eink-smart-display/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clock := mocks.NewMockClock(ctrl)
	c := service{
		c: newClient(),
		t: clock,
	}

	clock.EXPECT().UTC().Times(1).DoAndReturn(func() time.Time {
		currentTime, err := time.Parse(ctLayout, "2020-08-16 12:10:00")
		require.NoError(t, err)
		return currentTime
	})
	forecastTime, err := time.Parse(ctLayout, "2020-08-16 12:00:00")
	require.NoError(t, err)

	want := &Weather{
		Place: Place{
			Name: "Vilnius",
		},
		ForecastTimestamps: []Forecast{
			{
				ForecastTimeUTC: Time(forecastTime),
				AirTemperature:  29,
				WindSpeed:       3,
				WindGust:        7,
				WindDirection:   259,
				ConditionCode:   "clear",
				Icon:            'B',
			},
		},
	}

	res, err := c.Forecast("vilnius")
	assert.NoError(t, err)
	assert.Equal(t, want, res)
}

func TestDropPastTimestamps(t *testing.T) {
	in := []Forecast{
		{ForecastTimeUTC: ts(t, "2020-09-01 12:00:00")},
		{ForecastTimeUTC: ts(t, "2020-09-01 13:00:00")},
		{ForecastTimeUTC: ts(t, "2020-09-01 14:00:00")},
	}

	testCases := map[string]struct {
		now  string
		want []Forecast
	}{
		"current_time_is_prior_all_timestamps": {
			now:  "2020-09-01 11:00:00",
			want: in,
		},
		"current_time_is_equal_to_the_first_timestamp": {
			now:  "2020-09-01 12:00:00",
			want: in,
		},
		"current_time_is_in_the_same_hour_of_the_first_timestamp": {
			now:  "2020-09-01 12:01:00",
			want: in,
		},
		"current_time_after_first_timestamp": {
			now: "2020-09-01 13:00:00",
			want: []Forecast{
				{ForecastTimeUTC: ts(t, "2020-09-01 13:00:00")},
				{ForecastTimeUTC: ts(t, "2020-09-01 14:00:00")},
			},
		},
		"current_time_after_all_timestamps": {
			now:  "2020-09-01 15:00:00",
			want: []Forecast{},
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			clock := mocks.NewMockClock(ctrl)
			clock.EXPECT().UTC().DoAndReturn(func() time.Time {
				ts, err := time.Parse(ctLayout, test.now)
				require.NoError(t, err)
				return ts
			})
			c := service{
				t: clock,
			}
			result := c.dropPastTimestamps(in)
			assert.Equal(t, test.want, result)
		})
	}
}

func ts(t *testing.T, timestamp string) Time {
	ts, err := time.Parse(ctLayout, timestamp)
	require.NoError(t, err)
	return Time(ts)
}

func TestGetIcon(t *testing.T) {
	assert.Equal(t, 'B', getIcon("clear"))
	assert.Equal(t, ')', getIcon("na"))
	assert.Equal(t, iconNotFound, getIcon("some-random-condition"))
}
