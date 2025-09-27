package list

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"strings"
)

func GetUserListById(c interfaces.APIContext, db interfaces.ListReader) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	listId := c.Param(constants.ListIdPathIdentifier)
	limit, finalId, tags, random := endpoints.GetReadQuestionsParamsWithDefaults(c)

	if strings.TrimSpace(userId) == "" || strings.TrimSpace(listId) == "" {
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
	retryLimit := 3
	if random && len(questions) == 0 {
		for retries := 0; retries < retryLimit; retries += 1 {
			logger.Info("Unable to find list questions using random param, trying again", "retry", retries, "finalId", finalId)
			limit, finalId, tags, _ = endpoints.GetReadQuestionsParamsWithDefaults(c)
			list, questions, err = db.ReadList(userId, listId, limit, finalId, tags)
			if err != nil {
				logger.Error("Error reading user list by ID", "error", err, "listId", listId)
				c.IndentedJSON(http.StatusNotFound, responses.NewError("unable to read list"))
				return
			}
			if len(questions) > 0 {
				logger.Info("Found list questions using random param", "retry", retries, "finalId", finalId)
				break
			}
		}
	}

	c.IndentedJSON(
		http.StatusOK,
		responses.NewGetUserListByIdResponse(list, questions),
	)
}
