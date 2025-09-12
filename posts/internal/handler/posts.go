package handler

import (
	"awesomeProject1/pkg/httputils"
	"awesomeProject1/posts/schema"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *Handler) createPost(c *gin.Context) {
	var input schema.CreatePostInput
	errResponse := httputils.NewErrorResponse

	fmt.Printf("Received POST request to create post\n")

	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("Failed to bind JSON: %v\n", err)
		errResponse(c, http.StatusBadRequest, "Invalid body request", err)
		return
	}

	fmt.Printf("Creating post: userID=%d, content=%s\n", input.UserID, input.Content)

	id, err := h.Service.CreatePost(c.Request.Context(), input.UserID, input.Content)
	if err != nil {
		fmt.Printf("Failed to create post: %v\n", err)
		errResponse(c, http.StatusInternalServerError, "Failed to create post", err)
		return
	}

	fmt.Printf("Post created successfully with ID: %d\n", id)

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
