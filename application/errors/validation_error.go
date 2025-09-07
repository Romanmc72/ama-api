package errors

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	msg             string
	ValidationErrCt int
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("There are %d validation errors: %s", v.ValidationErrCt, v.msg)
}

// Given an array of validation error messages, return the validation error
func NewValidationError(validationErrMsgs []string) *ValidationError {
	return &ValidationError{
		msg:             strings.Join(validationErrMsgs, "; "),
		ValidationErrCt: len(validationErrMsgs),
	}
}
