package handler

import (
	"net/http"

	"github.com/pegondo/starwars-service/internal/resources/swapi"
)

// Response represents the response of a handler.
type Response[T swapi.Resource] struct {
	// Data is the resource data.
	Data []T `json:"data"`
	// Count is the number of elements in the
	Count int `json:"count"`
}

// getStatusCode returns the HTTP status code to return regarding the number of
// elements in the response. getStatusCode may misbehave if resp.Count is
// negative.
func getStatusCode[T swapi.Resource](resp swapi.SwapiResponse[T]) int {
	if areAllResourcesPresent := resp.Count <= len(resp.Results); !areAllResourcesPresent {
		return http.StatusPartialContent
	}
	return http.StatusOK
}
