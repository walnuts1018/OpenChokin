package usecase

import (
	"errors"
	"fmt"
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
	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		return MoneyPoolsSummaryResponse{}, err
	}

	var pools []MoneyPoolSummary
	for _, pool := range moneyPools {
		// userIDとloginUserIDが一致、またはmoneyPoolがpublicであるかを確認
		if userID == loginUserID || pool.Type == "public" {
			sum, balanceErr := u.db.GetMoneyPoolBalance(pool.ID, false)
			if balanceErr != nil {
				return MoneyPoolsSummaryResponse{}, balanceErr // エラーを返す
			}
			pools = append(pools, MoneyPoolSummary{
				ID:   pool.ID,
				Name: pool.Name,
				Sum:  sum,
				Type: pool.Type,
			})
		} else if loginUserID != "" {
			// userIDとloginUserIDが異なる場合、共有状態を確認
			shared, shareErr := u.db.IsMoneyPoolSharedWithUser(pool.ID, loginUserID)
			if shareErr != nil {
				return MoneyPoolsSummaryResponse{}, shareErr // エラーを返す
			}
			if shared {
				sum, balanceErr := u.db.GetMoneyPoolBalance(pool.ID, false)
				if balanceErr != nil {
					return MoneyPoolsSummaryResponse{}, balanceErr // エラーを返す
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
	// Fetch the money pool by ID
	moneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		return MoneyPoolResponse{}, err
	}

	// Check if the login user is the owner or has access to the pool
	var hasAccess bool
	if userID == loginUserID {
		hasAccess = true
	} else if loginUserID != "" {
		// Check if the money pool is shared with the login user or if it's public
		shared, err := u.db.IsMoneyPoolSharedWithUser(moneyPoolID, loginUserID)
		if err != nil {
			return MoneyPoolResponse{}, err
		}
		hasAccess = shared || moneyPool.Type == domain.PublicTypePublic
	} else {
		// If loginUserID is empty, access is granted only if the pool is public
		hasAccess = moneyPool.Type == domain.PublicTypePublic
	}

	// If the user has no access, return an error
	if !hasAccess {
		return MoneyPoolResponse{}, fmt.Errorf("unauthorized access: user %s does not have access to the money pool %s", userID, moneyPoolID)
	}

	// Fetch payments associated with the money pool
	payments, err := u.db.GetPaymentsByMoneyPoolID(moneyPoolID)
	if err != nil {
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

	// Return the money pool response
	return MoneyPoolResponse{
		ID:          moneyPool.ID,
		Name:        moneyPool.Name,
		Description: moneyPool.Description,
		Type:        moneyPool.Type,
		Payments:    paymentSummaries,
	}, nil
}

func (u Usecase) AddMoneyPool(userID string, name string, description string, publicType domain.PublicType) (MoneyPoolResponse, error) {
	newMoneyPool := domain.MoneyPool{
		Name:        name,
		Description: description,
		Type:        publicType,
		OwnerID:     userID,
	}

	createdMoneyPool, err := u.db.NewMoneyPool(newMoneyPool)
	if err != nil {
		return MoneyPoolResponse{}, err
	}

	return MoneyPoolResponse{
		ID:          createdMoneyPool.ID,
		Name:        createdMoneyPool.Name,
		Description: createdMoneyPool.Description,
		Type:        createdMoneyPool.Type,
		Payments:    []PaymentSummary{}, // No payments right after creation
	}, nil
}

func (u Usecase) UpdateMoneyPool(userID string, moneyPoolID string, name string, description string, publicationType domain.PublicType) (MoneyPoolResponse, error) {
	// Fetch the existing money pool to check ownership
	existingMoneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		return MoneyPoolResponse{}, err
	}

	if existingMoneyPool.OwnerID != userID {
		return MoneyPoolResponse{}, errors.New("you are not authorized to update this money pool")
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
		return MoneyPoolResponse{}, err
	}

	return MoneyPoolResponse{
		ID:          updatedMoneyPool.ID,
		Name:        updatedMoneyPool.Name,
		Description: updatedMoneyPool.Description,
		Type:        updatedMoneyPool.Type,
	}, nil
}

func (u Usecase) DeleteMoneyPool(userID string, moneyPoolID string) error {
	moneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		return err
	}

	if moneyPool.OwnerID != userID {
		return errors.New("you are not authorized to delete this money pool")
	}

	return u.db.DeleteMoneyPool(moneyPoolID)
}

func (u *Usecase) ChangePublicationScope(userID string, moneyPoolID string, userGroupIDs []string) error {
	// Retrieve the MoneyPool by its ID to check its publication type.
	moneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		// Return error if the MoneyPool cannot be retrieved.
		return err
	}

	// Check if the owner of the MoneyPool is the user making the request.
	if moneyPool.OwnerID != userID {
		// Return an error if the user is not the owner.
		return fmt.Errorf("user %s is not the owner of MoneyPool %s", userID, moneyPoolID)
	}

	// Check if the MoneyPool's publication type is restricted.
	if moneyPool.Type != domain.PublicTypeRestricted {
		// Return an error if the publication type is not restricted.
		return fmt.Errorf("moneyPool %s publication type is not restricted", moneyPoolID)
	}

	// If the publication type is restricted, share the MoneyPool with user groups.
	err = u.db.ShareMoneyPoolWithUserGroups(moneyPoolID, userGroupIDs)
	if err != nil {
		// Return error if sharing fails.
		return err
	}

	// Return nil if sharing is successful.
	return nil
}
