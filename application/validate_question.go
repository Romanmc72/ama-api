package application

import (
	"ama/api/application/errors"
	"slices"
	"strings"
)

// Apply server side validation to the question data prior to writing it to the database.
func ValidateQuestion(q Question) error {
	errs := []string{}
	if strings.Join(q.Tags, "") == "" {
		errs = append(errs, `question "tags" field is required to have at least 1 tag`)
	}
	slices.Sort(q.Tags)
	dedupedTags := slices.Compact(q.Tags)
	if len(q.Tags) != len(dedupedTags) {
		errs = append(errs, `question "tags" cannot have duplicate tag values`)
	}
	if strings.TrimSpace(q.Prompt) == "" {
		errs = append(errs, `question "prompt" field cannot be blank`)
	}
	if len(errs) != 0 {
		return errors.NewValidationError(errs)
	}
	return nil
}
