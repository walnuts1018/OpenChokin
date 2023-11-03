package domain

func (d *dbImpl) NewUser() (User, error) {
	var newUser User
	tx, err := d.db.Beginx()
	if err != nil {
		return newUser, err
	}
	err = tx.QueryRow("INSERT INTO users DEFAULT VALUES RETURNING id").Scan(&newUser.ID)
	if err != nil {
		tx.Rollback()
		return newUser, err
	}
	err = tx.Commit()
	if err != nil {
		return newUser, err
	}
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
