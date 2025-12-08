package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"weather-cache/internal/service"
)

type WeatherHandler struct {
	service   *service.WeatherService
	templates *template.Template
}

type WeatherViewModel struct {
	*service.WeatherData 
	Mock bool            
}

func NewWeatherHandler(s *service.WeatherService, t *template.Template) *WeatherHandler {
	return &WeatherHandler{service: s, templates: t}
}

func (h *WeatherHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	h.templates.Execute(w, nil)
}

func (h *WeatherHandler) HandleWeather(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")
	mock := r.URL.Query().Get("mock") == "true"

	if lat == "" || lon == "" {
		http.Error(w, "Missing params", http.StatusBadRequest)
		return
	}

	data, err := h.service.GetWeather(r.Context(), lat, lon, mock)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		
		status := "MISS"
		if data.Source == "Redis Cache ⚡" {
			status = "HIT"
		}
		w.Header().Set("X-Cache", status)
		
		json.NewEncoder(w).Encode(data)

	} else {
		viewModel := WeatherViewModel{
			WeatherData: data,
			Mock:        mock,
		}
		
		h.templates.Execute(w, viewModel)
	}
}