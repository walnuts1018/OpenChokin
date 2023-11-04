package usecase

import (
	"fmt"
	"log"

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

// NewUser creates a new user in the database
func (u Usecase) NewUser(user domain.User) (domain.User, error) {
	log.Printf("新規ユーザー作成を開始します。ユーザーID: %s", user.ID)
	user, err := u.db.NewUser(user)
	if err != nil {
		log.Printf("新規ユーザー作成に失敗しました。ユーザーID: %s, エラー: %v", user.ID, err)
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	log.Printf("新規ユーザー作成に成功しました。ユーザーID: %s", user.ID)
	return user, nil
}

// GetUser retrieves a user by ID from the database
func (u Usecase) GetUser(id string) (domain.User, error) {
	log.Printf("ユーザーID %s によるユーザー情報の取得を試みます。", id)
	user, err := u.db.GetUser(id)
	if err != nil {
		log.Printf("ユーザーID %s の情報取得に失敗しました。エラー: %v", id, err)
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	log.Printf("ユーザーID %s の情報取得に成功しました。", id)
	return user, nil
}

// UpdateUser updates an existing user's information in the database
func (u Usecase) UpdateUser(user domain.User) error {
	log.Printf("ユーザーID %s の情報更新を開始します。", user.ID)
	err := u.db.UpdateUser(user)
	if err != nil {
		log.Printf("ユーザーID %s の情報更新に失敗しました。エラー: %v", user.ID, err)
		return fmt.Errorf("failed to update user: %w", err)
	}
	log.Printf("ユーザーID %s の情報更新に成功しました。", user.ID)
	return nil
}
