package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/walnuts1018/openchokin/back/config"
	"github.com/walnuts1018/openchokin/back/handler"
	"github.com/walnuts1018/openchokin/back/infra/psql"
	"github.com/walnuts1018/openchokin/back/usecase"
)

func main() {
	config.LoadConfig()

	db, err := psql.NewDB()
	if err != nil {
		slog.Error("failed to create db", "message", err)
		os.Exit(1)
	}
	defer db.Close()

	u := usecase.NewUsecase(db)

	h, err := handler.NewHandler(u)
	if err != nil {
		slog.Error("failed to create handler", "message", err)
		os.Exit(1)
	}

	err = h.Run(fmt.Sprintf(":%v", config.Config.ServerPort))
	if err != nil {
		slog.Error("failed to run handler", "message", err)
		os.Exit(1)
	}
}
