package coordinates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Define a structure for the JSON response
type GeocodingResponse struct {
	Features []struct {
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

// GetCityCoordinates fetches the coordinates of a city from the MapTiler Geocoding API
func GetCityCoordinates(city, apiKey string) (float64, float64, error) {
	// Construct the request URL
	baseURL := "https://api.maptiler.com/geocoding/"
	query := url.QueryEscape(city)
	fullURL := fmt.Sprintf("%s%s.json?key=%s", baseURL, query, apiKey)

	// Make the HTTP GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("received non-200 response status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	var geocodingResponse GeocodingResponse
	err = json.Unmarshal(body, &geocodingResponse)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	// Check if any features were returned
	if len(geocodingResponse.Features) == 0 {
		return 0, 0, fmt.Errorf("no coordinates found for city: %s", city)
	}

	// Extract the coordinates
	coordinates := geocodingResponse.Features[0].Geometry.Coordinates
	if len(coordinates) < 2 {
		return 0, 0, fmt.Errorf("incomplete coordinates data for city: %s", city)
	}

	return coordinates[1], coordinates[0], nil // lat, lng
}
