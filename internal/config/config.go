package config

import (
	"html/template"
	"log"
	"os"
)

type Config struct {
	Port      string
	RedisURL  string
	RedisTLS  bool
	Templates *template.Template
}

func Load() *Config {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Critical: Could not load templates: %v", err)
	}

	return &Config{
		Port:      getEnv("PORT", "8080"),
		RedisURL:  getEnv("REDIS_URL", "redis://localhost:6379"),
		RedisTLS:  os.Getenv("REDIS_TLS") == "true",
		Templates: tmpl,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}