package sources

import (
	"context"

	"github.com/debecerra/city-go/backend/internal/models"
)

type RoutingResult struct {
}

func GetRoutes(ctx context.Context, origin, destination models.LatLng) (RoutingResult, error) {
	// call routing API, return result
	return RoutingResult{}, nil
}
