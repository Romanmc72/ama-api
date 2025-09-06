package user

import (
	"ama/api/application/requests"
	"ama/api/application/responses"
	"ama/api/application/user"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
)

// Creates a new user in the database and returns that new user
func PostUser(c interfaces.APIContext, db interfaces.UserCreator) {
	logger := logging.GetLogger()
	var newUser requests.PostUserRequest
	err := c.BindJSON(&newUser)
	if err != nil {
		msg := "malformed input data for user"
		logger.Error(msg, "body", c.GetString("body"), "error", err)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(msg))
		return
	}

	if err = user.ValidateUser(newUser.BaseUser); err != nil {
		logger.Error("invalid input data for user", "body", newUser, "error", err)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(err.Error()))
		return
	}

	user, err := db.CreateUser(newUser.BaseUser)
	if err != nil {
		msg := "error creating user"
		logger.Error(msg, "error", err)
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError(msg))
		return
	}
	c.IndentedJSON(http.StatusCreated, user)
}
