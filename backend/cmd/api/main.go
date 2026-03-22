package main

import (
	"os"

	"github.com/debecerra/city-go/backend/internal/http"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := http.NewServer(port)
	srv.Run()
}
