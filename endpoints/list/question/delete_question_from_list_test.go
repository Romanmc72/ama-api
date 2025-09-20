package question_test

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/list/question"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"net/http"
	"testing"
)

func TestDeleteQuestionFromList(t *testing.T) {
	testCases := []struct {
		name     string
		db       *test.MockListManager
		ctx      *test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name: "Success",
			db:   test.NewMockListManager(test.MockListManagerConfig{}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier:     fixtures.UserId,
					constants.ListIdPathIdentifier:     fixtures.ListId,
					constants.QuestionIdPathIdentifier: fixtures.QuestionId,
				},
			}),
			wantCode: http.StatusOK,
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
			name: "Failure - Delete Error",
			db: test.NewMockListManager(test.MockListManagerConfig{
				RemoveQuestionFromList: func(userId, listId, questionId string) error {
					return errors.New("unable to delete that question")
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
			question.DeleteQuestionFromList(tc.ctx, tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("DeleteQuestionFromList() wanted code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); !ok && tc.wantErr && tc.ctx.ResponseData != nil {
				t.Errorf("DeleteQuestionFromList() wanted error Response, got = %v", tc.ctx.ResponseData)
			}
		})
	}
}
