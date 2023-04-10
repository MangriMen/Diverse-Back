// Package main loads the environment and start the web server
package main

import (
	"github.com/MangriMen/Diverse-Back/api/server"
	"github.com/MangriMen/Diverse-Back/internal/helpers"
)

func main() {
	if !helpers.IsRunningInContainer() {
		helpers.LoadEnvironment(".env")
	}

	server.SetupAPI()
}
