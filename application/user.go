package application

import (
	"ama/api/application/list"
	"ama/api/application/user"
	"encoding/json"
)

// Describes the shape of a user profile within the application.
type User struct {
	// The unique identifier for the user.
	ID string `json:"userId" firestore:"-" binding:"required"`
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

func (u User) String() string {
	v, err := json.Marshal(u)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	return string(v)
}
