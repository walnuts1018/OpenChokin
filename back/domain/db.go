package domain

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type dbImpl struct {
	db *sqlx.DB
}

// NewDB creates a new dbImpl and returns it as a DB interface.
func NewDB(db *sqlx.DB) DB {
	return &dbImpl{
		db: db,
	}
}

type DB interface {
	NewUser() (User, error)
	GetUser(id string) (User, error)
	UpdateUser(user User) error

	NewMoneyPool(moneyPool MoneyPool) (MoneyPool, error)
	GetMoneyPool(id string) (MoneyPool, error)
	GetMoneyPoolsByUserID(userID string) ([]MoneyPool, error)
	UpdateMoneyPool(moneyPool MoneyPool) error
	DeleteMoneyPool(id string) error
	ShareMoneyPoolWithUserGroups(moneyPoolID string, shareUserGruopIDs []string) error

	NewMoneyProvider(moneyProvider MoneyProvider) (MoneyProvider, error)
	GetMoneyProvider(id string) (MoneyProvider, error)
	GetMoneyProvidersByUserID(userID string) ([]MoneyProvider, error)
	UpdateMoneyProvider(moneyProvider MoneyProvider) error
	DeleteMoneyProvider(id string) error

	NewStore(store Store) (Store, error)
	GetStore(id string) (Store, error)
	GetStoresByUserID(userID string) ([]Store, error)
	UpdateStore(store Store) error

	NewItem(item Item) (Item, error)
	GetItem(id string) (Item, error)
	GetItemsByUserID(userID string) ([]Item, error)
	UpdateItem(item Item) error

	NewPayment(payment Payment) (Payment, error)
	GetPayment(id string) (Payment, error)
	GetPaymentsByMoneyPoolID(moneyPoolID string) ([]Payment, error)
	UpdatePayment(payment Payment) error
	DeletePayment(id string) error

	GetMoneyPoolBalance(moneyPoolID string, includeExpceted bool) (float64, error)                       // transactionからマネープールの残高を計算する
	GetMoneyPoolBalanceOfDate(moneyPoolID string, date time.Time, includeExpceted bool) (float64, error) // transactionからマネープールの残高を計算する（ある日までの）

	NewUserGroup(userGroup UserGroup) (UserGroup, error)
	GetUserGroups(userID string) ([]UserGroup, error)
	GetUserGroup(id string) (UserGroup, error)
	GetUserGroupMembers(id string) ([]User, error)
	UpdateUserGroup(id string, name string, memberIDs []string) (UserGroup, error)
	DeleteUserGroup(id string) error
}
