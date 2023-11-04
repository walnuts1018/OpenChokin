package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler function for creating a new MoneyProvider.
func createMoneyProviderHandler(c *gin.Context) {
	var req struct {
		Name    string  `json:"name"`
		Balance float64 `json:"balance"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string) // Assuming authentication middleware sets this.
	response, err := uc.AddMoneyProvider(userID, req.Name, req.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Handler function for updating an existing MoneyProvider.
func updateMoneyProviderHandler(c *gin.Context) {
	var req struct {
		Name    string  `json:"name"`
		Balance float64 `json:"balance"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)
	moneyProviderID := c.Param("moneyprovider_id")

	response, err := uc.UpdateMoneyProvider(userID, moneyProviderID, req.Name, req.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Handler function for deleting a MoneyProvider.
func deleteMoneyProviderHandler(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	moneyProviderID := c.Param("moneyprovider_id")

	if err := uc.DeleteMoneyProvider(userID, moneyProviderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
