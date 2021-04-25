package entities

type ForecastDay string

var (
	Today ForecastDay = "Today"
)

type ForecastResponse struct {
	Place    string     `json:"place"`
	Forecast []Forecast `json:"forecast"`
}

type Forecast struct {
	Day              ForecastDay `json:"day"`
	AirTemperature   float32     `json:"temp"`
	DayTemperature   int         `json:"daytemp"`
	NightTemperature int         `json:"nighttemp"`
	WindSpeed        int         `json:"wind"`
	WindGust         int         `json:"gust"`
	WindDirection    int         `json:"direction"`
	ConditionCode    string      `json:"condition"`
	Icon             string      `json:"icon"`
}
