package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler for getting user group details
func getUserGroups(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	response, err := uc.GetUserGroups(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Handler for creating a new user group
func createUserGroup(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	var requestBody struct {
		Name      string   `json:"name"`
		MemberIDs []string `json:"member_ids"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := uc.AddUserGroup(userID, requestBody.Name, requestBody.MemberIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response)
}

// Handler for updating a user group
func updateUserGroup(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	userGroupID := c.Param("usergroup_id")
	var requestBody struct {
		Name      string   `json:"name"`
		MemberIDs []string `json:"member_ids"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := uc.UpdateUserGroup(userID, userGroupID, requestBody.Name, requestBody.MemberIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// Handler for deleting a user group
func deleteUserGroup(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	userGroupID := c.Param("usergroup_id")
	err := uc.DeleteUserGroup(userID, userGroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
