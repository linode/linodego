package linodego

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
)

type testResultNestedType struct {
	NestedInt    int    `json:"nested_int"`
	NestedString string `json:"nested_string"`
}

type testResultType struct {
	Foo  string               `json:"foo"`
	Bar  int                  `json:"bar"`
	Cool testResultNestedType `json:"cool"`
}

func TestRequestHelpers_get(t *testing.T) {
	client := CreateMockClient(t)

	desiredResponse := testResultType{
		Foo: "test",
		Bar: 123,
		Cool: testResultNestedType{
			NestedInt:    456,
			NestedString: "test2",
		},
	}

	httpmock.RegisterRegexpResponder("GET", MockRequestURL(t, "/foo/bar"),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse))

	result, err := doGETRequest[testResultType](
		context.Background(),
		client,
		"/foo/bar",
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*result, desiredResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(result, desiredResponse))
	}
}

func TestRequestHelpers_post(t *testing.T) {
	client := CreateMockClient(t)

	desiredResponse := testResultType{
		Foo: "test",
		Bar: 123,
		Cool: testResultNestedType{
			NestedInt:    456,
			NestedString: "test2",
		},
	}

	httpmock.RegisterRegexpResponder("POST", MockRequestURL(t, "/foo/bar"),
		MockRequestBodyValidate(t, desiredResponse, desiredResponse))

	result, err := doPOSTRequest[testResultType](
		context.Background(),
		client,
		"/foo/bar",
		desiredResponse,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*result, desiredResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(result, desiredResponse))
	}
}

func TestRequestHelpers_put(t *testing.T) {
	client := CreateMockClient(t)

	desiredResponse := testResultType{
		Foo: "test",
		Bar: 123,
		Cool: testResultNestedType{
			NestedInt:    456,
			NestedString: "test2",
		},
	}

	httpmock.RegisterRegexpResponder("PUT", MockRequestURL(t, "/foo/bar"),
		MockRequestBodyValidate(t, desiredResponse, desiredResponse))

	result, err := doPUTRequest[testResultType](
		context.Background(),
		client,
		"/foo/bar",
		desiredResponse,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*result, desiredResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(result, desiredResponse))
	}
}

func TestRequestHelpers_delete(t *testing.T) {
	client := CreateMockClient(t)

	httpmock.RegisterRegexpResponder(
		"DELETE",
		MockRequestURL(t, "/foo/bar/foo%20bar"),
		httpmock.NewStringResponder(200, "{}"),
	)

	if err := doDELETERequest(
		context.Background(),
		client,
		formatAPIPath("/foo/bar/%s", "foo bar"),
	); err != nil {
		t.Fatal(err)
	}
}
