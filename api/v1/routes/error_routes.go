package routes

import (
	"net/http"

	"github.com/asparkoffire/whisper-speech-transcribe/api"
	"github.com/gin-gonic/gin"
)

func setupErrorRoutes(app *gin.Engine) {
	app.HandleMethodNotAllowed = true

	app.NoMethod(api.ErrorWrapper(func(c *gin.Context) error {
		return api.Response{
			Code: http.StatusMethodNotAllowed,
			Msg:  "Method Not Allowed",
		}
	}))

	app.NoRoute(api.ErrorWrapper(func(c *gin.Context) error {
		return api.Response{
			Code: http.StatusNotFound,
			Msg:  "Not Found",
		}
	}))
}
