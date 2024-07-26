package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"generic.tv/internal/utils"
)

// Video represents the structure of a video item in the YouTube API response
type Video struct {
	ID struct {
		VideoID string `json:"videoId"`
	} `json:"id"`
	Snippet struct {
		Title string `json:"title"`
	} `json:"snippet"`
}

// YouTubeResponse represents the structure of the response from the YouTube API
type YouTubeResponse struct {
	Items []Video `json:"items"`
}

// FetchLatestVideoLink returns the link to the latest video from the specified YouTube channel
func FetchLatestVideoLink() (string, string, error) {
	apiKey := utils.ReadStringFromFile("../storage/gapi.txt")
	channelID := "UCGlYKd-FR4g0Tp4wF6_wxig"
	baseURL := "https://www.googleapis.com/youtube/v3/search"
	params := url.Values{}
	params.Add("part", "snippet")
	params.Add("channelId", channelID)
	params.Add("order", "date")
	params.Add("type", "video")
	params.Add("maxResults", "1")
	params.Add("key", apiKey)

	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", "", fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("unexpected response status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %v", err)
	}

	var ytResponse YouTubeResponse
	if err := json.Unmarshal(bodyBytes, &ytResponse); err != nil {
		return "", "", fmt.Errorf("error unmarshalling JSON response: %v", err)
	}

	if len(ytResponse.Items) == 0 {
		return "", "", fmt.Errorf("no videos found for channel ID %s", channelID)
	}

	latestVideo := ytResponse.Items[0]
	videoTitle := latestVideo.Snippet.Title
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", latestVideo.ID.VideoID)
	return videoTitle, videoURL, nil
}
