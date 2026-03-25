package config

import (
	"log"
	"os"
)

type Config struct {
	Port                string
	OpenRouteServiceKey string
}

func Load() Config {
	return Config{
		Port:                getEnv("PORT", "8080"),
		OpenRouteServiceKey: mustGetEnv("ORS_API_KEY"),
	}
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
