package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func main() {
	db, err := sql.Open("sqlite3", "../data/towns.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS towns (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        latitude REAL NOT NULL,
        longitude REAL NOT NULL
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database and table create successfully")

	accessToken := readStringFromFile("../data/oauth") // Access token is OAuth token of the bot account
	client := twitch.NewClient("setanoier", accessToken)
	admin := "setanoier"
	response(client, admin)

	client.Join("setanoier")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
}

func readStringFromFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

type WeatherResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
	} `json:"current_weather"`
}

func getWeather() (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"latitude":        "55.754461",
			"longitude":       "48.742641",
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

func response(client *twitch.Client, admin string) {
	isDrawing := false
	players := make(map[string]bool)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Message)
		switch {
		case message.Message == "!ping":
			client.Say(message.Channel, "pong!\n")

		case message.Message == "!weather":
			weather, err := getWeather()
			if err != nil {
				client.Say(message.Channel, "Failed to get weather information.")
			} else {
				client.Say(message.Channel, weather)
			}

		case message.Message == "!drawing" && isDrawing && message.User.Name != admin:
			if !players[message.User.Name] {
				players[message.User.Name] = true
			} else {
				client.Say(message.Channel, "You are already participating!")
			}

		case strings.HasPrefix(message.Message, "!drawing"):
			if message.User.Name != admin {
				client.Say(message.Channel, "You do not have appropriate rights to start this event")
				return
			}

			tokens := strings.Split(message.Message, " ")
			if len(tokens) != 2 {
				client.Say(message.Channel, "Usage: !drawing minutes:seconds.")
				return
			}

			durationParts := strings.Split(tokens[1], ":")
			if len(durationParts) != 2 {
				client.Say(message.Channel, "Invalid time format: Use minutes:seconds.")
				return
			}

			minutes, err1 := strconv.Atoi(durationParts[0])
			seconds, err2 := strconv.Atoi(durationParts[1])
			if err1 != nil || err2 != nil {
				client.Say(message.Channel, "Invalid time format. Use integers for minutes and seconds.")
				return
			}

			duration := time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second

			client.Say(message.Channel, "The drawing has begun!")
			isDrawing = true

			timer := time.NewTimer(duration)
			go func() {
				<-timer.C
				client.Say(message.Channel, "The drawing has ended!")
				isDrawing = false
				keys := reflect.ValueOf(players).MapKeys()
				winner := keys[rand.Intn(len(keys))].String()
				client.Say(message.Channel, fmt.Sprintf("Winner: %s", winner))
				players = make(map[string]bool)
			}()
		}
	})
}

func insertTown(db *sql.DB, name string, latitude, longitude float64) {
	insertSQL := `INSERT INTO towns (name, latitude, longitude) VALUES (?, ?, ?)`
	_, err := db.Exec(insertSQL, name, latitude, longitude)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added: ", name)
}

func getCoordinatesByTown(db *sql.DB, townName string) (float64, float64, error) {
	var latitude, longitude float64
	query := `SELECT latitude, longitude FROM towns WHERE name = ?`
	row := db.QueryRow(query, townName)
	err := row.Scan(&latitude, &longitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, fmt.Errorf("no town found with name %s", townName)
		}
		return 0, 0, err
	}
	return latitude, longitude, nil
}
