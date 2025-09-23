package list

import (
	"ama/api/application/errors"
	"strings"
)

const LikedQuestionsListName = "Liked questions"

// No significance behind this ID, just a hard coded uuid for every like question list
const LikedQuestionsListID = "84e37c9a-3ff3-4643-82fd-49977fb35fe8"

// A list of questions
type List struct {
	// The unique identifier for the list.
	ID string `json:"listId" firestore:"id" binding:"required"`
	// The human readable name for the list
	Name string `json:"name" firestore:"name" binding:"required"`
}

func GetLikedQuestionList() List {
	return List{
		ID:   LikedQuestionsListID,
		Name: LikedQuestionsListName,
	}
}

func (l *List) String() string {
	return "List(" +
		"Id=" + l.ID +
		", Name=" + l.Name +
		")"
}

// Validates a list of questions and returns if it has the liked question
// list and an error if the list of lists contains any errors
func ValidateListOfLists(lol []List) (bool, error) {
	errs := []string{}
	idMap := map[string]bool{}
	hasLikedQuestions := false
	for _, l := range lol {
		idMap[l.ID] = true
		// allow the list to have 1 liked question list
		if !hasLikedQuestions && l.Name == LikedQuestionsListName {
			hasLikedQuestions = true
			if err := validateListID(l.ID); err != nil {
				errs = append(errs, err.Error())
			}
			continue
		}
		err := ValidateList(l)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(lol) != len(idMap) {
		errs = append(errs, "duplicate IDs found in list")
	}
	if len(errs) != 0 {
		return hasLikedQuestions, errors.NewValidationError(errs)
	}
	return hasLikedQuestions, nil
}

// Ensures that a list's data is valid, if not an error is returned
func ValidateList(l List) error {
	errs := []string{}
	if err := validateListID(l.ID); err != nil {
		errs = append(errs, err.Error())
	}
	if strings.TrimSpace(l.Name) == "" {
		errs = append(errs, "Name cannot be blank")
	}
	if strings.TrimSpace(l.Name) == LikedQuestionsListName {
		errs = append(errs, `Cannot have 2 "`+LikedQuestionsListName+`" list names`)
	}
	if len(errs) != 0 {
		return errors.NewValidationError(errs)
	}
	return nil
}

func validateListID(ID string) error {
	if strings.TrimSpace(ID) == "" {
		return errors.NewValidationError([]string{"ID cannot be blank"})
	}
	return nil
}
