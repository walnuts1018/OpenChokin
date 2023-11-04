package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/openchokin/back/domain"
)

// Handler function for creating a new MoneyPool
func createMoneyPool(c *gin.Context) {
	userID := c.MustGet("userID").(string) // Get the authenticated user's ID
	var request struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Type        domain.PublicType `json:"type"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.Type != domain.PublicTypePrivate && request.Type != domain.PublicTypePublic && request.Type != domain.PublicTypeRestricted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request type does not match any options"})
		return
	}
	response, err := uc.AddMoneyPool(userID, request.Name, request.Description, request.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Handler function for updating an existing MoneyPool
func updateMoneyPool(c *gin.Context) {
	userID := c.MustGet("userID").(string) // Get the authenticated user's ID
	moneyPoolID := c.Param("moneypool_id")
	var request struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Type        domain.PublicType `json:"type"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.Type != domain.PublicTypePrivate && request.Type != domain.PublicTypePublic && request.Type != domain.PublicTypeRestricted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request type does not match any options"})
		return
	}
	response, err := uc.UpdateMoneyPool(userID, moneyPoolID, request.Name, request.Description, request.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Handler function for deleting an existing MoneyPool
func deleteMoneyPool(c *gin.Context) {
	userID := c.MustGet("userID").(string) // Get the authenticated user's ID
	moneyPoolID := c.Param("moneypool_id")
	err := uc.DeleteMoneyPool(userID, moneyPoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func changePublicationScope(c *gin.Context) {
	userID := c.MustGet("userID").(string) // Get the authenticated user's ID
	moneyPoolID := c.Param("moneypool_id")

	// リクエストボディから
	var request struct {
		UserGroupIDs []string `json:"user_group_ids"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := uc.ChangePublicationScope(userID, moneyPoolID, request.UserGroupIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
