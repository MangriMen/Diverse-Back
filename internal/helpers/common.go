package helpers

import (
	"io"
)

// CloseQuietly close io.Closer object quietly without returning error.
func CloseQuietly[T io.Closer](entity T) {
	if err := entity.Close(); err != nil {
		_ = entity.Close()
	}
}
