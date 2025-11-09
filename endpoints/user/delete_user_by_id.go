package user

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"strings"
)

// DeleteUserById godoc
//
//	@Summary		Delete a user given their ID
//	@Description	Delete a user from the database given their user id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"User ID"
//	@Success		200		{object}	responses.DeleteUserResponse
//	@Failure		400		{object}	responses.ErrorResponse
//	@Failure		500		{object}	responses.ErrorResponse
//	@Router			/user/{userId} [delete]
func DeleteUserById(c interfaces.APIContext, deleter interfaces.UserDeleter) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	if strings.TrimSpace(userId) == "" {
		msg := "cannot have blank user id"
		logger.Error(msg, "error", msg, constants.UserIdPathIdentifier, userId)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(msg))
		return
	}
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
