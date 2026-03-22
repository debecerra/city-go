package models

type RecommendResponse struct {
	Best        string       `json:"best"`   // "bike", "transit", "walk", "rideshare"
	Reason      string       `json:"reason"` // human-readable
	Modes       []ModeOption `json:"modes"`
	Conditions  []Condition  `json:"conditions"` // weather, transit status etc.
	GeneratedAt string       `json:"generated_at"`
}

type ModeOption struct {
	Mode        string `json:"mode"`
	DurationMin int    `json:"duration_min"`
	Cost        string `json:"cost,omitempty"` // "$2.75", "Free", "$14–19"
	Summary     string `json:"summary"`
	Alert       string `json:"alert,omitempty"` // "8 min delay", "Surge 1.4×"
}

type Condition struct {
	Status  string `json:"status"` // "ok", "warn"
	Message string `json:"message"`
}
