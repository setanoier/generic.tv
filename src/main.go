package main

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
)

func readStringFromFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

type WeatherResponse struct {
	CurrentWeather struct {
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		Winddirection float64 `json:"winddirection"`
	} `json:"current_weather"`
}

func getWeather() (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"latitude":        "35.6895",  // example latitude for Tokyo
			"longitude":       "139.6917", // example longitude for Tokyo
			"current_weather": "true",
		}).
		Get("https://api.open-meteo.com/v1/forecast")

	if err != nil {
		return "", err
	}

	var weather WeatherResponse
	err = json.Unmarshal(resp.Body(), &weather)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Current temperature: %.1f°C, Windspeed: %.1f km/h, Winddirection: %.1f°", weather.CurrentWeather.Temperature, weather.CurrentWeather.Windspeed, weather.CurrentWeather.Winddirection), nil
}

func main() {
	accessToken := readStringFromFile("../data/oauth") // access token is oauth token of bot account
	client := twitch.NewClient("setanoier", accessToken)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Message)
		switch message.Message {
		case "!ping":
			client.Say(message.Channel, "pong!\n")
		case "!weather":
			weather, err := getWeather()
			if err != nil {
				client.Say(message.Channel, "Failed to get weather information.")
			}
			client.Say(message.Channel, weather)
		}
	})

	client.Join("setanoier")

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
}
