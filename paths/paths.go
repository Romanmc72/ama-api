// The paths that represent various API endpoints.
package paths

const (
	// This is the base API endpoint off of which all other question endpoints will listen.
	Base         = "/questions"
	// For working with a specific question.
	QuestionById = Base + "/:id"
	// For interacting with a user profile
	UserById     = "/user/:id"
)
