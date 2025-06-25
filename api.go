// Package main holds the code for the main program that is responsible for
// running the simple API that reads and writes question data.
package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ama/api/constants"
	"ama/api/database"
	"ama/api/endpoints/list"
	listQuestion "ama/api/endpoints/list/question"
	"ama/api/endpoints/question"
	"ama/api/endpoints/user"
	"ama/api/logging"
)

// main() is the main entrypoint for the program that starts and serves the API.
func main() {
	port := flag.Int("port", 8088, "The port to launch the server on. Default=8088.")
	flag.Parse()

	logger := logging.GetLogger()
	router := gin.Default()
	db, err := database.Connect()
	if err != nil {
		logger.Error("Could not connect to the database", "error", err)
	}
	defer db.Close()
	// Health check endpoint
	router.GET(
		constants.PingPath,
		func(c *gin.Context) { c.IndentedJSON(http.StatusOK, map[string]string{"ping": "pong"}) },
	)

	// Question endpoints
	router.GET(
		constants.QuestionBasePath,
		func(c *gin.Context) { question.GetQuestions(c, &db) },
	)
	router.POST(
		constants.QuestionBasePath,
		func(c *gin.Context) { question.PostQuestion(c, &db) },
	)
	router.GET(
		constants.QuestionByIdPath,
		func(c *gin.Context) { question.GetQuestionById(c, &db) },
	)
	router.DELETE(
		constants.QuestionByIdPath,
		func(c *gin.Context) { question.DeleteQuestionById(c, &db) },
	)
	router.PUT(
		constants.QuestionByIdPath,
		func(c *gin.Context) { question.PutQuestionById(c, &db) },
	)

	// User endpoints
	router.DELETE(
		constants.UserByIdPath,
		func(c *gin.Context) { user.DeleteUserById(c, &db) },
	)
	router.GET(
		constants.UserByIdPath,
		func(c *gin.Context) { user.GetUserByUserId(c, &db) },
	)
	router.PUT(
		constants.UserByIdPath,
		func(c *gin.Context) { user.PutUserByUserId(c, &db) },
	)
	router.POST(
		constants.UserBasePath,
		func(c *gin.Context) { user.PostUser(c, &db) },
	)

	// User list endpoints
	router.GET(
		constants.UserListPath,
		func(c *gin.Context) { list.GetUserLists(c, &db) },
	)
	router.POST(
		constants.UserListPath,
		func(c *gin.Context) { list.PostUserList(c, &db) },
	)
	router.GET(
		constants.UserListByIdPath,
		func(c *gin.Context) { list.GetUserListById(c, &db) },
	)
	router.DELETE(
		constants.UserListByIdPath,
		func(c *gin.Context) { list.DeleteUserListByID(c, &db) },
	)
	router.PUT(
		constants.UserListByIdPath,
		func(c *gin.Context) { list.PutUserListById(c, &db) },
	)

	// User list question endpoints
	router.POST(
		constants.UserListQuestionByIdPath,
		func(c *gin.Context) { listQuestion.PostQuestionToList(c, &db) },
	)
	router.DELETE(
		constants.UserListQuestionByIdPath,
		func(c *gin.Context) { listQuestion.DeleteQuestionFromList(c, &db) },
	)

	// User question list endpoints
	logger.Info("Starting the server")
	router.Run(fmt.Sprintf("0.0.0.0:%d", *port))
	logger.Info("Shutting down the server")
}
