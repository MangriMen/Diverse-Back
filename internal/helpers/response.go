package helpers

import "github.com/gofiber/fiber/v2"

const (
	DefaultError uint32 = 1 << iota
)

func GetResponse(err error, responseType uint32) fiber.Map {
	switch responseType {
	case DefaultError:
		return fiber.Map{
			"error":   true,
			"message": err.Error(),
		}
	default:
		return fiber.Map{}
	}
}
