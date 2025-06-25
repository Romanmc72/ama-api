package user

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
)

// Delete a user from the database given their user id
func DeleteUserById(c interfaces.APIContext, deleter interfaces.UserDeleter) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	deleteTime, err := deleter.DeleteUser(userId)
	if err != nil {
		msg := "error deleting user"
		logger.Error(msg, "error", err)
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError(msg))
		return
	}
	c.IndentedJSON(
		http.StatusOK,
		responses.NewDeleteUserResponse(true, deleteTime.Unix()),
	)
}
