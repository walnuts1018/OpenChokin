package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// updatePaymentHandler handles the PATCH request for updating a payment
func updatePaymentHandler(c *gin.Context) {
	userID := c.MustGet("userID").(string) // Assuming userID retrieval from middleware
	moneyPoolID := c.Param("moneypool_id")
	paymentID := c.Param("payment_id")

	var req struct {
		Date        time.Time `json:"date"`
		Title       string    `json:"title"`
		Amount      float64   `json:"amount"`
		Description string    `json:"description"`
		IsPlanned   bool      `json:"is_planned"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentResponse, err := uc.UpdatePayment(userID, moneyPoolID, paymentID, req.Date, req.Title, req.Amount, req.Description, req.IsPlanned)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentResponse)
}

// deletePaymentHandler handles the DELETE request for a payment
func deletePaymentHandler(c *gin.Context) {
	userID := c.MustGet("userID").(string) // Assuming userID retrieval from middleware
	paymentID := c.Param("payment_id")

	err := uc.DeletePayment(userID, paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// POST /moneypools/:moneypool_id/payments
// 指定されたマネープールに新しい支払いを追加する
func postPayment(c *gin.Context) {
	userID := c.MustGet("userID").(string) // 認証ユーザーのIDを取得
	moneyPoolID := c.Param("moneypool_id") // パスパラメータからマネープールIDを取得

	// リクエストボディの構造体
	var paymentRequest struct {
		Title       string  `json:"title"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		IsPlanned   bool    `json:"is_planned"`
	}
	if err := c.BindJSON(&paymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := uc.AddNewPayment(userID, moneyPoolID, paymentRequest.Title, paymentRequest.Amount, paymentRequest.Description, paymentRequest.IsPlanned)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GET /payments
// 指定された月の支払い情報を取得する
func getMonthlyPayments(c *gin.Context) {
	userID := c.MustGet("userID").(string) // 認証ユーザーのIDを取得
	monthStr := c.Query("month")           // クエリパラメータから月を取得

	// "YYYY-MM"の形式であることを確認し、time.Time型にパースする
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month format"})
		return
	}

	response, err := uc.GetMonthlyPayments(userID, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
