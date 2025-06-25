package list

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
)

func GetUserListById(c interfaces.APIContext, db interfaces.ListReader) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	listId := c.Param(constants.ListIdPathIdentifier)
	limit, finalId, tags := endpoints.GetReadQuestionsParamsWithDefaults(c)

	if userId == "" || listId == "" {
		logger.Error("Error reading user list by ID", "error", "userId or listId is empty")
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("userId or listId is empty"))
		return
	}

	list, questions, err := db.ReadList(userId, listId, limit, finalId, tags)
	if err != nil {
		logger.Error("Error reading user list by ID", "error", err, "listId", listId)
		c.IndentedJSON(http.StatusNotFound, responses.NewError("unable to read list"))
		return
	}

	response := responses.NewGetUserListByIdResponse(list, questions)

	c.IndentedJSON(http.StatusOK, response)
}
