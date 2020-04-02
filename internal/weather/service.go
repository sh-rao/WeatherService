package weather

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"../config"
	"../model"
)

type Service struct {
	cfg             *config.Configuration
	hc              *http.Client
	lastUpdatedTime *time.Time
	cachedWd        *model.WeatherDetails
}

func NewService(cfg *config.Configuration) *Service {
	return &Service{cfg: cfg, hc: http.DefaultClient, lastUpdatedTime: nil, cachedWd: nil}
}

func (s *Service) GetCurrentWeather() (*model.WeatherDetails, error) {
	if !s.isStale() {
		return s.cachedWd, nil
	} else {
		err := s.getWeatherDetails()
		// This will be true only if we fail during the first hit itself
		if err != nil && s.cachedWd == nil {
			return nil, err
		}
	}
	return s.cachedWd, nil
}

func (s *Service) isStale() bool {
	if s.lastUpdatedTime == nil {
		return true
	}
	return time.Now().Sub(*s.lastUpdatedTime).Seconds() > 3
}

func (s *Service) getWeatherDetails() error {
	weatherData, err := s.MakeHttpRequest(&s.cfg.PrimaryWeatherProvider)
	if err == nil {
		err = s.extractPrimaryCurrentWeatherData(weatherData)
	}
	if err != nil {
		weatherData, err := s.MakeHttpRequest(&s.cfg.SecondaryWeatherProvider)
		if err != nil {
			return err
		}
		err = s.extractSecondaryCurrentWeatherData(weatherData)
		if err != nil {
			return err
		}
	}
	currentTime := time.Now()
	s.lastUpdatedTime = &currentTime
	return nil
}

func (s *Service) extractPrimaryCurrentWeatherData(data []byte) error {
	var weatherData interface{}
	err := json.Unmarshal(data, &weatherData)
	if err != nil {
		return err
	}
	if s.cachedWd == nil {
		s.cachedWd = &model.WeatherDetails{}
	}
	weatherDataMap := weatherData.(map[string]interface{})
	currentDataMap := weatherDataMap["current"]
	currentData := currentDataMap.(map[string]interface{})
	s.cachedWd.TemperatureDegrees = currentData["temperature"].(float64)
	s.cachedWd.WindSpeed = currentData["wind_speed"].(float64)
	return nil
}

func (s *Service) extractSecondaryCurrentWeatherData(data []byte) error {
	var weatherData interface{}
	err := json.Unmarshal(data, &weatherData)
	if err != nil {
		return err
	}
	if s.cachedWd == nil {
		s.cachedWd = &model.WeatherDetails{}
	}
	weatherDataMap := weatherData.(map[string]interface{})
	mainDataMap := weatherDataMap["main"]
	mainData := mainDataMap.(map[string]interface{})
	s.cachedWd.TemperatureDegrees = mainData["temp"].(float64)
	windDataMap := weatherDataMap["wind"]
	windData := windDataMap.(map[string]interface{})
	s.cachedWd.WindSpeed = windData["speed"].(float64)
	return nil
}

func (s *Service) MakeHttpRequest(cfg *config.WeatherProviderConfiguration) ([]byte, error) {
	req, err := http.NewRequest("GET", cfg.BaseUrl, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	query := req.URL.Query()
	if cfg.Primary == true {
		query.Add("query", cfg.City)
		query.Add("access_key", cfg.AccessKey)
	} else {
		query.Add("q", cfg.City)
		query.Add("appid", cfg.AccessKey)
	}

	query.Add("units", cfg.Unit)
	req.URL.RawQuery = query.Encode()

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return data, nil
}
