package domain

func (d *dbImpl) NewUser(user User) (User, error) {
	var newUser User
	// トランザクションを開始
	tx, err := d.db.Beginx()
	if err != nil {
		return newUser, err
	}

	// user.IDを持つ行を挿入する。ここではidが文字列として定義されていると仮定します。
	query := "INSERT INTO users (id) VALUES ($1) RETURNING id"
	err = tx.QueryRow(query, user.ID).Scan(&newUser.ID)
	if err != nil {
		tx.Rollback() // エラーがあればロールバック
		return newUser, err
	}

	// トランザクションをコミット
	err = tx.Commit()
	if err != nil {
		return newUser, err
	}

	// 新しいユーザーオブジェクトを返す
	return newUser, nil
}

func (d *dbImpl) GetUser(id string) (User, error) {
	var user User
	err := d.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (d *dbImpl) UpdateUser(user User) error {
	return nil
}
