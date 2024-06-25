package linodego

func Pointer[T any](value T) *T {
	return &value
}
