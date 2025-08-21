package config
import (
"log"
"os"
)
type Config struct {
	Port string
	DatabaseURL string
	AppBaseURL string
}

func Load() *Config {
	cfg := &Config{
	Port: get("PORT", "8080"),
	DatabaseURL: get("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"),
	AppBaseURL: get("APP_BASE_URL", "http://localhost:8080"),
	}
	if cfg.DatabaseURL == "" {
	log.Fatal("DATABASE_URL is required")
	}
	return cfg
	}
	func get(key, def string) string {
	if v := os.Getenv(key); v != "" {
	return v
	}
	return def
}