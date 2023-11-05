package domain

import "fmt"

func (d *dbImpl) NewPayment(payment Payment) (Payment, error) {
	// クエリ文字列で位置パラメータを使用します。$1、$2...はそれぞれの値のプレースホルダーです。
	query := `INSERT INTO payment (money_pool_id, date, title, amount, description, is_planned, store_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  RETURNING id`
	// QueryRowを使用してSQLクエリを実行し、戻り値のIDを取得します。
	err := d.db.QueryRow(query, payment.MoneyPoolID, payment.Date, payment.Title, payment.Amount, payment.Description, payment.IsPlanned, payment.StoreID).Scan(&payment.ID)
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
	query := `SELECT id, money_pool_id, date, title, amount, description, is_planned, store_id FROM payment WHERE money_pool_id = $1 ORDER BY date DESC`
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

func (d *dbImpl) DeletePayment(id string) error {
	// DELETE SQL文を実行します。
	query := `DELETE FROM payment WHERE id = $1`
	result, err := d.db.Exec(query, id)
	if err != nil {
		// SQL実行エラーを返します。
		return fmt.Errorf("error deleting payment with id %s: %v", id, err)
	}

	// 影響を受けた行の数を確認します。
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// 影響を受けた行数の確認エラーを返します。
		return fmt.Errorf("error getting rows affected during deletion of payment with id %s: %v", id, err)
	}

	if rowsAffected == 0 {
		// 削除する行がなかった場合、エラーを返します。
		return fmt.Errorf("no payment found with id %s", id)
	}

	// 削除が成功した場合、nilを返します。
	return nil
}
