package database

import (
	"ama/api/application"
	"ama/api/application/list"
	"ama/api/constants"

	"errors"
)

func (db *Database) ReadList(userId string, listId string, limit int, finalId string, tags []string) (l list.List, questions []application.Question, err error) {
	user, err := db.ReadUser(userId)
	if err != nil {
		return l, questions, err
	}
	if userList, ok := user.GetList(listId); !ok {
		db.logger.Error("list not found", "userId", userId, "listId", listId)
		return l, questions, errors.New("list not found")
	} else {
		l = userList
	}
	listRef := db.client.Collection(constants.UserCollection).Doc(userId).Collection(listId)
	questions, err = db.readQuestionCollection(listRef, limit, finalId, tags)
	if err != nil {
		return l, questions, err
	}
	return l, questions, nil
}
