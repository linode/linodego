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

// TestDoublePointer tests the DoublePointer helper function with various types
func TestDoublePointer(t *testing.T) {
	// Test with an integer
	intValue := 42
	intDoublePtr := DoublePointer(intValue)
	if **intDoublePtr != intValue {
		t.Errorf("Expected %d, got %d", intValue, **intDoublePtr)
	}

	// Test nil pointer for int (should be nil if explicitly set)
	nilIntPtr := Pointer[*int](nil)
	if *nilIntPtr != nil {
		t.Errorf("Expected nil pointer, got %v", nilIntPtr)
	}

	// Test with a string
	strValue := "double"
	strDoublePtr := DoublePointer(strValue)
	if **strDoublePtr != strValue {
		t.Errorf("Expected %s, got %s", strValue, **strDoublePtr)
	}

	// Test with a boolean
	boolValue := false
	boolDoublePtr := DoublePointer(boolValue)
	if **boolDoublePtr != boolValue {
		t.Errorf("Expected %t, got %t", boolValue, **boolDoublePtr)
	}

	// Test with a struct
	type myStruct struct {
		Field int
	}
	structValue := myStruct{Field: 7}
	structDoublePtr := DoublePointer(structValue)
	if (**structDoublePtr).Field != structValue.Field {
		t.Errorf("Expected %+v, got %+v", structValue, **structDoublePtr)
	}
}
