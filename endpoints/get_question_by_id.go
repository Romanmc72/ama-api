package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ama/api/application"
	"ama/api/interfaces"
)

// GetQuestionById(c *gin.Context) will get one question using its unique id.
func GetQuestionById(c *gin.Context, db interfaces.QuestionReader) {
	id := c.Param("id")
	questionFetched, err := db.ReadQuestion(id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.IndentedJSON(
				http.StatusNotFound,
				application.NewError("Did not find a document with that Id"),
			)
			return
		}
		c.IndentedJSON(
			http.StatusInternalServerError,
			application.NewError("Something went wrong trying to get that document"),
		)
		return
	}
	c.IndentedJSON(http.StatusOK, questionFetched)
}
