package user

import (
	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/application/user"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"strings"
)

// Updates the user given the user id and the update data
func PutUserByUserId(c interfaces.APIContext, db interfaces.UserWriter) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	if strings.TrimSpace(userId) == "" {
		msg := "userId cannot be blank"
		logger.Error(msg, constants.UserIdPathIdentifier, userId, "error", "userId is blank", "request", c.GetString("body"))
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(msg))
		return
	}
	var u user.BaseUser
	if err := c.BindJSON(&u); err != nil {
		logger.Error("Failed to bind user data", constants.UserIdPathIdentifier, userId, "error", err, "request", c.GetString("body"))
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("invalid input"))
		return
	}
	err := user.ValidateUser(u)
	if err != nil {
		logger.Error("Input user data invalid", constants.UserIdPathIdentifier, userId, "error", err, "request", u)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(err.Error()))
		return
	}
	finalUser := &application.User{
		ID:       userId,
		BaseUser: u,
	}
	if err := db.UpdateUser(finalUser); err != nil {
		logger.Error("Failed to update user", constants.UserIdPathIdentifier, userId, "error", err, "request", c.GetString("body"))
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError("unable to update user"))
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}
