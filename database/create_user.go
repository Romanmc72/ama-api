package database

import (
	"ama/api/application"
	"ama/api/application/list"
	"ama/api/application/user"
	"ama/api/constants"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Creates a new user in the database given the input data and returns the user
func (db *Database) CreateUser(userData user.BaseUser) (user application.User, err error) {
	// TODO: Get the user id from the firebase auth token instead of in the payload
	// also verify the auth token (if possible verify it upstream in an api gateway)
	docRef := db.client.Collection(constants.UserCollection).Doc(userData.FirebaseID)
	docSnap, err := docRef.Get(db.ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			db.logger.Error("Error checking if user exists", "error", err)
			return user, err
		}
	} else {
		if docSnap.Exists() {
			err = errors.New("user already exists")
			db.logger.Error("Error creating user", "error", err)
			return user, err
		}
	}
	hasLikedQuestions := false
	for _, l := range userData.Lists {
		if l.Name == list.LikedQuestionsListName {
			hasLikedQuestions = true
			break
		}
	}
	if len(userData.Lists) == 0 || !hasLikedQuestions {
		userData.Lists = append(userData.Lists, list.List{
			ID:   db.client.NewID(),
			Name: list.LikedQuestionsListName,
		})
	}
	writeResult, err := docRef.Set(db.ctx, userData)
	if err != nil {
		db.logger.Error("Error creating user", "error", err)
		return user, err
	}
	user.BaseUser = userData
	user.ID = user.FirebaseID
	db.logger.Debug("write succeeded creating user", "message", writeResult, "userId", user.ID)
	return user, nil
}
