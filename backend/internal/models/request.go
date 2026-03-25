package models

type RecommendRequest struct {
	Origin      LatLng `json:"origin"`
	Destination LatLng `json:"destination"`
	DepartAt    string `json:"depart_at,omitempty"` // ISO8601, empty = now
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
