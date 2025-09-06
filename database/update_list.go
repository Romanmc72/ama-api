package database

import (
	"ama/api/application"
	"ama/api/application/list"
	"ama/api/constants"
	"context"
	"errors"

	"cloud.google.com/go/firestore"
)

// TODO: prevent updating the name of the "Liked questions" list
// Update the attributes of a list for a user
func (db *Database) UpdateList(userId string, updatedList list.List) error {
	if updatedList.ID == "" || userId == "" {
		return errors.New("list id and user id cannot be empty")
	}
	db.logger.Debug("Updating list", "user", userId, "list", updatedList.ID)
	return db.client.RunTransaction(db.ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		userDocRef := db.client.Collection(constants.UserCollection).Doc(userId)
		userDoc, err := tx.Get(userDocRef.Ref())
		if err != nil {
			db.logger.Error("Error getting user", "error", err, "user", userId)
			return err
		}
		var user application.User
		err = userDoc.DataTo(user)
		if err != nil {
			db.logger.Error("Error converting user doc to user", "error", err, "user", userId)
			return err
		}
		var existingList *list.List
		for _, l := range user.Lists {
			if l.ID == updatedList.ID {
				existingList = &updatedList
				break
			}
		}
		if existingList == nil {
			return errors.New("list not for user found")
		}
		if err := list.ValidateList(*existingList); err != nil {
			db.logger.Error("This list cannot be changed", "error", err, "user", userId, "list", existingList)
			return err
		}
		if err := list.ValidateList(updatedList); err != nil {
			db.logger.Error("This updated list is invalid", "error", err, "user", userId, "list", updatedList)
			return err
		}
		return tx.Set(userDocRef.Ref(), user)
	})
}
