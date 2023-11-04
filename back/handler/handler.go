package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/openchokin/back/usecase"
)

var (
	uc *usecase.Usecase
)

func NewHandler(usecase *usecase.Usecase) (*gin.Engine, error) {
	uc = usecase
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		// クエリパラメータtype=summary or detailでサマリーと詳細を分けられる。
		// 今回はsummaryだけを実装する
		// /moneypools?type=summary
		v1.GET("/moneypools", moneyPoolsHandler)

		// パスパラメータで指定されたIDのマネープール情報を返す
		v1.GET("/moneypools/:moneypool_id", moneyPoolHandler)

		// クエリパラメータtype=summary or detailでサマリと詳細を分けられる
		// 今回はsummaryだけを実装する
		// /moneyproviders?type=summary
		v1.GET("/moneyproviders", moneyProvidersHandler)

		// オプションパラメータdateを持つ
		// /moneyinformation?date=2023-05-15
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

// moneyPoolsHandler はモニープールのサマリ情報を取得するハンドラです。
func moneyPoolsHandler(c *gin.Context) {
	// クエリパラメータ 'type' を取得し、'summary' が指定されているかチェックします。
	queryType := c.DefaultQuery("type", "summary")

	// 今回は 'summary' のみを実装しているため、それ以外の値が来た場合はエラーを返します。
	if queryType != "summary" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
		return
	}

	// 認証ミドルウェアでuserIDを指定する
	userID := c.MustGet("userID").(string)

	// UsecaseのGetMoneyPoolsSummary関数を呼び出してサマリ情報を取得します。
	summaryResponse, err := uc.GetMoneyPoolsSummary(userID)
	if err != nil {
		// エラーが発生した場合は、クライアントにエラーメッセージとともにステータスコード500を返します。
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get money pools summary"})
		return
	}

	// サマリ情報をクライアントに返します。
	c.JSON(http.StatusOK, summaryResponse)
}

// MoneyPoolHandler handles GET requests for a specific money pool by ID.
func moneyPoolHandler(c *gin.Context) {
	userID := c.MustGet("userID").(string) // Assuming userID is set in some middleware
	moneyPoolID := c.Param("moneypool_id")

	response, err := uc.GetMoneyPool(userID, moneyPoolID)
	if err != nil {
		// Handle the error, e.g., by logging and returning an appropriate HTTP status code.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

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
	userID := c.MustGet("userID").(string)

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
	// ユーザーIDをコンテキストから取得
	userID := c.MustGet("userID").(string)

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
		response, err = uc.GetMoneyInformationOfDate(userID, date)
	} else {
		// 日付が指定されていない場合は現在の情報を計算
		response, err = uc.GetMoneyInformation(userID)
	}

	// エラーハンドリング
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusOK, response)
}
