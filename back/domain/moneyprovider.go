package domain

import "github.com/pkg/errors"

func (d *dbImpl) NewMoneyProvider(moneyProvider MoneyProvider) (MoneyProvider, error) {
	query := `INSERT INTO money_providers (name, creator_id, balance)
			  VALUES (:name, :creator_id, :balance)
			  RETURNING id`
	err := d.db.QueryRowx(query, moneyProvider).StructScan(&moneyProvider)
	if err != nil {
		return MoneyProvider{}, errors.Wrap(err, "Failed to create new MoneyProvider")
	}
	return moneyProvider, nil
}

// GetMoneyProvider retrieves a money provider by its ID.
func (d *dbImpl) GetMoneyProvider(id string) (MoneyProvider, error) {
	var moneyProvider MoneyProvider
	query := `SELECT id, name, creator_id, balance FROM money_provider WHERE id = $1`
	err := d.db.Get(&moneyProvider, query, id)
	return moneyProvider, err
}

// GetMoneyProvidersByUserID retrieves all money providers created by a specific user.
func (d *dbImpl) GetMoneyProvidersByUserID(userID string) ([]MoneyProvider, error) {
	var moneyProviders []MoneyProvider
	query := `SELECT id, name, creator_id, balance FROM money_provider WHERE creator_id = $1`
	err := d.db.Select(&moneyProviders, query, userID)
	return moneyProviders, err
}

// UpdateMoneyProvider updates an existing money provider in the database.
func (d *dbImpl) UpdateMoneyProvider(moneyProvider MoneyProvider) error {
	query := `UPDATE money_provider SET name = :name, balance = :balance WHERE id = :id`
	_, err := d.db.NamedExec(query, moneyProvider)
	return err
}
