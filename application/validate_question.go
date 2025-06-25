package application

import (
	"errors"
	"slices"
	"strings"
)

// Apply server side validation to the question data prior to writing it to the database.
func ValidateQuestion(q Question) error {
	var errorBuilder strings.Builder
	if strings.Join(q.Tags, "") == "" {
		errorBuilder.WriteString(`question "tags" field is required to have at least 1 tag`)
	}
	slices.Sort(q.Tags)
	dedupedTags := slices.Compact(q.Tags)
	if len(q.Tags) != len(dedupedTags) {
		errorBuilder.WriteString(`question "tags" cannot have duplicate tag values`)
	}
	if strings.TrimSpace(q.Prompt) == "" {
		errorBuilder.WriteString(`question "prompt" field cannot be blank`)
	}
	errorMessage := errorBuilder.String()
	if errorMessage != "" {
		return errors.New(errorMessage)
	}
	return nil
}
