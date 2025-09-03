package fixtures

import (
	"ama/api/application"
	"ama/api/application/list"
	"ama/api/application/user"
	"time"
)

var ValidUserSubscription = user.UserSubscription{
	RenewalDate: time.Now().AddDate(0, 1, 0),
	PayCadence:  user.PayCadenceMonthly,
}

var ValidUserSettings = user.UserSettings{
	ColorScheme: user.GetDefaultUserColorScheme(),
}

var ValidLists = []list.List{
	{
		ID:   "list1",
		Name: "List 1",
	},
}

var BaseUserValid = user.BaseUser{
	FirebaseID:   "valid-firebase-id",
	Email:        "test@test.com",
	Tier:         user.FreeTier,
	Subscription: ValidUserSubscription,
	Settings:     ValidUserSettings,
	Lists:        ValidLists,
}

var UserValid = application.User{
	ID:       "valid-user-id",
	BaseUser: BaseUserValid,
}
