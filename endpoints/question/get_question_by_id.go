package question

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
)

// GetQuestionById(c *gin.Context) will get one question using its unique id.
func GetQuestionById(c interfaces.APIContext, db interfaces.QuestionReader) {
	id := c.Param(constants.QuestionIdPathIdentifier)
	questionFetched, err := db.ReadQuestion(id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.IndentedJSON(
				http.StatusNotFound,
				responses.NewError("Did not find a document with that Id"),
			)
			return
		}
		c.IndentedJSON(
			http.StatusInternalServerError,
			responses.NewError("Something went wrong trying to get that document"),
		)
		return
	}
	c.IndentedJSON(http.StatusOK, questionFetched)
}
