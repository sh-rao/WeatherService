package config

// Configurations exported
type Configuration struct {
	PrimaryWeatherProvider   WeatherProviderConfiguration
	SecondaryWeatherProvider WeatherProviderConfiguration
	StaleTime                string
}

// WeatherProviderConfiguration exported
type WeatherProviderConfiguration struct {
	BaseUrl   string
	City      string
	Unit      string
	AccessKey string
	Primary   bool
}
