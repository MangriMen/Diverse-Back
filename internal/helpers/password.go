// Package helpers provides common functionality for application
package helpers

import (
	"github.com/MangriMen/Diverse-Back/configs"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a string as input and uses the bcrypt algorithm to generate a hash of the password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), configs.PasswordEncryptCost)
	return string(bytes), err
}

// CheckPasswordHash takes a strings as input, password and hash, and
// returns a boolean value indicating whether the hash matches the given password.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
