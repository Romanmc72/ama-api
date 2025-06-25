package database

import (
	"ama/api/constants"
	"time"
)

// Delete a user and all of their lists from the database.
func (db *Database) DeleteUser(id string) (time.Time, error) {
	db.logger.Debug("Deleting this user", "id", id)
	user, err := db.ReadUser(id)
	if err != nil {
		return time.Now(), err
	}
	for _, list := range user.Lists {
		err = db.DeleteList(id, list.ID)
		if err != nil {
			db.logger.Error("Error deleting list", "user", id, "list", list.Name, "error", err)
			return time.Now(), err
		}
	}
	deleteResult, err := db.client.Collection(constants.UserCollection).Doc(id).Delete(db.ctx)
	if err != nil {
		db.logger.Error("Error deleting user", "user", id, "error", err)
		return time.Now(), err
	}
	return deleteResult.UpdateTime, nil
}
