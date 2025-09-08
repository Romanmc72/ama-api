package list_test

import (
	"ama/api/application/list"
	"ama/api/test/fixtures"
	"testing"
)

func TestValidateList(t *testing.T) {
	testCases := []struct {
		name    string
		l       list.List
		wantErr bool
	}{
		{
			name:    "Success",
			l:       fixtures.ValidLists[0],
			wantErr: false,
		},
		{
			name: "Failure - Blank ID",
			l: list.List{
				ID:   "   ",
				Name: "something",
			},
			wantErr: true,
		},
		{
			name: "Failure - Blank Name",
			l: list.List{
				ID:   fixtures.ListId,
				Name: "  ",
			},
			wantErr: true,
		},
		{
			name: "Failure - Is Liked Question List",
			l: list.List{
				ID:   fixtures.ListId,
				Name: list.LikedQuestionsListName,
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := list.ValidateList(tc.l)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateList() wantedErr = %v; got = %v", tc.wantErr, err)
			}
		})
	}
}

func TestValidateListOfLists(t *testing.T) {
	testCases := []struct {
		name        string
		lol         []list.List
		hasLikeList bool
		wantErr     bool
	}{
		{
			name:        "Success",
			lol:         fixtures.ValidLists,
			hasLikeList: true,
			wantErr:     false,
		},
		{
			name: "Failure - Blank ID",
			lol: []list.List{{
				ID:   "   ",
				Name: "something",
			}},
			hasLikeList: false,
			wantErr:     true,
		},
		{
			name: "Failure - Liked Questions Blank ID",
			lol: []list.List{{
				ID:   "     ",
				Name: list.LikedQuestionsListName,
			}},
			hasLikeList: true,
			wantErr:     true,
		},
		{
			name: "Failure - Blank Name",
			lol: []list.List{{
				ID:   fixtures.ListId,
				Name: "  ",
			}},
			hasLikeList: false,
			wantErr:     true,
		},
		{
			name: "Failure - Duplicate Liked Question List",
			lol: []list.List{
				{
					ID:   fixtures.ListId,
					Name: list.LikedQuestionsListName,
				},
				list.GetLikedQuestionList(),
			},
			hasLikeList: true,
			wantErr:     true,
		},
		{
			name: "Failure - Duplicate List IDs",
			lol: []list.List{
				{
					ID:   fixtures.ListId,
					Name: "List 1",
				},
				{
					ID:   fixtures.ListId,
					Name: "List 2",
				},
			},
			hasLikeList: false,
			wantErr:     true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hasLikedQuestions, err := list.ValidateListOfLists(tc.lol)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateList() wantedErr = %v; got = %v", tc.wantErr, err)
			}
			if hasLikedQuestions != tc.hasLikeList {
				t.Errorf("ValidateList() hasLIkedQuestions want = %v; got = %v", tc.hasLikeList, hasLikedQuestions)
			}
		})
	}
}

func TestListString(t *testing.T) {
	name := "test list"
	input := list.List{
		ID:   fixtures.ListId,
		Name: name,
	}
	got := input.String()
	want := "List(Id=" + fixtures.ListId + ", Name=" + name + ")"
	if got != want {
		t.Errorf("List.String() wanted = %s; got = %s", want, got)
	}
}
