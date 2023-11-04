package handler

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/openchokin/back/config"
	"github.com/walnuts1018/openchokin/back/usecase"
)

var (
	uc *usecase.Usecase
)

func userMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("loginUserID", "1")
		c.Next()
	}
}

// ミドルウェア関数
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorizationヘッダーを取得する
		authHeader := c.GetHeader("Authorization")
		// "Bearer "で始まる場合、トークンを検証する
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// ここでトークン検証の処理を実行する
			// OIDCプロバイダーのURLとクライアント情報
			issuer := "https://auth.walnuts.dev"
			clientID := "238653199337193865@walnuts.dev"

			// OIDCプロバイダーの構成情報を取得する
			provider, err := oidc.NewProvider(context.Background(), issuer)
			if err != nil {
				// エラー処理はログ出力に留める
				log.Printf("failed to get provider: %v\n", err)
				c.Next()
				return
			}

			// 公開鍵セットを取得してトークンを検証する
			verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
			idToken, err := verifier.Verify(context.Background(), tokenString)
			if err != nil {
				// エラー処理はログ出力に留める
				log.Printf("failed to verify token: %v\n", err)
				c.Next()
				return
			}

			// IDトークンのクレームを取得するための構造体
			var claims struct {
				Sub string `json:"sub"` // "sub"はOIDCのユーザーIDクレーム
			}

			// クレームをデコードする
			if err := idToken.Claims(&claims); err != nil {
				log.Printf("failed to decode claims: %v\n", err)
				c.Next()
				return
			}

			// クレームの情報をコンテキストにセットする
			c.Set("loginUserID", claims.Sub)
		}

		// 次のハンドラーまたはミドルウェアを実行
		c.Next()
	}
}
func NewHandler(usecase *usecase.Usecase) (*gin.Engine, error) {
	uc = usecase
	r := gin.Default()
	if config.Config.ISDebugMode == "true" {
		r.Use(userMiddleware())
	} else {
		r.Use(authMiddleware())
	}

	v1 := r.Group("/v1")
	{
		// クエリパラメータtype=summary or detailでサマリーと詳細を分けられる。
		// 今回はsummaryだけを実装する
		// /moneypools?type=summary&user_id=204938384
		v1.GET("/moneypools", moneyPoolsHandler)

		// パスパラメータで指定されたIDのマネープール情報を返す
		// クエリパラメータuserIDが必要
		v1.GET("/moneypools/:moneypool_id", moneyPoolHandler)

		// クエリパラメータtype=summary or detailでサマリと詳細を分けられる
		// 今回はsummaryだけを実装する
		// /moneyproviders?type=summary
		v1.GET("/moneyproviders", moneyProvidersHandler)

		// オプションパラメータdateを持つ
		// /moneyinformation?date=2023-05-15
		// クエリパラメータuserIDが必要
		v1.GET("/moneyinformation", moneyInformationHandler)

		// クエリパラメータmonthが必須パラメータである
		// /payments?month=2023-05
		v1.GET("/payments", getMonthlyPayments)

		// Paymentの追加・修正・削除
		v1.POST("/moneypools/:moneypool_id/payments", postPayment)
		v1.PATCH("/moneypools/:moneypool_id/payments/:payment_id", updatePaymentHandler)
		v1.DELETE("/moneypools/:moneypool_id/payments/:payment_id", deletePaymentHandler)

		// MoneyProviderの追加・修正・削除
		v1.POST("/moneyproviders", createMoneyProviderHandler)
		v1.PATCH("/moneyproviders/:moneyprovider_id", updateMoneyProviderHandler)
		v1.DELETE("/moneyproviders/:moneyprovider_id", deleteMoneyProviderHandler)

		// MoneyPoolの追加・修正・削除
		v1.POST("/moneypools", createMoneyPool)
		v1.PATCH("/moneypools/:moneypool_id", updateMoneyPool)
		v1.DELETE("/moneypools/:moneypool_id", deleteMoneyPool)
		// 公開範囲の設定(対象となるマネープールに対して、リクエストのjsonで指定されたユーザーグループに対して)
		v1.POST("/moneypools/:moneypool_id/publicationscope", changePublicationScope)

		// ユーザーグループの編集
		// これだけで詳細情報を全部取得する
		v1.GET("/usergroups", getUserGroups)
		v1.POST("/usergroups", createUserGroup)
		v1.PATCH("/usergroups/:usergroup_id", updateUserGroup)
		v1.DELETE("/usergroups/:usergroup_id", deleteUserGroup)
	}
	return r, nil
}

func moneyPoolsHandler(c *gin.Context) {
	queryType := c.DefaultQuery("type", "summary")
	if queryType != "summary" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
		return
	}

	queryUserID := c.Query("user_id")
	if queryUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id parameter"})
		return
	}

	// Default loginUserID to an empty string to handle non-logged-in state.
	loginUserID := ""

	// Check if userID exists and overwrite loginUserID with the actual userID if it does.
	if userID, exists := c.Get("loginUserID"); exists {
		var ok bool
		loginUserID, ok = userID.(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is not a string"})
			return
		}
	}

	// Retrieve summary information using the userID and loginUserID.
	summaryResponse, err := uc.GetMoneyPoolsSummary(queryUserID, loginUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get money pools summary"})
		return
	}

	// Send the retrieved summary information in the response.
	c.JSON(http.StatusOK, summaryResponse)
}

func moneyPoolHandler(c *gin.Context) {
	queryUserID := c.Query("user_id")
	if queryUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id parameter"})
		return
	}

	moneyPoolID := c.Param("moneypool_id")
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

	// Call the use case with the userID and loginUserID to get the money pool.
	response, err := uc.GetMoneyPool(queryUserID, loginUserID, moneyPoolID)
	if err != nil {
		// Handle the error, e.g., by logging and returning an appropriate HTTP status code.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Return the response.
	c.JSON(http.StatusOK, response)
}

// MoneyProvidersHandler handles GET requests for a summary of money providers.
func moneyProvidersHandler(c *gin.Context) {
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

// moneyInformationHandler は、/moneyinformation エンドポイントのリクエストを処理するハンドラです。
func moneyInformationHandler(c *gin.Context) {
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
