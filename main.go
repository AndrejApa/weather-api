// http://127.0.0.1:8080/api/weather?q='your city'

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	weatherAPI = "http://api.openweathermap.org/data/2.5/weather"
	units      = "metric"
)

type apiConfig struct {
	ApiToken string `json:"ApiToken"`
}

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func load(filename string) (apiConfig, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfig{}, err
	}
	var c apiConfig
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfig{}, err
	}
	return c, nil

}
func main() {
	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("q")
		data, err := getWeatherByCity(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			return
		}

	})
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		return
	}

}
func getWeatherByCity(city string) (WeatherResponse, error) {
	req, err := http.NewRequest("GET", weatherAPI, nil)
	if err != nil {
		return WeatherResponse{}, err
	}
	apiConfig, err := load(".apiConfig")
	if err != nil {
		return WeatherResponse{}, err
	}
	values := url.Values{}
	values.Add("appid", apiConfig.ApiToken)
	values.Add("units", units)
	values.Add("q", city)
	req.URL.RawQuery = values.Encode()

	cli := http.Client{}

	resp, err := cli.Do(req)
	if err != nil {
		return WeatherResponse{}, err
	}

	var Weather WeatherResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&Weather)
	if err != nil {
		return WeatherResponse{}, err
	}
	return Weather, nil
}
