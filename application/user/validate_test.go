package user_test

import (
	"ama/api/application/list"
	"ama/api/application/user"
	"ama/api/test/fixtures"
	"testing"
)

func TestValidateUser(t *testing.T) {
	type testCase struct {
		name  string
		user  user.BaseUser
		valid bool
	}
	testCases := []testCase{
		{
			name: "Valid Settings",
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "test@test.com",
				Tier:         user.FreeTier,
				Lists:        fixtures.ValidLists,
				Subscription: fixtures.ValidUserSubscription,
				Settings:     fixtures.ValidUserSettings,
			},
			valid: true,
		},
		{
			name:  "Invalid Empty User",
			user:  user.BaseUser{},
			valid: false,
		},
		{
			name: "Invalid Tier",
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "test@test.com",
				Tier:         "applesauce",
				Lists:        fixtures.ValidLists,
				Subscription: fixtures.ValidUserSubscription,
				Settings:     fixtures.ValidUserSettings,
			},
			valid: false,
		},
		{
			name: "Invalid Lists",
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "test@test.com",
				Tier:         user.FreeTier,
				Lists:        []list.List{},
				Subscription: fixtures.ValidUserSubscription,
				Settings:     fixtures.ValidUserSettings,
			},
			valid: false,
		},
		{
			name: "Invalid Empty Email",
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "",
				Tier:         user.FreeTier,
				Lists:        fixtures.ValidLists,
				Subscription: fixtures.ValidUserSubscription,
				Settings:     fixtures.ValidUserSettings,
			},
			valid: false,
		},
		{
			name: "Invalid Email Too Short",
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "t@t.c",
				Tier:         user.FreeTier,
				Lists:        fixtures.ValidLists,
				Subscription: fixtures.ValidUserSubscription,
				Settings:     fixtures.ValidUserSettings,
			},
			valid: false,
		},
		{
			name: "Invalid Email Missing @",
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "testtest.com",
				Tier:         user.FreeTier,
				Lists:        fixtures.ValidLists,
				Subscription: fixtures.ValidUserSubscription,
				Settings:     fixtures.ValidUserSettings,
			},
			valid: false,
		},
		{
			name: "Invalid Email Domain Too Short",
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "test@t.c",
				Tier:         user.FreeTier,
				Lists:        fixtures.ValidLists,
				Subscription: fixtures.ValidUserSubscription,
				Settings:     fixtures.ValidUserSettings,
			},
			valid: false,
		},
		{
			name: `Invalid Email Domain Missing "."`,
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "test@testcom",
				Tier:         user.FreeTier,
				Lists:        fixtures.ValidLists,
				Subscription: fixtures.ValidUserSubscription,
				Settings:     fixtures.ValidUserSettings,
			},
			valid: false,
		},
		{
			name: "Invalid Email Name Too Short",
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "@test.com",
				Tier:         user.FreeTier,
				Lists:        fixtures.ValidLists,
				Subscription: fixtures.ValidUserSubscription,
				Settings:     fixtures.ValidUserSettings,
			},
			valid: false,
		},
		{
			name: "Invalid Blank FirebaseID",
			user: user.BaseUser{
				FirebaseID:   "",
				Email:        "test@test.com",
				Tier:         user.FreeTier,
				Lists:        fixtures.ValidLists,
				Subscription: fixtures.ValidUserSubscription,
				Settings:     fixtures.ValidUserSettings,
			},
			valid: false,
		},
		{
			name: "Invalid Empty Subscription",
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "test@test.com",
				Tier:         user.FreeTier,
				Lists:        fixtures.ValidLists,
				Subscription: user.UserSubscription{},
				Settings:     fixtures.ValidUserSettings,
			},
			valid: false,
		},
		{
			name: "Invalid Empty Settings",
			user: user.BaseUser{
				FirebaseID:   "firebase123",
				Email:        "test@test.com",
				Tier:         user.FreeTier,
				Lists:        fixtures.ValidLists,
				Subscription: fixtures.ValidUserSubscription,
				Settings:     user.UserSettings{},
			},
			valid: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := user.ValidateUser(tc.user)
			if (err == nil) != tc.valid {
				t.Errorf("Expected validity: %v, got error: %v", tc.valid, err)
			}
		})
	}
}
