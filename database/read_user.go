package database

import (
	"ama/api/application"
	"ama/api/application/user"
	"ama/api/constants"
)

// TODO: Write this
func (db *Database) ReadUser(id string) (application.User, error) {
	doc, err := db.client.Collection(constants.UserCollection).Doc(id).Get(db.ctx)
	if err != nil {
		return application.User{}, err
	}
	var baseUser user.BaseUser
	err = doc.DataTo(&baseUser)
	if err != nil {
		return application.User{}, err
	}
	return application.User{
		ID:       id,
		BaseUser: baseUser,
	}, nil
}
