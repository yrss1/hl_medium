package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	return r
}
