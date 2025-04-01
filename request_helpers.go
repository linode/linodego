package linodego

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

// paginatedResponse represents a single response from a paginated
// endpoint.
type paginatedResponse[T any] struct {
	Page    int `json:"page"    url:"page,omitempty"`
	Pages   int `json:"pages"   url:"pages,omitempty"`
	Results int `json:"results" url:"results,omitempty"`
	Data    []T `json:"data"`
}

// getPaginatedResults aggregates results from the given
// paginated endpoint using the provided ListOptions.
// nolint:funlen
func getPaginatedResults[T any](
	ctx context.Context,
	client *Client,
	endpoint string,
	opts *ListOptions,
) ([]T, error) {
	result := make([]T, 0)

	if opts == nil {
		opts = &ListOptions{PageOptions: &PageOptions{Page: 0}}
	}

	if opts.PageOptions == nil {
		opts.PageOptions = &PageOptions{Page: 0}
	}

	// Makes a request to a particular page and
	// appends the response to the result
	handlePage := func(page int) error {
		var resultType paginatedResponse[T]
		opts.Page = page

		params := requestParams{
			Response: &resultType,
		}

		// Create a mutator to all user-provided list options to the request
		mutator := createListOptionsToRequestMutator(opts)

		// Make the request using doRequest
		err := client.doRequest(ctx, http.MethodGet, endpoint, params, &mutator)
		if err != nil {
			return err
		}

		// Extract the result from the response
		opts.Page = page
		opts.Pages = resultType.Pages
		opts.Results = resultType.Results

		// Append the data to the result slice
		result = append(result, resultType.Data...)
		return nil
	}

	// Determine starting page
	startingPage := 1
	pageDefined := opts.Page > 0

	if pageDefined {
		startingPage = opts.Page
	}

	// Get the first page
	if err := handlePage(startingPage); err != nil {
		return nil, err
	}

	// If a specific page is defined, return the result
	if pageDefined {
		return result, nil
	}

	// Get the remaining pages
	for page := 2; page <= opts.Pages; page++ {
		if err := handlePage(page); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// doGETRequest runs a GET request using the given client and API endpoint,
// and returns the result
func doGETRequest[T any](
	ctx context.Context,
	client *Client,
	endpoint string,
) (*T, error) {
	var resultType T
	params := requestParams{
		Response: &resultType,
	}

	err := client.doRequest(ctx, http.MethodGet, endpoint, params, nil)
	if err != nil {
		return nil, err
	}

	return &resultType, nil
}

// doPOSTRequest runs a PUT request using the given client, API endpoint,
// and options/body.
func doPOSTRequest[T, O any](
	ctx context.Context,
	client *Client,
	endpoint string,
	options ...O,
) (*T, error) {
	var resultType T
	numOpts := len(options)
	if numOpts > 1 {
		return nil, fmt.Errorf("invalid number of options: %d", numOpts)
	}

	params := requestParams{
		Response: &resultType,
	}
	if numOpts > 0 && !isNil(options[0]) {
		body, err := json.Marshal(options[0])
		if err != nil {
			return nil, err
		}
		params.Body = bytes.NewReader(body)
	}

	err := client.doRequest(ctx, http.MethodPost, endpoint, params, nil)
	if err != nil {
		return nil, err
	}
	return &resultType, nil
}

// doPOSTRequestNoResponseBody runs a POST request using the given client, API endpoint,
// and options/body. It expects only empty response from the endpoint.
func doPOSTRequestNoResponseBody[T any](
	ctx context.Context,
	client *Client,
	endpoint string,
	options ...T,
) error {
	_, err := doPOSTRequest[any, T](ctx, client, endpoint, options...)
	return err
}

// doPOSTRequestNoRequestBody runs a POST request using the given client and API endpoint.
// It does not expect a request body but does expect a response from the endpoint.
func doPOSTRequestNoRequestBody[T any](
	ctx context.Context,
	client *Client,
	endpoint string,
) (*T, error) {
	return doPOSTRequest[T, any](ctx, client, endpoint)
}

// doPOSTRequestNoRequestResponseBody runs a POST request where no request body is needed and no response body
// is expected from the endpoints.
func doPOSTRequestNoRequestResponseBody(
	ctx context.Context,
	client *Client,
	endpoint string,
) error {
	return doPOSTRequestNoResponseBody(ctx, client, endpoint, struct{}{})
}

// doPUTRequest runs a PUT request using the given client, API endpoint,
// and options/body.
func doPUTRequest[T, O any](
	ctx context.Context,
	client *Client,
	endpoint string,
	options ...O,
) (*T, error) {
	var resultType T
	numOpts := len(options)
	if numOpts > 1 {
		return nil, fmt.Errorf("invalid number of options: %d", numOpts)
	}

	params := requestParams{
		Response: &resultType,
	}
	if numOpts > 0 && !isNil(options[0]) {
		body, err := json.Marshal(options[0])
		if err != nil {
			return nil, err
		}
		params.Body = bytes.NewReader(body)
	}

	err := client.doRequest(ctx, http.MethodPut, endpoint, params, nil)
	if err != nil {
		return nil, err
	}
	return &resultType, nil
}

// doDELETERequest runs a DELETE request using the given client
// and API endpoint.
func doDELETERequest(
	ctx context.Context,
	client *Client,
	endpoint string,
) error {
	params := requestParams{}
	err := client.doRequest(ctx, http.MethodDelete, endpoint, params, nil)
	return err
}

// formatAPIPath allows us to safely build an API request with path escaping
func formatAPIPath(format string, args ...any) string {
	escapedArgs := make([]any, len(args))
	for i, arg := range args {
		if typeStr, ok := arg.(string); ok {
			arg = url.PathEscape(typeStr)
		}

		escapedArgs[i] = arg
	}

	return fmt.Sprintf(format, escapedArgs...)
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	// Check for nil pointers
	v := reflect.ValueOf(i)
	return v.Kind() == reflect.Ptr && v.IsNil()
}
