# City-go — project design doc

Ambient transport mode recommender for Seattle.

---

## Product

An ambient mobile app that recommends the best way to get somewhere right now —
walking, biking, transit, or rideshare — by factoring in weather, real-time transit
delays, elevation, and cost. Optimises for the right choice in the moment, not just
the fastest route.

### Core value proposition

- Opinionated recommendation with a plain-English reason ("Bike — it's dry and
  downhill, saves you $4")
- Shows trade-offs: time vs. cost vs. effort vs. reliability
- Seattle-specific: hill weighting for bikes, route reliability by Metro/Link line
- Conditions panel explains why the recommendation was made

### Key differentiators vs. Google Maps

- Factors in weather, not just time
- Real-time transit reliability (not just scheduled)
- Elevation-aware bike routing (Dexter Ave vs. brutal climb)
- Cost-aware including rideshare surge pricing
- Single clear recommendation, not a list of options

---

## Data sources

All free, no billing risk at early stage.

| Source | Used for |
|---|---|
| OpenMeteo | Weather — temperature, precipitation, wind. Free, no API key. |
| GTFS-RT (King County Metro / Sound Transit) | Real-time transit vehicle positions and delay data. |
| OpenRouteService | Walking and bike routing with elevation profile. Free tier. |
| OpenStreetMap | Map tiles and base geodata via flutter_map. |
| Uber / Lyft API | Rideshare price and time estimates (requires API access). |

---

## Tech stack

### Mobile app — Flutter

Flutter with Dart. Cross-platform: ships to iOS App Store and Android Play Store
from one codebase. Dart's type system and structure are closer to Go than JavaScript,
reducing context-switching.

| Package | Used for |
|---|---|
| flutter_map | OpenStreetMap tile rendering. Open-source, no billing. |
| geolocator | Device GPS and permissions handling. |
| Riverpod | State management for async API data. |
| Dio | HTTP client for Go API calls. |
| Hive / SharedPreferences | Local caching of user preferences and recent responses. |

### Backend — Go

Go REST API. The core operation is fanning out to 4–5 external APIs concurrently
using goroutines, aggregating results, applying the recommendation logic, and
returning a single JSON response in under 500ms. Deployed as a Docker container
on Fly.io.

| Component | Detail |
|---|---|
| HTTP router | chi — stays close to stdlib, handlers are plain http.HandlerFunc |
| Concurrency | sync.WaitGroup + channels for parallel API fan-out |
| Caching | In-memory or Redis for short-lived weather/transit responses |
| Deployment | Fly.io — Go binary in Docker, global edge, cheap at small scale |

### Project structure
```
city-go/
  backend/
    cmd/api/
      main.go
    internal/
      http/
        handlers/
          recommend.go
          health.go
        server.go
      recommendation/
        engine.go
      sources/
        weather.go
        transit.go
        routing.go
      models/
        request.go
        response.go
      cache/
        cache.go
    config/
      config.go
    go.mod
    Dockerfile
    fly.toml
  app/
    lib/
      main.dart
      features/
        recommendation/
        map/
        settings/
      core/
        api_client.dart
        models/
        theme/
        router.dart
      shared/
    pubspec.yaml
  docs/
    architecture.md
    api.md
    data-sources.md
    adr/
  CLAUDE.md
  Makefile
  docker-compose.yml
  .gitignore
```

### Web (deferred)

Flutter web (`flutter build web`) is available as a low-effort phase 2 — same
codebase, no backend changes. If web becomes a priority, a thin SvelteKit or
Next.js frontend calling the same Go API would give better load performance and
SEO. The Go backend is unchanged in either case.

---

## Maps decision

OpenStreetMap via flutter_map. Chosen over Google Maps Platform and Mapbox:

- Free with no MAU caps or monthly billing
- Mapbox free tier: 25,000 MAU/month on mobile — fine at early stage but a future
  constraint
- Google Maps mobile SDK is free for rendering but Directions API costs ~$5–10 per
  1k calls after 5k/month free
- OSM data quality in Seattle is excellent and actively maintained
- flutter_map API is clean and well-documented for custom tile sources and route
  overlays

Revisit Mapbox if custom map styling becomes important post-launch.

---

## API contract

### POST /v1/recommend

**Request**
```json
{
  "origin": { "lat": 47.6062, "lng": -122.3321 },
  "destination": { "lat": 47.6490, "lng": -122.3577 },
  "depart_at": ""
}
```

**Response**
```json
{
  "best": "bike",
  "reason": "Dry and 52°F, mostly downhill on Eastlake. Saves you $4 vs. transit.",
  "modes": [
    { "mode": "bike", "duration_min": 24, "cost": "Free", "summary": "Downhill on Eastlake Ave", "alert": "" },
    { "mode": "transit", "duration_min": 38, "cost": "$2.75", "summary": "Route 49", "alert": "8 min delay" },
    { "mode": "walk", "duration_min": 52, "cost": "Free", "summary": "Eastlake path", "alert": "" },
    { "mode": "rideshare", "duration_min": 18, "cost": "$14–19", "summary": "", "alert": "Surge 1.4×" }
  ],
  "conditions": [
    { "status": "ok", "message": "Weather: Dry, 52°F, light clouds" },
    { "status": "warn", "message": "Transit: Route 49 running 8 min late" },
    { "status": "ok", "message": "Bike lanes: Clear on Eastlake Ave" }
  ],
  "generated_at": "2026-03-24T08:42:00Z"
}
```

---

## Phased delivery

- **Phase 1** — Flutter mobile app (iOS + Android) + Go backend. Core recommendation
  flow: location → fan-out to APIs → single opinionated result with reason.
- **Phase 2** — Flutter web build if a browser-accessible version is needed. Same
  codebase, low effort.
- **Phase 3** — Thin SvelteKit/Next.js web frontend if load performance or SEO
  matters. Go backend unchanged.
- **Funding** — validate locally in Seattle first. If it gains traction, OSM/GTFS
  data generalises to any city.

---

## Open questions

- Recommendation weighting logic: rule-based thresholds (temp < 40°F = don't suggest
  biking) vs. learned user preferences
- User accounts: needed for preference storage, or anonymous with local prefs only?
- Uber/Lyft API access: requires partnership approval — assess early
- Elevation data source for bike routing: OpenRouteService includes it, confirm
  accuracy on Seattle hills
- Notification strategy: push alerts for transit delays on a saved commute route?