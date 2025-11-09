package list

import (
	"ama/api/application/list"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"fmt"
	"net/http"
	"strings"
)

// DeleteUserListByID godoc
//
//	@Summary		Delete a list
//	@Description	Deletes the list from the user profile and removes all questions from the list. If the list does not exist, returns successfully.
//	@Tags			list
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"User ID"
//	@Param			listId	path		string	true	"List ID"
//	@Success		200		{object}	nil
//	@Failure		400		{object}	responses.ErrorResponse
//	@Failure		500		{object}	responses.ErrorResponse
//	@Router			/user/{userId}/list/{listId} [delete]
func DeleteUserListByID(c interfaces.APIContext, db interfaces.ListDeleter) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	listId := c.Param(constants.ListIdPathIdentifier)
	if strings.TrimSpace(userId) == "" || strings.TrimSpace(listId) == "" {
		msg := "userId and listId cannot be blank"
		logger.Error(msg, constants.UserIdPathIdentifier, userId, constants.ListIdPathIdentifier, listId)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(msg))
		return
	}
	u, err := db.ReadUser(userId)
	if err != nil {
		logger.Error("user not found", constants.UserIdPathIdentifier, userId, constants.ListIdPathIdentifier, listId)
		c.IndentedJSON(http.StatusNotFound, responses.NewError("userId not found"))
		return
	}
	l, ok := u.GetList(listId)
	if !ok {
		logger.Warn("user list not found", constants.UserIdPathIdentifier, userId, constants.ListIdPathIdentifier, listId)
		c.IndentedJSON(http.StatusOK, nil)
		return
	}
	if l.Name == list.LikedQuestionsListName {
		msg := fmt.Sprintf("cannot delete the '%s' list", list.LikedQuestionsListName)
		logger.Warn(msg, constants.UserIdPathIdentifier, userId, constants.ListIdPathIdentifier, listId)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(msg))
		return
	}
	if err := db.DeleteList(userId, listId); err != nil {
		logger.Error("Unable to delete list", constants.UserIdPathIdentifier, userId, constants.ListIdPathIdentifier, listId, "error", err)
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError("unable to delete list"))
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}
