package application_test

import (
	"ama/api/application"
	"testing"
)

func TestValidateQuestion(t *testing.T) {
	testCases := []struct {
		name     string
		question application.Question
		valid    bool
	}{
		{
			name: "Valid Question",
			question: application.Question{
				Prompt: "What is the meaning of life?",
				Tags:   []string{"philosophy", "life"},
			},
			valid: true,
		},
		{
			name:     "Invalid Empty Question",
			question: application.Question{},
			valid:    false,
		},
		{
			name: "Invalid No Tags",
			question: application.Question{
				ID:     "q1",
				Prompt: "What is the meaning of life?",
			},
			valid: false,
		},
		{
			name: "Invalid Duplicate Tags",
			question: application.Question{
				Prompt: "What is the meaning of life?",
				Tags:   []string{"philosophy", "philosophy"},
			},
			valid: false,
		},
		{
			name: "Invalid No Prompt",
			question: application.Question{
				Tags: []string{"philosophy"},
			},
			valid: false,
		},
		{
			name: "Invalid Blank Prompt",
			question: application.Question{
				Prompt: "   ",
				Tags:   []string{"philosophy"},
			},
			valid: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := application.ValidateQuestion(tc.question)
			if (err == nil) != tc.valid {
				t.Errorf("Expected validity: %v, got error: %v", tc.valid, err)
			}
		})
	}
}
