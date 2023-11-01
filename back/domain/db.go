package domain

type DB interface {
	NewUser(user User) error
	GetUser(id string) (User, error)
	UpdateUser(user User) error
	NewMoneyPool(moneyPool MoneyPool) error
	GetMoneyPool(id string) (MoneyPool, error)
	GetMoneyPoolsByUsers(user User) ([]MoneyPool, error)
	UpdateMoneyPool(moneyPool MoneyPool) error
	NewTransaction(transaction Transaction) error
	GetTransaction(transactionID string, hint GetTransactionHints) (Transaction, error)
	GetTransactionsByMoneyPool(moneyPoolID string, hint GetTransactionHints) ([]Transaction, error)
	GetTransactionsByUser(userID string, hint GetTransactionHints) ([]Transaction, error)
	UpdateTransaction(transaction Transaction) error
	NewItem(item Item) error
	GetItem(id string) (Item, error)
	GetItemsByUser(userID string) ([]Item, error)
	UpdateItem(item Item) error
	IsItemExist(itemID string) (bool, error)
	NewStore(store Store) error
	GetStore(id string) (Store, error)
	GetStoresByUser(userID string) ([]Store, error)
	UpdateStore(store Store) error
	IsStoreExist(storeID string) (bool, error)
}
