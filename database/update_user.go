package database

import (
	"ama/api/application"
	"ama/api/interfaces"
)

// TODO: Write this
func (db * Database) UpdateUser(userData interfaces.UserConverter) application.User {
	return userData.User()
}
