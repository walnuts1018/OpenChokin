/*
THIS IS NOT UP TO DATE
TODO: change
see @/back/readme.md
*/

package psql

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"database/sql"

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
	defer db.Close()

    var dbName string
    err = db.Get(&dbName, "SELECT datname FROM pg_database WHERE datname = $1", config.Config.PostgresDb)
    if err != nil && err != sql.ErrNoRows {
        return fmt.Errorf("error checking for database existence: %w", err)
    }

    // If the database does not exist, create it
    if dbName == "" {
        _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %v", config.Config.PostgresDb))
		if err != nil {
			return fmt.Errorf("failed to create db: %w", err)
		}
    }

    return nil
}

func NewDB() (*DB, error) {
	err := dbInit()
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

    db, err := sqlx.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v sslmode=%v", config.Config.PostgresHost, config.Config.PostgresPort, config.Config.PostgresAdminUser, config.Config.PostgresAdminPassword, sslmode))
    if err != nil {
        return nil, fmt.Errorf("failed to open db: %w", err)
    }

	// SQLファイルからテーブルを作成
	err = executeSQLFile(db, "/app/infra/psql/init.sql")
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}
func executeSQLFile(db *sqlx.DB, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open SQL file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sqlStatement string
	var inDOBlock bool

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// コメントを無視
		if strings.HasPrefix(trimmedLine, "--") {
			continue
		}

		// DOブロックの開始を検出
		if strings.HasPrefix(trimmedLine, "DO") {
			inDOBlock = true
		}

		// DOブロック内では、END; まで読み込む
		if inDOBlock && strings.HasPrefix(trimmedLine, "END;") {
			inDOBlock = false
		}

		sqlStatement += line + "\n" // SQLステートメントを行ごとに追加

		// SQLステートメントが終わったかどうか（セミコロンかDOブロックの終わり）
		if (!inDOBlock && strings.HasSuffix(trimmedLine, ";")) || (!inDOBlock && strings.HasPrefix(trimmedLine, "END;")) {
			_, err = db.Exec(sqlStatement)
			if err != nil {
				return fmt.Errorf("failed to exec SQL statement: %w", err)
			}
			sqlStatement = "" // ステートメントをリセット
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading SQL file: %w", err)
	}

	return nil
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
