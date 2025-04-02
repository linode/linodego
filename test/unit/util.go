package unit

import (
	"encoding/json"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// assertJSONObjectsSimilar returns whether the given two structures implementing JSON
// are equivalent, only accounting for fields with shared JSON keys.
//
// This is primarily used to ensure that the GetCreateOptions() and GetUpdateOptions()
// functions are implemented correctly.
func assertJSONObjectsSimilar[TA, TB any](t testing.TB, a TA, b TB) {
	// Encoding and decoding JSON here is hacky, but it
	// lets us avoid some ugly type reflection pointer logic
	aJSON, err := json.Marshal(a)
	require.NoError(t, err)

	bJSON, err := json.Marshal(b)
	require.NoError(t, err)

	var aParsed, bParsed map[string]any

	require.NoError(t, json.Unmarshal(aJSON, &aParsed))
	require.NoError(t, json.Unmarshal(bJSON, &bParsed))

	assertJSONObjectsSimilarInner(t, []string{}, aParsed, bParsed)
}

func assertJSONObjectsSimilarInner(t testing.TB, path []string, a, b any) {
	a = normalizeEmptyValues(a)
	b = normalizeEmptyValues(b)

	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)

	aKind := aValue.Kind()
	bKind := bValue.Kind()

	require.Equalf(
		t,
		aKind,
		bKind,
		"%s type mismatch: %s != %s",
		strings.Join(path, "."),
		aKind,
		bKind,
	)

	switch aValue.Kind() {
	case reflect.Map:
		for _, key := range aValue.MapKeys() {
			aFieldValue := aValue.MapIndex(key)
			bFieldValue := bValue.MapIndex(key)

			if !bFieldValue.IsValid() {
				// This key is not shared so we can ignore it
				continue
			}

			assertJSONObjectsSimilarInner(
				t,
				slices.Concat(path, []string{key.String()}),
				aFieldValue.Interface(),
				bFieldValue.Interface(),
			)
		}
	case reflect.Slice:
		require.Equalf(
			t,
			aValue.Len(),
			bValue.Len(),
			"%s slice length mismatch: %d != %d",
			strings.Join(path, "."),
			aValue.Len(),
			bValue.Len(),
		)

		for index := range aValue.Len() {
			aFieldValue := aValue.Index(index)
			bFieldValue := bValue.Index(index)

			assertJSONObjectsSimilarInner(
				t,
				slices.Concat(path, []string{strconv.Itoa(index)}),
				aFieldValue.Interface(),
				bFieldValue.Interface(),
			)
		}
	default:
		require.Equal(
			t,
			a,
			b,
			"%s value mismatch: %v != %v",
			strings.Join(path, "."),
			a,
			b,
		)
	}
}

// normalizeEmptyValues normalizes the given value for use in JSON object diffing,
// primarily replacing any map, slice, or array values with nil.
//
// This is necessary because an empty length-having type is functionally equivalent
// to a nil value when using GetCreateOptions(...) and GetUpdateOptions(...).
func normalizeEmptyValues(v any) any {
	vValue := reflect.ValueOf(v)
	vKind := vValue.Kind()

	if (vKind == reflect.Map || vKind == reflect.Slice || vKind == reflect.Array) && vValue.Len() < 1 {
		return nil
	}

	return v
}
