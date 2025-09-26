package server

import (
	"github.com/gin-gonic/gin"	// Gin, a third-party web framework, is a very popular and high-performance HTTP framework in the Go community.
	"net/http"	// HTTP tools from the Go standard library (using this to import the http.StatusOK constant = 200)

)

func SetupRouter() *gin.Engine {	// *gin.Engine Returns a Gin HTTP server instance (you can call .Run() to start the server).
	router := gin.Default()		// The default HTTP router instance includes log/recovery, which can be used to quickly enable the server. And this line is equals to "router := gin.New(), router.Use(gin.Logger()), router.Use(gin.Recovery())"

	router.GET("/", func(c *gin.Context) {	// c is *gin.Context, which represents the context of this HTTP request. You can read parameters, write responses, set cookies, etc. from it.
		c.String(http.StatusOK, "Welcome to the Chatrtoom!!")	// health check endpoint
	})

	return  router		// return Gin HTTP server instance
}