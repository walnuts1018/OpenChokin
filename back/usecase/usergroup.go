package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/walnuts1018/openchokin/back/domain"
)

type UserGroupMember struct {
	ID string `json:"id"`
}
type UserGroupResponse struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Members []UserGroupMember `json:"members"`
}

// AddUserGroup creates a new user group with the given members
func (u Usecase) AddUserGroup(userID string, name string, memberIDs []string) (UserGroupResponse, error) {
	// Log the action of adding a new user group
	log.Printf("ユーザーID %s による新しいユーザーグループ %s の追加を開始します。", userID, name)

	// Create a new UserGroup object
	newGroup := domain.UserGroup{
		CreatorID: userID,
		Name:      name,
	}

	// Add the new user group using the DB interface
	addedGroup, err := u.db.NewUserGroup(newGroup)
	if err != nil {
		log.Printf("ユーザーグループ %s の追加中にエラーが発生しました: %v", name, err)
		return UserGroupResponse{}, err
	}

	// Log the successful creation of the user group
	log.Printf("ユーザーグループ %s の追加に成功しました。グループID: %s", name, addedGroup.ID)

	// Create a response object
	response := UserGroupResponse{
		ID:   addedGroup.ID,
		Name: addedGroup.Name,
	}

	// For each member ID, create a UserGroupMember and add it to the response
	for _, memberID := range memberIDs {
		response.Members = append(response.Members, UserGroupMember{ID: memberID})
	}

	return response, nil
}

// GetUserGroups retrieves all user groups for a given user and constructs responses including group members.
func (u *Usecase) GetUserGroups(userID string) ([]UserGroupResponse, error) {
	// Log the action of retrieving user groups
	log.Printf("ユーザーID %s のユーザーグループの取得を開始します。", userID)

	// Retrieve all user groups created by the given userID.
	userGroups, err := u.db.GetUserGroups(userID)
	if err != nil {
		log.Printf("ユーザーグループの取得中にエラーが発生しました: %v", err)
		return nil, err
	}

	// Prepare a slice to hold the response data.
	var responseGroups []UserGroupResponse

	// Iterate over each user group to fetch its members and validate the creator.
	for _, group := range userGroups {
		// Check if the creator ID matches the userID; if not, return an error.
		if group.CreatorID != userID {
			errMessage := fmt.Sprintf("ユーザー %s はユーザーグループ %s の作成者ではありません。", userID, group.ID)
			log.Print(errMessage)
			return nil, fmt.Errorf(errMessage)
		}

		// Get the members of the current user group.
		members, err := u.db.GetUserGroupMembers(group.ID)
		if err != nil {
			log.Printf("ユーザーグループ %s のメンバー取得中にエラーが発生しました: %v", group.ID, err)
			return nil, err
		}

		// Map members data to UserGroupMember slice.
		memberResponses := make([]UserGroupMember, len(members))
		for i, member := range members {
			memberResponses[i] = UserGroupMember{ID: member.ID}
		}

		// Add the constructed UserGroupResponse to the response slice.
		responseGroups = append(responseGroups, UserGroupResponse{
			ID:      group.ID,
			Name:    group.Name,
			Members: memberResponses,
		})
	}

	// Log the successful retrieval of user groups
	log.Printf("ユーザーID %s のユーザーグループの取得が完了しました。", userID)

	return responseGroups, nil
}

// UpdateUserGroup updates an existing user group's members
func (u Usecase) UpdateUserGroup(userID string, userGroupID string, name string, memberIDs []string) (UserGroupResponse, error) {
	log.Printf("ユーザーID %s によるユーザーグループID %s の更新開始。", userID, userGroupID)

	// Retrieve and validate the user group
	userGroup, err := u.db.GetUserGroup(userGroupID)
	if err != nil {
		log.Printf("ユーザーグループID %s の取得中にエラー: %v", userGroupID, err)
		return UserGroupResponse{}, err
	}

	// Validate the userID against the CreatorID of the UserGroup
	if userID != userGroup.CreatorID {
		log.Printf("ユーザーID %s はユーザーグループID %s を更新する権限がありません。", userID, userGroupID)
		return UserGroupResponse{}, errors.New("user is not authorized to update this user group")
	}

	// Update the user group with the new member IDs
	updatedGroup, err := u.db.UpdateUserGroup(userGroupID, name, memberIDs)
	if err != nil {
		log.Printf("ユーザーグループID %s の更新中にエラー: %v", userGroupID, err)
		return UserGroupResponse{}, err
	}

	log.Printf("ユーザーグループID %s の更新が成功しました。", userGroupID)

	// Create a response object
	response := UserGroupResponse{
		ID:   updatedGroup.ID,
		Name: updatedGroup.Name,
	}

	// Populate the response with updated member information
	for _, id := range memberIDs {
		response.Members = append(response.Members, UserGroupMember{ID: id})
	}

	return response, nil
}

// DeleteUserGroup deletes an existing user group
func (u Usecase) DeleteUserGroup(userID string, userGroupID string) error {
	log.Printf("ユーザーID %s によるユーザーグループID %s の削除開始。", userID, userGroupID)

	// Retrieve and validate the user group
	userGroup, err := u.db.GetUserGroup(userGroupID)
	if err != nil {
		log.Printf("ユーザーグループID %s の取得中にエラー: %v", userGroupID, err)
		return err
	}

	// Validate the userID against the CreatorID of the UserGroup
	if userID != userGroup.CreatorID {
		log.Printf("ユーザーID %s はユーザーグループID %s を削除する権限がありません。", userID, userGroupID)
		return errors.New("user is not authorized to delete this user group")
	}

	// Delete the user group using the DB interface
	err = u.db.DeleteUserGroup(userGroupID)
	if err != nil {
		log.Printf("ユーザーグループID %s の削除中にエラー: %v", userGroupID, err)
		return err
	}

	log.Printf("ユーザーグループID %s の削除が成功しました。", userGroupID)
	return nil
}
