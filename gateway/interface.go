package gateway

import "github.com/Pawka/esp32-eink-smart-display/entities"

// WeatherInterface defines an API for weather forecasting service.
type WeatherInterface interface {
	Forecast(place string) (*entities.ForecastResponse, error)
}
