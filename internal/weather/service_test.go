package weather_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"math"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

	"../config"

	w "../weather"
)

func TestGetCurrentWeather(t *testing.T) {

	t.Run("primary weather provider success", func(t *testing.T) {
		cfg := config.Configuration{
			PrimaryWeatherProvider:   config.WeatherProviderConfiguration{Primary: true},
			SecondaryWeatherProvider: config.WeatherProviderConfiguration{Primary: false},
			StaleTime:                "3s",
		}
		ws := w.NewService(&cfg)
		monkey.PatchInstanceMethod(reflect.TypeOf(ws), "MakeHttpRequest", func(_ *w.Service, cfg *config.WeatherProviderConfiguration) ([]byte, error) {
			data, _ := ioutil.ReadFile("../../test/primary_success.json")
			return data, nil
		})
		defer monkey.UnpatchAll()

		wd, err := ws.GetCurrentWeather()
		assert.Nil(t, err)
		assert.Equal(t, float64(20), wd.TemperatureDegrees)
		assert.Equal(t, float64(15), wd.WindSpeed)
	})

	t.Run("secondary weather provider success", func(t *testing.T) {
		cfg := config.Configuration{
			PrimaryWeatherProvider:   config.WeatherProviderConfiguration{Primary: true},
			SecondaryWeatherProvider: config.WeatherProviderConfiguration{Primary: false},
			StaleTime:                "3s",
		}
		ws := w.NewService(&cfg)
		monkey.PatchInstanceMethod(reflect.TypeOf(ws), "MakeHttpRequest", func(_ *w.Service, cfg *config.WeatherProviderConfiguration) ([]byte, error) {
			if cfg.Primary == true {
				return nil, errors.New("primary will return error, so we can failover to secondary")
			}
			data, _ := ioutil.ReadFile("../../test/secondary_success.json")
			return data, nil
		})
		defer monkey.UnpatchAll()

		wd, err := ws.GetCurrentWeather()
		assert.Nil(t, err)
		assert.Equal(t, 16.63, wd.TemperatureDegrees)
		assert.Equal(t, 3.1, wd.WindSpeed)
	})

	t.Run("primary weather bad json, so fails over to secondary", func(t *testing.T) {
		cfg := config.Configuration{
			PrimaryWeatherProvider:   config.WeatherProviderConfiguration{Primary: true},
			SecondaryWeatherProvider: config.WeatherProviderConfiguration{Primary: false},
			StaleTime:                "3s",
		}
		ws := w.NewService(&cfg)
		monkey.PatchInstanceMethod(reflect.TypeOf(ws), "MakeHttpRequest", func(_ *w.Service, cfg *config.WeatherProviderConfiguration) ([]byte, error) {
			if cfg.Primary == true {
				var buf bytes.Buffer
				binary.Write(&buf, binary.BigEndian, math.Inf(1))
				return buf.Bytes(), nil
			} else {
				data, _ := ioutil.ReadFile("../../test/secondary_success.json")
				return data, nil
			}
		})
		defer monkey.UnpatchAll()

		wd, err := ws.GetCurrentWeather()
		assert.Nil(t, err)
		assert.Equal(t, 16.63, wd.TemperatureDegrees)
		assert.Equal(t, 3.1, wd.WindSpeed)
	})
	t.Run("failure - primary weather bad json, econdary weather bad json", func(t *testing.T) {
		cfg := config.Configuration{
			PrimaryWeatherProvider:   config.WeatherProviderConfiguration{Primary: true},
			SecondaryWeatherProvider: config.WeatherProviderConfiguration{Primary: false},
			StaleTime:                "3s",
		}
		ws := w.NewService(&cfg)
		monkey.PatchInstanceMethod(reflect.TypeOf(ws), "MakeHttpRequest", func(_ *w.Service, cfg *config.WeatherProviderConfiguration) ([]byte, error) {
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, math.Inf(1))
			return buf.Bytes(), nil
		})
		defer monkey.UnpatchAll()
		wd, err := ws.GetCurrentWeather()
		assert.NotNil(t, err)
		assert.Nil(t, wd)
	})
}
