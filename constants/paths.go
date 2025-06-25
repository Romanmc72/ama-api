// The paths that represent various API endpoints.
package constants

// Paths defining API routes.
const (
	// The path that allows for a health check ping.
	PingPath = "/ping"
	// This is the base API endpoint off of which all other question endpoints will listen.
	QuestionBasePath = "/question"
	// For working with a specific question.
	QuestionByIdPath = QuestionBasePath + "/:" + QuestionIdPathIdentifier
	// The root user endpoint.
	UserBasePath = "/user"
	// For interacting with a user profile.
	UserByIdPath = UserBasePath + "/:" + UserIdPathIdentifier
	// The root user question list endpoint.
	UserListPath = UserByIdPath + "/list"
	// For working with a specific user question list.
	UserListByIdPath = UserListPath + "/:" + ListIdPathIdentifier
	// For working with a questions in a user question list.
	UserListQuestionPath = UserListPath + "/:" + ListIdPathIdentifier + QuestionBasePath
	// For working with a specific question in a user question list.
	UserListQuestionByIdPath = UserListQuestionPath + "/:" + QuestionIdPathIdentifier
)
