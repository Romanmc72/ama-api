package user_test

import (
	"ama/api/application/user"
	"testing"
	"time"
)

func TestValidateUserSubscription(t *testing.T) {
	type testCase struct {
		name  string
		sub   user.UserSubscription
		valid bool
	}
	testCases := []testCase{
		{
			name: "Valid Subscription",
			sub: user.UserSubscription{
				PayCadence:  user.PayCadenceMonthly,
				RenewalDate: time.Now().Add(24 * time.Hour),
			},
			valid: true,
		},
		{
			name:  "Invalid Empty Subscription",
			sub:   user.UserSubscription{},
			valid: false,
		},
		{
			name: "Invalid No Cadence Subscription",
			sub: user.UserSubscription{
				PayCadence:  "",
				RenewalDate: time.Now().Add(24 * time.Hour),
			},
			valid: false,
		},
		{
			name: "Invalid Renewal Date Subscription",
			sub: user.UserSubscription{
				PayCadence:  "",
				RenewalDate: time.Now().Add(-24 * time.Hour),
			},
			valid: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := user.ValidateUserSubscription(tc.sub)
			if (err == nil) != tc.valid {
				t.Errorf("Expected validity: %v, got error: %v", tc.valid, err)
			}
		})
	}
}
