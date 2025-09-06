package question_test

import (
	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/question"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetQuestionById(t *testing.T) {
	testCases := []struct {
		name     string
		db       test.MockQuestionManager
		ctx      test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name:     "Success - Found Questions",
			wantCode: http.StatusOK,
			db: *test.NewMockQuestionManager(test.MockQuestionManagerConfig{
				ReadQuestion: func(id string) (application.Question, error) {
					return fixtures.ValidQuestion, nil
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.QuestionIdPathIdentifier: fixtures.QuestionId,
				},
			}),
			wantErr: false,
		},
		{
			name:     "Failure - Not Found",
			wantCode: http.StatusNotFound,
			db: *test.NewMockQuestionManager(test.MockQuestionManagerConfig{
				ReadQuestion: func(id string) (application.Question, error) {
					return application.Question{}, status.Error(codes.NotFound, "not found")
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.QuestionIdPathIdentifier: fixtures.QuestionId,
				},
			}),
			wantErr: true,
		},
		{
			name:     "Failure - Internal Server Error",
			wantCode: http.StatusInternalServerError,
			db: *test.NewMockQuestionManager(test.MockQuestionManagerConfig{
				ReadQuestion: func(id string) (application.Question, error) {
					return application.Question{}, errors.New("failed to read")
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.QuestionIdPathIdentifier: fixtures.QuestionId,
				},
			}),
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			question.GetQuestionById(&tc.ctx, &tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("GetQuestionById() wanted Code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); tc.wantErr && !ok {
				t.Errorf("GetQuestionById() wanted error, got %v", tc.ctx.ResponseData)
			}
			if q, ok := tc.ctx.ResponseData.(application.Question); !tc.wantErr && (!ok || !reflect.DeepEqual(q, fixtures.ValidQuestion)) {
				t.Errorf("GetQuestionById() wanted %v, got %v", fixtures.ValidQuestion, tc.ctx.ResponseData)
			}
		})
	}
}
