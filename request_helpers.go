package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
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
	var resultType paginatedResponse[T]

	result := make([]T, 0)

	req := client.R(ctx).SetResult(resultType)

	var options ListOptions

	// Apply all user-provided list options to the request
	// if applicable
	if opts != nil {
		options = *opts

		if err := applyListOptionsToRequest(&options, req); err != nil {
			return nil, err
		}
	}

	numPages := 0

	// Makes a request to a particular page and
	// appends the response to the result
	handlePage := func(page int) error {
		req.SetQueryParam("page", strconv.Itoa(page))

		res, err := coupleAPIErrors(req.Get(endpoint))
		if err != nil {
			return err
		}

		response := res.Result().(*paginatedResponse[T])

		// Only update the number of pages if it hasn't been set yet
		if numPages == 0 {
			numPages = response.Pages
		}

		result = append(result, response.Data...)
		return nil
	}

	// This helps simplify the logic below
	startingPage := 1
	pageDefined := options.PageOptions != nil && options.PageOptions.Page > 0

	if pageDefined {
		startingPage = options.PageOptions.Page
	}

	// Get the first page
	if err := handlePage(startingPage); err != nil {
		return nil, err
	}

	// If the user has explicitly specified a page, we don't
	// need to get any other pages.
	if pageDefined {
		return result, nil
	}

	// Get the rest of the pages
	for page := 2; page <= numPages; page++ {
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

	req := client.R(ctx).SetResult(&resultType)
	r, err := coupleAPIErrors(req.Get(endpoint))
	if err != nil {
		return nil, err
	}

	return r.Result().(*T), nil
}

// doPOSTRequest runs a PUT request using the given client, API endpoint,
// and options/body.
func doPOSTRequest[T, O any](
	ctx context.Context,
	client *Client,
	endpoint string,
	options O,
) (*T, error) {
	var resultType T

	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	req := client.R(ctx).SetResult(&resultType).SetBody(string(body))
	r, err := coupleAPIErrors(req.Post(endpoint))
	if err != nil {
		return nil, err
	}

	return r.Result().(*T), nil
}

// doPUTRequest runs a PUT request using the given client, API endpoint,
// and options/body.
func doPUTRequest[T, O any](
	ctx context.Context,
	client *Client,
	endpoint string,
	options O,
) (*T, error) {
	var resultType T

	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	req := client.R(ctx).SetResult(&resultType).SetBody(string(body))
	r, err := coupleAPIErrors(req.Put(endpoint))
	if err != nil {
		return nil, err
	}

	return r.Result().(*T), nil
}

// doDELETERequest runs a DELETE request using the given client
// and API endpoint.
func doDELETERequest(
	ctx context.Context,
	client *Client,
	endpoint string,
) error {
	req := client.R(ctx)
	_, err := coupleAPIErrors(req.Delete(endpoint))
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
