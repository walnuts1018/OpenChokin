package usecase

import (
	"log"
	"time"

	"github.com/walnuts1018/openchokin/back/domain"
)

type MoneySumResponse struct {
	MoneyProviderSum       float64
	ActualMoneyPoolSum     float64
	ForecastedMoneyPoolSum float64
}

// GetMoneyInformation retrieves the sum of money information for a user.
func (u Usecase) GetMoneyInformation(userID string, loginUserID string) (MoneySumResponse, error) {
	var response MoneySumResponse
	log.Printf("ユーザーID %s のマネー情報取得を開始します。ログインユーザーID: %s", userID, loginUserID)

	// Retrieve all MoneyPools associated with the user.
	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		log.Printf("ユーザーID %s のマネープール取得時にエラー: %v", userID, err)
		return response, err
	}
	log.Printf("ユーザーID %s に関連するマネープールを %d 件取得しました。", userID, len(moneyPools))

	// Process each MoneyPool based on access rights.
	for _, pool := range moneyPools {
		var hasAccess bool // Flag to check access permissions

		// Direct access for the same user or public access for others.
		if userID == loginUserID {
			hasAccess = true // Users have full access to their own pools
			log.Printf("ユーザーID %s はマネープールID %s に直接アクセスできます。", userID, pool.ID)
		} else if loginUserID != "" {
			// Check shared access or public type for logged-in users.
			shared, err := u.db.IsMoneyPoolSharedWithUser(pool.ID, loginUserID)
			if err != nil {
				log.Printf("マネープールID %s とユーザーID %s の共有状況確認エラー: %v", pool.ID, loginUserID, err)
				return response, err // Error checking shared status
			}
			hasAccess = shared || pool.Type == domain.PublicTypePublic
			log.Printf("ログインユーザーID %s はマネープールID %s へのアクセス権が %v です。", loginUserID, pool.ID, hasAccess)
		} else {
			// No login user ID provided; only include public pools.
			hasAccess = pool.Type == domain.PublicTypePublic
			if hasAccess {
				log.Printf("ログインしていないユーザーへのマネープールID %s の公開アクセスが許可されています。", pool.ID)
			}
		}

		// If the user has access, sum up the actual and forecasted balances.
		if hasAccess {
			balance, err := u.db.GetMoneyPoolBalance(pool.ID, false)
			if err != nil {
				log.Printf("マネープールID %s の実際の残高取得エラー: %v", pool.ID, err)
				return response, err
			}
			response.ActualMoneyPoolSum += balance

			forecastedBalance, err := u.db.GetMoneyPoolBalance(pool.ID, true)
			if err != nil {
				log.Printf("マネープールID %s の予測残高取得エラー: %v", pool.ID, err)
				return response, err
			}
			response.ForecastedMoneyPoolSum += forecastedBalance
			log.Printf("マネープールID %s の残高を追加: 実際の残高 %f, 予測残高 %f", pool.ID, balance, forecastedBalance)
		}
	}

	// Retrieve all MoneyProviders for the user and calculate the sum.
	moneyProviders, err := u.db.GetMoneyProvidersByUserID(userID)
	if err != nil {
		log.Printf("ユーザーID %s のマネープロバイダー取得エラー: %v", userID, err)
		return response, err
	}
	for _, provider := range moneyProviders {
		response.MoneyProviderSum += provider.Balance
	}
	log.Printf("ユーザーID %s のマネープロバイダー合計を計算: %f", userID, response.MoneyProviderSum)

	log.Printf("ユーザーID %s のマネー情報取得が完了しました。", userID)
	return response, nil
}

func (u Usecase) GetMoneyInformationOfDate(userID string, loginUserID string, date time.Time) (MoneySumResponse, error) {
	var response MoneySumResponse
	log.Printf("特定日の金銭情報取得を開始: ユーザーID: %s, ログインユーザーID: %s, 日付: %v", userID, loginUserID, date)

	// ユーザーに関連する全てのMoneyProvidersを取得し、そのバランスの合計を計算する。
	moneyProviders, err := u.db.GetMoneyProvidersByUserID(userID)
	if err != nil {
		log.Printf("MoneyProviders取得エラー: ユーザーID: %s, エラー: %v", userID, err)
		return response, err
	}
	for _, provider := range moneyProviders {
		response.MoneyProviderSum += provider.Balance
	}

	// アクセス権に基づいてMoneyPoolsを取得
	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		log.Printf("MoneyPools取得エラー: ユーザーID: %s, エラー: %v", userID, err)
		return response, err
	}
	for _, pool := range moneyPools {
		var hasAccess bool // アクセス権を確認するフラグ
		if userID == loginUserID {
			hasAccess = true // ユーザーは自分自身のプールに全アクセス権を持つ
		} else if loginUserID != "" {
			// MoneyPoolが共有または公開されているかチェック
			shared, err := u.db.IsMoneyPoolSharedWithUser(pool.ID, loginUserID)
			if err != nil {
				log.Printf("MoneyPool共有状態チェックエラー: プールID: %s, ログインユーザーID: %s, エラー: %v", pool.ID, loginUserID, err)
				return response, err
			}
			hasAccess = shared || pool.Type == domain.PublicTypePublic
		} else {
			// ログインユーザーIDが提供されていない場合、公開プールのみを含む
			hasAccess = pool.Type == domain.PublicTypePublic
		}

		if hasAccess {
			// 指定された日付までの実際のバランスを計算
			balance, err := u.db.GetMoneyPoolBalanceOfDate(pool.ID, date, false)
			if err != nil {
				log.Printf("実際のバランス計算エラー: プールID: %s, 日付: %v, エラー: %v", pool.ID, date, err)
				return response, err
			}
			response.ActualMoneyPoolSum += balance

			// 指定された日付までの予測バランスを計算
			balance, err = u.db.GetMoneyPoolBalanceOfDate(pool.ID, date, true)
			if err != nil {
				log.Printf("予測バランス計算エラー: プールID: %s, 日付: %v, エラー: %v", pool.ID, date, err)
				return response, err
			}
			response.ForecastedMoneyPoolSum += balance
		}
	}

	log.Printf("特定日の金銭情報取得完了: ユーザーID: %s, 実際の合計: %f, 予測合計: %f", userID, response.ActualMoneyPoolSum, response.ForecastedMoneyPoolSum)
	return response, nil
}
