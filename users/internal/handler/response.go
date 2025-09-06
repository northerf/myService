package handler

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

func newErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	// Теперь мы логируем и сообщение для контекста, и саму ошибку
	logrus.WithFields(logrus.Fields{
		"error": err.Error(),
	}).Error(message)

	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
