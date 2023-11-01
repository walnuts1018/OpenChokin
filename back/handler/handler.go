package handler

import (
	"fmt"

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
		v1.GET("/slack/profile", func(ctx *gin.Context) { fmt.Printf("%v", uc) })
	}
	return r, nil
}
