package client

import (
	goErrors "errors"
	"fmt"
	"net/url"
	"strconv"

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
	pageNumber int
	// pageSize is the size of the pagination page.
	pageSize int
	// search is the search criteria. If "", no search criteria should be
	// applied.
	search string
	// sort is the sort criteria. If nil, no sort criteria should be applied.
	sort *swapi.SortCriteria
}

// NewRequestOpts creates and returns a request options structure with the given
// data.
func NewRequestOpts(
	pageNumber,
	pageSize int,
	search string,
	sort *swapi.SortCriteria,
) requestOpts {
	return requestOpts{
		pageNumber: pageNumber,
		pageSize:   pageSize,
		search:     search,
		sort:       sort,
	}
}

// buildUrl builds a URL to request the addr in c with the given endpoint and
// options.
func (c *Client) buildUrl(endpoint string, opts requestOpts) string {
	// TODO: Consider using this library when building urls in the code.
	reqUrl, _ := url.Parse(fmt.Sprintf("%s/%s", c.addr, endpoint))

	query := url.Values{}
	if opts.pageNumber != 0 {
		query.Add("page", strconv.Itoa(opts.pageNumber))
	}
	if opts.pageSize != 0 {
		query.Add("pageSize", strconv.Itoa(opts.pageSize))
	}
	if opts.search != "" {
		query.Add("search", opts.search)
	}
	if opts.sort != nil {
		if opts.sort.Field != "" {
			query.Add("sortField", string(opts.sort.Field))
		}
		if opts.sort.Order != "" {
			query.Add("sortOrder", string(opts.sort.Order))
		}
	}
	reqUrl.RawQuery = query.Encode()

	return reqUrl.String()
}
