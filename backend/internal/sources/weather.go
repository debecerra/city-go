package sources

import (
	"context"

	"github.com/debecerra/city-go/backend/internal/models"
)

type WeatherResult struct {
}

func GetWeather(ctx context.Context, location models.LatLng) (WeatherResult, error) {
	// call weather API, return result
	return WeatherResult{}, nil
}
