package domain

import (
	"fmt"

	"github.com/pkg/errors"
)

func (d *dbImpl) NewItem(item Item) (Item, error) {
	// クエリ文字列で名前付きパラメータを位置パラメータに置き換えます
	query := `INSERT INTO item (name, creator_id)
			  VALUES ($1, $2)
			  RETURNING id`
	// QueryRowを使ってクエリを実行し、結果のIDをスキャンします
	err := d.db.QueryRow(query, item.Name, item.CreatorID).Scan(&item.ID)
	if err != nil {
		// errors.Wrapを使って、エラーのコンテキストを提供します
		return Item{}, errors.Wrap(err, "Failed to create new Item")
	}
	return item, nil
}

// GetItem retrieves an item by its ID.
func (d *dbImpl) GetItem(id string) (Item, error) {
	var item Item
	err := d.db.Get(&item, "SELECT * FROM item WHERE id = $1", id)
	if err != nil {
		return Item{}, fmt.Errorf("error fetching item: %v", err)
	}
	return item, nil
}

// GetItemsByUserID retrieves all items created by a specific user.
func (d *dbImpl) GetItemsByUserID(userID string) ([]Item, error) {
	var items []Item
	err := d.db.Select(&items, "SELECT * FROM item WHERE creator_id = $1", userID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateItem updates an existing item.
func (d *dbImpl) UpdateItem(item Item) error {
	// Use positional parameters with $1, $2, etc.
	query := `UPDATE item SET name = $1, creator_id = $2 WHERE id = $3`
	// Use the Exec function with the struct's fields passed in the order of the parameters.
	_, err := d.db.Exec(query, item.Name, item.CreatorID, item.ID)
	if err != nil {
		return fmt.Errorf("error updating item: %v", err)
	}
	return nil
}
