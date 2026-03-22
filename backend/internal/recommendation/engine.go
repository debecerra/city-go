package recommendation

import (
	"context"
	"sync"

	"github.com/debecerra/city-go/backend/internal/models"
	"github.com/debecerra/city-go/backend/internal/sources"
)

func GetRecommendation(ctx context.Context, req models.RecommendRequest) (models.RecommendResponse, error) {
	var (
		wg      sync.WaitGroup
		weather sources.WeatherResult
		transit sources.TransitResult
		routing sources.RoutingResult
		errs    = make([]error, 3)
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		weather, errs[0] = sources.GetWeather(ctx, req.Origin)
	}()

	go func() {
		defer wg.Done()
		transit, errs[1] = sources.GetTransitStatus(ctx, req.Origin, req.Destination)
	}()

	go func() {
		defer wg.Done()
		routing, errs[2] = sources.GetRoutes(ctx, req.Origin, req.Destination)
	}()

	wg.Wait()

	// handle errs, then weigh and pick best mode
	return weigh(weather, transit, routing), nil
}

func weigh(weather sources.WeatherResult, transit sources.TransitResult, routing sources.RoutingResult) models.RecommendResponse {
	// simple logic for now, can be more complex later
	return models.RecommendResponse{}
}
