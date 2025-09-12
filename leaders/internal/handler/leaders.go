package handler

import (
	"awesomeProject1/leaders/schema"
	"awesomeProject1/pkg/httputils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getLeaders(c *gin.Context) {
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}

	leaders, err := h.Service.GetTopLeaders(c.Request.Context(), limit)
	if err != nil {
		httputils.NewErrorResponse(c, http.StatusInternalServerError, "failed to get leaders", err)
		return
	}

	var responseData []schema.LeaderOutput
	for _, leader := range leaders {
		responseData = append(responseData, schema.LeaderOutput{
			UserID:   leader.UserID,
			Score:    int(leader.Score),
			Rank:     leader.Rank,
			Username: leader.Username,
		})
	}

	if responseData == nil {
		responseData = []schema.LeaderOutput{}
	}

	c.JSON(http.StatusOK, responseData)
}
