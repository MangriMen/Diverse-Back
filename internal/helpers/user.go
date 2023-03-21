package helpers

func GetNotEmpty[T comparable](a T, b T) T {
	if a == *new(T) {
		return b
	}
	return a
}

func ptr[T any](obj T) *T { return &obj }
