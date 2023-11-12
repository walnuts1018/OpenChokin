package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/openchokin/back/config"
	"github.com/walnuts1018/openchokin/back/domain"
	"github.com/walnuts1018/openchokin/back/usecase"
)

var (
	uc *usecase.Usecase
)

func userMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("loginUserID", "1")
		if _, err := uc.GetUser("1"); err != nil {
			log.Printf("created new user %s\n", "1")
			uc.NewUser(domain.User{ID: "1"})
		}
		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorizationヘッダーを取得する
		authHeader := c.GetHeader("Authorization")
		// "Bearer "で始まる場合、トークンを検証する
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// OIDCプロバイダーのURLとクライアント情報
			issuer := "https://auth.walnuts.dev"
			clientID := "238653199337193865@walnuts.dev"

			// OIDCプロバイダーの構成情報を取得する
			provider, err := oidc.NewProvider(context.Background(), issuer)
			if err != nil {
				log.Printf("OIDCプロバイダーの取得に失敗しました: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "内部サーバーエラー"})
				return
			}

			// 公開鍵セットを取得してトークンを検証する
			verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
			idToken, err := verifier.Verify(context.Background(), tokenString)
			if err != nil {
				log.Printf("トークンの検証に失敗しました: %v", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証エラー"})
				return
			}

			// IDトークンのクレームを取得するための構造体
			var claims struct {
				Sub string `json:"sub"` // "sub"はOIDCのユーザーIDクレーム
			}

			// クレームをデコードする
			if err := idToken.Claims(&claims); err != nil {
				log.Printf("クレームのデコードに失敗しました: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "クレームデコードエラー"})
				return
			}

			// クレームの情報をコンテキストにセットする
			c.Set("loginUserID", claims.Sub)

			// ユーザーの存在確認と新規作成のロジックは省略...
			if _, err := uc.GetUser(claims.Sub); err != nil {
				uc.NewUser(domain.User{ID: claims.Sub})
			}

			log.Printf("ユーザー認証成功: ユーザーID %s", claims.Sub)
		} else {
			log.Printf("AuthorizationヘッダーがBearer形式ではありません。")
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
		v1.GET("/moneypools", getMoneyPools)

		// パスパラメータで指定されたIDのマネープール情報を返す
		// クエリパラメータuserIDが必要
		v1.GET("/moneypools/:moneypool_id", getMoneyPool)

		// クエリパラメータtype=summary or detailでサマリと詳細を分けられる
		// 今回はsummaryだけを実装する
		// /moneyproviders?type=summary
		v1.GET("/moneyproviders", getMoneyProviders)

		// オプションパラメータdateを持つ
		// /moneyinformation?date=2023-05-15
		// クエリパラメータuserIDが必要
		v1.GET("/moneyinformation", getMoneyInformation)

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
