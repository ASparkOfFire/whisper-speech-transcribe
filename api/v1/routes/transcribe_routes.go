package routes

import (
	"github.com/asparkoffire/whisper-speech-transcribe/api"
	"github.com/asparkoffire/whisper-speech-transcribe/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func setupTranscribeRoutes(router *gin.RouterGroup) {
	transcribePrefix := router.Group("/transcribe")

	transcribePrefix.GET("/ws", handlers.HandleWebSocketTranscribe)
	transcribePrefix.POST("/", api.ErrorWrapper(handlers.HandleFileUploadTranscribe))
}
