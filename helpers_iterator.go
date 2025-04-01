package linodego

import (
	"iter"
	"slices"
)

// Map returns a new iterator of the values in the given iterator transformed using the given transform function.
func Map[I, O any](values iter.Seq[I], transform func(I) O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for value := range values {
			if !yield(transform(value)) {
				return
			}
		}
	}
}

// MapSlice returns a new slice of the values in the given slice transformed using the given transform function.
func MapSlice[I, O any](values []I, transform func(I) O) []O {
	return slices.Collect(Map(slices.Values(values), transform))
}
