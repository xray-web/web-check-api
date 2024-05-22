package routers

import (
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
)

var header controllers.HeaderController
var archive controllers.ArchivesController
var cookies controllers.CookiesController
var carbon controllers.CarbonController
var blockLists controllers.BlockListsController

func WebCheckRoutes(route *gin.Engine) {
	api := route.Group("/api")
	{
		api.GET("/headers", header.GetHeaders)
		api.GET("/archives", archive.ArchivesHandler)
		api.GET("/cookies", cookies.CookiesHandler)
		api.GET("/carbon", carbon.CarbonHandler)
		api.GET("/block-lists", blockLists.BlockListsHandler)
	}
}
