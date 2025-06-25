package database

import (
	"ama/api/application"
	"ama/api/application/list"
	"ama/api/constants"

	"errors"
)

func (db *Database) ReadList(userId string, listId string, limit int, finalId string, tags []string) (list list.List, questions []application.Question, err error) {
	user, err := db.ReadUser(userId)
	if err != nil {
		return list, questions, err
	}
	if list, ok := user.GetList(listId); !ok {
		db.logger.Error("list not found", "userId", userId, "listId", listId)
		return list, questions, errors.New("list not found")
	}
	listRef := db.client.Collection(constants.UserCollection).Doc(userId).Collection(listId)
	questions, err = db.readQuestionCollection(listRef, limit, finalId, tags)
	if err != nil {
		return list, questions, err
	}
	return list, questions, nil
}
