package handler

import (
	"github.com/pegondo/starwars/service/internal/resources/swapi"
)

// Response represents the response of a handler.
type Response[T swapi.Resource] struct {
	// Data is the resource data.
	Data []T `json:"data"`
}
