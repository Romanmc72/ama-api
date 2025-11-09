package list

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"strings"

	"net/http"
)

// GetUserLists godoc
//
//	@Summary		Get lists
//	@Description	Get all of the lists for a user.
//	@Tags			list
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"User ID"
//	@Success		201		{object}	[]list.List
//	@Failure		400		{object}	responses.ErrorResponse
//	@Failure		404		{object}	responses.ErrorResponse
//	@Router			/user/{userId}/list [get]
func GetUserLists(c interfaces.APIContext, db interfaces.UserReader) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	if strings.TrimSpace(userId) == "" {
		logger.Error("Error reading user lists", "error", "userId is empty")
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("userId is empty"))
		return
	}
	user, err := db.ReadUser(userId)
	if err != nil {
		logger.Error("Error reading user lists", "error", err, "userId", userId)
		c.IndentedJSON(http.StatusNotFound, responses.NewError("unable to read user"))
		return
	}
	c.IndentedJSON(http.StatusOK, user.Lists)
}
