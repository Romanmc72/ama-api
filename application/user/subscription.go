package user

import (
	"errors"
	"strings"
	"time"
)

const (
	PayCadenceMonthly = "monthly"
)

// Details about a user's subscription.
type UserSubscription struct {
	// How often a user makes payments for their subscription.
	// Currently "monthly" is the only supported value.
	PayCadence string `json:"payCadence" firestore:"payCadence"`
	// When the next payment is to be made to renew the subscription.
	RenewalDate time.Time `json:"renewalDate" firestore:"renewalDate" binding:"required"`
}

func ValidateUserSubscription(sub UserSubscription) error {
	errorList := make([]string, 0)
	if sub.PayCadence != PayCadenceMonthly {
		errorList = append(errorList, `subscription "payCadence" field is required and cannot be empty`)
	}
	if sub.RenewalDate.IsZero() || sub.RenewalDate.Before(time.Now()) {
		errorList = append(errorList, `subscription "renewalDate" field is required and must be in the future`)
	}
	if len(errorList) > 0 {
		return errors.New("subscription validation errors: " + strings.Join(errorList, "; "))
	}
	return nil
}
