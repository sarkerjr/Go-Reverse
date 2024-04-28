package initializer

import (
	"log"

	"github.com/joho/godotenv"
)

func Initialize() {
	// initialize environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
}
