package user

import (
	"errors"
	"strings"
)

const DefaultColor = "default"

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

func isHexColorCode(s string) bool {
	if len(s) != 7 || s[0] != '#' {
		return false
	}
	for _, c := range s[1:] {
		if (c < '0' || c > '9') && (c < 'A' || c > 'F') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

func ValidateUserColorScheme(scheme UserColorScheme) error {
	errorList := make([]string, 0)
	if scheme.Background != DefaultColor && !isHexColorCode(scheme.Background) {
		errorList = append(errorList, "background must be 'default' or a hex color code like '#RRGGBB'")
	}
	if scheme.Foreground != DefaultColor && !isHexColorCode(scheme.Foreground) {
		errorList = append(errorList, "background must be 'default' or a hex color code like '#RRGGBB'")
	}
	if scheme.HighlightedBackground != DefaultColor && !isHexColorCode(scheme.HighlightedBackground) {
		errorList = append(errorList, "background must be 'default' or a hex color code like '#RRGGBB'")
	}
	if scheme.HighlightedForeground != DefaultColor && !isHexColorCode(scheme.HighlightedForeground) {
		errorList = append(errorList, "background must be 'default' or a hex color code like '#RRGGBB'")
	}
	if len(errorList) > 0 {
		return errors.New("color scheme validation errors: " + strings.Join(errorList, "; "))
	}
	return nil
}

func GetDefaultUserColorScheme() UserColorScheme {
	return UserColorScheme{
		Background:            DefaultColor,
		Foreground:            DefaultColor,
		HighlightedBackground: DefaultColor,
		HighlightedForeground: DefaultColor,
	}
}
