package linodego

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	expected := map[string]any{
		"vcpus": map[string]any{
			"+gte": 12,
		},
		"class": "standard",
	}

	expectedStr, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("failed to marshal expected json: %v", err)
	}

	f := Filter{}
	f.AddField(Gte, "vcpus", 12)
	f.AddField(Eq, "class", "standard")

	result, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshal filter: %v", err)
	}

	if !reflect.DeepEqual(result, expectedStr) {
		t.Fatal(string(result), " doesn't match ", string(expectedStr))
	}
}

func TestFilterAscending(t *testing.T) {
	expected := map[string]any{
		"vcpus": map[string]any{
			"+gte": 12,
		},
		"class":     "standard",
		"+order_by": "class",
		"+order":    "asc",
	}

	expectedStr, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("failed to marshal expected json: %v", err)
	}

	f := Filter{
		Order:   Ascending,
		OrderBy: "class",
	}
	f.AddField(Gte, "vcpus", 12)
	f.AddField(Eq, "class", "standard")

	result, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshal filter: %v", err)
	}

	if !reflect.DeepEqual(result, expectedStr) {
		t.Fatal(string(result), " doesn't match ", string(expectedStr))
	}
}

func TestFilterDescending(t *testing.T) {
	expected := map[string]any{
		"vcpus": map[string]any{
			"+gte": 12,
		},
		"class":     "standard",
		"+order_by": "class",
		"+order":    "desc",
	}

	expectedStr, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("failed to marshal expected json: %v", err)
	}

	f := Filter{
		Order:   Descending,
		OrderBy: "class",
	}
	f.AddField(Gte, "vcpus", 12)
	f.AddField(Eq, "class", "standard")

	result, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshal filter: %v", err)
	}

	if !reflect.DeepEqual(result, expectedStr) {
		t.Fatal(string(result), " doesn't match ", string(expectedStr))
	}
}

func TestFilterAnd(t *testing.T) {
	expected := map[string]any{
		"+and": []map[string]any{
			{
				"vcpus": map[string]any{
					"+gte": 12,
				},
			},
			{
				"class": "standard",
			},
		},
	}

	expectedStr, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("failed to marshal expected json: %v", err)
	}

	c1 := &Comp{"vcpus", Gte, 12}
	c2 := &Comp{"class", Eq, "standard"}
	out := And("", "", c1, c2)

	result, err := out.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshal filter: %v", err)
	}

	if !reflect.DeepEqual(result, expectedStr) {
		t.Fatal(string(result), " doesn't match ", string(expectedStr))
	}
}

func TestFilterOr(t *testing.T) {
	expected := map[string]any{
		"+or": []map[string]any{
			{
				"vcpus": map[string]any{
					"+gte": 12,
				},
			},
			{
				"class": "standard",
			},
		},
	}

	expectedStr, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("failed to marshal expected json: %v", err)
	}

	c1 := &Comp{"vcpus", Gte, 12}
	c2 := &Comp{"class", Eq, "standard"}
	out := Or("", "", c1, c2)

	result, err := out.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshal filter: %v", err)
	}

	if !reflect.DeepEqual(result, expectedStr) {
		t.Fatal(string(result), " doesn't match ", string(expectedStr))
	}
}
