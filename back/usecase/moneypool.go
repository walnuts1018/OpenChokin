package usecase

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/walnuts1018/openchokin/back/domain"
)

type MoneyPoolSummary struct {
	ID   string
	Name string
	// このIDのMoneyPoolに紐づくPlanではない実際の支払いの総額
	Sum  float64
	Type domain.PublicType
}

type MoneyPoolsSummaryResponse struct {
	Pools []MoneyPoolSummary
}

// GetMoneyPoolsSummary メソッドは、指定されたuserIDのMoneyPoolsの要約を返します。
func (u *Usecase) GetMoneyPoolsSummary(userID string, loginUserID string) (MoneyPoolsSummaryResponse, error) {
	log.Printf("ユーザーのMoneyPoolsの概要取得開始: ユーザーID: %s, ログインユーザーID: %s", userID, loginUserID)
	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		log.Printf("ユーザーのMoneyPoolsの取得に失敗: ユーザーID: %s, エラー: %v", userID, err)
		return MoneyPoolsSummaryResponse{}, err
	}

	var pools []MoneyPoolSummary
	for _, pool := range moneyPools {
		if userID == loginUserID || pool.Type == "public" {
			sum, balanceErr := u.db.GetMoneyPoolBalance(pool.ID, false)
			if balanceErr != nil {
				log.Printf("MoneyPoolのバランス取得に失敗: Pool ID: %s, エラー: %v", pool.ID, balanceErr)
				return MoneyPoolsSummaryResponse{}, balanceErr
			}
			pools = append(pools, MoneyPoolSummary{
				ID:   pool.ID,
				Name: pool.Name,
				Sum:  sum,
				Type: pool.Type,
			})
		} else if loginUserID != "" {
			shared, shareErr := u.db.IsMoneyPoolSharedWithUser(pool.ID, loginUserID)
			if shareErr != nil {
				log.Printf("MoneyPoolの共有状態確認に失敗: Pool ID: %s, ログインユーザーID: %s, エラー: %v", pool.ID, loginUserID, shareErr)
				return MoneyPoolsSummaryResponse{}, shareErr
			}
			if shared {
				sum, balanceErr := u.db.GetMoneyPoolBalance(pool.ID, false)
				if balanceErr != nil {
					log.Printf("MoneyPoolのバランス取得に失敗: Pool ID: %s, エラー: %v", pool.ID, balanceErr)
					return MoneyPoolsSummaryResponse{}, balanceErr
				}
				pools = append(pools, MoneyPoolSummary{
					ID:   pool.ID,
					Name: pool.Name,
					Sum:  sum,
					Type: pool.Type,
				})
			}
		}
	}

	log.Printf("ユーザーのMoneyPoolsの概要取得完了: ユーザーID: %s", userID)
	return MoneyPoolsSummaryResponse{Pools: pools}, nil
}

type PaymentSummary struct {
	ID          string
	Date        time.Time
	Title       string
	Amount      float64
	Description string
	IsPlanned   bool
}
type MoneyPoolResponse struct {
	ID          string
	Name        string
	Description string
	Type        domain.PublicType
	Payments    []PaymentSummary
}

func (u Usecase) GetMoneyPool(userID string, loginUserID string, moneyPoolID string) (MoneyPoolResponse, error) {
	log.Printf("ユーザーID: %sのためのMoneyPoolID: %sの取得を試みます。", userID, moneyPoolID)

	// Fetch the money pool by ID
	moneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		log.Printf("MoneyPoolID: %sの取得に失敗しました。エラー: %v", moneyPoolID, err)
		return MoneyPoolResponse{}, err
	}

	// Check access rights
	hasAccess := false
	if userID == loginUserID {
		hasAccess = true
	} else if loginUserID != "" {
		shared, err := u.db.IsMoneyPoolSharedWithUser(moneyPoolID, loginUserID)
		if err != nil {
			log.Printf("MoneyPoolID: %sの共有状態の確認に失敗しました。エラー: %v", moneyPoolID, err)
			return MoneyPoolResponse{}, err
		}
		hasAccess = shared || moneyPool.Type == domain.PublicTypePublic
	} else {
		hasAccess = moneyPool.Type == domain.PublicTypePublic
	}

	if !hasAccess {
		log.Printf("ユーザーID: %sはMoneyPoolID: %sへのアクセス権がありません。", userID, moneyPoolID)
		return MoneyPoolResponse{}, fmt.Errorf("unauthorized access: user %s does not have access to the money pool %s", userID, moneyPoolID)
	}

	// Fetch payments associated with the money pool
	payments, err := u.db.GetPaymentsByMoneyPoolID(moneyPoolID)
	if err != nil {
		log.Printf("MoneyPoolID: %sに関連する支払いの取得に失敗しました。エラー: %v", moneyPoolID, err)
		return MoneyPoolResponse{}, err
	}

	// Map payments to payment summaries
	var paymentSummaries []PaymentSummary
	for _, payment := range payments {
		paymentSummaries = append(paymentSummaries, PaymentSummary{
			ID:          payment.ID,
			Date:        payment.Date,
			Title:       payment.Title,
			Amount:      payment.Amount,
			Description: payment.Description,
			IsPlanned:   payment.IsPlanned,
		})
	}

	log.Printf("MoneyPoolID: %sに関する情報を正常に取得しました。", moneyPoolID)
	return MoneyPoolResponse{
		ID:          moneyPool.ID,
		Name:        moneyPool.Name,
		Description: moneyPool.Description,
		Type:        moneyPool.Type,
		Payments:    paymentSummaries,
	}, nil
}

// AddMoneyPool adds a new money pool to the database and logs the process in Japanese.
func (u Usecase) AddMoneyPool(userID string, name string, description string, publicType domain.PublicType) (MoneyPoolResponse, error) {
	log.Printf("ユーザーID: %sによる新しいマネープールの作成を開始します。名前: %s", userID, name)

	newMoneyPool := domain.MoneyPool{
		Name:        name,
		Description: description,
		Type:        publicType,
		OwnerID:     userID,
	}

	createdMoneyPool, err := u.db.NewMoneyPool(newMoneyPool)
	if err != nil {
		log.Printf("マネープールの作成中にエラーが発生しました: %v", err)
		return MoneyPoolResponse{}, err
	}

	log.Printf("マネープールが正常に作成されました。ID: %s", createdMoneyPool.ID)
	return MoneyPoolResponse{
		ID:          createdMoneyPool.ID,
		Name:        createdMoneyPool.Name,
		Description: createdMoneyPool.Description,
		Type:        createdMoneyPool.Type,
		Payments:    []PaymentSummary{}, // No payments right after creation
	}, nil
}

// UpdateMoneyPool updates an existing money pool and logs the process in Japanese.
func (u Usecase) UpdateMoneyPool(userID string, moneyPoolID string, name string, description string, publicationType domain.PublicType) (MoneyPoolResponse, error) {
	log.Printf("ユーザーID: %sがマネープールID: %sを更新しようとしています。", userID, moneyPoolID)

	existingMoneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		log.Printf("マネープールID: %sの取得中にエラーが発生しました: %v", moneyPoolID, err)
		return MoneyPoolResponse{}, err
	}

	if existingMoneyPool.OwnerID != userID {
		log.Printf("ユーザーID: %sはマネープールID: %sを更新する権限がありません。", userID, moneyPoolID)
		return MoneyPoolResponse{}, errors.New("更新権限がありません")
	}

	updatedMoneyPool := domain.MoneyPool{
		ID:          moneyPoolID,
		Name:        name,
		Description: description,
		Type:        publicationType,
		OwnerID:     userID,
	}

	err = u.db.UpdateMoneyPool(updatedMoneyPool)
	if err != nil {
		log.Printf("マネープールID: %sの更新中にエラーが発生しました: %v", moneyPoolID, err)
		return MoneyPoolResponse{}, err
	}

	log.Printf("マネープールID: %sが正常に更新されました。", moneyPoolID)
	return MoneyPoolResponse{
		ID:          updatedMoneyPool.ID,
		Name:        updatedMoneyPool.Name,
		Description: updatedMoneyPool.Description,
		Type:        updatedMoneyPool.Type,
	}, nil
}

// DeleteMoneyPool deletes an existing money pool and logs the process in Japanese.
func (u Usecase) DeleteMoneyPool(userID string, moneyPoolID string) error {
	log.Printf("ユーザーID: %sがマネープールID: %sの削除を試みます。", userID, moneyPoolID)

	moneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		log.Printf("削除するマネープールID: %sの取得中にエラーが発生しました: %v", moneyPoolID, err)
		return err
	}

	if moneyPool.OwnerID != userID {
		log.Printf("ユーザーID: %sにはマネープールID: %sを削除する権限がありません。", userID, moneyPoolID)
		return errors.New("削除権限がありません")
	}

	err = u.db.DeleteMoneyPool(moneyPoolID)
	if err != nil {
		log.Printf("マネープールID: %sの削除中にエラーが発生しました: %v", moneyPoolID, err)
		return err
	}

	log.Printf("マネープールID: %sが正常に削除されました。", moneyPoolID)
	return nil
}

// ChangePublicationScope changes the scope of publication for a money pool and logs the process in Japanese.
func (u *Usecase) ChangePublicationScope(userID string, moneyPoolID string, userGroupIDs []string) error {
	log.Printf("ユーザーID: %sによるマネープールID: %sの公開範囲変更を試みます。", userID, moneyPoolID)

	// Retrieve the MoneyPool by its ID to check its publication type.
	moneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		log.Printf("マネープールID: %sの取得に失敗しました: %v", moneyPoolID, err)
		// Return error if the MoneyPool cannot be retrieved.
		return err
	}

	// Check if the owner of the MoneyPool is the user making the request.
	if moneyPool.OwnerID != userID {
		errMsg := fmt.Sprintf("ユーザーID: %sはマネープールID: %sの所有者ではありません。", userID, moneyPoolID)
		log.Println(errMsg)
		// Return an error if the user is not the owner.
		return errors.New(errMsg)
	}

	// Check if the MoneyPool's publication type is restricted.
	if moneyPool.Type != domain.PublicTypeRestricted {
		errMsg := fmt.Sprintf("マネープールID: %sの公開タイプは制限されていません。", moneyPoolID)
		log.Println(errMsg)
		// Return an error if the publication type is not restricted.
		return errors.New(errMsg)
	}

	// If the publication type is restricted, share the MoneyPool with user groups.
	err = u.db.ShareMoneyPoolWithUserGroups(moneyPoolID, userGroupIDs)
	if err != nil {
		log.Printf("ユーザーグループにマネープールID: %sの共有に失敗しました: %v", moneyPoolID, err)
		// Return error if sharing fails.
		return err
	}

	log.Printf("マネープールID: %sをユーザーグループに正常に共有しました。", moneyPoolID)
	// Return nil if sharing is successful.
	return nil
}
