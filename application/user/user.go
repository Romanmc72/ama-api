package user

import (
	"ama/api/application/list"
)

// The base user object, without an ID. Can be used for various implementations
// of the user without duplicating struct types or creating import cycles.
type BaseUser struct {
	// The unique user identifier from firebase.
	FirebaseID string `json:"firebaseId" firestore:"firebaseId" binding:"required"`
	// The user's username.
	Name string `json:"name" firestore:"name"`
	// The user's username.
	Email string `json:"email" firestore:"email" binding:"required"`
	// The user's subscription tier. One of "free" | "lite" | "premium"
	Tier string `json:"tier" firestore:"tier" binding:"required"`
	// The user's subscription information.
	Subscription UserSubscription `json:"subscription" firestore:"subscription"`
	// The users settings and preferences.
	Settings UserSettings `json:"settings" firestore:"settings" binding:"required"`
	// The list of question lists that the user has created.
	Lists []list.List `json:"lists" firestore:"lists" binding:"required"`
}
