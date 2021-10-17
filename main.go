// http://127.0.0.1:8080/api/weather?city='your city'

package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	weatherAPI = "http://api.openweathermap.org/data/2.5/weather"
	units      = "metric"
)

type MyResponse struct {
	City string  `json:"city"`
	Temp float64 `json:"temp"`
}

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

		resp := MyResponse{
			City: data.Name,
			Temp: data.Main.Temp,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}

	})
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
}

func init() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
func getWeatherByCity(city string) (WeatherResponse, error) {
	req, err := http.NewRequest("GET", weatherAPI, nil)
	if err != nil {
		return WeatherResponse{}, err
	}
	APIKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		fmt.Println("API_KEY is not present")
	} else {
		fmt.Printf("API_KEY: %t\n", ok)
	}
	values := url.Values{}
	values.Add("appid", APIKey)
	values.Add("units", units)
	values.Add("q", city)
	req.URL.RawQuery = values.Encode()

	cli := http.Client{}

	resp, err := cli.Do(req)
	if err != nil {
		return WeatherResponse{}, err
	}

	var weather WeatherResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&weather)
	if err != nil {
		return WeatherResponse{}, err
	}
	return weather, nil
}
