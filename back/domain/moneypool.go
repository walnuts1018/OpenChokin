package domain

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

func (d *dbImpl) NewMoneyPool(moneyPool MoneyPool) (MoneyPool, error) {
	// クエリ文字列で位置パラメータを使用します。
	query := `INSERT INTO money_pool (name, description, type, owner_id, emoji, is_deleted)
			  VALUES ($1, $2, $3, $4, $5, false)
			  RETURNING id`
	// QueryRowを使用してIDを取得します。
	var returnedID int64
	err := d.db.QueryRow(query, moneyPool.Name, moneyPool.Description, moneyPool.Type, moneyPool.OwnerID, moneyPool.Emoji).Scan(&returnedID)
	if err != nil {
		return MoneyPool{}, errors.Wrap(err, "新規MoneyPoolの作成とIDの返却に失敗しました")
	}

	// 返却されたIDをmoneyPool構造体のIDフィールドに割り当てます。
	moneyPool.ID = fmt.Sprintf("%d", returnedID)

	return moneyPool, nil
}

func (d *dbImpl) GetMoneyPool(id string) (MoneyPool, error) {
	var moneyPool MoneyPool
	query := `SELECT * FROM money_pool WHERE id = $1 AND is_deleted = false`
	err := d.db.Get(&moneyPool, query, id)
	if err != nil {
		return MoneyPool{}, fmt.Errorf("could not find money pool: %v", err)
	}
	return moneyPool, nil
}

func (d *dbImpl) GetMoneyPoolsByUserID(userID string) ([]MoneyPool, error) {
	var moneyPools []MoneyPool
	query := `SELECT * FROM money_pool WHERE owner_id = $1 AND is_deleted = false`
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
	// 公開タイプを取得
	err = tx.Get(&currentType, "SELECT type FROM money_pool WHERE id = $1", moneyPool.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// if the current type is restricted and the new type is not restricted, delete the restricted publication scope table
	if currentType == "restricted" && moneyPool.Type != "restricted" {
		_, err := tx.Exec("DELETE FROM restricted_publication_scope WHERE pool_id = $1", moneyPool.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 名前付きパラメータを位置パラメータに置き換えたクエリを作成します
	query := `UPDATE money_pool SET name = $2, description = $3, type = $4, owner_id = $5, emoji = $6 WHERE id = $1`
	// Execを使用して更新を実行し、パラメータを順番にバインドします
	_, err = tx.Exec(query, moneyPool.ID, moneyPool.Name, moneyPool.Description, moneyPool.Type, moneyPool.OwnerID, moneyPool.Emoji)
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
	query := `UPDATE money_pool SET is_deleted = true, deleted_at = $2 WHERE id = $1`
	result, err := d.db.Exec(query, id, time.Now())
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
	log.Printf("ユーザー共有マネープールのチェック開始: MoneyPoolID=%s, UserID=%s", id, userID)

	// マネープールのタイプをチェックします。
	var poolType string
	query := `SELECT type FROM MoneyPool WHERE id = $1`
	err := d.db.Get(&poolType, query, id)
	if err != nil {
		log.Printf("マネープールのタイプの取得に失敗しました: %v", err)
		return false, err
	}
	log.Printf("マネープールのタイプ: %s", poolType)

	// マネープールがRestrictedタイプでない場合は、共有されていないと判断します。
	if poolType != "Restricted" {
		log.Println("マネープールはRestrictedタイプではありません。共有されていません。")
		return false, nil
	}

	// マネープールが特定のユーザーと共有されているかどうかを確認します。
	query = `
		SELECT COUNT(*) FROM UserGroupMembership ugm
		INNER JOIN RestrictedPublicationScope rps ON ugm.GroupID = rps.GroupID
		WHERE rps.PoolID = $1 AND ugm.UserID = $2
	`
	var count int
	err = d.db.Get(&count, query, id, userID)
	if err != nil {
		log.Printf("共有状態の確認中にエラーが発生しました: %v", err)
		return false, err
	}

	if count > 0 {
		log.Println("マネープールは指定ユーザーと共有されています。")
	} else {
		log.Println("マネープールは指定ユーザーと共有されていません。")
	}

	return count > 0, nil
}
