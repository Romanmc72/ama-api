package user

// The settings unique to this user's experience and profile.
type UserSettings struct {
	// The color scheme that the user has chose.
	ColorScheme UserColorScheme `json:"colorScheme" firestore:"colorScheme" binding:"required"`
}

func ValidateUserSettings(settings UserSettings) error {
	return ValidateUserColorScheme(settings.ColorScheme)
}
