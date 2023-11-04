package usecase

import (
	"fmt"

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

func (u Usecase) NewUser(user domain.User) (domain.User, error) {
	user, err := u.db.NewUser(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (u Usecase) GetUser(id string) (domain.User, error) {
	user, err := u.db.GetUser(id)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (u Usecase) UpdateUser(user domain.User) error {
	err := u.db.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
