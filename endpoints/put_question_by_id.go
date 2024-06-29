package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ama/api/application"
	"ama/api/interfaces"
	"ama/api/logging"
)

// PutQuestionsById(c *gin.Context) will update an existing question in the database.
// If the question does not exist it will be created.
func PutQuestionById(c *gin.Context, db interfaces.QuestionWriter) {
	id := c.Param("id")
	var newQuestion application.NewQuestion
	logger := logging.GetLogger()
	
	if err := c.BindJSON(&newQuestion); err != nil {
		c.IndentedJSON(http.StatusBadRequest, application.NewError("invalid data"))
		logger.Error("Invalid input data provided of", "body", c.GetString("body"), "error", err)
		return
	}

	updatedQuestion, err := db.UpdateQuestion(id, &newQuestion)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, application.NewError("unable to update question"))
		logger.Error("Update failed", "question", newQuestion, "error", err)
		return
	}
	c.IndentedJSON(http.StatusOK, updatedQuestion)
}
