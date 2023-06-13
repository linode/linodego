package linodego

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestFlattenQueryStruct(t *testing.T) {
	type TestStruct struct {
		TestInt      int    `query:"test_int"`
		TestString   string `query:"test_string"`
		TestString2  string `query:"test_string_2"`
		TestInt64    int64  `query:"test_int64"`
		TestIn64Ptr  *int64 `query:"test_int64_ptr"`
		TestBool     bool   `query:"test_bool"`
		TestUntagged string
	}

	testInt64 := int64(789)

	inst := TestStruct{
		TestInt:      123,
		TestString:   "test+string",
		TestString2:  "",
		TestInt64:    567,
		TestIn64Ptr:  &testInt64,
		TestBool:     true,
		TestUntagged: "cool",
	}

	expectedOutput := map[string]string{
		"test_int":       "123",
		"test_string":    "test+string",
		"test_int64":     "567",
		"test_int64_ptr": "789",
		"test_bool":      "true",
	}

	result, err := flattenQueryStruct(inst)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(result, expectedOutput) {
		t.Fatalf("diff in result: %v", cmp.Diff(result, expectedOutput))
	}
}
