package handler

import (
	"awesomeProject1/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Auth
}

func NewHandler(service service.Auth) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)

	}
	return router
}
