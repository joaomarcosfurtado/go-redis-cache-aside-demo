package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type WeatherData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Temp float64 `json:"temperature_2m"`
		Wind float64 `json:"wind_speed_10m"`
	} `json:"current"`
	Source string `json:"source,omitempty"`
}

type WeatherService struct {
	rdb    *redis.Client
	client *http.Client
}

func NewWeatherService(rdb *redis.Client) *WeatherService {
	return &WeatherService{
		rdb:    rdb,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, lat, lon string, mock bool) (*WeatherData, error) {
	cacheKey := fmt.Sprintf("weather:%s:%s", lat, lon)
	var data WeatherData

	// 1. Try Cache
	val, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &data) == nil {
			data.Source = "Redis Cache ⚡"
			return &data, nil
		}
	}

	// 2. Try API (Miss)
	data, err = s.fetchFromAPI(ctx, lat, lon, mock)
	if err != nil {
		return nil, err
	}

	// 3. Save Cache (Async)
	go func() {
		bytes, _ := json.Marshal(data)
		s.rdb.Set(context.Background(), cacheKey, bytes, 10*time.Second)
	}()

	return &data, nil
}

func (s *WeatherService) fetchFromAPI(ctx context.Context, lat, lon string, mock bool) (WeatherData, error) {
	var data WeatherData

	if mock {
		time.Sleep(500 * time.Millisecond)
		raw := []byte(`{"latitude":52.52,"longitude":13.41,"current":{"temperature_2m":25.0,"wind_speed_10m":10.0}}`)
		json.Unmarshal(raw, &data)
		data.Source = "Mock API (Slow) 🐢"
		return data, nil
	}

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,wind_speed_10m", lat, lon)
	
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		return data, err
	}

	data.Source = "Open-Meteo API (Slow) 🐢"
	return data, nil
}