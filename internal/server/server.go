package server

import (
	"github.com/pegondo/starwars-service/internal/errors"
	"github.com/pegondo/starwars-service/internal/handler"
	"github.com/pegondo/starwars-service/internal/logger"
	"github.com/pegondo/starwars-service/internal/request"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// router is the router instance.
var router *gin.Engine

// Init initializes the local router instance.
func Init() {
	router = gin.Default()

	router.Use(errors.RecoveryMiddleware(), request.RequestIdMiddleware(), logger.Middleware())

	api := router.Group("/api")
	api.GET(handler.PeopleEndpoint, handler.RetrievePeople)
	api.GET(handler.PlanetEndpoint, handler.RetrievePlanets)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

// Run runs the router.
func Run() {
	router.Run()
}
