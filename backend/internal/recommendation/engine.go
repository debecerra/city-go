package recommendation

import (
	"context"
	"sync"

	"github.com/debecerra/city-go/backend/config"
	"github.com/debecerra/city-go/backend/internal/models"
	"github.com/debecerra/city-go/backend/internal/sources"
)

func GetRecommendation(ctx context.Context, req models.RecommendRequest) (models.RecommendResponse, error) {
	var (
		wg      sync.WaitGroup
		routing []sources.RoutingResult
		weather sources.WeatherResult
		errs    = make([]error, 2)
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		routing, errs[2] = sources.GetRoutes(ctx, req.Origin, req.Destination, config.Load().OpenRouteServiceKey)
	}()

	go func() {
		defer wg.Done()
		weather, errs[0] = sources.GetWeather(ctx, req.Origin)
	}()

	wg.Wait()

	// handle errs, then weigh and pick best mode
	return weigh(weather, routing), nil
}

func weigh(weather sources.WeatherResult, routing []sources.RoutingResult) models.RecommendResponse {
	// simple logic for now, can be more complex later
	return models.RecommendResponse{}
}
