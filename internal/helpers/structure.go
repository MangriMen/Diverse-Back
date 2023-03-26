package helpers

// GetNotEmpty returns the non-empty value between a and b.
func GetNotEmpty[T comparable](a T, b T) T {
	if a == *new(T) {
		return b
	}
	return a
}

// Ptr returns pointer to an object.
func Ptr[T any](obj T) *T { return &obj }
