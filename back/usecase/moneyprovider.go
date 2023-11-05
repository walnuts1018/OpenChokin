package usecase

import (
	"fmt"
	"log"

	"github.com/walnuts1018/openchokin/back/domain"
)

type MoneyProviderSummary struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
type MoneyProvidersSummaryResponse struct {
	Providers []MoneyProviderSummary `json:"provider"`
}

func (u Usecase) GetMoneyProvidersSummary(userID string) (MoneyProvidersSummaryResponse, error) {
	log.Printf("ユーザーID %s のMoneyProvidersの概要を取得を開始します。", userID)
	moneyProviders, err := u.db.GetMoneyProvidersByUserID(userID)
	if err != nil {
		log.Printf("ユーザーID %s のMoneyProvidersの取得中にエラーが発生しました: %v", userID, err)
		return MoneyProvidersSummaryResponse{}, err
	}

	var providersSummary []MoneyProviderSummary
	for _, provider := range moneyProviders {
		balance, err := u.db.GetMoneyPoolBalance(provider.ID, false)
		if err != nil {
			log.Printf("MoneyProvider ID %s の残高取得時にエラーが発生しました。エラーを記録し、次のプロバイダーに続けます: %v", provider.ID, err)
			continue
		}
		providersSummary = append(providersSummary, MoneyProviderSummary{
			ID:      provider.ID,
			Name:    provider.Name,
			Balance: balance,
		})
		log.Printf("MoneyProvider ID %s: 名前：%s, 残高：%f", provider.ID, provider.Name, balance)
	}

	log.Printf("ユーザーID %s のMoneyProvidersの概要取得が完了しました。", userID)
	return MoneyProvidersSummaryResponse{Providers: providersSummary}, nil
}

type MoneyProviderResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	CreatorID string  `json:"creator_id"`
	Balance   float64 `json:"balance"`
}

func (u Usecase) UpdateMoneyProvider(userID string, moneyProviderID string, name string, balance float64) (MoneyProviderResponse, error) {
	log.Printf("MoneyProvider ID %s の更新を開始します。ユーザーID: %s", moneyProviderID, userID)

	existingProvider, err := u.db.GetMoneyProvider(moneyProviderID)
	if err != nil {
		log.Printf("MoneyProvider ID %s のデータ取得中にエラーが発生しました。エラー: %v", moneyProviderID, err)
		return MoneyProviderResponse{}, err
	}

	if existingProvider.CreatorID != userID {
		log.Printf("ユーザーID %s はMoneyProvider ID %s の更新が許可されていません。", userID, moneyProviderID)
		return MoneyProviderResponse{}, fmt.Errorf("unauthorized to update money provider: %s", moneyProviderID)
	}

	updatedProvider := domain.MoneyProvider{
		ID:        moneyProviderID,
		Name:      name,
		CreatorID: userID,
		Balance:   balance,
	}

	err = u.db.UpdateMoneyProvider(updatedProvider)
	if err != nil {
		log.Printf("MoneyProvider ID %s の更新中にエラーが発生しました。エラー: %v", moneyProviderID, err)
		return MoneyProviderResponse{}, err
	}

	log.Printf("MoneyProvider ID %s の更新が完了しました。", moneyProviderID)
	return MoneyProviderResponse{
		ID:        updatedProvider.ID,
		Name:      updatedProvider.Name,
		CreatorID: updatedProvider.CreatorID,
		Balance:   updatedProvider.Balance,
	}, nil
}

func (u Usecase) AddMoneyProvider(userID string, name string, balance float64) (MoneyProviderResponse, error) {
	log.Printf("新しいMoneyProviderの追加を開始します。ユーザーID: %s", userID)

	newProvider := domain.MoneyProvider{
		Name:      name,
		CreatorID: userID,
		Balance:   balance,
	}

	createdProvider, err := u.db.NewMoneyProvider(newProvider)
	if err != nil {
		log.Printf("新しいMoneyProviderの作成中にエラーが発生しました。エラー: %v", err)
		return MoneyProviderResponse{}, err
	}

	log.Printf("新しいMoneyProviderが作成されました。ID: %s", createdProvider.ID)
	return MoneyProviderResponse{
		ID:        createdProvider.ID,
		Name:      createdProvider.Name,
		CreatorID: createdProvider.CreatorID,
		Balance:   createdProvider.Balance,
	}, nil
}

func (u Usecase) DeleteMoneyProvider(userID string, moneyProviderID string) error {
	log.Printf("MoneyProvider ID %s の削除を試みます。ユーザーID: %s", moneyProviderID, userID)

	provider, err := u.db.GetMoneyProvider(moneyProviderID)
	if err != nil {
		log.Printf("MoneyProvider ID %s のデータ取得中にエラーが発生しました。エラー: %v", moneyProviderID, err)
		return err
	}

	if provider.CreatorID != userID {
		log.Printf("ユーザーID %s はMoneyProvider ID %s の削除が許可されていません。", userID, moneyProviderID)
		return fmt.Errorf("unauthorized to delete money provider: %s", moneyProviderID)
	}

	err = u.db.DeleteMoneyProvider(moneyProviderID)
	if err != nil {
		log.Printf("MoneyProvider ID %s の削除中にエラーが発生しました。エラー: %v", moneyProviderID, err)
		return err
	}

	log.Printf("MoneyProvider ID %s の削除が完了しました。", moneyProviderID)
	return nil
}
