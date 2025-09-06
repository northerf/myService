package handler

import (
	"awesomeProject1/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input schema.SignUpInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid body request", err)
		return
	}
	id, err := h.service.CreateUser(c.Request.Context(), input.Name, input.Password, input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input schema.SignInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid body request", err)
		return
	}

	token, err := h.service.GenerateToken(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to create token", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
