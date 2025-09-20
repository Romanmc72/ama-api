package requests

import (
	"ama/api/application/errors"
	"ama/api/application/list"
	"ama/api/application/user"
)

type PostUserRequest struct {
	// The unique user identifier from firebase.
	FirebaseID string `json:"firebaseId" firestore:"firebaseId" binding:"required" validate:"required"`
	// The user's username.
	Name string `json:"name" firestore:"name"`
	// The user's username.
	Email string `json:"email" firestore:"email" binding:"required" validate:"required"`
	// The user's subscription tier. One of "free" | "lite" | "premium"
	Tier string `json:"tier" firestore:"tier" binding:"required" validate:"required"`
	// The user's subscription information.
	Subscription user.UserSubscription `json:"subscription" firestore:"subscription" binding:"required" validate:"required"`
	// The users settings and preferences.
	Settings user.UserSettings `json:"settings" firestore:"settings" binding:"required" validate:"required"`
	// The list of question lists that the user has created.
	Lists []list.List `json:"lists" firestore:"lists"`
}

// Converts the Post request to a base user object
// or returns an error if it fails validation
func (p *PostUserRequest) BaseUser() (user.BaseUser, error) {
	errs := []string{}
	lists := p.Lists
	hasLikedQuestions, err := list.ValidateListOfLists(p.Lists)
	if err != nil {
		errs = append(errs, err.Error())
	}
	if !hasLikedQuestions {
		lists = append(
			lists,
			list.List{
				ID:   list.LikedQuestionsListID,
				Name: list.LikedQuestionsListName,
			},
		)
	}
	u := user.BaseUser{
		FirebaseID:   p.FirebaseID,
		Email:        p.Email,
		Name:         p.Name,
		Tier:         p.Tier,
		Settings:     p.Settings,
		Subscription: p.Subscription,
		Lists:        lists,
	}
	err = user.ValidateUser(u)
	if err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) != 0 {
		return user.BaseUser{}, errors.NewValidationError(errs)
	}
	return u, nil
}
