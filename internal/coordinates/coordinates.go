package coordinates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type GeoFeature struct {
	Geometry GeoGeometry `json:"geometry"`
}

type GeoGeometry struct {
	Coordinates []float64 `json:"coordinates"`
}

func GetCoordinatesByTown(town, apiKey string) (float64, float64, error) {
	encodedTownName := url.QueryEscape(town)

	apiURL := fmt.Sprintf("https://api.maptiler.com/geocoding/%s.json?key=%s", encodedTownName, apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, 0, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("error reading response body: %v", err)
	}

	var geoFeature GeoFeature
	err = json.Unmarshal(body, &geoFeature)
	if err != nil {
		return 0, 0, fmt.Errorf("error decoding JSON: %v", err)
	}

	coordinates := geoFeature.Geometry.Coordinates
	if len(coordinates) < 2 {
		return 0, 0, fmt.Errorf("no coordinates found in response")
	}

	fmt.Println(1)

	latitude := coordinates[1]
	longitude := coordinates[0]

	return latitude, longitude, nil
}
