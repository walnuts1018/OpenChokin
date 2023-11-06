package domain

import "time"

// getMoneyPoolBalanceInternal is a helper function that constructs the SQL query for retrieving the money pool balance.
// It is used to avoid repetition in public methods.
func (d *dbImpl) getMoneyPoolBalanceInternal(moneyPoolID string, date *time.Time, includePlanned bool) (float64, error) {
	var balance float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM payment WHERE money_pool_id = $1 AND is_deleted = false`
	args := []interface{}{moneyPoolID}

	// Add date condition if it is provided
	if date != nil {
		query += ` AND date <= $2`
		args = append(args, *date)
	}

	// Exclude planned payments if not included
	if !includePlanned {
		query += ` AND is_planned = false`
	}

	err := d.db.Get(&balance, query, args...)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

// GetMoneyPoolBalance calculates the total amount of payments associated with the specified moneyPoolID.
// If includePlanned is true, it includes the planned payments in the calculation.
func (d *dbImpl) GetMoneyPoolBalance(moneyPoolID string, includePlanned bool) (float64, error) {
	return d.getMoneyPoolBalanceInternal(moneyPoolID, nil, includePlanned)
}

// GetMoneyPoolBalanceOfDate calculates the total amount of payments for a moneyPoolID up to a certain date.
// If includePlanned is true, it includes the planned payments in the calculation.
func (d *dbImpl) GetMoneyPoolBalanceOfDate(moneyPoolID string, date time.Time, includePlanned bool) (float64, error) {
	return d.getMoneyPoolBalanceInternal(moneyPoolID, &date, includePlanned)
}
