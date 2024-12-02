package initializer

import (
	"log"

	"github.com/joho/godotenv"
	"sarkerjr.com/go-reverse/config"
)

func Initialize() {
	// initialize environment variable package
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	// initialize configs
	config.LoadConfig()

}
