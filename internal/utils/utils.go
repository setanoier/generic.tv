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
