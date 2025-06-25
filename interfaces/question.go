// Defines the shapes of methods that need to be associated to various objects
// in order for the app to work.
package interfaces

import (
	"time"

	"ama/api/application"
)

// Creates a new question without any additional input.
type QuestionConverter interface {
	// Anything that can create a question from itself with no additional input.
	Question(questionId string) application.Question
}

// Writes questions to the database
type QuestionWriter interface {
	// Creates a question using any object that can convert itself to a question.
	CreateQuestion(questionData QuestionConverter) (application.Question, error)
	// Updates a question using anything that is capable of converting itself
	// into a question.
	UpdateQuestion(id string, questionData QuestionConverter) (application.Question, error)
}

// Reads data from the database
type QuestionReader interface {
	// Given a question id, will return that question if it exists.
	ReadQuestion(id string) (application.Question, error)
	// Read several questions at a time and get the matching results based on a query.
	ReadQuestions(limit int, finalId string, tags []string) ([]application.Question, error)
}

// Reads data from the database
type QuestionDeleter interface {
	// Given a question id, delete the question and return when the delete was complete.
	DeleteQuestion(id string) (time.Time, error)
}

// Can perform both the reading and writing functionality for managing questions.
type QuestionManager interface {
	QuestionDeleter
	QuestionReader
	QuestionWriter
}
