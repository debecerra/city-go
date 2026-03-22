package sources

import (
	"context"

	"github.com/debecerra/city-go/backend/internal/models"
)

type TransitResult struct {
}

func GetTransitStatus(ctx context.Context, origin, destination models.LatLng) (TransitResult, error) {
	// call transit API, return result
	return TransitResult{}, nil
}
