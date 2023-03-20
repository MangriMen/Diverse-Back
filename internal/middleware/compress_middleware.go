package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func Compress(a *fiber.App) {
	a.Use(compress.New(compress.ConfigDefault))
}
