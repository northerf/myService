package handler

import (
	"awesomeProject1/pkg/httputils"
	"awesomeProject1/posts/schema"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) listRecent(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	posts, err := h.Service.ListRecent(c.Request.Context(), limit)
	if err != nil {
		httputils.NewErrorResponse(c, http.StatusBadRequest, "Invalid body request", err)
		return
	}

	var responseData []schema.PostsOutput

	for _, metaInfo := range posts {
		responseData = append(responseData, schema.PostsOutput{
			ID:        metaInfo.ID,
			UserID:    metaInfo.UserID,
			Content:   metaInfo.Content,
			CreatedAt: metaInfo.CreatedAt,
		})
	}
	if responseData == nil {
		responseData = []schema.PostsOutput{}
	}

	c.JSON(http.StatusOK, responseData)
}
