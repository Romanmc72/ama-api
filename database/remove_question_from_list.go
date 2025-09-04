package database

import (
	"ama/api/constants"
	"errors"
)

// Given a user, a list, and a question, remove that question from the list
func (db *Database) RemoveQuestionFromList(userId string, listId string, questionId string) error {
	db.logger.Debug("Removing question from list", "user", userId, "list", listId, "question", questionId)
	if userId == "" || listId == "" || questionId == "" {
		db.logger.Error("missing required parameters", "user", userId, "list", listId, "question", questionId)
		return errors.New("missing required parameter")
	}
	result, err := db.client.
		Collection(constants.UserCollection).
		Doc(userId).
		Collection(listId).
		Doc(questionId).
		Delete(db.ctx)
	if err != nil {
		db.logger.Error("Error removing question from list", "error", err, "user", userId, "list", listId, "question", questionId)
		return err
	}
	db.logger.Debug("Question removed from list", "user", userId, "list", listId, "question", questionId, "removedAt", result.UpdateTime.Unix())
	return nil
}
