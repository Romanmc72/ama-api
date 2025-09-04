package database

import (
	"ama/api/application"
	"ama/api/application/user"
	"ama/api/constants"
)

func (db *Database) ReadUser(id string) (u application.User, err error) {
	doc, err := db.client.Collection(constants.UserCollection).Doc(id).Get(db.ctx)
	if err != nil {
		return u, err
	}
	var baseUser user.BaseUser
	err = doc.DataTo(&baseUser)
	if err != nil {
		return u, err
	}
	u.ID = id
	u.BaseUser = baseUser
	return u, nil
}
