package usecase

import (
	"encoding/base64"
	"fmt"
	"math/rand"
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

func (u Usecase) NewUser(userid string) (domain.User, error) {
	user := domain.User{
		ID: userid,
	}
	err := u.db.NewUser(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (u Usecase) NewMoneyPool(moneyPoolName, moneyPoolColor, userID string, isWorldPublic bool) (domain.MoneyPool, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	moneyPoolID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v%v%v", moneyPoolName, userID, r.Int63())))
	moneyPool := domain.MoneyPool{
		ID:            moneyPoolID,
		Name:          moneyPoolName,
		Color:         moneyPoolColor,
		IsWorldPublic: isWorldPublic,
		ShareUserIDs:  []string{userID},
	}

	err := u.db.NewMoneyPool(moneyPool)
	if err != nil {
		return domain.MoneyPool{}, fmt.Errorf("failed to create money pool: %w", err)
	}

	user, err := u.db.GetUser(userID)
	if err != nil {
		return domain.MoneyPool{}, fmt.Errorf("failed to get user: %w", err)
	}

	user.MoneyPoolIDs = append(user.MoneyPoolIDs, moneyPool.ID)
	err = u.db.UpdateUser(user)
	if err != nil {
		return domain.MoneyPool{}, fmt.Errorf("failed to update user: %w", err)
	}
	return moneyPool, nil
}
