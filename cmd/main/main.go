package main

import (
	"log"

	"sarkerjr.com/go-reverse/internal/server"
	"sarkerjr.com/go-reverse/pkg/initializer"
)

func main() {
	initializer.Initialize()

	if err := server.StartServer(); err != nil {
		log.Fatal(err)
	}
}
