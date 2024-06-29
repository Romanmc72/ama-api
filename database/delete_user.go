package database

import (
	"time"
)

// TODO: Write this
func (db *Database) DeleteUser(id string) (time.Time, error) {
	db.logger.Debug("Deleting this user", "id", id)
	return time.Now(), nil
}
