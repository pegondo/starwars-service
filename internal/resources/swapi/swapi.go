package swapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Source: https://swapi.dev/documentation
const (
	// swapiBaseUrl is the base URL of the SWAPI.
	swapiBaseUrl = "https://swapi.dev/api/"
	// swapiPageSize is the SWAPI fixed page size.
	swapiPageSize = 10
)

// request performs a HTTP request to the given URL returns its response.
func request[T PersonResponse](url string) (response T, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return response, fmt.Errorf("error while performing the request :: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error while reading the response body :: %v", err)
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("error while parsing the response to a JSON :: %v", err)
	}

	return response, err
}
