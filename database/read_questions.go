package database

import (
	"ama/api/application"
	"ama/api/constants"
)

// Retrieve all of the questions from the database one page at a time.
// limit = the number of items per page. 0 or less will not paginate but return
// an iterator for all documents.
// finalId = the Id field of the last document on the previous page. It will
// start the next query off after that one.
// tag = the tag to search for within the set of questions.
// tags = the set of tags to search for an inclusive match of.
func (db *Database) ReadQuestions(limit int, finalId string, tags []string) ([]application.Question, error) {
	collection := db.client.Collection(constants.QuestionCollection)
	return db.readQuestionCollection(collection, limit, finalId, tags)
}
