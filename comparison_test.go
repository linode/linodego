package linodego

import "testing"

func TestFilter(t *testing.T) {
	expected := `"vcpus": {"+gte": 12}, {"class": "standard"}`
	f := Filter{}
	f.Add(&Comp{"vcpus", Gte, 12})
	f.Add(&Comp{"class", Eq, "standard"})
	if f.JSON() != expected {
		t.Fatal(f.JSON(), " doesn't match ", expected)
	}
}

func TestAscending(t *testing.T) {
	expected := `{"vcpus": {"+gte": 12}}, {"class": "standard"}, "+order_by": "class", "+order": "asc"`
	f := Filter{
		Order:   Ascending,
		OrderBy: "class",
	}
	f.Add(&Comp{"vcpus", Gte, 12})
	f.Add(&Comp{"class", Eq, "standard"})
	if f.JSON() != expected {
		t.Fatal(f.JSON(), " doesn't match ", expected)
	}
}

func TestDescending(t *testing.T) {
	expected := `{"vcpus": {"+gte": 12}}, {"class": "standard"}, "+order_by": "class", "+order": "desc"`
	f := Filter{
		Order:   Descending,
		OrderBy: "class",
	}
	f.Add(&Comp{"vcpus", Gte, 12})
	f.Add(&Comp{"class", Eq, "standard"})
	if f.JSON() != expected {
		t.Fatal(f.JSON(), " doesn't match ", expected)
	}
}

func TestAnd(t *testing.T) {
	expected := `"+and": [{"vcpus": {"+gte": 12}}, {"class": "standard"}]`
	c1 := &Comp{"vcpus", Gte, 12}
	c2 := &Comp{"class", Eq, "standard"}
	out := And("", "", c1, c2)
	if out.JSON() != expected {
		t.Fatal(out.JSON(), " doesn't match ", expected)
	}
}

func TestOr(t *testing.T) {
	expected := `"+or": [{"vcpus": {"+gte": 12}}, {"class": "standard"}]`
	c1 := &Comp{"vcpus", Gte, 12}
	c2 := &Comp{"class", Eq, "standard"}
	out := Or("", "", c1, c2)
	if out.JSON() != expected {
		t.Fatal(out.JSON(), " doesn't match ", expected)
	}
}
