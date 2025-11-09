package question

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"strings"
)

// GetQuestionFromList godoc
//
//	@Summary		Get a question from a list
//	@Description	Get a question from a list
//	@Tags			list question
//	@Accept			json
//	@Produce		json
//	@Param			userId		path		string	true	"User ID"
//	@Param			listId		path		string	true	"List ID"
//	@Param			questionId	path		string	true	"Question ID"
//	@Success		200			{object}	application.Question
//	@Failure		400			{object}	responses.ErrorResponse
//	@Failure		404			{object}	responses.ErrorResponse
//	@Router			/user/{userId}/list/{listId}/question/{questionId} [get]
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
