package domain

import (
	"fmt"

	"github.com/pkg/errors"
)

func (d *dbImpl) NewMoneyPool(moneyPool MoneyPool) (MoneyPool, error) {
	// クエリ文字列で位置パラメータを使用します。
	query := `INSERT INTO money_pool (name, description, type, owner_id)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id`
	// QueryRowを使用してIDを取得します。
	var returnedID int64
	err := d.db.QueryRow(query, moneyPool.Name, moneyPool.Description, moneyPool.Type, moneyPool.OwnerID).Scan(&returnedID)
	if err != nil {
		return MoneyPool{}, errors.Wrap(err, "新規MoneyPoolの作成とIDの返却に失敗しました")
	}

	// 返却されたIDをmoneyPool構造体のIDフィールドに割り当てます。
	moneyPool.ID = fmt.Sprintf("%d", returnedID)

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

	var currentType string
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

	var poolType string
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

func (d *dbImpl) IsMoneyPoolSharedWithUser(id string, userID string) (bool, error) {
	// Check if the MoneyPool with the provided ID is of type Restricted.
	var poolType string
	query := `SELECT type FROM MoneyPool WHERE id = ?`
	err := d.db.Get(&poolType, query, id)
	if err != nil {
		return false, err // or return false, nil to ignore error handling
	}

	// If not Restricted, it's not shared with specific users, so return an error or false as per your error handling policy.
	if poolType != PublicTypeRestricted {
		return false, nil // Pool is not restricted, hence not explicitly shared
	}

	query = `
		SELECT EXISTS (
			SELECT 1 FROM UserGroupMembership ugm
			INNER JOIN RestrictedPublicationScope rps ON ugm.GroupID = rps.GroupID
			WHERE rps.PoolID = ? AND ugm.UserID = ?
		)`

	var exists bool
	err = d.db.Get(&exists, query, id, userID)
	if err != nil {
		return false, err
	}

	return exists, nil
}
