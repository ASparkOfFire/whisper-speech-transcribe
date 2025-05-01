package routes

import (
	"github.com/asparkoffire/whisper-speech-transcribe/api"
	"github.com/asparkoffire/whisper-speech-transcribe/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func setupHealthCheckRoutes(router *gin.RouterGroup) {
	router.GET("/healthz", api.ErrorWrapper((handlers.HandleHealthCheck)))
}
