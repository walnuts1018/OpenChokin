package usecase

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/walnuts1018/openchokin/back/domain"
)

type Usecase struct {
	db domain.DB
}

func NewUsecase(db domain.DB) *Usecase {
	return &Usecase{
		db: db,
	}
}

func (u Usecase) NewUser(userName, email string) (domain.User, error) {
	userid := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v%v", userName, email)))
	user := domain.User{
		ID: userid,
	}
	err := u.db.NewUser(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (u Usecase) GetUser(userID string) (domain.User, error) {
	return u.db.GetUser(userID)
}

func (u Usecase) NewMoneyPool(moneyPool domain.MoneyPool, user domain.User) error {
	err := u.db.NewMoneyPool(moneyPool)
	if err != nil {
		return fmt.Errorf("failed to create money pool: %w", err)
	}

	user.MoneyPoolIDs = append(user.MoneyPoolIDs, moneyPool.ID)
	err = u.db.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil

}

func (u Usecase) GetMoneyPool(moneyPoolID string) (domain.MoneyPool, error) {
	return u.db.GetMoneyPool(moneyPoolID)
}

func (u Usecase) GetMoneyPoolsByUser(user domain.User) ([]domain.MoneyPool, error) {
	return u.db.GetMoneyPoolsByUsers(user)
}

func (u Usecase) UpdateMoneyPool(moneyPool domain.MoneyPool) error {
	return u.db.UpdateMoneyPool(moneyPool)
}

func (u Usecase) NewTransaction(transaction domain.Transaction, user domain.User) error {

}

func (u Usecase) GetTransactionsByTimeRange(UserID string, moneyPoolID string, from time.Time, to time.Time) ([]domain.Transaction, error) {
	if moneyPoolID != "" {
		PartitioningKeys := []string{}
		for t := from; t.Before(to); t = t.AddDate(0, 1, 0) {
			PartitioningKeys = append(PartitioningKeys, t.Format("2006-01"))
		}
		return u.db.GetTransactionsByMoneyPool(moneyPoolID, domain.GetTransactionHints{
			PartitioningKeys: PartitioningKeys,
		})
	} else {
		PartitioningKeys := []string{}
		for t := from; t.Before(to); t = t.AddDate(0, 1, 0) {
			PartitioningKeys = append(PartitioningKeys, t.Format("2006-01"))
		}
		return u.db.GetTransactionsByUser(UserID, domain.GetTransactionHints{
			PartitioningKeys: PartitioningKeys,
		})
	}
}

func (u Usecase) GetTransaction(transactionID string) (domain.Transaction, error) {
	return u.db.GetTransaction(transactionID, domain.GetTransactionHints{})
}

func (u Usecase) UpdateTransaction(transaction domain.Transaction) error {
	return u.db.UpdateTransaction(transaction)
}

func (u Usecase) NewItem(item domain.Item, user domain.User) error {
	return u.db.NewItem(item)
}
