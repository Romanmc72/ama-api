package user_test

import (
	"ama/api/application/user"
	"testing"
)

func TestValidateUserColorScheme(t *testing.T) {
	type testCase struct {
		name   string
		scheme user.UserColorScheme
		valid  bool
	}
	testCases := []testCase{
		{
			name: "Invalid Empty Scheme",
			scheme: user.UserColorScheme{
				Background:            "",
				Foreground:            "",
				HighlightedBackground: "",
				HighlightedForeground: "",
			},
			valid: false,
		},
		{
			name:   "Default Scheme",
			scheme: user.GetDefaultUserColorScheme(),
			valid:  true,
		},
		{
			name: "Hex Codes",
			scheme: user.UserColorScheme{
				Background:            "#FFFFFF",
				Foreground:            "#f8f9f0",
				HighlightedBackground: "#abcdef",
				HighlightedForeground: "#123456",
			},
			valid: true,
		},
		{
			name: "Invalid Hex Codes",
			scheme: user.UserColorScheme{
				Background:            "#FFFgFF",
				Foreground:            "#fff",
				HighlightedBackground: "0000000",
				HighlightedForeground: "#zzzzzz",
			},
			valid: false,
		},
		{
			name: "Invalid Hex Codes",
			scheme: user.UserColorScheme{
				Background:            "#FFFgFF",
				Foreground:            "#",
				HighlightedBackground: "#000",
				HighlightedForeground: "#00z000",
			},
			valid: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := user.ValidateUserColorScheme(tc.scheme)
			if (err == nil) != tc.valid {
				t.Errorf("Expected valid=%v, got error=%v", tc.valid, err)
			}
		})
	}
}
