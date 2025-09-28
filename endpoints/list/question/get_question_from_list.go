package question

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"strings"
)

func GetQuestionFromList(c interfaces.APIContext, db interfaces.ListReader) {
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	listId := c.Param(constants.ListIdPathIdentifier)
	questionId := c.Param(constants.QuestionIdPathIdentifier)
	if strings.TrimSpace(userId) == "" || strings.TrimSpace(listId) == "" || strings.TrimSpace(questionId) == "" {
		logger.Error(
			"list id, question id, or user id was blank",
			"userId", userId,
			"listId", listId,
			"questionId", questionId,
		)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("userId, listId, and questionId cannot be blank"))
		return
	}
	question, err := db.ReadListQuestion(userId, listId, questionId)
	if err != nil {
		logger.Error(
			"Error reading question",
			"error", err,
			"userId", userId,
			"listId", listId,
			"questionId", questionId,
		)
		c.IndentedJSON(http.StatusNotFound, responses.NewError("Unable to find list question"))
		return
	}
	c.IndentedJSON(http.StatusOK, question)
}
