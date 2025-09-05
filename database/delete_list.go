package database

import (
	"ama/api/application"
	"ama/api/application/list"
	"ama/api/constants"
	"context"

	"cloud.google.com/go/firestore"
)

// TODO: prevent deleting the "Liked questions" list
// Erase a list from a user in the database
func (db *Database) DeleteList(userId string, listId string) error {
	userDocRef := db.client.Collection(constants.UserCollection).Doc(userId)
	listCollection := userDocRef.Collection(listId)
	err := db.deleteCollection(listCollection, 500)
	if err != nil {
		return err
	}
	// yay race transactions! The update operation to remove from an array is not
	// supported in golang firestore client. So we have to do a read, then
	// delete the list by updating the field. A transaction is needed to prevent
	// race conditions.
	var user application.User
	return db.client.RunTransaction(db.ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		document, err := tx.Get(userDocRef.Ref())
		if err != nil {
			return err
		}
		err = document.DataTo(&user)
		if err != nil {
			return err
		}
		newList := []list.List{}
		for _, l := range user.Lists {
			if l.ID == listId {
				continue
			}
			newList = append(newList, l)
		}
		return tx.Set(
			userDocRef.Ref(),
			map[string][]list.List{"lists": newList},
			firestore.MergeAll,
		)
	})
}
