package commands

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"generic.tv/internal/utils"
	"generic.tv/internal/weather"
	"generic.tv/internal/youtube"

	"github.com/gempir/go-twitch-irc/v4"
)

func Response(client *twitch.Client, admin string) {
	isDrawing := false
	players := make(map[string]bool)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		switch {
		case message.Message == "!ping":
			client.Say(message.Channel, "pong!\n")

		case strings.HasPrefix(message.Message, "!weather"):
			tokens := strings.Split(message.Message, " ")

			if len(tokens) == 1 {
				tokens = append(tokens, "Innopolis")
			}

			weather, err := weather.GetWeather(tokens[1])

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
		case message.Message == "!yt":
			title, link, err := youtube.FetchLatestVideoLink()
			if err != nil {
				client.Say(message.Channel, "Something went wrong...")
				return
			} else {
				client.Say(message.Channel, title+": "+link)
			}
		}
	})
}

func Start() error {
	if _, err := os.Stat("../storage/oauth.txt"); os.IsNotExist(err) {
		// Handle case when oauth.txt does not exist
	}
	accessToken := utils.ReadStringFromFile("../storage/oauth.txt") // Access token is OAuth token of the bot account
	client := twitch.NewClient("setanoier", accessToken)
	admin := "setanoier"
	Response(client, admin)

	client.Join("setanoier")

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	return nil
}
