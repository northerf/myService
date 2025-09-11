package httputils

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	logrus.WithFields(logrus.Fields{
		"error": err.Error(),
	}).Error(message)

	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
