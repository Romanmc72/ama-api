package endpoints

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ama/api/application"
	"ama/api/interfaces"
	"ama/api/logging"
)

// DeleteQuestionById(c *gin.Context) will delete one question using its id.
func DeleteQuestionById(c *gin.Context, db interfaces.QuestionDeleter) {
	logger := logging.GetLogger()
	id := c.Param("id")
	deleteTime, err := db.DeleteQuestion(id)
	if err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			application.NewError("Something went wrong trying to delete that document"),
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
