package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseSourceModel struct {
	Data    any    `json:"data" default:"[]"`
	Message string `json:"message" default:""`
	Error   string `json:"error" default:""`
}

func Success(c *gin.Context, data any, message string, error string) {
	var response = ResponseSourceModel{
		Data:    data,
		Message: message,
		Error:   error,
	}

	c.JSON(http.StatusOK, gin.H{"result": response})

}

func Fail(c *gin.Context, data any, message string, error string) {
	var response = ResponseSourceModel{
		Data:    data,
		Message: message,
		Error:   error,
	}

	c.JSON(http.StatusInternalServerError, gin.H{"result": response})
}
