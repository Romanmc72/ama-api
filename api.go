// Package main holds the code for the main program that is responsible for
// running the simple API that reads and writes question data.
package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ama/api/auth"
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
	jwtVerifier, err := auth.NewAuthClient()
	if err != nil {
		logger.Error("Could not connect to Firebase", "error", err)
	}
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

	// Routes requiring the authorization header to be set
	authorizedGroup := router.Group("/")
	// TODO: Something is wrong here and the token is not parsing correctly
	authorizedGroup.Use(func(c *gin.Context) { auth.VerifyToken(c, jwtVerifier, logger) })

	// Question endpoints
	authorizedGroup.GET(
		constants.QuestionBasePath,
		func(c *gin.Context) { question.GetQuestions(c, &db) },
	)
	authorizedGroup.GET(
		constants.QuestionByIdPath,
		func(c *gin.Context) { question.GetQuestionById(c, &db) },
	)

	// admin only question endpoints
	adminOnlyGroup := authorizedGroup.Group(constants.QuestionBasePath)
	adminOnlyGroup.Use(func(c *gin.Context) { auth.VerifyRequiredScope(c, logger, constants.GetAdminScopes()) })
	adminOnlyGroup.POST(
		constants.QuestionBasePath,
		func(c *gin.Context) { question.PostQuestion(c, &db) },
	)
	adminOnlyGroup.DELETE(
		constants.QuestionByIdPath,
		func(c *gin.Context) { question.DeleteQuestionById(c, &db) },
	)
	adminOnlyGroup.PUT(
		constants.QuestionByIdPath,
		func(c *gin.Context) { question.PutQuestionById(c, &db) },
	)

	// Create user endpoint (not a part of the user-validated group)
	authorizedGroup.POST(
		constants.UserBasePath,
		func(c *gin.Context) { user.PostUser(c, &db) },
	)

	validatedUserGroup := authorizedGroup.Group(constants.UserByIdPath)
	validatedUserGroup.Use(func(c *gin.Context) { auth.VerifyUserID(c, logger) })

	// Validated user endpoints
	validatedUserGroup.DELETE(
		constants.NoPath,
		func(c *gin.Context) { user.DeleteUserById(c, &db) },
	)
	validatedUserGroup.GET(
		constants.NoPath,
		func(c *gin.Context) { user.GetUserByUserId(c, &db) },
	)
	validatedUserGroup.PUT(
		constants.NoPath,
		func(c *gin.Context) { user.PutUserByUserId(c, &db) },
	)

	// User list endpoints
	validatedUserGroup.GET(
		constants.UserListPath,
		func(c *gin.Context) { list.GetUserLists(c, &db) },
	)
	validatedUserGroup.POST(
		constants.UserListPath,
		func(c *gin.Context) { list.PostUserList(c, &db) },
	)
	validatedUserGroup.GET(
		constants.UserListByIdPath,
		func(c *gin.Context) { list.GetUserListById(c, &db) },
	)
	validatedUserGroup.DELETE(
		constants.UserListByIdPath,
		func(c *gin.Context) { list.DeleteUserListByID(c, &db) },
	)
	validatedUserGroup.PUT(
		constants.UserListByIdPath,
		func(c *gin.Context) { list.PutUserListById(c, &db) },
	)

	// User list question endpoints
	validatedUserGroup.POST(
		constants.UserListQuestionByIdPath,
		func(c *gin.Context) { listQuestion.PostQuestionToList(c, &db) },
	)
	validatedUserGroup.DELETE(
		constants.UserListQuestionByIdPath,
		func(c *gin.Context) { listQuestion.DeleteQuestionFromList(c, &db) },
	)

	// User question list endpoints
	logger.Info("Starting the server")
	router.Run(fmt.Sprintf("0.0.0.0:%d", *port))
	logger.Info("Shutting down the server")
}
