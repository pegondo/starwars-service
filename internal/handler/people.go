package handler

import (
	"net/http"
	"starwars/service/internal/errors"
	"starwars/service/internal/logger"
	"starwars/service/internal/resources/swapi"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	// retrievePeopleHandlerName is the name of the retrieve people handler.
	retrievePeopleHandlerName = "retrieve people"

	// pageParamKey is the key to get the page query parameter.
	pageParamKey = "page"

	// defaultPageSize is the default page size used to request the API.
	defaultPageSize = 15
	// pageSizeParamKey is the key to get the page size query parameter.
	pageSizeParamKey = "pageSize"
)

// PeopleResponse represents the response the retrieve people handler returns.
type PeopleResponse struct {
	// Data is the person data.
	Data []swapi.Person `json:"data"`
}

// getNumericParam returns the parameter with the given key from the context. If
// the param isn't defined, getNumericParam returns the default value provided.
func getNumericParam(c *gin.Context, key string, defaultValue int) (value int, err error) {
	valueStr := c.DefaultQuery(key, strconv.Itoa(defaultValue))
	return strconv.Atoi(valueStr)
}

func RetrievePeople(c *gin.Context) {
	l := logger.Logger(c)
	l.Info().Msgf("received request to the %s endpoint", retrievePeopleHandlerName)

	page, err := getNumericParam(c, pageParamKey, 1)
	if err != nil {
		// If the page number is invalid, return a 404.
		l.Warn().Msgf("invalid page number :: %v", err)
		err = errors.New(errors.InvalidPageErrorCode, errors.InvalidPageErrorMsg)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pageSize, err := getNumericParam(c, pageSizeParamKey, defaultPageSize)
	if err != nil {
		// If the page number is invalid, return a 404.
		l.Warn().Msgf("invalid page size :: %v", err)
		err = errors.New(errors.InvalidPageSizeErrorCode, errors.InvalidPageSizeErrorMsg)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	people, err := swapi.RetrievePeople(page, pageSize)
	if err != nil {
		// If there is an issue while requesting for the people, return a 500.
		l.Error().Msg(err.Error())
		err = errors.New(errors.InternalServerErrorCode, errors.InternalServerErrorMsg)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	statusCode := http.StatusOK
	if areAllResources := people.Count == len(people.Results); !areAllResources {
		statusCode = http.StatusPartialContent
	}
	c.JSON(statusCode, PeopleResponse{
		Data: people.Results,
	})
}
