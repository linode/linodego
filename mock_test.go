package linodego

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var validTestAPIKey = "NOTANAPIKEY"

func MockRequestURL(t *testing.T, path string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("/[a-zA-Z0-9]+/%s", strings.TrimPrefix(path, "/")))
}

func MockRequestBodyValidate(t *testing.T, expected interface{}, response interface{}) httpmock.Responder {
	return func(request *http.Request) (*http.Response, error) {
		eType := reflect.TypeOf(expected)
		result := reflect.New(eType)

		i := result.Interface()

		data, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatal(err)
		}

		if err := json.Unmarshal(data, &i); err != nil {
			t.Fatal(err)
		}

		// Deref the pointer if necessary
		if result.Kind() == reflect.Pointer {
			result = result.Elem()
		}

		resultValue := result.Interface()

		if !reflect.DeepEqual(expected, resultValue) {
			t.Fatalf("request body does not match request options: %s", cmp.Diff(expected, resultValue))
		}

		return httpmock.NewJsonResponse(200, response)
	}
}

func CreateMockClient(t *testing.T) *Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: validTestAPIKey})

	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}
	httpmock.ActivateNonDefault(client)

	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	c := NewClient(client)
	return &c
}
