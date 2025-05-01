package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(app *gin.Engine) {
	setupErrorRoutes(app)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	setupHealthCheckRoutes(v1)
	setupTranscribeRoutes(v1)
}
