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

// GetQuestions godoc
//
//	@Summary		Get zero or more questions
//	@Description	Retrieves all of the questions from the database given input query parameters.
//	@Tags			question
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int						false	"Limit"
//	@Param			finalId	query		string					false	"Final ID from previous page"
//	@Param			tag		query		int						false	"Tag to match (specify multiple times for && match)"
//	@Param			random	query		bool					false	"Get a random question"
//	@Success		201		{object}	[]application.Question	"The questions that were found (or weren't)"
//	@Failure		500		{object}	responses.ErrorResponse
//	@Router			/question [get]
func GetQuestions(c interfaces.APIContext, db interfaces.QuestionReader) {
	logger := logging.GetLogger()
	limit, finalId, tags, random := endpoints.GetReadQuestionsParamsWithDefaults(c)
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
	retryLimit := 3
	if random && len(questions) == 0 {
		for retries := 0; retries < retryLimit; retries += 1 {
			logger.Info("Did not find questions with random params set, trying again", "retry", retries, "finalId", finalId)
			limit, finalId, tags, _ := endpoints.GetReadQuestionsParamsWithDefaults(c)
			questions, err = db.ReadQuestions(limit, finalId, tags)
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
				logger.Info("Found questions with random param", "retry", retries, "finalId", finalId)
				break
			}
		}
	}
	if len(questions) > 0 {
		c.IndentedJSON(http.StatusOK, questions)
	} else {
		c.IndentedJSON(http.StatusOK, []application.Question{})
	}
}
