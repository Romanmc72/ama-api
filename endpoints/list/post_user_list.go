package list

import (
	"ama/api/application/errors"
	"ama/api/application/list"
	"ama/api/application/requests"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"strings"
)

func PostUserList(c interfaces.APIContext, db interfaces.ListCreator) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	if strings.TrimSpace(userId) == "" {
		msg := "userId cannot be blank"
		logger.Error("userID is blank", "error", msg, "userId", userId)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(msg))
		return
	}
	var listRequest requests.PostUserListRequest
	if err := c.BindJSON(&listRequest); err != nil {
		logger.Error("Error binding JSON", "error", err, "userId", userId)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("invalid request"))
		return
	}
	logger.Info("Creating list for user", "userId", userId, "listName", listRequest.Name)
	l := list.List{
		Name: listRequest.Name,
	}
	l, err := db.CreateList(userId, l)
	if err != nil {
		if ve, ok := err.(*errors.ValidationError); ok && ve.ValidationErrCt > 0 {
			logger.Error("Data validation error", "error", err, "userId", userId, "listName", listRequest.Name)
			c.IndentedJSON(http.StatusBadRequest, responses.NewError(ve.Error()))
			return
		}
		logger.Error("Error creating list", "error", err, "userId", userId, "listName", listRequest.Name)
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError("unable to create list"))
		return
	}
	c.IndentedJSON(http.StatusCreated, l)
}
