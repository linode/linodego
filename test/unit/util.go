package unit

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertJSONObjectsSimilar returns whether the given two structures implementing JSON
// are equivalent, only accounting for fields with shared JSON keys.
//
// This is primarily used to ensure that the GetCreateOptions() and GetUpdateOptions()
// functions are implemented correctly.
func assertJSONObjectsSimilar[TA, TB any](t testing.TB, a TA, b TB) {
	assertJSONObjectsSimilarInner(t, []string{}, a, b)
}

func assertJSONObjectsSimilarInner[TA, TB any](t testing.TB, path []string, a TA, b TB) {
	aValue := derefValueRecursive(reflect.ValueOf(a))
	bValue := derefValueRecursive(reflect.ValueOf(b))

	aFields := aggregateJSONFields(reflect.ValueOf(a))
	bFields := aggregateJSONFields(reflect.ValueOf(b))

	require.Equalf(
		t,
		aValue.Kind(),
		bValue.Kind(),
		"%s kind mismatch: %s != %s",
		path,
		aValue.Kind(),
		bValue.Kind(),
	)

	switch aValue.Kind() {
	case reflect.Slice:
		assert.Equalf(
			t,
			aValue.Len(),
			bValue.Len(),
			"%s slice length mismatch: %d != %d",
			path,
			aValue.Len(),
			bValue.Len(),
		)

		for index := range aValue.Len() {
			assertJSONObjectsSimilarInner(
				t,
				slices.Concat(path, []string{strconv.Itoa(index)}),
				aValue.Index(index),
				bValue.Index(index),
			)
		}
	case reflect.Map:
		aKeys := aValue.MapKeys()
		bKeys := bValue.MapKeys()

		assert.Equalf(
			t,
			aKeys,
			bKeys,
			"%s map keys mismatch: %v != %v",
			path,
			aKeys,
			bKeys,
		)

		for _, key := range aKeys {
			assertJSONObjectsSimilarInner(
				t,
				slices.Concat(path, []string{key.String()}),
				aValue.MapIndex(key),
				bValue.MapIndex(key),
			)
		}
	case reflect.Struct:
		for key, aFieldValue := range aFields {
			bFieldValue, ok := bFields[key]
			if !ok {
				// This key isn't shared, nothing to do here
				continue
			}

			assertJSONObjectsSimilarInner(t, slices.Concat(path, []string{key}), aFieldValue, bFieldValue)
		}
	default:
		assert.Equal(
			t,
			aValue.Interface(),
			bValue.Interface(),
			"%s value mismatch: %s != %s",
			path,
			aValue.Interface(),
			bValue.Interface(),
		)
	}
}

func aggregateJSONFields(v reflect.Value) map[string]reflect.Value {
	vType := derefTypeRecursive(v.Type())

	result := make(map[string]reflect.Value, vType.NumField())

	for fieldNum := range vType.NumField() {
		field := vType.Field(fieldNum)

		jsonTag, jsonTagOk := field.Tag.Lookup("json")
		if !jsonTagOk {
			// No JSON tag is defined, nothing to do here
			continue
		}

		if jsonTag == "-" {
			continue
		}

		jsonTagKey := strings.Split(jsonTag, ",")[0]
		result[jsonTagKey] = derefValueRecursive(derefValueRecursive(v).FieldByName(field.Name))
	}

	return result
}

func derefTypeRecursive(v reflect.Type) reflect.Type {
	if v.Kind() == reflect.Ptr {
		return derefTypeRecursive(v.Elem())
	}

	return v
}

func derefValueRecursive(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return derefValueRecursive(v.Elem())
	}

	return v
}
