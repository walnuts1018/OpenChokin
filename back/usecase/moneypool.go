package usecase

import (
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

func (u Usecase) GetMoneyPoolsSummary(userID string) (MoneyPoolsSummaryResponse, error) {
	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		return MoneyPoolsSummaryResponse{}, err
	}

	var pools []MoneyPoolSummary
	for _, pool := range moneyPools {
		sum, err := u.db.GetMoneyPoolBalance(pool.ID, false)
		if err != nil {
			return MoneyPoolsSummaryResponse{}, err
		}

		pools = append(pools, MoneyPoolSummary{
			ID:   pool.ID,
			Name: pool.Name,
			Sum:  sum,
			Type: pool.Type,
		})
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

func (u Usecase) GetMoneyPool(userID string, moneyPoolID string) (MoneyPoolResponse, error) {
	moneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		return MoneyPoolResponse{}, err
	}

	payments, err := u.db.GetPaymentsByMoneyPoolID(moneyPoolID)
	if err != nil {
		return MoneyPoolResponse{}, err
	}

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

	return MoneyPoolResponse{
		ID:          moneyPool.ID,
		Name:        moneyPool.Name,
		Description: moneyPool.Description,
		Type:        moneyPool.Type,
		Payments:    paymentSummaries,
	}, nil
}
