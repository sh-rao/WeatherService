package internal

import (
	"fmt"
	"net/http"

	"./netutil"
	"./weather"
)

type Handler struct {
	weatherService weather.Service
}

func NewHandler(weatherService weather.Service) *Handler {
	return &Handler{
		weatherService: weatherService,
	}
}

func (h *Handler) Execute(w http.ResponseWriter, r *http.Request) {
	currentWeatherData, err := h.weatherService.GetCurrentWeather()
	if err != nil {
		errorMsg := fmt.Sprintf("{\"message\": \"A system error occurred. Details: %+v \"} ", err)
		netutil.WriteResponse(errorMsg, http.StatusInternalServerError, w)
	}
	netutil.WriteResponse(currentWeatherData, http.StatusOK, w)
}
