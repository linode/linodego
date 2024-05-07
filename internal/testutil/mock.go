package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"golang.org/x/oauth2"
)

var validTestAPIKey = "NOTANAPIKEY"

func MockRequestURL(path string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("/[a-zA-Z0-9]+/%s", strings.TrimPrefix(path, "/")))
}

func MockRequestBodyValidate(t *testing.T, expected any, response any) httpmock.Responder {
	t.Helper()

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

		return httpmock.NewJsonResponse(http.StatusOK, response)
	}
}

func MockRequestBodyValidateNoBody(t *testing.T, response any) httpmock.Responder {
	t.Helper()

	return func(request *http.Request) (*http.Response, error) {
		if request.Body != nil {
			t.Fatal("got request body when no request body was expected")
		}

		return httpmock.NewJsonResponse(http.StatusOK, response)
	}
}

// CreateMockClient is generic because importing the linodego package will result
// in a cyclic dependency. This pattern isn't ideal but works for now.
func CreateMockClient[T any](t *testing.T, createFunc func(*http.Client) T) *T {
	t.Helper()

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

	result := createFunc(client)
	return &result
}

type Logger interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

func CreateLogger() *TestLogger {
	l := &TestLogger{L: log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)}
	return l
}

var _ Logger = (*TestLogger)(nil)

type TestLogger struct {
	L *log.Logger
}

func (l *TestLogger) Errorf(format string, v ...interface{}) {
	l.output("ERROR RESTY "+format, v...)
}

func (l *TestLogger) Warnf(format string, v ...interface{}) {
	l.output("WARN RESTY "+format, v...)
}

func (l *TestLogger) Debugf(format string, v ...interface{}) {
	l.output("DEBUG RESTY "+format, v...)
}

func (l *TestLogger) output(format string, v ...interface{}) {
	if len(v) == 0 {
		l.L.Print(format)
		return
	}
	l.L.Printf(format, v...)
}
