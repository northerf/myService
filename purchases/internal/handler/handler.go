package handler

import (
	"awesomeProject1/pkg/auth"
	"awesomeProject1/purchases/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service      service.Purchase
	tokenManager auth.TokenManager
}

func NewHandler(service service.Purchase, tm *auth.TokenManager) *Handler {
	return &Handler{service: service, tokenManager: *tm}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api", h.tokenManager.Middleware())
	{
		api.POST("/purchases", h.createPurchase)
		api.GET("/purchases", h.listUserPurchases)
	}
	return router
}
