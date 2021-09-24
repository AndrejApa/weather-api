package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	weatherAPI = "http://api.openweathermap.org/data/2.5/weather"
	apiToken   = "2ae6e50490e35649304462a5d8f6cb29"
	units      = "metric"
)

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func main() {
	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("city")
		data, err := getWeatherByCity(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)

	})
	_ = http.ListenAndServe("127.0.0.1:8080", nil)

}

func getWeatherByCity(city string) (WeatherResponse, error) {
	req, err := http.NewRequest("GET", weatherAPI, nil)
	if err != nil {
		return WeatherResponse{}, err
	}
	values := url.Values{}
	values.Add("appid", apiToken)
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
