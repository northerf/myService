package handler

import (
	"awesomeProject1/pkg/auth"
	"awesomeProject1/users/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service      service.Auth
	tokenManager auth.TokenManager
}

func NewHandler(service service.Auth, tm *auth.TokenManager) *Handler {
	return &Handler{
		service:      service,
		tokenManager: *tm,
	}
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
