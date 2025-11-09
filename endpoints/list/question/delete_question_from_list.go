package question

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"strings"
)

// DeleteQuestionFromList godoc
//
//	@Summary		Delete a question from a list
//	@Description	Delete a question from a list. If the question does not exist in the list, returns successfully.
//	@Tags			list question
//	@Accept			json
//	@Produce		json
//	@Param			userId		path		string	true	"User ID"
//	@Param			listId		path		string	true	"List ID"
//	@Param			questionId	path		string	true	"Question ID"
//	@Success		200			{object}	responses.SuccessResponse
//	@Failure		400			{object}	responses.ErrorResponse
//	@Failure		500			{object}	responses.ErrorResponse
//	@Router			/user/{userId}/list/{listId}/question/{questionId} [delete]
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
