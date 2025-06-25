package user

// A color scheme that a user prefers for their app.
type UserColorScheme struct {
	// The hex code for the color to use in the background OR the word "default"
	// to use the default setting on the user's device.
	Background string `json:"background" firestore:"background" binding:"required"`
	// The hex code for the color to use in the foreground OR the word "default"
	// to use the default setting on the user's device.
	Foreground string `json:"foreground" firestore:"foreground" binding:"required"`
	// The hex code for the color to use in the highlighted background OR the
	// word "default" to use the default setting on the user's device.
	HighlightedBackground string `json:"highlightedBackground" firestore:"highlightedBackground" binding:"required"`
	// The hex code for the color to use in the highlighted background OR the
	// word "default" to use the default setting on the user's device.
	HighlightedForeground string `json:"highlightedForeground" firestore:"highlightedForeground" binding:"required"`
}
