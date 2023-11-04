package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MoneyProvidersHandler handles GET requests for a summary of money providers.
func getMoneyProviders(c *gin.Context) {
	// クエリパラメータ 'type' を取得し、'summary' が指定されているかチェックします。
	queryType := c.DefaultQuery("type", "summary")

	// 今回は 'summary' のみを実装しているため、それ以外の値が来た場合はエラーを返します。
	if queryType != "summary" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
		return
	}

	// 認証ミドルウェアでuserIDを指定する
	userID := c.MustGet("loginUserID").(string)

	response, err := uc.GetMoneyProvidersSummary(userID)
	if err != nil {
		// Handle the error, e.g., by logging and returning an appropriate HTTP status code.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

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

	userID := c.MustGet("loginUserID").(string) // Assuming authentication middleware sets this.
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

	userID := c.MustGet("loginUserID").(string)
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
	userID := c.MustGet("loginUserID").(string)
	moneyProviderID := c.Param("moneyprovider_id")

	if err := uc.DeleteMoneyProvider(userID, moneyProviderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
