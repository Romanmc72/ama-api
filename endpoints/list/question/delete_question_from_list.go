package question

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"strings"
)

func DeleteQuestionFromList(c interfaces.APIContext, db interfaces.ListUpdater) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	listId := c.Param(constants.ListIdPathIdentifier)
	questionId := c.Param(constants.QuestionIdPathIdentifier)
	if strings.TrimSpace(userId) == "" || strings.TrimSpace(listId) == "" || strings.TrimSpace(questionId) == "" {
		logger.Error(
			"Missing path parameters",
			"userId", userId,
			"listId", listId,
			"questionId", questionId,
		)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("Missing path parameters"))
		return
	}
	err := db.RemoveQuestionFromList(userId, listId, questionId)
	if err != nil {
		logger.Error(
			"Error removing question from list",
			"error", err,
			"userId", userId,
			"listId", listId,
			"questionId", questionId,
		)
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError("Error removing question from list"))
		return
	}
	c.IndentedJSON(http.StatusOK, responses.NewSuccessResponse(true))
}
