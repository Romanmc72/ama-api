package question

import (
	"net/http"

	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
)

// PutQuestionsById(c *gin.Context) will update an existing question in the database.
// If the question does not exist it will be created.
func PutQuestionById(c interfaces.APIContext, db interfaces.QuestionWriter) {
	id := c.Param(constants.QuestionIdPathIdentifier)
	var newQuestion application.NewQuestion
	logger := logging.GetLogger()

	if err := c.BindJSON(&newQuestion); err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("invalid data"))
		logger.Error("Invalid input data provided of", "body", c.GetString("body"), "error", err)
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
