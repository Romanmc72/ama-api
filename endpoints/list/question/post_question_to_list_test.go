package question_test

import (
	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/list/question"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"net/http"
	"testing"
)

func TestPostQuestionToList(t *testing.T) {
	testCases := []struct {
		name     string
		db       *test.MockListManager
		ctx      *test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name: "Success",
			db: test.NewMockListManager(test.MockListManagerConfig{
				ReadQuestion: func(id string) (application.Question, error) {
					return fixtures.ValidQuestion, nil
				},
				AddQuestionToList: func(userId, listId string, question application.Question) error {
					return nil
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier:     fixtures.UserId,
					constants.ListIdPathIdentifier:     fixtures.ListId,
					constants.QuestionIdPathIdentifier: fixtures.QuestionId,
				},
			}),
			wantCode: http.StatusCreated,
			wantErr:  false,
		},
		{
			name: "Failure - Blank Path Identifiers",
			db:   test.NewMockListManager(test.MockListManagerConfig{}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier:     "     ",
					constants.ListIdPathIdentifier:     "     ",
					constants.QuestionIdPathIdentifier: "     ",
				},
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Read Error",
			db: test.NewMockListManager(test.MockListManagerConfig{
				ReadQuestion: func(id string) (application.Question, error) {
					return application.Question{}, errors.New("read error")
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier:     fixtures.UserId,
					constants.ListIdPathIdentifier:     fixtures.ListId,
					constants.QuestionIdPathIdentifier: fixtures.QuestionId,
				},
			}),
			wantCode: http.StatusNotFound,
			wantErr:  true,
		},
		{
			name: "Failure - Write Error",
			db: test.NewMockListManager(test.MockListManagerConfig{
				ReadQuestion: func(id string) (application.Question, error) {
					return fixtures.ValidQuestion, nil
				},
				AddQuestionToList: func(userId, listId string, question application.Question) error {
					return errors.New("could not write to list")
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier:     fixtures.UserId,
					constants.ListIdPathIdentifier:     fixtures.ListId,
					constants.QuestionIdPathIdentifier: fixtures.QuestionId,
				},
			}),
			wantCode: http.StatusInternalServerError,
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			question.PostQuestionToList(tc.ctx, tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("PostQuestionToList() wanted code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); !ok && tc.wantErr && tc.ctx.ResponseData != nil {
				t.Errorf("PostQuestionToList() wanted error Response, got = %v", tc.ctx.ResponseData)
			}
		})
	}
}
