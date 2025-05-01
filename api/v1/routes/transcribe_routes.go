package routes

import (
	"github.com/asparkoffire/whisper-speech-transcribe/api/v1/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func setupTranscribeRoutes(router *gin.RouterGroup) {
	transcribePrefix := router.Group("/transcribe")

	transcribeHandlers, err := handlers.NewTranscribeHandler()
	if err != nil {
		log.Fatalf("Failed to create transcribe handler: %v\n", err)
	}
	transcribePrefix.GET("/ws", transcribeHandlers.HandleWebSocketTranscribe)
}
