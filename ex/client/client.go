package client

import (
	goErrors "errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pegondo/starwars/service/internal/handler"
	"github.com/pegondo/starwars/service/internal/resources/swapi"
)

// ErrInvalidResponseFormat is the error returned when the response format is
// invalid.
var ErrInvalidResponseFormat = goErrors.New("invalid response format")

// Client a HTTP client to call the endpoints of the server.
type Client struct {
	addr string
}

// New creates and returns a new client.
func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

// Response is the response returned by the client.
type Response[T swapi.Resource] struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int
	// Response is the content of the response.
	Response handler.Response[T]
}

// requestOpts represents the options of the request.
type requestOpts struct {
	// pageNumber is the number of the pagination page requested.
	pageNumber string
	// pageSize is the size of the pagination page.
	pageSize string
	// search is the search criteria. If "", no search criteria should be
	// applied.
	search string
	// sortField is the field to sort by. If "", no sort criteria should be
	// applied.
	sortField string
	// sortOrder is the order to sort by.
	sortOrder string
}

// NewRequestOpts creates and returns a request options structure with the given
// data.
func NewRequestOpts(
	pageNumber,
	pageSize,
	search,
	sortField,
	sortOrder string,
) requestOpts {
	return requestOpts{
		pageNumber: pageNumber,
		pageSize:   pageSize,
		search:     search,
		sortField:  sortField,
		sortOrder:  sortOrder,
	}
}

// buildUrl builds a URL to request the addr in c with the given endpoint and
// options.
func (c *Client) buildUrl(endpoint string, opts requestOpts) string {
	// TODO: Consider using this library when building urls in the code.
	reqUrl, _ := url.Parse(fmt.Sprintf("%s/%s", c.addr, endpoint))

	query := url.Values{}
	if opts.pageNumber != "" {
		query.Add("page", opts.pageNumber)
	}
	if opts.pageSize != "" {
		query.Add("pageSize", opts.pageSize)
	}
	if opts.search != "" {
		query.Add("search", opts.search)
	}
	if opts.sortField != "" {
		query.Add("sortField", opts.sortField)
	}
	if opts.sortOrder != "" {
		query.Add("sortOrder", opts.sortOrder)
	}
	reqUrl.RawQuery = query.Encode()

	return reqUrl.String()
}

// isStatusOk returns whether the status code corresponds to a 2xx or not.
func isStatusOk(statusCode int) bool {
	return http.StatusOK <= statusCode && statusCode < http.StatusMultipleChoices
}
