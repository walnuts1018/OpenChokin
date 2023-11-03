package domain

type DB interface {
	NewUser(user User) error
	GetUser(id string) (User, error)
	UpdateUser(user User) error

	NewMoneyPool(moneyPool MoneyPool) error
	GetMoneyPool(id string) (MoneyPool, error)
	GetMoneyPoolsByUserID(userID string) ([]MoneyPool, error)
	UpdateMoneyPool(moneyPool MoneyPool) error
	UpdateMoneyPoolShareUserIDs(moneyPoolID string, shareUserIDs []string) error

	NewMoneyProvider(moneyProvider MoneyProvider) error
	GetMoneyProvider(id string) (MoneyProvider, error)
	GetMoneyProvidersByUserID(userID string) ([]MoneyProvider, error)
	UpdateMoneyProvider(moneyProvider MoneyProvider) error

	NewStore(store Store) error
	GetStore(id string) (Store, error)
	GetStoresByUserID(userID string) ([]Store, error)
	UpdateStore(store Store) error

	NewItem(item Item) error
	GetItem(id string) (Item, error)
	GetItemsByUserID(userID string) ([]Item, error)
	UpdateItem(item Item) error
}
