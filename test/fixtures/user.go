package fixtures

import (
	"ama/api/application"
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

var ValidBaseUser = user.BaseUser{
	FirebaseID:   UserId,
	Email:        "test@test.com",
	Tier:         user.FreeTier,
	Subscription: ValidUserSubscription,
	Settings:     ValidUserSettings,
	Lists:        ValidLists,
}

var ValidUser = application.User{
	ID:       UserId,
	BaseUser: ValidBaseUser,
}
