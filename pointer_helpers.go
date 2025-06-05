package linodego

/*
Pointer takes a value of any type T and returns a pointer to that value.
Go does not allow directly creating pointers to literals, so Pointer enables
abstraction away the pointer logic.

Example:

		booted := true

		createOpts := linodego.InstanceCreateOptions{
			Booted: &booted,
		}

		can be replaced with

		createOpts := linodego.InstanceCreateOptions{
			Booted: linodego.Pointer(true),
		}
*/

func Pointer[T any](value T) *T {
	return &value
}

// DoublePointer creates a double pointer to a value of type T.
// It returns a **T, where a nil double pointer (**T == nil) represents a null field,
// and omitting the field entirely indicates that the field won't be included in the request body.
// This is useful for APIs that distinguish between null and omitted fields.
//
// Example:
//
//	// For a field that should be null in the API payload:
//	value := DoublePointer(42) // Returns **int pointing to 42
//
//	// Omit the field in the struct to indicate it won't be included in the request body.
//	nullValue := DoublePointer[int](nil) // Returns **int that is nil
func DoublePointer[T any](value T) **T {
	valuePtr := &value
	return &valuePtr
}
