package database

import (
	"ama/api/application"
	"ama/api/constants"
)

// Retrieve a particular question from the database
func (db *Database) ReadQuestion(id string) (q application.Question, err error) {
	collection := db.client.Collection(constants.QuestionCollection)
	return db.readQuestionFromCollection(collection, id)
}
