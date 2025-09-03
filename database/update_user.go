package database

import (
	"ama/api/application/user"
	"ama/api/constants"
	"ama/api/interfaces"
)

// Updates the user in the database with the provided user data.
func (db *Database) UpdateUser(userData interfaces.UserConverter) error {
	u := userData.User()
	err := user.ValidateUser(u.BaseUser)
	if err != nil {
		return err
	}
	db.logger.Debug("updating user", "user", u)
	result, err := db.client.Collection(constants.UserCollection).
		Doc(u.ID).
		Set(db.ctx, u)
	if err != nil {
		db.logger.Error("Error updating user", "error", err, "user", u.ID)
		return err
	}
	db.logger.Debug("User updated", "user", u.ID, "updatedAt", result.UpdateTime.Unix())
	return nil
}
