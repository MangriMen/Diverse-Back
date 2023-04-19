package jwthelpers

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims is a struct that extends the JWT Registered Claim Names
// by adding a user id.
type CustomClaims struct {
	jwt.RegisteredClaims

	ID string `json:"id,omitempty"`
}

// GetTokenClaims verifying the token and extracting it JWT claims.
func GetTokenClaims(c *fiber.Ctx) (*CustomClaims, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func getToken(c *fiber.Ctx) string {
	const maxTokenParts = 2

	bearerToken := strings.Split(c.Get("Authorization"), " ")

	if len(bearerToken) == maxTokenParts {
		return bearerToken[1]
	}

	return ""
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := getToken(c)

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
