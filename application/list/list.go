package list

import (
	"errors"
	"strings"
)

const LikedQuestionsListName = "Liked questions"

// A list of questions
type List struct {
	// The unique identifier for the list.
	ID string `json:"id" firestore:"id" binding:"required"`
	// The human readable name for the list
	Name string `json:"name" firestore:"name" binding:"required"`
}

func (l *List) String() string {
	return "List(" +
		"Id=" + l.ID +
		", Name=" + l.Name +
		")"
}

func ValidateList(l List) error {
	errs := []string{}
	if strings.TrimSpace(l.ID) == "" {
		errs = append(errs, "ID cannot be blank")
	}
	if strings.TrimSpace(l.Name) == "" {
		errs = append(errs, "Name cannot be blank")
	}
	if strings.TrimSpace(l.Name) == LikedQuestionsListName {
		errs = append(errs, `Cannot have 2 "`+LikedQuestionsListName+`" list names`)
	}
	if len(errs) != 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}
