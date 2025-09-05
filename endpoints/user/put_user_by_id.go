package user

import (
	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
)

// Updates the user given the user id and the update data
func PutUserByUserId(c interfaces.APIContext, db interfaces.UserWriter) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	var user application.User
	if err := c.BindJSON(&user); err != nil {
		logger.Error("Failed to bind user data", constants.UserIdPathIdentifier, userId, "error", err, "request", c.GetString("body"))
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("invalid input"))
		return
	}
	user.ID = userId
	if err := db.UpdateUser(user); err != nil {
		logger.Error("Failed to update user", constants.UserIdPathIdentifier, userId, "error", err, "request", c.GetString("body"))
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError("unable to update user"))
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}
