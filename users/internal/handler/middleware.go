package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "No authorization header", errors.New("no authorization header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header", errors.New("invalid authorization header"))
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "Token is empty", errors.New("token is empty"))
	}

	userID, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Token is invalid", err)
	}
	c.Set(userCtx, userID)

	c.Next()
}

func (h *Handler) getUserID(c *gin.Context) (int64, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "No user id", errors.New("no user id"))
	}
	idInt, ok := id.(int64)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid user id", errors.New("invalid user id"))
	}

	return idInt, nil
}
