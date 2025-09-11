package handler

import (
	"awesomeProject1/pkg/auth"
	"awesomeProject1/pkg/httputils"
	"awesomeProject1/purchases/schema"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) (int64, error) {
	idInterface, exists := c.Get("userId")
	if !exists {
		return 0, errors.New("user id not found in context")
	}

	id, ok := idInterface.(int64)
	if !ok {
		return 0, errors.New("user id is of invalid type in context")
	}

	return id, nil
}

func (h *Handler) createPurchase(c *gin.Context) {
	var purchase schema.CreatePuchaseInput
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := c.BindJSON(&purchase); err != nil {
		httputils.NewErrorResponse(c, http.StatusBadRequest, "Invalid body request", err)
		return
	}

	id, err := h.service.CreatePurchase(c.Request.Context(), userID, purchase.ItemName, purchase.Amount)
	if err != nil {
		httputils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to create purchase", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) listUserPurchases(c *gin.Context) {
	userID, ok := auth.UserIDFromContext(c)
	if !ok {
		httputils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", http.ErrNoCookie)
		return
	}

	list, err := h.service.ListPurchases(c.Request.Context(), userID)
	if err != nil {
		httputils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to get purchases", err)
		return
	}

	var responseData []schema.PurchaseOutput

	for _, purchase := range list {
		responseData = append(responseData, schema.PurchaseOutput{
			ID:          purchase.ID,
			ItemName:    purchase.ItemName,
			Amount:      purchase.Amount,
			PurchasedAt: purchase.PurchasedAt,
		})
	}

	if responseData == nil {
		responseData = []schema.PurchaseOutput{}
	}

	c.JSON(http.StatusOK, responseData)
}
