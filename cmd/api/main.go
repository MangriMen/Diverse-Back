package main

import (
	"github.com/MangriMen/Value-Back/api/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	server.InitApi()
}
