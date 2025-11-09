// Package main holds the code for the main program that is responsible for
// running the simple API that reads and writes question data.
package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ama/api/auth"
	"ama/api/constants"
	"ama/api/database"
	_ "ama/api/docs"
	"ama/api/endpoints"
	"ama/api/endpoints/list"
	listQuestion "ama/api/endpoints/list/question"
	"ama/api/endpoints/question"
	"ama/api/endpoints/user"
	"ama/api/logging"
)

type PingPong struct {
	Ping string `json:"ping"`
}

// Ping godoc
//
//	@Summary		Check API health
//	@Description	see if the server is up and running with this health check endpoint
//	@Tags			healthcheck
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	PingPong	"{\"ping\": \"pong\"}"
//	@Router			/ping [get]
func Ping(c *gin.Context) { c.IndentedJSON(http.StatusOK, PingPong{Ping: "pong"}) }

//	@BasePath	/

// main() is the main entrypoint for the program that starts and serves the API.
//
//	@title	Ask Me Anything API
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

	router.Use(func(c *gin.Context) { auth.CORSHeaders(endpoints.NewAPIContext(c)) })

	router.GET(constants.PingPath, Ping)

	// Routes requiring the authorization header to be set
	authorizedGroup := router.Group("/")
	authorizedGroup.Use(func(c *gin.Context) { auth.VerifyToken(endpoints.NewAPIContext(c), jwtVerifier, logger) })

	// Question endpoints
	authorizedGroup.GET(
		constants.QuestionBasePath,
		// PingExample godoc
		//	@Summary	ping health check
		//	@Schemes
		//	@Description	do ping to see if the server is up
		//	@Tags			healthcheck
		//	@Accept			json
		//	@Produce		json
		//	@Success		200	{map[string]string{"ping": "pong"}}
		//	@Router			/ping [get]
		func(c *gin.Context) { question.GetQuestions(endpoints.NewAPIContext(c), &db) },
	)
	authorizedGroup.GET(
		constants.QuestionByIdPath,
		func(c *gin.Context) {
			question.GetQuestionById(endpoints.NewAPIContext(c), &db)
		},
	)

	// admin only question endpoints
	adminOnlyGroup := authorizedGroup.Group(constants.QuestionBasePath)
	adminOnlyGroup.Use(func(c *gin.Context) {
		auth.VerifyRequiredScope(endpoints.NewAPIContext(c), logger, constants.GetAdminScopes())
	})
	adminOnlyGroup.POST(
		constants.NoPath,
		func(c *gin.Context) { question.PostQuestion(endpoints.NewAPIContext(c), &db) },
	)
	adminOnlyGroup.DELETE(
		constants.QuestionIdPathSegment,
		func(c *gin.Context) { question.DeleteQuestionById(endpoints.NewAPIContext(c), &db) },
	)
	adminOnlyGroup.PUT(
		constants.QuestionIdPathSegment,
		func(c *gin.Context) { question.PutQuestionById(endpoints.NewAPIContext(c), &db) },
	)

	// Create user endpoint (not a part of the user-validated group)
	authorizedGroup.POST(
		constants.UserBasePath,
		func(c *gin.Context) { user.PostUser(endpoints.NewAPIContext(c), &db) },
	)

	validatedUserGroup := authorizedGroup.Group(constants.UserByIdPath)
	validatedUserGroup.Use(func(c *gin.Context) { auth.VerifyUserID(endpoints.NewAPIContext(c), logger) })

	// Validated user endpoints
	validatedUserGroup.DELETE(
		constants.NoPath,
		func(c *gin.Context) { user.DeleteUserById(endpoints.NewAPIContext(c), &db) },
	)
	validatedUserGroup.GET(
		constants.NoPath,
		func(c *gin.Context) { user.GetUserByUserId(endpoints.NewAPIContext(c), &db) },
	)
	validatedUserGroup.PUT(
		constants.NoPath,
		func(c *gin.Context) { user.PutUserByUserId(endpoints.NewAPIContext(c), &db) },
	)

	// User list endpoints
	validatedUserGroup.GET(
		constants.UserListPath,
		func(c *gin.Context) { list.GetUserLists(endpoints.NewAPIContext(c), &db) },
	)
	validatedUserGroup.POST(
		constants.UserListPath,
		func(c *gin.Context) { list.PostUserList(endpoints.NewAPIContext(c), &db) },
	)
	validatedUserGroup.GET(
		constants.UserListByIdPath,
		func(c *gin.Context) { list.GetUserListById(endpoints.NewAPIContext(c), &db) },
	)
	validatedUserGroup.DELETE(
		constants.UserListByIdPath,
		func(c *gin.Context) { list.DeleteUserListByID(endpoints.NewAPIContext(c), &db) },
	)
	validatedUserGroup.PUT(
		constants.UserListByIdPath,
		func(c *gin.Context) { list.PutUserListById(endpoints.NewAPIContext(c), &db) },
	)

	// User list question endpoints
	validatedUserGroup.GET(
		constants.UserListQuestionByIdPath,
		func(c *gin.Context) { listQuestion.GetQuestionFromList(endpoints.NewAPIContext(c), &db) },
	)
	validatedUserGroup.POST(
		constants.UserListQuestionByIdPath,
		func(c *gin.Context) { listQuestion.PostQuestionToList(endpoints.NewAPIContext(c), &db) },
	)
	validatedUserGroup.DELETE(
		constants.UserListQuestionByIdPath,
		func(c *gin.Context) { listQuestion.DeleteQuestionFromList(endpoints.NewAPIContext(c), &db) },
	)

	// User question list endpoints
	logger.Info("Starting the server")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run(fmt.Sprintf("0.0.0.0:%d", *port))
	logger.Info("Shutting down the server")
}
