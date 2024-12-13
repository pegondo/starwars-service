package handler

import (
	"net/http"

	"github.com/pegondo/starwars-service/internal/errors"
	"github.com/pegondo/starwars-service/internal/logger"
	"github.com/pegondo/starwars-service/internal/request"
	"github.com/pegondo/starwars-service/internal/resources/swapi"

	"github.com/gin-gonic/gin"
)

// retrievePlanetsHandlerName is the name of the retrieve planets handler.
const retrievePlanetsHandlerName = "retrive planets"

func RetrievePlanets(c *gin.Context) {
	l := logger.Logger(c)
	l.Info().Msgf("received request to the %s endpoint", retrievePlanetsHandlerName)

	params, err := request.Params(c)
	if err != nil {
		l.Warn().Msgf("invalid request parameters :: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	planets, err := swapi.RetrievePlanets(params)
	if err != nil {
		// If there is an issue while requesting for the planets, return a 500.
		l.Error().Msg(err.Error())
		err = errors.New(errors.InternalServerErrorCode, errors.InternalServerErrorMsg)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	statusCode := getStatusCode(planets)
	c.JSON(statusCode, Response[swapi.Planet]{
		Data: planets.Results,
	})
}
