package linodego

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/linode/linodego/internal/testutil"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
)

type testResultNestedType struct {
	NestedInt    int    `json:"nested_int"`
	NestedString string `json:"nested_string"`
}

type testResultType struct {
	ID   int                  `json:"id"`
	Bar  *string              `json:"bar"`
	Foo  string               `json:"foo"`
	Cool testResultNestedType `json:"cool"`
}

var testResponse = testResultType{
	Foo: "test",
	ID:  123,
	Cool: testResultNestedType{
		NestedInt:    456,
		NestedString: "test2",
	},
}

func TestRequestHelpers_get(t *testing.T) {
	client := testutil.CreateMockClient(t, NewClient)

	httpmock.RegisterRegexpResponder("GET", testutil.MockRequestURL("/foo/bar"),
		httpmock.NewJsonResponderOrPanic(200, &testResponse))

	result, err := doGETRequest[testResultType](
		context.Background(),
		client,
		"/foo/bar",
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*result, testResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(result, testResponse))
	}
}

func TestRequestHelpers_post(t *testing.T) {
	client := testutil.CreateMockClient(t, NewClient)

	httpmock.RegisterRegexpResponder("POST", testutil.MockRequestURL("/foo/bar"),
		testutil.MockRequestBodyValidate(t, testResponse, testResponse))

	result, err := doPOSTRequest[testResultType](
		context.Background(),
		client,
		"/foo/bar",
		testResponse,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*result, testResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(result, testResponse))
	}
}

func TestRequestHelpers_postNoOptions(t *testing.T) {
	client := testutil.CreateMockClient(t, NewClient)

	httpmock.RegisterRegexpResponder("POST", testutil.MockRequestURL("/foo/bar"),
		testutil.MockRequestBodyValidateNoBody(t, testResponse))

	result, err := doPOSTRequest[testResultType, any](
		context.Background(),
		client,
		"/foo/bar",
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*result, testResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(result, testResponse))
	}
}

func TestRequestHelpers_put(t *testing.T) {
	client := testutil.CreateMockClient(t, NewClient)

	httpmock.RegisterRegexpResponder("PUT", testutil.MockRequestURL("/foo/bar"),
		testutil.MockRequestBodyValidate(t, testResponse, testResponse))

	result, err := doPUTRequest[testResultType](
		context.Background(),
		client,
		"/foo/bar",
		testResponse,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*result, testResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(result, testResponse))
	}
}

func TestRequestHelpers_putNoOptions(t *testing.T) {
	client := testutil.CreateMockClient(t, NewClient)

	httpmock.RegisterRegexpResponder("PUT", testutil.MockRequestURL("/foo/bar"),
		testutil.MockRequestBodyValidateNoBody(t, testResponse))

	result, err := doPUTRequest[testResultType, any](
		context.Background(),
		client,
		"/foo/bar",
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*result, testResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(result, testResponse))
	}
}

func TestRequestHelpers_delete(t *testing.T) {
	client := testutil.CreateMockClient(t, NewClient)

	httpmock.RegisterRegexpResponder("DELETE", testutil.MockRequestURL("/foo/bar/foo%20bar"),
		httpmock.NewStringResponder(200, "{}"))

	if err := doDELETERequest(
		context.Background(),
		client,
		formatAPIPath("/foo/bar/%s", "foo bar"),
	); err != nil {
		t.Fatal(err)
	}
}

func TestRequestHelpers_paginateAll(t *testing.T) {
	const totalResults = 4123

	client := testutil.CreateMockClient(t, NewClient)

	numRequests := 0

	httpmock.RegisterRegexpResponder("GET", testutil.MockRequestURL("/foo/bar"),
		mockPaginatedResponse(buildPaginatedEntries(totalResults), &numRequests))

	response, err := getPaginatedResults[testResultType](
		context.Background(),
		client,
		"/foo/bar",
		&ListOptions{
			PageSize: 500,
			Filter:   "{\"foo\": \"bar\"}",
		},
	)
	require.NoError(t, err)

	require.Equal(t, 9, numRequests)
	require.Len(t, response, totalResults)

	for i := 0; i < totalResults; i++ {
		entry := response[i]

		require.Equal(t, i, entry.ID)
		require.Equal(t, fmt.Sprintf("test-%d", i), *entry.Bar)
	}
}

func TestRequestHelpers_paginateSingle(t *testing.T) {
	client := testutil.CreateMockClient(t, NewClient)

	numRequests := 0

	httpmock.RegisterRegexpResponder("GET", testutil.MockRequestURL("/foo/bar"),
		mockPaginatedResponse(buildPaginatedEntries(12), &numRequests))

	response, err := getPaginatedResults[testResultType](
		context.Background(),
		client,
		"/foo/bar",
		&ListOptions{
			PageOptions: &PageOptions{
				Page: 3,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if numRequests != 1 {
		t.Fatalf("expected 1 request, got %d", numRequests)
	}

	if len(response) != 3 {
		t.Fatalf("expected 3 results, got %d", len(response))
	}

	for i := 0; i < 3; i++ {
		entry := response[i]
		desiredID := i + 6
		if entry.ID != desiredID {
			t.Fatalf("expected id %d, got %d", desiredID, entry.ID)
		}
	}
}

func buildPaginatedEntries(numEntries int) []testResultType {
	result := make([]testResultType, numEntries)

	for i := 0; i < numEntries; i++ {
		bar := fmt.Sprintf("test-%d", i)
		result[i] = testResultType{
			Bar: &bar,
			Foo: "foo",
			ID:  i,
		}
	}

	return result
}

func mockPaginatedResponse(
	entries []testResultType, numRequests *int,
) httpmock.Responder {
	return func(request *http.Request) (*http.Response, error) {
		*numRequests++

		// Default page size for testing purposes
		pageSize := 3

		page, err := strconv.Atoi(request.URL.Query().Get("page"))
		if err != nil {
			return nil, err
		}

		if request.URL.Query().Has("page_size") {
			pageSize, err = strconv.Atoi(request.URL.Query().Get("page_size"))
			if err != nil {
				return nil, err
			}
		}

		// Clamp the top index to prevent out of bounds issues
		lastEntryIdx := pageSize * page
		if lastEntryIdx > len(entries) {
			lastEntryIdx = len(entries)
		}

		return httpmock.NewJsonResponse(
			200,
			paginatedResponse[testResultType]{
				Page:    page,
				Pages:   int(math.Ceil(float64(len(entries)) / float64(pageSize))),
				Results: pageSize,
				Data:    entries[pageSize*(page-1) : lastEntryIdx],
			},
		)
	}
}
