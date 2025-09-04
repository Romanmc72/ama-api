package database

import (
	"ama/api/application"
	"ama/api/constants"
	"ama/api/interfaces"
)

// Given the input data, create a new question in the database.
func (db *Database) CreateQuestion(questionData interfaces.QuestionConverter) (application.Question, error) {
	question := questionData.Question("")
	err := application.ValidateQuestion(question)
	if err != nil {
		return application.Question{}, err
	}
	databaseQuestion := question.DatabaseQuestion()
	docRef, writeResult, err := db.client.
		Collection(constants.QuestionCollection).
		Add(db.ctx, databaseQuestion)
	if err != nil {
		db.logger.Error("Encountered an error writing that question", "question", err)
		return application.Question{}, err
	}
	question = questionData.Question(docRef.ID())
	db.logger.Debug("Write succeeded", "message", writeResult)
	return question, nil
}
