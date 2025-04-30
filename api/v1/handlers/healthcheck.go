package handlers

import (
	"net/http"
	"time"

	"github.com/asparkoffire/whisper-speech-transcribe/api"
	"github.com/gin-gonic/gin"
)

func HandleHealthCheck(c *gin.Context) error {
	start := time.Now()
	end := time.Since(start)

	return api.Response{
		Code: http.StatusOK,
		Msg:  "Service is Up.",
		Data: gin.H{
			"response_time": end.String(),
		},
	}
}
