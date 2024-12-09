package server

import (
	"starwars/service/internal/errors"
	"starwars/service/internal/handler"
	"starwars/service/internal/logger"
	"starwars/service/internal/request"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// router is the router instance.
var router *gin.Engine

// Init initializes the local router instance.
func Init() {
	router = gin.Default()

	router.Use(errors.RecoveryMiddleware(), request.RequestIdMiddleware(), logger.Middleware())

	router.GET("/people", handler.RetrievePeople)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

// Run runs the router.
func Run() {
	router.Run()
}
