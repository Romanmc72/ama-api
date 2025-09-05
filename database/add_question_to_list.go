package database

import (
	"ama/api/application"
	"ama/api/constants"
	"errors"
)

func (db *Database) AddQuestionToList(userId string, listId string, question application.Question) error {
	if userId == "" || listId == "" || question.ID == "" {
		err := errors.New("userId, listId, and question.ID cannot be empty")
		db.logger.Error(err.Error(), "userId", userId, "listId", listId, "question", question)
		return err
	}
	db.logger.Debug("Adding question to list", "user", userId, "list", listId, "question", question)
	result, err := db.client.
		Collection(constants.UserCollection).
		Doc(userId).
		Collection(listId).
		Doc(question.ID).
		Set(db.ctx, question)
	if err != nil {
		db.logger.Error("Error adding question to list", "error", err, "user", userId, "list", listId, "question", question)
		return err
	}
	db.logger.Debug("Question added to list", "user", userId, "list", listId, "question", question, "result", result)
	return nil
}
