package swapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

// Source: https://swapi.dev/documentation
const (
	// swapiBaseUrl is the base URL of the SWAPI.
	swapiBaseUrl = "https://swapi.dev/api/"
	// swapiPageSize is the SWAPI fixed page size.
	swapiPageSize = 10
)

// Resource represents a SWAPI resource the API serves.
type Resource interface {
	Person | Planet
}

// SwapiResponse represents the SWAPI response for a resource T.
type SwapiResponse[T Resource] struct {
	// Count is the number of resources in the collection.
	Count int `json:"count"`
	// Resoults are the paginated elements in the collection.
	Results []T `json:"results"`
}

// request performs a HTTP request to the given URL returns its response.
func request[T Resource](url string) (response SwapiResponse[T], err error) {
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

// retrievePage is a recursive solution to request SWAPI given a variable page
// size.
func retrievePage[T Resource](
	resource SwapiResponse[T],
	endpoint,
	search string,
	remainingResources,
	pageNumber,
	offset int,
) (
	swapiResp SwapiResponse[T],
	err error,
) {
	if remainingResources <= 0 {
		return resource, nil
	}

	url := fmt.Sprintf("%s/%s?page=%d", swapiBaseUrl, endpoint, pageNumber)
	if search != "" {
		url = fmt.Sprintf("%s&search=%s", url, search)
	}
	swapiResp, err = request[T](url)
	if err != nil {
		return swapiResp, fmt.Errorf("error while requesting the %s endpoint :: %v", endpoint, err)
	}
	if swapiResp.Count == 0 {
		// If there are no results, return an empty array.
		return SwapiResponse[T]{
			Count:   0,
			Results: []T{},
		}, nil
	}

	remainingResources = int(math.Min(float64(remainingResources), float64(swapiResp.Count)))

	minIdx := offset
	maxIdx := int(math.Min(swapiPageSize, float64(minIdx+remainingResources)))
	swapiResp.Results = append(resource.Results, swapiResp.Results[minIdx:maxIdx]...)
	numElementsAdded := maxIdx - minIdx
	return retrievePage(swapiResp, endpoint, search, remainingResources-numElementsAdded, pageNumber+1, 0)
}
