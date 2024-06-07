package routers

import (
	"github.com/xray-web/web-check-api/middleware"

	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	RegisterRoutes(router) //routes register

	return router
}
