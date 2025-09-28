package database

import (
	"ama/api/application"
	"ama/api/constants"
)

// Retrieve a particular question from a list in the database
func (db *Database) ReadListQuestion(userId string, listId string, questionId string) (q application.Question, err error) {
	collection := db.client.Collection(constants.UserCollection).Doc(userId).Collection(listId)
	return db.readQuestionFromCollection(collection, questionId)
}
