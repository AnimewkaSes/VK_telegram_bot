package util

import (
	"github.com/subosito/gotenv"
	"log"
)

func LoadConfig() {
	if err := gotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	return
}
