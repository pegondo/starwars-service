package handler

import (
	"strconv"

	"github.com/pegondo/starwars/service/internal/resources/swapi"

	"github.com/gin-gonic/gin"
)

const (
	// pageParamKey is the key to get the page query parameter.
	pageParamKey = "page"

	// defaultPageSize is the default page size used to request the API.
	defaultPageSize = 15
	// pageSizeParamKey is the key to get the page size query parameter.
	pageSizeParamKey = "pageSize"

	// searchParamKey is the key to get the search query parameter.
	searchParamKey = "search"
)

// Response represents the response of a handler.
type Response[T swapi.Resource] struct {
	// Data is the resource data.
	Data []T `json:"data"`
}

// getNumericParam returns the parameter with the given key from the context. If
// the param isn't defined, getNumericParam returns the default value provided.
func getNumericParam(c *gin.Context, key string, defaultValue int) (value int, err error) {
	valueStr := c.DefaultQuery(key, strconv.Itoa(defaultValue))
	return strconv.Atoi(valueStr)
}
