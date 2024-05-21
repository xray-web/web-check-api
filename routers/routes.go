package routers

import (
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
)

var header controllers.HeaderController

func WebCheckRoutes(route *gin.Engine) {
	api := route.Group("/api")
	{
		api.GET("/headers", header.GetHeaders)
	}
}
