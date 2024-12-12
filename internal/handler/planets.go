package handler

import (
	"net/http"
	"strings"

	"github.com/pegondo/starwars/service/internal/errors"
	"github.com/pegondo/starwars/service/internal/logger"
	"github.com/pegondo/starwars/service/internal/resources/swapi"

	"github.com/gin-gonic/gin"
)

// retrievePlanetsHandlerName is the name of the retrieve planets handler.
const retrievePlanetsHandlerName = "retrive planets"

func RetrievePlanets(c *gin.Context) {
	l := logger.Logger(c)
	l.Info().Msgf("received request to the %s endpoint", retrievePlanetsHandlerName)

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

	search := c.DefaultQuery(searchParamKey, "")
	search = strings.ToLower(search)

	sortCriteria := swapi.GetSortCriteria(c)

	planets, err := swapi.RetrievePlanets(page, pageSize, search, sortCriteria)
	if err != nil {
		// If there is an issue while requesting for the planets, return a 500.
		l.Error().Msg(err.Error())
		err = errors.New(errors.InternalServerErrorCode, errors.InternalServerErrorMsg)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	statusCode := http.StatusOK
	if areAllResources := planets.Count == len(planets.Results); !areAllResources {
		statusCode = http.StatusPartialContent
	}
	c.JSON(statusCode, Response[swapi.Planet]{
		Data: planets.Results,
	})
}
