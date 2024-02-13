package linodego

import (
	"context"
	"strconv"
)

type PaginationResponse[T any] struct {
	Page    int `json:"page"    url:"page,omitempty"`
	Pages   int `json:"pages"   url:"pages,omitempty"`
	Results int `json:"results" url:"results,omitempty"`
	Data    []T `json:"data"`
}

func aggregatePaginatedResults[T any](
	ctx context.Context,
	client *Client,
	endpoint string,
	opts *ListOptions,
) ([]T, error) {
	var resultType PaginationResponse[T]

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

		response := res.Result().(*PaginationResponse[T])

		// Only update the number of pages if it hasn't been set yet
		if numPages == 0 {
			numPages = response.Pages
		}

		result = append(result, response.Data...)
		return nil
	}

	// This helps simplify the logic below
	startingPage := 0

	if options.PageOptions != nil && options.PageOptions.Page > 0 {
		startingPage = options.PageOptions.Page
	}

	// Get the first page
	if err := handlePage(startingPage); err != nil {
		return nil, err
	}

	// If the user has explicitly specified a page, we don't
	// need to get any other pages.
	if startingPage > 1 {
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
