package database

import (
	"ama/api/application"
	"ama/api/interfaces"
)

// TODO: Write this
func (db *Database) CreateUser(userData interfaces.UserConverter) application.User {
	user := userData.User()
	return user
}
