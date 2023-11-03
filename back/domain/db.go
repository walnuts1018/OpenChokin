package domain

type DB interface {
	NewUser(user User) error
	GetUser(id string) (User, error)
	UpdateUser(user User) error

	NewMoneyPool(moneyPool MoneyPool) error
	GetMoneyPool(id string) (MoneyPool, error)
	GetMoneyPoolsByUserID(userID string) ([]MoneyPool, error)
	UpdateMoneyPool(moneyPool MoneyPool) error

	NewMoneyProvider(moneyProvider MoneyProvider) error
	GetMoneyProvider(id string) (MoneyProvider, error)
	GetMoneyProvidersByUserID(userID string) ([]MoneyProvider, error)
	UpdateMoneyProvider(moneyProvider MoneyProvider) error
}
