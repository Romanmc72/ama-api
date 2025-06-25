package question

import (
	"net/http"

	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints"
	"ama/api/interfaces"
	"ama/api/logging"
)

// GetQuestions(c *gin.Context) retrieves all of the questions from the database.
func GetQuestions(c interfaces.APIContext, db interfaces.QuestionReader) {
	logger := logging.GetLogger()
	limit, finalId, tags := endpoints.GetReadQuestionsParamsWithDefaults(c)
	questions, err := db.ReadQuestions(limit, finalId, tags)
	if err != nil {
		logger.Error(
			"Something went wrong getting the questions",
			"error", err,
			constants.LimitParam, limit,
			constants.FinalIdParam, finalId,
			constants.TagParam, tags,
		)
		c.IndentedJSON(
			http.StatusInternalServerError,
			responses.NewError("Could not retrieve questions"),
		)
		return
	}
	if len(questions) > 0 {
		c.IndentedJSON(http.StatusOK, questions)
	} else {
		c.IndentedJSON(http.StatusOK, []application.Question{})
	}
}
