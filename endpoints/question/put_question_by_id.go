package question

import (
	"net/http"
	"strings"

	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
)

// PutQuestionsById(c *gin.Context) will update an existing question in the database.
// If the question does not exist it will be created.
func PutQuestionById(c interfaces.APIContext, db interfaces.QuestionWriter) {
	logger := logging.GetLogger()
	id := c.Param(constants.QuestionIdPathIdentifier)
	if strings.TrimSpace(id) == "" {
		msg := "questionId cannot be blank"
		logger.Error(msg, "body", c.GetString("body"), "questionId", id)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("missing questionId in path"))
		return
	}
	var newQuestion application.NewQuestion

	if err := c.BindJSON(&newQuestion); err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("invalid data"))
		logger.Error("Invalid input data provided of", "body", c.GetString("body"), "error", err, constants.QuestionIdPathIdentifier, id)
		return
	}

	if err := application.ValidateQuestion(newQuestion.Question(id)); err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(err.Error()))
		logger.Error("Invalid input data provided of", "body", c.GetString("body"), "error", err, constants.QuestionIdPathIdentifier, id)
		return
	}

	updatedQuestion, err := db.UpdateQuestion(id, &newQuestion)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError("unable to update question"))
		logger.Error("Update failed", "question", newQuestion, "error", err)
		return
	}
	c.IndentedJSON(http.StatusOK, updatedQuestion)
}
