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
		TestInt64    int64  `query:"test_int64"`
		TestBool     bool   `query:"test_bool"`
		TestUntagged string
	}

	inst := TestStruct{
		TestInt:      123,
		TestString:   "test+string",
		TestInt64:    567,
		TestBool:     true,
		TestUntagged: "cool",
	}

	expectedOutput := map[string]string{
		"test_int":    "123",
		"test_string": "test+string",
		"test_int64":  "567",
		"test_bool":   "true",
	}

	result, err := flattenQueryStruct(inst)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(result, expectedOutput) {
		t.Fatalf("diff in result: %v", cmp.Diff(result, expectedOutput))
	}
}
