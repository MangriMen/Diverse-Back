package server

import (
	"fmt"
	"os"

	"github.com/MangriMen/Value-Back/api/db"
	"github.com/gofiber/fiber/v2"
)

func InitApi() {
	app := fiber.New()

	db, err := db.OpenDBConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	app.Listen(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
