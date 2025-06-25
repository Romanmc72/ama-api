package question

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
)

// DeleteQuestionById(c *gin.Context) will delete one question using its id.
func DeleteQuestionById(c interfaces.APIContext, db interfaces.QuestionDeleter) {
	logger := logging.GetLogger()
	id := c.Param(constants.QuestionIdPathIdentifier)
	deleteTime, err := db.DeleteQuestion(id)
	if err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			responses.NewError("Something went wrong trying to delete that document"),
		)
		return
	}
	logger.Debug(
		"Successfully deleted the question",
		"id", id,
	)
	c.IndentedJSON(
		http.StatusOK,
		gin.H{
			"message": fmt.Sprintf(
				"Deleted %s at %s",
				id,
				deleteTime.String(),
			),
		},
	)
}
