package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	weatherAPI = "http://api.openweathermap.org/data/2.5/weather"
	apiToken   = "2ae6e50490e35649304462a5d8f6cb29"
	units      = "metric"
)

type Item struct {
	City  string  `json:"city"`
	Tempa float64 `json:"temp"`
}

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	}
}

func main() {
	var a string
	fmt.Print("Enter city: ")
	_, _ = fmt.Scan(&a)
	city := a
	req, err := http.NewRequest("GET", weatherAPI, nil)
	if err != nil {
		log.Fatal(err)
	}
	values := url.Values{}
	values.Add("appid", apiToken)
	values.Add("units", units)
	values.Add("q", city)
	req.URL.RawQuery = values.Encode()

	cli := http.Client{}

	resp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var Weather WeatherResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&Weather)
	if err != nil {
		log.Fatal(err)
	}
	item := Item{Tempa: 18, City: city}
	jitem, err := json.Marshal(item)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(string(jitem))

	http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%s\n", jitem)

	})
	fmt.Println("Server is listening...")
	_ = http.ListenAndServe("127.0.0.1:8080", nil)

}
