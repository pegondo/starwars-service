package server

import "github.com/gin-gonic/gin"

// router is the router instance.
var router *gin.Engine

// Init initializes the local router instance.
func Init() {
	router = gin.Default()
}

// Run runs the router.
func Run() {
	router.Run()
}
