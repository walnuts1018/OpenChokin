package domain

import (
	"fmt"
)

func (d *dbImpl) NewStore(store Store) (Store, error) {
	// クエリ文字列で位置パラメータを使って、:nameと:creator_idを$1と$2に置き換えます。
	query := `INSERT INTO store (name, creator_id)
			  VALUES ($1, $2)
			  RETURNING id`
	// QueryRowを使ってSQLクエリを実行し、戻り値のIDをスキャンします。
	err := d.db.QueryRow(query, store.Name, store.CreatorID).Scan(&store.ID)
	if err != nil {
		return Store{}, fmt.Errorf("failed to create new Store: %v", err)
	}
	return store, nil
}

// GetStore retrieves a single store by its ID.
func (d *dbImpl) GetStore(id string) (Store, error) {
	var store Store
	query := `SELECT id, name, creator_id FROM store WHERE id = $1`
	err := d.db.Get(&store, query, id)
	if err != nil {
		return Store{}, fmt.Errorf("error fetching store: %v", err)
	}
	return store, nil
}

// GetStoresByUserID retrieves all stores created by a specific user.
func (d *dbImpl) GetStoresByUserID(userID string) ([]Store, error) {
	var stores []Store
	query := `SELECT id, name, creator_id FROM store WHERE creator_id = $1`
	err := d.db.Select(&stores, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching stores: %v", err)
	}
	return stores, nil
}

// UpdateStore updates an existing store's details.
func (d *dbImpl) UpdateStore(store Store) error {
	query := `UPDATE store SET name = $1, creator_id = $2 WHERE id = $3`
	_, err := d.db.Exec(query, store.Name, store.CreatorID, store.ID)
	if err != nil {
		return fmt.Errorf("error updating store: %v", err)
	}
	return nil
}
