package usecase

import (
	"fmt"
	"time"

	"github.com/walnuts1018/openchokin/back/domain"
)

// AddNewPayment adds a new payment to the specified MoneyPool for a given user.
func (u *Usecase) AddNewPayment(userID string, moneyPoolID string, title string, amount float64, description string, isPlanned bool) error {
	// Retrieve the MoneyPool to ensure it exists and belongs to the user
	moneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		return err // MoneyPool retrieval failed
	}
	if moneyPool.OwnerID != userID {
		return fmt.Errorf("error: user unauthorized")
	}

	// Create the Payment entity
	payment := domain.Payment{
		ID:          "",
		MoneyPoolID: moneyPoolID,
		Date:        time.Now(), // Use current time for payment date
		Title:       title,
		Amount:      amount,
		Description: description,
		IsPlanned:   isPlanned,
		// StoreID can be nil or set based on additional logic if required
	}

	// Persist the new payment
	_, err = u.db.NewPayment(payment)
	return err // Will be nil if the operation was successful
}

type DailyPaymentItem struct {
	ID          string
	MoneyPoolID string
	Title       string
	Amount      float64
	IsPlanned   bool
}
type DailyPayments struct {
	Payments []DailyPaymentItem
}
type MonthlyPaymentsResponse struct {
	DailyPayments map[int]DailyPayments
}

func (u *Usecase) GetMonthlyPayments(userID string, month time.Time) (MonthlyPaymentsResponse, error) {
	response := MonthlyPaymentsResponse{
		DailyPayments: make(map[int]DailyPayments),
	}

	daysInMonth := time.Date(month.Year(), month.Month()+1, 0, 0, 0, 0, 0, month.Location()).Day()
	for day := 1; day <= daysInMonth; day++ {
		response.DailyPayments[day] = DailyPayments{Payments: []DailyPaymentItem{}}
	}

	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		return MonthlyPaymentsResponse{}, err
	}

	for _, pool := range moneyPools {
		payments, err := u.db.GetPaymentsByMoneyPoolID(pool.ID)
		if err != nil {
			return MonthlyPaymentsResponse{}, err
		}

		for _, payment := range payments {
			if payment.Date.Month() == month.Month() && payment.Date.Year() == month.Year() {
				day := payment.Date.Day()

				item := DailyPaymentItem{
					ID:          payment.ID,
					MoneyPoolID: payment.MoneyPoolID,
					Title:       payment.Title,
					Amount:      payment.Amount,
					IsPlanned:   payment.IsPlanned,
				}

				dailyPayments := response.DailyPayments[day]

				dailyPayments.Payments = append(dailyPayments.Payments, item)

				response.DailyPayments[day] = dailyPayments
			}
		}
	}

	return response, nil
}
