// Package main loads the environment and start the web server
package main

import (
	"log"

	"github.com/MangriMen/Diverse-Back/api/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("environment variables were not loaded")
	}

	server.InitAPI()
}
