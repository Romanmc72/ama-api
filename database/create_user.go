package database

import (
	"ama/api/application"
	"ama/api/application/list"
	"ama/api/application/user"
	"ama/api/constants"
)

// Creates a new user in the database given the input data and returns the user
func (db *Database) CreateUser(userData user.BaseUser) (user application.User, err error) {
	user.BaseUser = userData
	// TODO: Get the user id from the firebase auth token instead of in the payload
	// also verify the auth token (if possible verify it upstream in an api gateway)
	writeResult, err := db.client.
		Collection(constants.UserCollection).
		Doc(user.FirebaseID).
		Set(db.ctx, userData)
	if err != nil {
		db.logger.Error("Error creating user", "error", err)
		return user, err
	}
	user.ID = user.FirebaseID
	db.logger.Debug("write succeeded creating user", "message", writeResult, "userId", user.ID)
	if len(user.Lists) > 0 {
		for _, l := range user.Lists {
			err = db.CreateList(user.ID, l)
			if err != nil {
				db.logger.Error("Error creating list during user creation", "error", err)
				return user, err
			}
		}
	} else {
		err = db.CreateList(user.ID, list.List{Name: "Liked Questions"})
		if err != nil {
			db.logger.Error("Error creating default list during user creation", "error", err)
			return user, err
		}
	}
	return user, nil
}
