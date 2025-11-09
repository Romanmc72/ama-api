package question

import (
	"net/http"

	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
)

//	DeleteQuestionById godoc
//
// @Summary		Delete one question.
// @Description	Deletes one question using its id. If the question already does not exist, this will return successfully.
// @Tags			question
// @Accept			json
// @Produce		json
// @Param			questionId	path		string	true	"Question ID"
// @Success		200			{object}	responses.SuccessResponse
// @Failure		500			{object}	responses.ErrorResponse
// @Router			/question/{questionId} [delete]
func DeleteQuestionById(c interfaces.APIContext, db interfaces.QuestionDeleter) {
	logger := logging.GetLogger()
	id := c.Param(constants.QuestionIdPathIdentifier)
	deleteTime, err := db.DeleteQuestion(id)
	if err != nil {
		logger.Error("Unable to delete question", "error", err, "questionId", id)
		c.IndentedJSON(
			http.StatusInternalServerError,
			responses.NewError("Something went wrong trying to delete that document"),
		)
		return
	}
	logger.Debug(
		"Successfully deleted the question",
		"id", id,
	)
	c.IndentedJSON(
		http.StatusOK,
		responses.SuccessResponse{
			Success: true,
			Time:    deleteTime.Unix(),
		},
	)
}
