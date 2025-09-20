// Defines the NewQuestion data type and any methods accessible to it.
package application

import (
	"strings"
)

// `NewQuestion` is the base question but without an id attached as the
// application will create one.
type NewQuestion struct {
	// The actual question itself.
	Prompt string `json:"prompt" binding:"required" validate:"required"`
	// The tags identifying what kind of question this is.
	Tags []string `json:"tags" binding:"required" validate:"required"`
}

// `Question.String` converts the question to a string format.
func (q *NewQuestion) String() string {
	return "NewQuestion(" +
		"Prompt=" + q.Prompt +
		", Tags=[" + strings.Join(q.Tags, ", ") +
		"])"
}

// `Question` will convert the NewQuestion into a question object which can
// be written to the database.
func (q *NewQuestion) Question(questionId string) Question {
	return Question{
		ID:     questionId,
		Prompt: q.Prompt,
		Tags:   q.Tags,
	}
}
