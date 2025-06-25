package user

import (
	"time"
)

// Details about a user's subscription.
type UserSubscription struct {
	// How often a user makes payments for their subscription.
	// Currently "monthly" is the only supported value.
	PayCadence string `json:"payCadence" firestore:"payCadence"`
	// When the next payment is to be made to renew the subscription.
	RenewalDate time.Time `json:"renewalDate" firestore:"renewalDate" binding:"required"`
}
