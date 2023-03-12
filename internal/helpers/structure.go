package helpers

func GetNotEmpty[T comparable](a T, b T) T {
	if a == *new(T) {
		return b
	}
	return a
}
