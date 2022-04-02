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

func weather(w http.ResponseWriter, r *http.Request) {
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

}
func handleRequest() {
	http.HandleFunc("/api/weather", weather)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal("Error starting server", err)
	}
}

func main() {
	handleRequest()
}

func init() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	/*
		l:= logrus.New()
		l.SetReportCaller(true)
		l.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				filename := path.Base(frame.File)
				return fmt.Sprintf("%s()",frame.Function),fmt.Sprintf("%s:%d",filename,frame.Line)
			},
			DisableColors: false,
			FullTimestamp: true,
		}
		err := os.MkdirAll("logs",0664)
		if err != nil {
			panic(err)
		}
		allfile,err :=os.OpenFile("logs/all.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0640)
		if err != nil{
			panic(err)
		}
		l.SetOutput(io.Discard)
		l.AddHook(&writerHook{
			Writer: []io.Writer{allfile,os.Stdout},
			Loglevels: logrus.AllLevels,

		})
		l.SetLevel(logrus.TraceLevel)

	*/
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

/*type writerHook struct {
	Writer []io.Writer
	Loglevels []logrus.Level

}
func (hook *writerHook) Fire(entry *logrus.Entry) error{
	line , err := entry.String()
	if err != nil{
		return err
	}
	for _, v := range hook.Writer{
		v.Write([]byte(line))
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level{
	return hook.Loglevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func (l *Logger) GetLoggerWithField(k string,v interface{}) Logger {
	return Logger{l.WithField(k,v)}

}

*/
