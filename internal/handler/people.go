package handler

import (
	"net/http"

	"github.com/pegondo/starwars/service/internal/errors"
	"github.com/pegondo/starwars/service/internal/logger"
	"github.com/pegondo/starwars/service/internal/request"
	"github.com/pegondo/starwars/service/internal/resources/swapi"

	"github.com/gin-gonic/gin"
)

// retrievePeopleHandlerName is the name of the retrieve people handler.
const retrievePeopleHandlerName = "retrieve people"

func RetrievePeople(c *gin.Context) {
	l := logger.Logger(c)
	l.Info().Msgf("received request to the %s endpoint", retrievePeopleHandlerName)

	params, err := request.Params(c)
	if err != nil {
		l.Warn().Msgf("invalid request parameters :: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	people, err := swapi.RetrievePeople(params)
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
	c.JSON(statusCode, Response[swapi.Person]{
		Data: people.Results,
	})
}
