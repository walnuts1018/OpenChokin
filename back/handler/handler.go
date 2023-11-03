package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/openchokin/back/domain"
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

		// リクエストボディの構造体を適切に定義してください
		v1.POST("/moneypools/:moneypool_id/payments", postPayment)

		// クエリパラメータmonthが必須パラメータである
		// /payments?month=2023-05
		v1.GET("/payments", getMonthlyPayments)

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

func AddNewUser(c *gin.Context) {
	user, err := uc.NewUser()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"user": user})
}

func GetUser(c *gin.Context) {
	userID := c.Param("userid")
	user, err := uc.GetUser(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("userid")
	user := domain.User{ID: userID}
	err := uc.UpdateUser(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"user": user})
}
