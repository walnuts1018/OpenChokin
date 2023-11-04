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

type PaymentResponse struct {
	ID          string
	MoneyPoolID string
	Date        time.Time
	Title       string
	Amount      float64
	Description string
	IsPlanned   bool
}

func (u *Usecase) UpdatePayment(userID string, moneyPoolID string, paymentID string, date time.Time, title string, amount float64, description string, isPlanned bool) (PaymentResponse, error) {
	// Get the payment details from the DB.
	payment, err := u.db.GetPayment(paymentID)
	if err != nil {
		return PaymentResponse{}, err
	}

	// Get the associated MoneyPool to check if the user is the owner.
	moneyPool, err := u.db.GetMoneyPool(payment.MoneyPoolID)
	if err != nil {
		return PaymentResponse{}, err
	}

	// Check if the user is the owner of the MoneyPool.
	if moneyPool.OwnerID != userID || moneyPool.ID != moneyPoolID {
		return PaymentResponse{}, fmt.Errorf("unauthorized: user %s is not the owner of the MoneyPool %s", userID, payment.MoneyPoolID)
	}

	// Update the payment details.
	payment.Date = date
	payment.Title = title
	payment.Amount = amount
	payment.Description = description
	payment.IsPlanned = isPlanned

	// Persist the updated payment in the DB.
	err = u.db.UpdatePayment(payment)
	if err != nil {
		return PaymentResponse{}, err
	}

	// Return the updated payment as a response.
	return PaymentResponse{
		ID:          payment.ID,
		MoneyPoolID: payment.MoneyPoolID,
		Date:        payment.Date,
		Title:       payment.Title,
		Amount:      payment.Amount,
		Description: payment.Description,
		IsPlanned:   payment.IsPlanned,
	}, nil
}

func (u *Usecase) DeletePayment(userID string, paymentID string) error {
	// Get the payment to check ownership.
	payment, err := u.db.GetPayment(paymentID)
	if err != nil {
		return err
	}

	// Get the associated MoneyPool to check if the user is the owner.
	moneyPool, err := u.db.GetMoneyPool(payment.MoneyPoolID)
	if err != nil {
		return err
	}

	// Check if the user is the owner of the MoneyPool.
	if moneyPool.OwnerID != userID {
		return fmt.Errorf("unauthorized: user %s is not the owner of the MoneyPool %s", userID, payment.MoneyPoolID)
	}

	// Use the DB interface method to delete the payment.
	err = u.db.DeletePayment(paymentID)
	if err != nil {
		return err
	}

	return nil
}
