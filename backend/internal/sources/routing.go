package sources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/debecerra/city-go/backend/internal/models"
)

const orsBaseURL = "https://api.openrouteservice.org/v2/directions"

type orsRequest struct {
	Coordinates [][]float64 `json:"coordinates"`
}

type orsResponse struct {
	Routes []struct {
		Summary struct {
			Distance float64 `json:"distance"` // metres
			Duration float64 `json:"duration"` // seconds
		} `json:"summary"`
		Segments []struct {
			Steps []struct {
				Instruction string  `json:"instruction"`
				Distance    float64 `json:"distance"`
				Duration    float64 `json:"duration"`
			} `json:"steps"`
		} `json:"segments"`
	} `json:"routes"`
}

type RoutingResult struct {
	Mode        string
	DurationMin int
	DistanceKm  float64
}

func GetRoutes(ctx context.Context, origin, dest models.LatLng, apiKey string) ([]RoutingResult, error) {
	profiles := []string{"foot-walking", "cycling-regular", "driving-car"}
	results := make([]RoutingResult, 0, len(profiles))

	body := orsRequest{
		Coordinates: [][]float64{
			{origin.Lng, origin.Lat}, // ORS takes [lng, lat] not [lat, lng]
			{dest.Lng, dest.Lat},
		},
	}

	for _, profile := range profiles {
		result, err := fetchRoute(ctx, profile, body, apiKey)
		if err != nil {
			return nil, fmt.Errorf("routing %s: %w", profile, err)
		}
		results = append(results, result)
	}

	return results, nil
}

func fetchRoute(ctx context.Context, profile string, body orsRequest, apiKey string) (RoutingResult, error) {
	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s/json", orsBaseURL, profile),
		bytes.NewReader(b),
	)
	if err != nil {
		return RoutingResult{}, err
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return RoutingResult{}, err
	}
	defer resp.Body.Close()

	var orsResp orsResponse
	if err := json.NewDecoder(resp.Body).Decode(&orsResp); err != nil {
		return RoutingResult{}, err
	}

	if len(orsResp.Routes) == 0 {
		return RoutingResult{}, fmt.Errorf("no routes returned for profile %s", profile)
	}

	summary := orsResp.Routes[0].Summary
	return RoutingResult{
		Mode:        profile,
		DurationMin: int(summary.Duration / 60),
		DistanceKm:  summary.Distance / 1000,
	}, nil
}
