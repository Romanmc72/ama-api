package application

import (
	"ama/api/application/list"
	"ama/api/application/user"
)

// TODO remove the firestore bindings from these objects if we have database objects, otherwise remove the database objects themselves
// TODO create a user with non-required bindings to enable updates more easily with partial user objects
// TODO create a validation function that can validate input user objects

// Describes the shape of a user profile within the application.
type User struct {
	// The unique identifier for the user.
	ID string `json:"id" firestore:"-" binding:"required"`
	user.BaseUser
}

func (u User) User() User {
	return u
}

func (u User) GetList(listId string) (list.List, bool) {
	for _, list := range u.Lists {
		if list.ID == listId {
			return list, true
		}
	}

	return list.List{}, false
}
