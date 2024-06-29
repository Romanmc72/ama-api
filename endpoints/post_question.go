package endpoints

import (
	"net/http"

	"ama/api/application"
	"ama/api/interfaces"
	"ama/api/logging"
)

// PostQuestions(c *gin.Context) will create a brand new question in the database.
func PostQuestion(c interfaces.APIContext, db interfaces.QuestionWriter) {
	logger := logging.GetLogger()
	var newQuestion application.NewQuestion

	if err := c.BindJSON(&newQuestion); err != nil {
		c.IndentedJSON(http.StatusBadRequest, application.NewError("invalid data"))
		logger.Error("Invalid input data provided of", "body", c.GetString("body"), "error", err)
		return
	}
	question, err := db.CreateQuestion(&newQuestion)
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			application.NewError("encountered an error writing that data"),
		)
		logger.Error("Encountered an error writing that question", "question", err)
		return
	}
	c.IndentedJSON(http.StatusCreated, question)
}
