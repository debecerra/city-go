package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/debecerra/city-go/backend/internal/models"
	"github.com/debecerra/city-go/backend/internal/recommendation"
)

func Recommend(w http.ResponseWriter, r *http.Request) {
	var req models.RecommendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := recommendation.GetRecommendation(r.Context(), req)
	if err != nil {
		http.Error(w, "failed to get recommendation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
