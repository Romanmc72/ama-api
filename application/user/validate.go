package user

import (
	"ama/api/application/list"
	"errors"
	"slices"
	"strings"
)

func ValidateUser(user BaseUser) error {
	errorList := make([]string, 0)
	if len(user.Lists) == 0 {
		errorList = append(errorList, `user "lists" field is required and cannot be empty`)
	}
	hasLikedQuestions := false
	likeQuestionLists := 0
	listMap := map[string]bool{}
	for _, l := range user.Lists {
		listMap[l.ID] = true
		if l.Name == list.LikedQuestionsListName {
			hasLikedQuestions = true
			likeQuestionLists += 1
		}
	}
	if !hasLikedQuestions {
		errorList = append(errorList, "missing '"+list.LikedQuestionsListName+"' List")
	}
	if likeQuestionLists > 1 {
		errorList = append(errorList, "can only have 1 '"+list.LikedQuestionsListName+"' List")
	}
	if len(listMap) < len(user.Lists) {
		errorList = append(errorList, "there are lists with duplicate IDs")
	}
	emailParts := strings.Split(user.Email, "@")
	if len(user.Email) < 6 || len(emailParts) != 2 || len(emailParts[0]) == 0 || len(emailParts[1]) < 4 || !strings.Contains(emailParts[1], ".") {
		errorList = append(errorList, `user "email" field is required and cannot be empty and must be a valid email address on the web`)
	}
	if user.FirebaseID == "" {
		errorList = append(errorList, `user "firebaseId" field is required and cannot be empty`)
	}
	tiers := Tiers()
	if !slices.Contains(tiers, user.Tier) {
		errorList = append(errorList, `user "tier" field must be one of [`+strings.Join(tiers, ", ")+"]")
	}
	if err := ValidateUserSettings(user.Settings); err != nil {
		errorList = append(errorList, err.Error())
	}
	if err := ValidateUserSubscription(user.Subscription); err != nil {
		errorList = append(errorList, err.Error())
	}
	if len(errorList) > 0 {
		return errors.New("validation errors: " + strings.Join(errorList, "; "))
	}
	return nil
}
