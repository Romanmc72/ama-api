package database

import (
	"ama/api/application"
	"ama/api/firestoreobjects"
	"ama/api/interfaces"
)

// Change a question on the backend
func (db *Database) UpdateQuestion(id string, questionData interfaces.QuestionConverter) (application.Question, error) {
	question := questionData.Question(id)
	err := application.ValidateQuestion(question)
	if err != nil {
		return application.Question{}, err
	}
	writeResult, err := db.client.
		Collection(firestoreobjects.QuestionCollection).
		Doc(question.ID).
		Set(db.ctx, question)
	if err != nil {
		db.logger.Error("Encountered an error updating that question", "question", err)
		return application.Question{}, err
	}
	db.logger.Debug("Write succeeded", "message", writeResult)
	return question, nil
}
