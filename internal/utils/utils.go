package utils

import (
	"log"
	"os"
)

func ReadStringFromFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func GetWeatherEmoji(temp int) string {
	switch {
	case temp <= 13:
		return "❄️"
	case temp > 13 && temp <= 20:
		return "🧥"
	case temp > 20 && temp <= 30:
		return "👕"
	case temp > 30:
		return "💀"
	default:
		return "???"
	}
}
