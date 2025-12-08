package main

import (
	"log/slog"
	"net/http"
	"os"

	"weather-cache/internal/config"
	"weather-cache/internal/handler"
	"weather-cache/internal/platform"
	"weather-cache/internal/service"
)

func main() {
	// 1. Config
	cfg := config.Load()

	// 2. Infra (Redis)
	rdb, err := platform.NewRedisClient(cfg)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}

	// 3. Service (Lógica)
	weatherSvc := service.NewWeatherService(rdb)

	// 4. Handlers (HTTP)
	weatherHdl := handler.NewWeatherHandler(weatherSvc, cfg.Templates)

	// 5. Router
	mux := http.NewServeMux()
	mux.HandleFunc("/", weatherHdl.HandleHome)
	mux.HandleFunc("/weather", weatherHdl.HandleWeather)

	// 6. Start
	slog.Info("Server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		slog.Error("Server failed", "error", err)
	}
}