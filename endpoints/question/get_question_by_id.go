package question

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
)

// GetQuestionById(c *gin.Context) will get one question using its unique id.
func GetQuestionById(c interfaces.APIContext, db interfaces.QuestionReader) {
	logger := logging.GetLogger()
	id := c.Param(constants.QuestionIdPathIdentifier)
	questionFetched, err := db.ReadQuestion(id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			logger.Error("Question not found", "error", err, "userId", id)
			c.IndentedJSON(
				http.StatusNotFound,
				responses.NewError("Did not find a document with that Id"),
			)
			return
		}
		logger.Error("Could not read user data", "error", err, "userId", id)
		c.IndentedJSON(
			http.StatusInternalServerError,
			responses.NewError("Something went wrong trying to get that document"),
		)
		return
	}
	c.IndentedJSON(http.StatusOK, questionFetched)
}
