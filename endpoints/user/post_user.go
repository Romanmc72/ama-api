package user

import (
	"ama/api/application/requests"
	"ama/api/application/responses"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
)

// PostUser godoc
//
//	@Summary		Create a user
//	@Description	Creates a new user in the database and returns that new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		requests.PostUserRequest	true	"User data"
//	@Success		201		{object}	application.User
//	@Failure		400		{object}	responses.ErrorResponse
//	@Failure		500		{object}	responses.ErrorResponse
//	@Router			/user [post]
func PostUser(c interfaces.APIContext, db interfaces.UserCreator) {
	logger := logging.GetLogger()
	var postReq requests.PostUserRequest
	err := c.BindJSON(&postReq)
	if err != nil {
		msg := "malformed input data for user"
		logger.Error(msg, "body", c.GetString("body"), "error", err)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(msg))
		return
	}

	newUser, err := postReq.BaseUser()
	if err != nil {
		logger.Error("invalid input data for user", "body", newUser, "error", err)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(err.Error()))
		return
	}

	user, err := db.CreateUser(newUser)
	if err != nil {
		msg := "error creating user"
		logger.Error(msg, "error", err)
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError(msg))
		return
	}
	c.IndentedJSON(http.StatusCreated, user)
}
