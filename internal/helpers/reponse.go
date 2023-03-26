package helpers

import (
	"github.com/MangriMen/Diverse-Back/internal/responses"
	"github.com/gofiber/fiber/v2"
)

// Response used to compact return default error response.
func Response(c *fiber.Ctx, responseStatus int, err interface{}) error {
	return c.Status(responseStatus).JSON(responses.BaseResponseBody{
		Error:   true,
		Message: &err,
	})
}
