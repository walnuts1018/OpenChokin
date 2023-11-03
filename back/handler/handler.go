package handler

import (
	"fmt"

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
		v1.GET("/slack/profile", func(ctx *gin.Context) { fmt.Printf("%v", uc) })
		// クエリパラメータtype=summary or detailでサマリーと詳細を分けられる
		v1.GET("/moneypools")
		v1.GET("/moneypools/:moneypool_id")
		v1.GET("/moneyproviders")
		v1.GET("/moneyinformation")
		v1.POST("/moneypools/:moneypool_id/payments")
		// クエリパラメータで日付を指定する
		v1.GET("/payments")
		// 今回は実装しない
		// v1.POST("/stores", createStore)
		// v1.PATCH("/stores/:store_id", updateStore)
		// v1.POST("/items", createItem)
		// v1.PATCH("/items/:item_id", updateItem)
	}
	return r, nil
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
