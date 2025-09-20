package database

import (
	"ama/api/application/list"
	"ama/api/constants"
	"ama/api/interfaces"
	"errors"
	"strings"

	"google.golang.org/api/iterator"
)

// TODO This still created a list with a blank id
// (is this still a needed TODO?)

// "Creates" a list for a given user, without a question in that list, the list does not "exist" per se
// but it is a placeholder for the user to add questions to it
func (db *Database) CreateList(userId string, l list.List) (list.List, error) {
	if strings.Trim(userId, " ") == "" {
		return list.List{}, errors.New("user id cannot be empty")
	}
	user, err := db.ReadUser(userId)
	if err != nil {
		db.logger.Error(
			"error getting user document",
			"error", err,
			"userId", userId,
		)
		return list.List{}, err
	}
	if l.ID == "" {
		l.ID = db.client.NewID()
		db.logger.Debug(
			"list has no ID, generated a new one",
			"list", l,
			"userId", userId,
		)
	}
	if err := list.ValidateList(l); err != nil {
		db.logger.Error(
			"List invalid",
			"list", l,
			"userId", userId,
		)
		return list.List{}, err
	}
	db.logger.Debug(
		"Checking if list already exists",
		"listId", l.ID,
		"userId", userId,
	)
	userDocumentRef := db.client.
		Collection(constants.UserCollection).
		Doc(userId)
	exists, err := db.checkIfListExists(userDocumentRef, l.ID)
	if err != nil {
		db.logger.Error(
			"error checking if list exists",
			"error", err,
			"listId", l.ID,
			"userId", userId,
		)
		return list.List{}, err
	}
	if exists {
		db.logger.Debug(
			"list already exists",
			"listId", l.ID,
			"userId", userId,
		)
	} else {
		db.logger.Debug(
			"list does not exist, list creation is a no-op, doing nothing",
			"listId", l.ID,
			"userId", userId,
		)
	}
	for _, eachList := range user.Lists {
		if l.ID == eachList.ID {
			db.logger.Debug(
				"list already exists on user record",
				"listId", l.ID,
				"userId", userId,
			)
			return list.List{}, nil
		}
	}
	user.Lists = append(user.Lists, l)
	err = db.UpdateUser(user)
	if err != nil {
		return list.List{}, err
	}
	return l, err
}

func (db *Database) checkIfListExists(userDocRef interfaces.DocumentRef, listId string) (bool, error) {
	iter := userDocRef.
		Collection(listId).
		Limit(1).
		Documents(db.ctx)
	_, err := iter.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
