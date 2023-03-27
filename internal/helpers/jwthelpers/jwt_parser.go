package jwthelpers

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// TokenMetadata is a struct that represents the metadata associated with a JWT access token.
// It contains expires and id fields.
type TokenMetadata struct {
	Expires int64
	ID      string
}

// CustomClaims is a struct that represents the custom claims associated with a JWT access token.
// It contains expires and id fields.
type CustomClaims struct {
	Expires float64 `json:"exp"`
	ID      string  `json:"id"`

	jwt.StandardClaims
}

// GetTokenMetadata verifying the token and extracting its expiration time and ID from the JWT claims.
// Returns the TokenMetadata struct.
func GetTokenMetadata(c *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(c)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(CustomClaims)

	if ok && token.Valid {
		expires := int64(claims.Expires)
		id := claims.ID

		return &TokenMetadata{
			Expires: expires,
			ID:      id,
		}, nil
	}

	return nil, err
}

func getToken(c *fiber.Ctx) string {
	const maxTokenParts = 2

	bearToken := c.Get("Authorization")

	onlyToken := strings.Split(bearToken, " ")

	if len(onlyToken) == maxTokenParts {
		return onlyToken[1]
	}

	return ""
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := getToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
