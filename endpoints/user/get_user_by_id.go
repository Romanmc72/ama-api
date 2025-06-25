package user

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
)

// Given the user id, grab the user from the database
func GetUserByUserId(c interfaces.APIContext, db interfaces.UserReader) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	user, err := db.ReadUser(userId)
	if err != nil {
		logger.Error("could not find user by id", "error", err, "user", userId)
		c.IndentedJSON(http.StatusNotFound, responses.NewError("could not find user"))
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
