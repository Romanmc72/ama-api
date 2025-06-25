package list

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"net/http"
)

// Deletes the list from the user profile and removes all questions from the list
func DeleteUserListByID(c interfaces.APIContext, db interfaces.ListDeleter) {
	userId := c.Param(constants.UserIdPathIdentifier)
	listId := c.Param(constants.ListIdPathIdentifier)
	if err := db.DeleteList(userId, listId); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError("unable to delete list"))
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}
