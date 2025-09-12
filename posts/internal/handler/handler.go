package handler

import (
	"awesomeProject1/posts/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.PostsService
}

func NewHandler(service *service.PostsService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.GET("/posts", h.listRecent)
	}

	return router
}
