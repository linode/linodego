package linodego

import "testing"

func TestComparisonOperator(t *testing.T) {
	var opTests = []struct {
		in  ComparisonOperator
		out string
	}{
		{Eq, "+eq"},
		{Neq, "+neq"},
		{Gt, "+gt"},
		{Gte, "+gte"},
		{Lt, "+lt"},
		{Lte, "+lte"},
		{Contains, "+contains"},
	}
	for _, tests := range opTests {
		out := tests.in.String()
		if out != tests.out {
			t.Fatal(out, " doesn't match ", tests.out)
		}
	}
}

func TestLogicalOperator(t *testing.T) {
	var opTests = []struct {
		in  LogicalOperator
		out string
	}{
		{LogicalOr, "+or"},
		{LogicalAnd, "+and"},
	}
	for _, tests := range opTests {
		out := tests.in.String()
		if out != tests.out {
			t.Fatal(out, " doesn't match ", tests.out)
		}
	}
}

func TestFilter(t *testing.T) {
	expected := `"+and": [{"vcpus": {"+gte": 12}, {"class": "standard"}]`
	c1 := &Comparison{
		Column:   "vcpus",
		Operator: Gte,
		Value:    12,
	}
	c2 := &Comparison{
		Column:   "class",
		Operator: Eq,
		Value:    "standard",
	}
	out := And(c1, c2)
	if out.JSON() != expected {
		t.Fatal(out.JSON(), " doesn't match ", expected)
	}
}
