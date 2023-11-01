package psql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/walnuts1018/openchokin/back/config"
	"github.com/walnuts1018/openchokin/back/domain"
)

const (
	sslmode = "disable"
)

type DB struct {
	db *sqlx.DB
}

func dbInit() error {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v sslmode=%v", config.Config.PostgresHost, config.Config.PostgresPort, config.Config.PostgresAdminUser, config.Config.PostgresAdminPassword, sslmode))
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v", config.Config.PostgresDb))
	if err != nil {
		return fmt.Errorf("failed to create db: %w", err)
	}

	return nil
}

func NewDB() (*DB, error) {
	err := dbInit()
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v", config.Config.PostgresHost, config.Config.PostgresPort, config.Config.PostgresUser, config.Config.PostgresPassword, config.Config.PostgresDb, sslmode))
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	//Create Table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id text PRIMARY KEY, money_pool_ids text[])")
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS money_pools (id text PRIMARY KEY, name text, color varchar(6), is_world_public boolean, share_user_ids text[])")
	if err != nil {
		return nil, fmt.Errorf("failed to create money_pools table: %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS transactions (id text PRIMARY KEY, money_pool_id text, transaction_date date, title text, amount float8, labels text[], is_world_public boolean, share_user_ids text[], expectation boolean, store_id text, item_ids text[], description text) PARTITION BY RANGE (transaction_date_id)")
	if err != nil {
		return nil, fmt.Errorf("failed to create transactions table: %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS items (id text PRIMARY KEY, name text, price_per_unit float8, user_id text")
	if err != nil {
		return nil, fmt.Errorf("failed to create items table: %w", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS stores (id text PRIMARY KEY, name text, user_id text)")
	if err != nil {
		return nil, fmt.Errorf("failed to create stores table: %w", err)
	}

	return &DB{db: db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) NewUser(user domain.User) error {
	_, err := db.db.NamedExec(`INSERT INTO users (id, money_pool_ids) VALUES (:id, :money_pool_ids)`, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (db *DB) GetUser(id string) (domain.User, error) {
	user := domain.User{}
	err := db.db.Get(&user, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (db *DB) UpdateUser(user domain.User) error {
	_, err := db.db.NamedExec(`UPDATE users SET money_pool_ids = :money_pool_ids WHERE id = :id`, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (db *DB) NewMoneyPool(moneyPool domain.MoneyPool) error {
	_, err := db.db.NamedExec(`INSERT INTO money_pools (id, name, color, is_world_public, share_user_ids) VALUES (:id, :name, :color, :is_world_public, :share_user_ids)`, moneyPool)
	if err != nil {
		return fmt.Errorf("failed to insert money pool: %w", err)
	}
	return nil
}

func (db *DB) GetMoneyPool(id string) (domain.MoneyPool, error) {
	moneyPool := domain.MoneyPool{}
	err := db.db.Get(&moneyPool, `SELECT * FROM money_pools WHERE id = $1`, id)
	if err != nil {
		return domain.MoneyPool{}, fmt.Errorf("failed to get money pool: %w", err)
	}
	return moneyPool, nil
}

func (db *DB) GetMoneyPoolsByUsers(user domain.User) ([]domain.MoneyPool, error) {
	moneyPools := []domain.MoneyPool{}
	err := db.db.Select(&moneyPools, `SELECT * FROM money_pools WHERE id = ANY($1)`, user.MoneyPoolIDs)
	if err != nil {
		return []domain.MoneyPool{}, fmt.Errorf("failed to get money pools: %w", err)
	}
	return moneyPools, nil
}

func (db *DB) UpdateMoneyPool(moneyPool domain.MoneyPool) error {
	_, err := db.db.NamedExec(`UPDATE money_pools SET name = :name, color = :color, is_world_public = :is_world_public, share_user_ids = :share_user_ids WHERE id = :id`, moneyPool)
	if err != nil {
		return fmt.Errorf("failed to update money pool: %w", err)
	}
	return nil
}

func (db *DB) NewTransaction(transaction domain.Transaction) error {
	partitioningKey := transaction.TransactionDate.Format("2006-01")

	_, err := db.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS transactions_%v PARTITION OF transactions FOR VALUES FROM ('%v-01') TO ('%v-01')`, partitioningKey, partitioningKey, partitioningKey))
	if err != nil {
		return fmt.Errorf("failed to create partition: %w", err)
	}

	_, err = db.db.NamedExec(fmt.Sprintf(`INSERT INTO transactions_%v (id, money_pool_id, transaction_date, title, amount, labels, is_world_public, share_user_ids, expectation, store_id, item_ids, description) VALUES (:id, :money_pool_id, :transaction_date, :title, :amount, :labels, :is_world_public, :share_user_ids, :expectation, :store_id, :item_ids, :description)`, partitioningKey), transaction)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}
	return nil
}

func (db *DB) GetTransaction(transactionID string, hint domain.GetTransactionHints) (domain.Transaction, error) {
	transaction := domain.Transaction{}
	if len(hint.PartitioningKeys) > 0 {
		for _, partitioningKey := range hint.PartitioningKeys {
			err := db.db.Get(&transaction, fmt.Sprintf(`SELECT * FROM transactions_%v WHERE id = $1`, partitioningKey), transactionID)
			if err != nil {
				return domain.Transaction{}, fmt.Errorf("failed to get transaction: %w", err)
			}
		}
	} else {
		err := db.db.Get(&transaction, `SELECT * FROM transactions WHERE id = $1`, transactionID)
		if err != nil {
			return domain.Transaction{}, fmt.Errorf("failed to get transaction: %w", err)
		}
	}
	return transaction, nil
}

func (db *DB) GetTransactionsByMoneyPool(moneyPoolID string, hint domain.GetTransactionHints) ([]domain.Transaction, error) {
	transactions := []domain.Transaction{}
	if len(hint.PartitioningKeys) > 0 {
		for _, partitioningKey := range hint.PartitioningKeys {
			tmpTransactions := []domain.Transaction{}
			err := db.db.Select(&tmpTransactions, fmt.Sprintf(`SELECT * FROM transactions_%v WHERE money_pool_id = $1`, partitioningKey), moneyPoolID)
			if err == nil {
				transactions = append(transactions, tmpTransactions...)
			}
		}
	} else {
		err := db.db.Select(&transactions, `SELECT * FROM transactions WHERE money_pool_id = $1`, moneyPoolID)
		if err != nil {
			return []domain.Transaction{}, fmt.Errorf("failed to get transactions: %w", err)
		}
	}
	return transactions, nil
}

func (db *DB) GetTransactionsByUser(userID string, hint domain.GetTransactionHints) ([]domain.Transaction, error) {
	transactions := []domain.Transaction{}
	if len(hint.PartitioningKeys) > 0 {
		for _, partitioningKey := range hint.PartitioningKeys {
			tmpTransactions := []domain.Transaction{}
			err := db.db.Select(&tmpTransactions, fmt.Sprintf(`SELECT * FROM transactions_%v WHERE $1 = ANY(share_user_ids)`, partitioningKey), userID)
			if err == nil {
				transactions = append(transactions, tmpTransactions...)
			}
		}
	} else {
		err := db.db.Select(&transactions, `SELECT * FROM transactions WHERE $1 = ANY(share_user_ids)`, userID)
		if err != nil {
			return []domain.Transaction{}, fmt.Errorf("failed to get transactions: %w", err)
		}
	}
	return transactions, nil
}

func (db *DB) UpdateTransaction(transaction domain.Transaction) error {
	partitioningKey := transaction.TransactionDate.Format("2006-01")

	_, err := db.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS transactions_%v PARTITION OF transactions FOR VALUES FROM ('%v-01') TO ('%v-01')`, partitioningKey, partitioningKey, partitioningKey))
	if err != nil {
		return fmt.Errorf("failed to create partition: %w", err)
	}

	_, err = db.db.NamedExec(fmt.Sprintf(`UPDATE transactions_%v SET money_pool_id = :money_pool_id, transaction_date = :transaction_date, title = :title, amount = :amount, labels = :labels, is_world_public = :is_world_public, share_user_ids = :share_user_ids, expectation = :expectation, store_id = :store_id, item_ids = :item_ids, description = :description WHERE id = :id`, partitioningKey), transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

func (db *DB) NewItem(item domain.Item) error {
	_, err := db.db.NamedExec(`INSERT INTO items (id, name, price_per_unit) VALUES (:id, :name, :price_per_unit, :user_id)`, item)
	if err != nil {
		return fmt.Errorf("failed to insert item: %w", err)
	}
	return nil
}

func (db *DB) GetItem(id string) (domain.Item, error) {
	item := domain.Item{}
	err := db.db.Get(&item, `SELECT * FROM items WHERE id = $1`, id)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to get item: %w", err)
	}
	return item, nil
}

func (db *DB) GetItemsByUser(userID string) ([]domain.Item, error) {
	items := []domain.Item{}
	err := db.db.Select(&items, `SELECT * FROM items WHERE user_id = $1`, userID)
	if err != nil {
		return []domain.Item{}, fmt.Errorf("failed to get items: %w", err)
	}
	return items, nil
}

func (db *DB) UpdateItem(item domain.Item) error {
	_, err := db.db.NamedExec(`UPDATE items SET name = :name, price_per_unit = :price_per_unit, user_id = :user_id WHERE id = :id`, item)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}
	return nil
}

func (db *DB) IsItemExist(itemID string) (bool, error) {
	var count int
	err := db.db.Get(&count, `SELECT COUNT(*) FROM items WHERE id = $1`, itemID)
	if err != nil {
		return false, fmt.Errorf("failed to get item: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (db *DB) NewStore(store domain.Store) error {
	_, err := db.db.NamedExec(`INSERT INTO stores (id, name, user_id) VALUES (:id, :name, :user_id)`, store)
	if err != nil {
		return fmt.Errorf("failed to insert store: %w", err)
	}
	return nil
}

func (db *DB) GetStore(id string) (domain.Store, error) {
	store := domain.Store{}
	err := db.db.Get(&store, `SELECT * FROM stores WHERE id = $1`, id)
	if err != nil {
		return domain.Store{}, fmt.Errorf("failed to get store: %w", err)
	}
	return store, nil
}

func (db *DB) GetStoresByUser(userID string) ([]domain.Store, error) {
	stores := []domain.Store{}
	err := db.db.Select(&stores, `SELECT * FROM stores WHERE user_id = $1`, userID)
	if err != nil {
		return []domain.Store{}, fmt.Errorf("failed to get stores: %w", err)
	}
	return stores, nil
}

func (db *DB) UpdateStore(store domain.Store) error {
	_, err := db.db.NamedExec(`UPDATE stores SET name = :name, user_id = :user_id WHERE id = :id`, store)
	if err != nil {
		return fmt.Errorf("failed to update store: %w", err)
	}
	return nil
}

func (db *DB) IsStoreExist(storeID string) (bool, error) {
	var count int
	err := db.db.Get(&count, `SELECT COUNT(*) FROM stores WHERE id = $1`, storeID)
	if err != nil {
		return false, fmt.Errorf("failed to get store: %w", err)
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
