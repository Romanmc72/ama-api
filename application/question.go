// Defines the Question data type and any methods accessible to it.
package application

import (
	"strings"
)

const SearchTagDelimiter = "|"

// `Question` is the fundamental piece of data stored within the app. It
// represents a short question number which can be easily remembered by a
// human and the prompt which is the main text of the question itself.
// Tags will help questions be found and categorized.
type Question struct {
	// TODO: Look at this for marshalling and unmarshalling the id
	// https://stackoverflow.com/questions/68870654/how-to-get-and-store-the-document-id-using-a-struct
	// The database identifier for this question.
	ID string `json:"id" firestore:"-" binding:"required"`
	// The actual question itself.
	Prompt string `json:"prompt" firestore:"prompt" binding:"required"`
	// The tags identifying what kind of question this is.
	Tags []string `json:"tags" firestore:"tags" binding:"required"`
}

// `Question.String` converts the question to a string format.
func (q *Question) String() string {
	return "Question(" +
		"Id=" + q.ID +
		", Prompt=" + q.Prompt +
		", Tags=[" + strings.Join(q.Tags, ", ") +
		"])"
}

// A question that returns itself to satisfy the interface
func (q *Question) Question() Question {
	return *q
}

// Deduplicates, sorts, and combines the search tags into a list of tags that
// are searchable by firebase.
func (q *Question) createSearchTags() []string {
	return Combine(q.Tags, SearchTagDelimiter)
}

// Converts the question to its database question representation
func (q *Question) DatabaseQuestion() DatabaseQuestion {
	return DatabaseQuestion{
		Prompt: q.Prompt,
		Tags: q.Tags,
		SearchTags: q.createSearchTags(),
	}
}
