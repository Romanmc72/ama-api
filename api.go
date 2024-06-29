// Package main holds the code for the main program that is responsible for
// running the simple API that reads and writes question data.
package main

import (
	"github.com/gin-gonic/gin"

	"ama/api/database"
	"ama/api/endpoints"
	"ama/api/logging"
	"ama/api/paths"
)

// main() is the main entrypoint for the program that starts and serves the API.
func main() {
	logger := logging.GetLogger()
	router := gin.Default()
	db, err := database.Connect()
	if err != nil {
		logger.Error("Could not connect to the database", "error", err)
	}
	defer db.Close()
	router.GET(
		paths.Base,
		func(c *gin.Context) { endpoints.GetQuestions(c, &db) },
	)
	router.POST(
		paths.Base,
		func(c *gin.Context) { endpoints.PostQuestion(c, &db) },
	)
	router.GET(
		paths.QuestionById,
		func(c *gin.Context) { endpoints.GetQuestionById(c, &db) },
	)
	router.DELETE(
		paths.QuestionById,
		func(c *gin.Context) { endpoints.DeleteQuestionById(c, &db) },
	)
	router.PUT(
		paths.QuestionById,
		func(c *gin.Context) { endpoints.PutQuestionById(c, &db) },
	)
	logger.Info("Starting the server")
	router.Run("localhost:8081")
	logger.Info("Shutting down the server")
}
