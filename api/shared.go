package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func (r Response) Error() string {
	return r.Msg
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
	Data any    `json:"data,omitempty"`
}

func (e Error) Error() string {
	return e.Err
}

func ErrorWrapper(fn func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiError Error
		var apiResponse Response
		err := fn(c)
		if err != nil {
			switch {
			case errors.As(err, &apiError):
				c.AbortWithStatusJSON(apiError.Code, apiError)
				return
			case errors.As(err, &apiResponse):
				c.AbortWithStatusJSON(apiResponse.Code, apiResponse)
				return
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
					Code: http.StatusInternalServerError,
					Err:  "Internal Server Error",
				})
				return
			}
		}
	}
}
