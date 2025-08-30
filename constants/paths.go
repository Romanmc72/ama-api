// The paths that represent various API endpoints.
package constants

const NoPath = ""
const BasePath = "/"

const (
	// The path that allows for a health check ping.
	PingPath = "/ping"
	// This is the base API endpoint off of which all other question endpoints will listen.
	QuestionBasePath = "/question"
	QuestionIdPathSegment = "/:" + QuestionIdPathIdentifier
	// For working with a specific question.
	QuestionByIdPath = QuestionBasePath + QuestionIdPathSegment
)

const (
	// The root user endpoint.
	UserBasePath = "/user"
	// For interacting with a user profile.
	UserByIdPath = UserBasePath + "/:" + UserIdPathIdentifier
	// The root user question list endpoint.
	// NOTE: This is not chained to the user by ID path because the user ID is
	// passed in via the route grouping and is automatically chained in the api
	// definition as a result.
	UserListPath = "/list"
	// For working with a specific user question list.
	UserListByIdPath = UserListPath + "/:" + ListIdPathIdentifier
	// For working with a questions in a user question list.
	UserListQuestionPath = UserListPath + "/:" + ListIdPathIdentifier + QuestionBasePath
	// For working with a specific question in a user question list.
	UserListQuestionByIdPath = UserListQuestionPath + "/:" + QuestionIdPathIdentifier
)
