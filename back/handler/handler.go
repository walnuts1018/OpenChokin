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
