package domain

import (
	"fmt"

	"github.com/pkg/errors"
)

func (d *dbImpl) NewMoneyPool(moneyPool MoneyPool) (MoneyPool, error) {
	query := `INSERT INTO money_pool (name, description, type, owner_id)
			  VALUES (:name, :description, :type, :owner_id)
			  RETURNING id`
	err := d.db.QueryRowx(query, moneyPool).StructScan(&moneyPool)
	if err != nil {
		return MoneyPool{}, errors.Wrap(err, "Failed to create new MoneyPool")
	}
	return moneyPool, nil
}

func (d *dbImpl) GetMoneyPool(id string) (MoneyPool, error) {
	var moneyPool MoneyPool
	query := `SELECT * FROM money_pool WHERE id = $1`
	err := d.db.Get(&moneyPool, query, id)
	if err != nil {
		return MoneyPool{}, fmt.Errorf("could not find money pool: %v", err)
	}
	return moneyPool, nil
}

func (d *dbImpl) GetMoneyPoolsByUserID(userID string) ([]MoneyPool, error) {
	var moneyPools []MoneyPool
	query := `SELECT * FROM money_pool WHERE owner_id = $1`
	err := d.db.Select(&moneyPools, query, userID)
	if err != nil {
		return nil, fmt.Errorf("could not find money pools for user: %v", err)
	}
	return moneyPools, nil
}

func (d *dbImpl) UpdateMoneyPool(moneyPool MoneyPool) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	var currentType PublicType
	err = tx.Get(&currentType, "SELECT type FROM money_pool WHERE id = $1", moneyPool.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if currentType == PublicTypeRestricted && moneyPool.Type != PublicTypeRestricted {
		_, err := tx.Exec("DELETE FROM restricted_publication_scope WHERE pool_id = $1", moneyPool.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	query := `UPDATE money_pool SET name = :name, description = :description, type = :type, owner_id = :owner_id WHERE id = :id`
	_, err = tx.NamedExec(query, moneyPool)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (d *dbImpl) ShareMoneyPoolWithUserGroups(moneyPoolID string, shareUserGroupIDs []string) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	var poolType PublicType
	err = tx.Get(&poolType, "SELECT type FROM money_pool WHERE id = $1", moneyPoolID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if poolType != PublicTypeRestricted {
		tx.Rollback()
		return errors.New("money pool must be of type 'restricted' to share with user groups")
	}

	_, err = tx.Exec("DELETE FROM restricted_publication_scope WHERE pool_id = $1", moneyPoolID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, groupID := range shareUserGroupIDs {
		_, err := tx.Exec("INSERT INTO restricted_publication_scope (pool_id, group_id) VALUES ($1, $2)", moneyPoolID, groupID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *dbImpl) DeleteMoneyPool(id string) error {
	query := `DELETE FROM money_pool WHERE id = $1`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not delete money pool: %v", err)
	}

	// Execの結果から影響を受けた行の数を確認します。Deleteが実行されなかった場合にはエラーを返すことも可能です。
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not determine rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, nothing to delete")
	}

	return nil
}
