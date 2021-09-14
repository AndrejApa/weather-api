package main

import (
	"encoding/json"
	"fmt"
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

type CityResp struct {
	City string  `json:"city"`
	Temp float64 `json:"temp"`
}

func main() {
	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("city")
		fmt.Println("city: ", city)
		// перепроверить
		m, _ := getWeatherByCity(city)
		_, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	})
	fmt.Println("Server is listening...")
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
