package weather

import (
	"encoding/json"
	"fmt"
	"generic.tv/internal/coordinates"
	"github.com/go-resty/resty/v2"
)

type WeatherResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
	} `json:"current_weather"`
}

func GetWeather(town []string) (string, error) {
	apiKey := "lM9XWQOYlh3z0uk3T5px"
	latitude, longitude, err := coordinates.GetCoordinatesByTown("London", apiKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(latitude, longitude)

	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"latitude":        fmt.Sprintf("%.6f", latitude),
			"longitude":       fmt.Sprintf("%.6f", longitude),
			"current_weather": "true",
		}).Get("https://api.open-meteo.com/v1/forecast")

	if err != nil {
		return "", err
	}

	var weather WeatherResponse
	err = json.Unmarshal(resp.Body(), &weather)
	if err != nil {
		return "", err
	}

	temperature := int(weather.CurrentWeather.Temperature)

	return fmt.Sprintf("Current temperature: %d°C", temperature), nil
}
