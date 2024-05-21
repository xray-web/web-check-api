package config

import (
	"os"

	"github.com/gin-gonic/gin"
)

func ServerConfig(router *gin.Engine) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		return err
	}
	return err
}
