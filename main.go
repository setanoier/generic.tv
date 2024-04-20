package main

import (
	"fmt"
	"github.com/gempir/go-twitch-irc/v4"
	"log"
	"os"
)

func readStringFromFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	strData := string(data)
	return strData
}

func main() {
	accessToken := readStringFromFile(".client") // access token is oauth token of bot account
	client := twitch.NewClient("yourTwitch", accessToken)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Message)
		if message.Message == "!ping" {
			client.Say(message.Channel, "pong!\n")
		}
	})

	client.Join("yourTwitch")

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
}
