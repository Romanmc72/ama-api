package question

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"strings"
)

// PostQuestionToList godoc
//
//	@Summary		Add a question to a list
//	@Description	Add a question to a list
//	@Tags			list question
//	@Accept			json
//	@Produce		json
//	@Param			userId		path		string	true	"User ID"
//	@Param			listId		path		string	true	"List ID"
//	@Param			questionId	path		string	true	"Question ID"
//	@Success		200			{object}	responses.SuccessResponse
//	@Failure		400			{object}	responses.ErrorResponse
//	@Failure		404			{object}	responses.ErrorResponse
//	@Failure		500			{object}	responses.ErrorResponse
//	@Router			/user/{userId}/list/{listId}/question/{questionId} [post]
func PostQuestionToList(c interfaces.APIContext, db interfaces.ListUpdater) {
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
	question, err := db.ReadQuestion(questionId)
	if err != nil {
		logger.Error(
			"Error reading question",
			"error", err,
			"userId", userId,
			"listId", listId,
			"questionId", questionId,
		)
		c.IndentedJSON(http.StatusNotFound, responses.NewError("Invalid question id"))
		return
	}
	err = db.AddQuestionToList(userId, listId, question)
	if err != nil {
		logger.Error(
			"Error adding question to list",
			"error", err,
			"userId", userId,
			"listId", listId,
			"question", question,
		)
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError("Error adding question to list"))
		return
	}
	c.IndentedJSON(http.StatusCreated, responses.NewSuccessResponse(true))
}
