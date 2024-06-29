package endpoints

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ama/api/application"
	"ama/api/interfaces"
	"ama/api/logging"
)

const (
	limitParam = "limit"
	defaultLimit = 0
	finalIdParam = "finalId"
	defaultFinalId = ""
	tagParam = "tags"
	defaultTag = ""
)

// GetQuestions(c *gin.Context) retrieves all of the questions from the database.
func GetQuestions(c *gin.Context, db interfaces.QuestionReader) {
	logger := logging.GetLogger()
	limit := GetQueryParamToInt(c, limitParam, defaultLimit)
	finalId := c.DefaultQuery(finalIdParam, defaultFinalId)
	rawTags := c.DefaultQuery(tagParam, defaultTag)
	tags := strings.Split(rawTags, application.SearchTagDelimiter)
	questions, err := db.ReadQuestions(limit, finalId, tags)
	if err != nil {
		logger.Error(
			"Something went wrong getting the questions",
			"error", err,
			limitParam, limit,
			finalIdParam, finalId,
			tagParam, tags,
		)
		c.IndentedJSON(
			http.StatusInternalServerError,
			application.NewError("Could not retrieve questions"),
		)
		return
	}
	if len(questions) > 0 {
		c.IndentedJSON(http.StatusOK, questions)
	} else {
		c.IndentedJSON(http.StatusOK, []application.Question{})
	}
}
