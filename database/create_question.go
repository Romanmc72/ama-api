package database

import (
	"ama/api/application"
	"ama/api/firestoreobjects"
	"ama/api/interfaces"
)

// Given the input data, create a new question in the database.
func (db *Database) CreateQuestion(questionData interfaces.QuestionConverter) (application.Question, error) {
	err := application.ValidateQuestion(questionData.Question(""))
	if err != nil {
		return application.Question{}, err
	}
	// questionId, err := db.incrementIdentifier(firestoreobjects.QuestionIdDoc)
	// if err != nil {
	// 	db.logger.Error("Encountered an error converting that data to a question", "question", err)
	// 	return application.Question{}, err
	// }
	question := questionData.Question("")
	databaseQuestion := question.DatabaseQuestion()
	docRef, writeResult, err := db.client.
		Collection(firestoreobjects.QuestionCollection).
		Add(db.ctx, databaseQuestion)
	question = questionData.Question(docRef.ID)
	if err != nil {
		db.logger.Error("Encountered an error writing that question", "question", err)
		return application.Question{}, err
	}
	db.logger.Debug("Write succeeded", "message", writeResult)
	return question, nil
}
