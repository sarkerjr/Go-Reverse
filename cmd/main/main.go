package main

import (
	"sarkerjr.com/go-reverse/internal/server"
	"sarkerjr.com/go-reverse/pkg/initializer"
)

func main() {
	initializer.Initialize()

	server.StartServer()
}
