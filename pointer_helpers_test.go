package linodego

import (
	"testing"
)

// TestPointer tests the Pointer helper function with various types
func TestPointer(t *testing.T) {
	// Test with an integer
	intValue := 11
	intPtr := Pointer(intValue)
	if *intPtr != intValue {
		t.Errorf("Expected %d, got %d", intValue, *intPtr)
	}

	// Test with a float
	floatValue := 1.23
	floatPtr := Pointer(floatValue)
	if *floatPtr != floatValue {
		t.Errorf("Expected %f, got %f", floatValue, *floatPtr)
	}

	// Test with a string
	stringValue := "hello world"
	stringPtr := Pointer(stringValue)
	if *stringPtr != stringValue {
		t.Errorf("Expected %s, got %s", stringValue, *stringPtr)
	}

	// Test with a boolean
	boolValue := true
	boolPtr := Pointer(boolValue)
	if *boolPtr != boolValue {
		t.Errorf("Expected %t, got %t", boolValue, *boolPtr)
	}

	// Test with a struct
	type myStruct struct {
		Field1 int
		Field2 string
	}
	structValue := myStruct{Field1: 1, Field2: "test"}
	structPtr := Pointer(structValue)
	if structPtr.Field1 != structValue.Field1 || structPtr.Field2 != structValue.Field2 {
		t.Errorf("Expected %+v, got %+v", structValue, *structPtr)
	}
}
