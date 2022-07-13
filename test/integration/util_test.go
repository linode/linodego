package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
)

func assertDateSet(t *testing.T, compared *time.Time) {
	emptyTime := time.Time{}
	if compared == nil || *compared == emptyTime {
		t.Errorf("Expected date to be set, got %v", compared)
	}
}

func mockRequestBodyValidate(t *testing.T, expected interface{}, response interface{}) httpmock.Responder {
	return func(request *http.Request) (*http.Response, error) {
		eType := reflect.TypeOf(expected)
		result := reflect.New(eType)

		i := result.Interface()

		data, err := ioutil.ReadAll(request.Body)
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

func mockRequestURL(t *testing.T, path string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("/[a-zA-Z0-9]+/%s", strings.TrimPrefix(path, "/")))
}
