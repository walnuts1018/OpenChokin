package usecase

import (
	"fmt"
	"log"
	"time"

	"github.com/walnuts1018/openchokin/back/domain"
)

// AddNewPayment adds a new payment to the specified MoneyPool for a given user.
func (u *Usecase) AddNewPayment(userID string, moneyPoolID string, title string, amount float64, description string, isPlanned bool) error {
	log.Printf("ユーザーID %s のための新規支払い追加を開始します。マネープールID: %s, タイトル: %s", userID, moneyPoolID, title)
	// Retrieve the MoneyPool to ensure it exists and belongs to the user
	moneyPool, err := u.db.GetMoneyPool(moneyPoolID)
	if err != nil {
		log.Printf("マネープールID %s の取得に失敗しました。エラー: %v", moneyPoolID, err)
		return err // MoneyPool retrieval failed
	}
	if moneyPool.OwnerID != userID {
		log.Printf("エラー: ユーザーID %s はマネープールID %s の支払い追加に対して権限がありません。", userID, moneyPoolID)
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
	if err != nil {
		log.Printf("新規支払いの保存に失敗しました。マネープールID: %s, タイトル: %s, エラー: %v", moneyPoolID, title, err)
		return err
	}
	log.Printf("新規支払いを保存しました。マネープールID: %s, タイトル: %s", moneyPoolID, title)
	return nil // Will be nil if the operation was successful
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

// GetMonthlyPayments retrieves payments for a given user and month.
func (u *Usecase) GetMonthlyPayments(userID string, month time.Time) (MonthlyPaymentsResponse, error) {
	log.Printf("ユーザーID %s の月間支払い情報取得を開始します。対象月: %s", userID, month.Format("2006-01"))
	response := MonthlyPaymentsResponse{
		DailyPayments: make(map[int]DailyPayments),
	}

	daysInMonth := time.Date(month.Year(), month.Month()+1, 0, 0, 0, 0, 0, month.Location()).Day()
	for day := 1; day <= daysInMonth; day++ {
		response.DailyPayments[day] = DailyPayments{Payments: []DailyPaymentItem{}}
	}

	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		log.Printf("ユーザーID %s のマネープール取得に失敗しました。エラー: %v", userID, err)
		return MonthlyPaymentsResponse{}, err
	}

	for _, pool := range moneyPools {
		payments, err := u.db.GetPaymentsByMoneyPoolID(pool.ID)
		if err != nil {
			log.Printf("マネープールID %s の支払い情報取得に失敗しました。エラー: %v", pool.ID, err)
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

	log.Printf("ユーザーID %s の月間支払い情報を取得しました。対象月: %s", userID, month.Format("2006-01"))
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

// UpdatePayment updates a payment's details.
func (u *Usecase) UpdatePayment(userID string, moneyPoolID string, paymentID string, date time.Time, title string, amount float64, description string, isPlanned bool) (PaymentResponse, error) {
	log.Printf("支払いID %s の更新処理を開始します。ユーザーID: %s", paymentID, userID)

	// Get the payment details from the DB.
	payment, err := u.db.GetPayment(paymentID)
	if err != nil {
		log.Printf("支払いの詳細取得に失敗しました。支払いID: %s, エラー: %v", paymentID, err)
		return PaymentResponse{}, err
	}

	// Get the associated MoneyPool to check if the user is the owner.
	moneyPool, err := u.db.GetMoneyPool(payment.MoneyPoolID)
	if err != nil {
		log.Printf("マネープールの詳細取得に失敗しました。支払いID: %s, エラー: %v", paymentID, err)
		return PaymentResponse{}, err
	}

	// Check if the user is the owner of the MoneyPool.
	if moneyPool.OwnerID != userID || moneyPool.ID != moneyPoolID {
		log.Printf("不正アクセス：ユーザーID %s はマネープールID %s の所有者ではありません。", userID, moneyPoolID)
		return PaymentResponse{}, fmt.Errorf("unauthorized: user %s is not the owner of the MoneyPool %s", userID, moneyPoolID)
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
		log.Printf("支払いの更新に失敗しました。支払いID: %s, エラー: %v", paymentID, err)
		return PaymentResponse{}, err
	}

	log.Printf("支払いID %s の更新が完了しました。ユーザーID: %s", paymentID, userID)
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

// DeletePayment deletes a payment.
func (u *Usecase) DeletePayment(userID string, paymentID string) error {
	log.Printf("支払いID %s の削除処理を開始します。ユーザーID: %s", paymentID, userID)

	// Get the payment to check ownership.
	payment, err := u.db.GetPayment(paymentID)
	if err != nil {
		log.Printf("支払いの詳細取得に失敗しました。支払いID: %s, エラー: %v", paymentID, err)
		return err
	}

	// Get the associated MoneyPool to check if the user is the owner.
	moneyPool, err := u.db.GetMoneyPool(payment.MoneyPoolID)
	if err != nil {
		log.Printf("マネープールの詳細取得に失敗しました。支払いID: %s, エラー: %v", paymentID, err)
		return err
	}

	// Check if the user is the owner of the MoneyPool.
	if moneyPool.OwnerID != userID {
		log.Printf("不正アクセス：ユーザーID %s はマネープールID %s の所有者ではありません。", userID, payment.MoneyPoolID)
		return fmt.Errorf("unauthorized: user %s is not the owner of the MoneyPool %s", userID, payment.MoneyPoolID)
	}

	// Use the DB interface method to delete the payment.
	err = u.db.DeletePayment(paymentID)
	if err != nil {
		log.Printf("支払いの削除に失敗しました。支払いID: %s, エラー: %v", paymentID, err)
		return err
	}

	log.Printf("支払いID %s の削除が完了しました。ユーザーID: %s", paymentID, userID)
	return nil
}
