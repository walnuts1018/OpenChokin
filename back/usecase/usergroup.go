package usecase

import (
	"errors"
	"fmt"

	"github.com/walnuts1018/openchokin/back/domain"
)

type UserGroupMember struct {
	ID string
}
type UserGruopResponse struct {
	ID      string `db:"id"`
	Name    string `db:"name"`
	Members []UserGroupMember
}

// AddUserGroup creates a new user group with the given members
func (u Usecase) AddUserGroup(userID string, name string, memberIDs []string) (UserGruopResponse, error) {
	// Create a new UserGroup object
	newGroup := domain.UserGroup{
		CreatorID: userID,
		Name:      name,
	}

	// Add the new user group using the DB interface
	addedGroup, err := u.db.NewUserGroup(newGroup)
	if err != nil {
		return UserGruopResponse{}, err
	}

	// Create a response object
	response := UserGruopResponse{
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
func (u *Usecase) GetUserGroups(userID string) ([]UserGruopResponse, error) {
	// Retrieve all user groups created by the given userID.
	userGroups, err := u.db.GetUserGroups(userID)
	if err != nil {
		return nil, err
	}

	// Prepare a slice to hold the response data.
	var responseGroups []UserGruopResponse

	// Iterate over each user group to fetch its members and validate the creator.
	for _, group := range userGroups {
		// Check if the creator ID matches the userID; if not, return an error.
		if group.CreatorID != userID {
			return nil, fmt.Errorf("user %s is not the creator of user group %s", userID, group.ID)
		}

		// Get the members of the current user group.
		members, err := u.db.GetUserGroupMembers(group.ID)
		if err != nil {
			return nil, err
		}

		// Map members data to UserGroupMember slice.
		memberResponses := make([]UserGroupMember, len(members))
		for i, member := range members {
			memberResponses[i] = UserGroupMember{ID: member.ID}
		}

		// Add the constructed UserGruopResponse to the response slice.
		responseGroups = append(responseGroups, UserGruopResponse{
			ID:      group.ID,
			Name:    group.Name,
			Members: memberResponses,
		})
	}

	return responseGroups, nil
}

// UpdateUserGroup updates an existing user group's members
func (u Usecase) UpdateUserGroup(userID string, userGroupID string, name string, memberIDs []string) (UserGruopResponse, error) {
	// Retrieve and validate the user group
	userGroup, err := u.db.GetUserGroup(userGroupID)
	if err != nil {
		return UserGruopResponse{}, err
	}

	// Validate the userID against the CreatorID of the UserGroup
	if userID != userGroup.CreatorID {
		return UserGruopResponse{}, errors.New("user is not authorized to update this user group")
	}

	// Update the user group with the new member IDs
	updatedGroup, err := u.db.UpdateUserGroup(userGroupID, name, memberIDs)
	if err != nil {
		return UserGruopResponse{}, err
	}

	// Create a response object
	response := UserGruopResponse{
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
	// Retrieve and validate the user group
	userGroup, err := u.db.GetUserGroup(userGroupID)
	if err != nil {
		return err
	}

	// Validate the userID against the CreatorID of the UserGroup
	if userID != userGroup.CreatorID {
		return errors.New("user is not authorized to delete this user group")
	}

	// Delete the user group using the DB interface
	err = u.db.DeleteUserGroup(userGroupID)
	if err != nil {
		return err
	}

	return nil
}
