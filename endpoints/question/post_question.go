package question

import (
	"net/http"

	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/interfaces"
	"ama/api/logging"
)

// PostQuestion godoc
//
//	@Summary		Create a question
//	@Description	Will create a brand new question in the database.
//	@Tags			question
//	@Accept			json
//	@Produce		json
//	@Param			question	body		application.NewQuestion	true	"Question data"
//	@Success		201			{object}	application.Question
//	@Failure		400			{object}	responses.ErrorResponse
//	@Failure		500			{object}	responses.ErrorResponse
//	@Router			/question [post]
func PostQuestion(c interfaces.APIContext, db interfaces.QuestionWriter) {
	logger := logging.GetLogger()
	var newQuestion application.NewQuestion

	if err := c.BindJSON(&newQuestion); err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("invalid data"))
		logger.Error("Invalid input data provided of", "body", c.GetString("body"), "error", err)
		return
	}

	if err := application.ValidateQuestion(newQuestion.Question("")); err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.NewError(err.Error()))
		logger.Error("Invalid input data provided of", "body", newQuestion, "error", err)
		return
	}

	question, err := db.CreateQuestion(&newQuestion)
	if err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			responses.NewError("encountered an error writing that data"),
		)
		logger.Error("Encountered an error writing that question", "question", err)
		return
	}
	c.IndentedJSON(http.StatusCreated, question)
}
