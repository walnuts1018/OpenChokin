package domain

import "fmt"

func (d *dbImpl) NewUserGroup(userGroup UserGroup) (UserGroup, error) {
	query := `INSERT INTO user_groups (name, creator_id) VALUES ($1, $2) RETURNING id`
	err := d.db.QueryRow(query, userGroup.Name, userGroup.CreatorID).Scan(&userGroup.ID)
	if err != nil {
		return UserGroup{}, fmt.Errorf("failed to create user group: %v", err)
	}
	return userGroup, nil
}

func (d *dbImpl) GetUserGroups(userID string) ([]UserGroup, error) {
	var userGroups []UserGroup
	query := `SELECT id, name, creator_id FROM user_groups WHERE creator_id = $1`
	err := d.db.Select(&userGroups, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups for user %s: %v", userID, err)
	}
	return userGroups, nil
}

func (d *dbImpl) GetUserGroup(id string) (UserGroup, error) {
	var userGroup UserGroup
	query := `SELECT id, name, creator_id FROM user_groups WHERE id = $1`
	err := d.db.Get(&userGroup, query, id)
	if err != nil {
		return UserGroup{}, fmt.Errorf("failed to get user group with id %s: %v", id, err)
	}
	return userGroup, nil
}

func (d *dbImpl) GetUserGroupMembers(groupID string) ([]User, error) {
	var users []User
	query := `SELECT u.id FROM users u 
              JOIN user_group_membership ugm ON u.id = ugm.user_id 
              WHERE ugm.group_id = $1`
	err := d.db.Select(&users, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users for group %s: %v", groupID, err)
	}
	return users, nil
}

func (d *dbImpl) UpdateUserGroup(id string, name string, memberIDs []string) (UserGroup, error) {
	// Transaction start
	tx, err := d.db.Beginx()
	if err != nil {
		return UserGroup{}, fmt.Errorf("failed to start transaction: %v", err)
	}

	// Update the user group name
	query := `UPDATE user_groups SET name = $2 WHERE id = $1`
	_, err = tx.Exec(query, id, name)
	if err != nil {
		tx.Rollback()
		return UserGroup{}, fmt.Errorf("failed to update user group name: %v", err)
	}

	// Clear existing membership
	_, err = tx.Exec(`DELETE FROM user_group_membership WHERE group_id = $1`, id)
	if err != nil {
		tx.Rollback()
		return UserGroup{}, fmt.Errorf("failed to clear user group membership: %v", err)
	}

	// Add new members
	for _, userID := range memberIDs {
		_, err = tx.Exec(`INSERT INTO user_group_membership (group_id, user_id) VALUES ($1, $2)`, id, userID)
		if err != nil {
			tx.Rollback()
			return UserGroup{}, fmt.Errorf("failed to add user to group: %v", err)
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return UserGroup{}, fmt.Errorf("failed to commit user group update: %v", err)
	}

	return d.GetUserGroup(id) // Fetch and return the updated user group
}

func (d *dbImpl) DeleteUserGroup(id string) error {
	// Transaction start
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Delete entries from user_group_membership table first to avoid foreign key constraint violation
	_, err = tx.Exec(`DELETE FROM user_group_membership WHERE group_id = $1`, id)
	if err != nil {
		tx.Rollback() // rollback if any error occurs
		return fmt.Errorf("failed to delete user group memberships for group %s: %v", id, err)
	}

	// Delete the user group
	_, err = tx.Exec(`DELETE FROM user_groups WHERE id = $1`, id)
	if err != nil {
		tx.Rollback() // rollback if any error occurs
		return fmt.Errorf("failed to delete user group with id %s: %v", id, err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit delete operation for user group %s: %v", id, err)
	}

	return nil
}
