package main

import (
	v1 "github.com/asparkoffire/whisper-speech-transcribe/api/v1/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"time"
)

func main() {
	app := gin.Default()
	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Change to specific origins in production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1.SetupRoutes(app)

	app.Run("0.0.0.0:8010")
}
