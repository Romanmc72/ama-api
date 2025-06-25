package list

import (
	"ama/api/application/list"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"

	"github.com/google/uuid"
)

type postUserListRequest struct {
	ListName string `json:"name"`
}

func PostUserList(c interfaces.APIContext, db interfaces.ListCreator) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	var listRequest *postUserListRequest
	if err := c.BindJSON(&listRequest); err != nil {
		logger.Error("Error binding JSON", "error", err, "userId", userId)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("invalid request"))
		return
	}
	listName := listRequest.ListName
	logger.Info("Creating list for user", "userId", userId, "listName", listName)
	if listName == "" {
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("listName is required"))
		return
	}
	listId := uuid.NewString()
	list := list.List{
		ID:   listId,
		Name: listName,
	}
	err := db.CreateList(userId, list)
	if err != nil {
		logger.Error("Error creating list", "error", err, "userId", userId, "listName", listName)
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError("unable to create list"))
		return
	}
	c.IndentedJSON(http.StatusCreated, list)
}
