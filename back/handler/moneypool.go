package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/openchokin/back/domain"
)

// getMoneyPools APIのコメント
// @Summary マネープールの要約情報を取得
// @Description ユーザーIDに基づいたマネープールの要約情報を取得します。クエリパラメータとしてtypeとuser_idを受け取ります。
// @Tags moneypools
// @Accept  json
// @Produce  json
// @Param   type query string false "リクエストタイプ (summary または detail)" Enums(summary, detail) default(summary)
// @Param   user_id query string true "ユーザーID"
// @Success 200 {object} MoneyPoolsSummaryResponse "成功したレスポンス"
// @Failure 400 {object} map[string]string "不正なリクエストパラメータ"
// @Failure 500 {object} map[string]string "サーバ内部エラー"
// @Router /v1/moneypools [get]
func getMoneyPools(c *gin.Context) {
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

// getMoneyPool APIのコメント
// @Summary 特定のマネープールの情報を取得
// @Description ユーザーIDをクエリパラメータとして受け取り、指定されたマネープールIDの情報を返す。
// @Tags moneypools
// @Accept  json
// @Produce  json
// @Param   user_id       query    string  true  "ユーザーID"
// @Param   moneypool_id  path     string  true  "マネープールID"
// @Success 200 {object}  MoneyPoolResponse "成功時にマネープール情報を返す"
// @Failure 400 {object}  map[string]string      "ユーザーIDが不正である場合のエラーメッセージを返す"
// @Failure 500 {object}  map[string]string      "サーバー内部エラーが発生した場合のエラーメッセージを返す"
// @Router /v1/moneypools/{moneypool_id} [get]
func getMoneyPool(c *gin.Context) {
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

// createMoneyPool は新しいマネープールを作成します。
// @Summary 新しいマネープールを作成
// @Description 認証済みユーザーのための新しいマネープールを作成します。
// @Tags moneypools
// @Accept  json
// @Produce  json
// @Param loginUserID header string true "ログインユーザーID"
// @Param body body struct{Name string `json:"name"; Description string `json:"description"; Type domain.PublicType `json:"type"`} true "マネープール情報"
// @Success 200 {object} MoneyPoolResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /moneypools [post]
func createMoneyPool(c *gin.Context) {
	userID := c.MustGet("loginUserID").(string) // Get the authenticated user's ID
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Emoji       string `json:"emoji"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.Type != domain.PublicTypePrivate && request.Type != domain.PublicTypePublic && request.Type != domain.PublicTypeRestricted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request type does not match any options"})
		return
	}
	response, err := uc.AddMoneyPool(userID, request.Name, request.Description, request.Type, request.Emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// updateMoneyPool は指定されたIDのマネープールを更新します。
// @Summary 指定されたマネープールを更新
// @Description 認証済みユーザーが指定したIDのマネープールの情報を更新します。
// @Tags moneypools
// @Accept  json
// @Produce  json
// @Param loginUserID header string true "ログインユーザーID"
// @Param moneypool_id path string true "マネープールID"
// @Param body body struct{Name string `json:"name"; Description string `json:"description"; Type domain.PublicType `json:"type"`} true "更新するマネープール情報"
// @Success 200 {object} MoneyPoolResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /moneypools/{moneypool_id} [patch]
func updateMoneyPool(c *gin.Context) {
	userID := c.MustGet("loginUserID").(string) // Get the authenticated user's ID
	moneyPoolID := c.Param("moneypool_id")
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Emoji       string `json:"emoji"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.Type != domain.PublicTypePrivate && request.Type != domain.PublicTypePublic && request.Type != domain.PublicTypeRestricted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request type does not match any options"})
		return
	}
	response, err := uc.UpdateMoneyPool(userID, moneyPoolID, request.Name, request.Description, request.Type, request.Emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// deleteMoneyPool deletes a money pool.
// @Summary マネープールの削除
// @Description 認証されたユーザーのマネープールを削除します。
// @Tags moneypools
// @Accept  json
// @Produce  json
// @Param   moneypool_id   path      string  true  "マネープールID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error: Execution failure"
// @Router /v1/moneypools/{moneypool_id} [delete]
func deleteMoneyPool(c *gin.Context) {
	userID := c.MustGet("loginUserID").(string) // Get the authenticated user's ID
	moneyPoolID := c.Param("moneypool_id")
	err := uc.DeleteMoneyPool(userID, moneyPoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// changePublicationScope changes the publication scope of a money pool.
// @Summary マネープールの公開範囲変更
// @Description 認証されたユーザーが指定したマネープールの公開範囲を変更します。
// @Tags moneypools
// @Accept  json
// @Produce  json
// @Param   moneypool_id   path      string  true  "マネープールID"
// @Param   user_group_ids body      []string true  "ユーザーグループIDの配列"
// @Success 200 {string} string "OK"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error: Execution failure"
// @Router /v1/moneypools/{moneypool_id}/publicationscope [post]
func changePublicationScope(c *gin.Context) {
	userID := c.MustGet("loginUserID").(string) // Get the authenticated user's ID
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
