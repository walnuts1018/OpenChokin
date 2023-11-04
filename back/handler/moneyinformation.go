package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/openchokin/back/usecase"
)

// getMoneyInformation は、/moneyinformation エンドポイントのリクエストを処理するハンドラです。
func getMoneyInformation(c *gin.Context) {
	queryUserID := c.Query("user_id")
	if queryUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id parameter"})
		return
	}

	// ユーザーIDをコンテキストから取得
	loginUserID := "" // Default to empty string to indicate no user is logged in.

	// Check if userID exists in the context, indicating a logged-in state.
	userID, exists := c.Get("loginUserID")
	if exists {
		// Type assert to string to make sure we have the correct format for userID.
		var ok bool
		loginUserID, ok = userID.(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not a string"})
			return
		}
	}

	// オプショナルなクエリパラメータ 'date' を解析
	dateParam := c.DefaultQuery("date", "")
	var response usecase.MoneySumResponse
	var err error

	// 日付が指定されている場合はその日付で計算
	if dateParam != "" {
		var date time.Time
		date, err = time.Parse("2006-01-02", dateParam)
		if err != nil {
			// 日付のフォーマットが不正な場合はエラーレスポンスを返す
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		response, err = uc.GetMoneyInformationOfDate(queryUserID, loginUserID, date)
	} else {
		// 日付が指定されていない場合は現在の情報を計算
		response, err = uc.GetMoneyInformation(queryUserID, loginUserID)
	}

	// エラーハンドリング
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusOK, response)
}
