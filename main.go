package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"

	h "./internal"
	cfg "./internal/config"
	w "./internal/weather"
)

func main() {
	weatherService := w.NewService(initConfig())
	handler := h.NewHandler(*weatherService)
	http.HandleFunc("/v1/weather/", handler.Execute)
	http.ListenAndServe(":8080", nil)
}

func initConfig() *cfg.Configuration {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	var configuration cfg.Configuration

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return &configuration
}
