package database

import (
	"ama/api/constants"
	"ama/api/interfaces"
)

// Updates the user in the database with the provided user data.
func (db *Database) UpdateUser(userData interfaces.UserConverter) error {
	user := userData.User()

	db.logger.Debug("updating user", "user", user)
	result, err := db.client.Collection(constants.UserCollection).
		Doc(user.ID).
		Set(db.ctx, user)
	if err != nil {
		db.logger.Error("Error updating user", "error", err, "user", user.ID)
		return err
	}
	db.logger.Debug("User updated", "user", user.ID, "updatedAt", result.UpdateTime.Unix())
	return nil
}
