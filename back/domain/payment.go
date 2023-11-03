package domain

import "fmt"

func (d *dbImpl) NewPayment(payment Payment) (Payment, error) {
	query := `INSERT INTO payment (money_pool_id, date, title, amount, description, is_planned, store_id)
			  VALUES (:money_pool_id, :transaction_date, :title, :amount, :description, :is_expectation, :store_id)
			  RETURNING id`
	err := d.db.QueryRowx(query, payment).StructScan(&payment)
	if err != nil {
		return Payment{}, fmt.Errorf("failed to create new Payment: %v", err)
	}
	return payment, nil
}

// GetPayment retrieves a single payment by its ID.
func (d *dbImpl) GetPayment(id string) (Payment, error) {
	var payment Payment
	query := `SELECT id, money_pool_id, date, title, amount, description, is_planned, store_id FROM payment WHERE id = $1`
	err := d.db.Get(&payment, query, id)
	if err != nil {
		return Payment{}, fmt.Errorf("error fetching payment: %v", err)
	}
	return payment, nil
}

// GetPaymentsByMoneyPoolID retrieves all payments associated with a specific money pool.
func (d *dbImpl) GetPaymentsByMoneyPoolID(moneyPoolID string) ([]Payment, error) {
	var payments []Payment
	query := `SELECT id, money_pool_id, date, title, amount, description, is_planned, store_id FROM payment WHERE money_pool_id = $1`
	err := d.db.Select(&payments, query, moneyPoolID)
	if err != nil {
		return nil, fmt.Errorf("error fetching payments: %v", err)
	}
	return payments, nil
}

// UpdatePayment updates an existing payment's details.
func (d *dbImpl) UpdatePayment(payment Payment) error {
	query := `UPDATE payment SET money_pool_id = $1, date = $2, title = $3, amount = $4, description = $5, is_planned = $6, store_id = $7 WHERE id = $8`
	_, err := d.db.Exec(query, payment.MoneyPoolID, payment.Date, payment.Title, payment.Amount, payment.Description, payment.IsPlanned, payment.StoreID, payment.ID)
	if err != nil {
		return fmt.Errorf("error updating payment: %v", err)
	}
	return nil
}
