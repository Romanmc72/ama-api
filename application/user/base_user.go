package user

import (
	"ama/api/application/list"
)

const (
	// You like the app but don't want all of the features and you don't mind ads
	FreeTier      = "free"
	// You like the app and want all of the features, but don't mind ads
	LiteTier      = "lite"
	// You like the app and want all of the features without ads
	PremiumTier   = "premium"
	// Looks like you missed a payment lmao, now you're temporarily in free tier
	// til you pay up or downgrade
	SuspendedTier = "suspended"
)

// Get the list of valid user tiers.
func Tiers() []string {
	return []string{FreeTier, LiteTier, PremiumTier, SuspendedTier}
}

// The base user object, without an ID. Can be used for various implementations
// of the user without duplicating struct types or creating import cycles.
type BaseUser struct {
	// The unique user identifier from firebase.
	FirebaseID string `json:"firebaseId" firestore:"firebaseId" binding:"required" validate:"required"`
	// The user's username.
	Name string `json:"name" firestore:"name"`
	// The user's username.
	Email string `json:"email" firestore:"email" binding:"required" validate:"required"`
	// The user's subscription tier. One of "free" | "lite" | "premium"
	Tier string `json:"tier" firestore:"tier" binding:"required" validate:"required"`
	// The user's subscription information.
	Subscription UserSubscription `json:"subscription" firestore:"subscription" binding:"required" validate:"required"`
	// The users settings and preferences.
	Settings UserSettings `json:"settings" firestore:"settings" binding:"required" validate:"required"`
	// The list of question lists that the user has created.
	Lists []list.List `json:"lists" firestore:"lists" binding:"required" validate:"required"`
}
