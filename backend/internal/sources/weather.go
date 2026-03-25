// sources/weather.go
package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/debecerra/city-go/backend/internal/models"
)

const openMeteoURL = "https://api.open-meteo.com/v1/forecast"

type WeatherResult struct {
	TempC        float64
	PrecipMM     float64 // precipitation in last hour
	WindSpeedKmh float64
	Description  string
}

type openMeteoResponse struct {
	Current struct {
		Temperature2m float64 `json:"temperature_2m"`
		Precipitation float64 `json:"precipitation"`
		WindSpeed10m  float64 `json:"wind_speed_10m"`
		WeatherCode   int     `json:"weather_code"`
	} `json:"current"`
}

func GetWeather(ctx context.Context, location models.LatLng) (WeatherResult, error) {
	url := fmt.Sprintf(
		"%s?latitude=%.6f&longitude=%.6f&current=temperature_2m,precipitation,wind_speed_10m,weather_code",
		openMeteoURL, location.Lat, location.Lng,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return WeatherResult{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return WeatherResult{}, err
	}
	defer resp.Body.Close()

	var omResp openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&omResp); err != nil {
		return WeatherResult{}, err
	}

	c := omResp.Current
	return WeatherResult{
		TempC:        c.Temperature2m,
		PrecipMM:     c.Precipitation,
		WindSpeedKmh: c.WindSpeed10m,
		Description:  describeWeather(c.WeatherCode),
	}, nil
}

// WMO weather code to human description
// Full table: open-meteo.com/en/docs#weathervariables
func describeWeather(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code <= 3:
		return "Partly cloudy"
	case code <= 49:
		return "Foggy"
	case code <= 59:
		return "Drizzle"
	case code <= 69:
		return "Rain"
	case code <= 79:
		return "Snow"
	case code <= 82:
		return "Rain showers"
	case code <= 99:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}
