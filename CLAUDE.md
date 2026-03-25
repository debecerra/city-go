# city-go

Ambient transport mode recommender for Seattle.

## Stack
- Backend: Go, chi router, internal/http/ transport layer
- Mobile: Flutter + Dart
- Maps: flutter_map + OpenStreetMap tiles
- Routing API: OpenRouteService (ORS)
- Weather API: OpenMeteo (no key needed)

## Structure
- backend/ — Go API server
- app/ — Flutter mobile app
- docs/ — architecture notes and ADRs

## Key conventions
- Go module: github.com/debecerra/city-go/backend
- HTTP handlers live in internal/http/handlers/
- Business logic lives in internal/recommendation/
- External API calls live in internal/sources/
- Models shared across layers in internal/models/