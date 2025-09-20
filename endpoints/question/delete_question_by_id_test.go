package question_test

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/question"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestDeleteQuestionById(t *testing.T) {
	successTime := time.Now()
	resp := responses.SuccessResponse{
		Success: true,
		Time:    successTime.Unix(),
	}
	testCases := []struct {
		name     string
		db       test.MockQuestionManager
		ctx      test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name:     "Success",
			wantCode: http.StatusOK,
			db: *test.NewMockQuestionManager(test.MockQuestionManagerConfig{
				DeleteQuestion: func(id string) (time.Time, error) {
					return successTime, nil
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
			name: "Failure - Internal Server Error",
			db: *test.NewMockQuestionManager(test.MockQuestionManagerConfig{
				DeleteQuestion: func(id string) (time.Time, error) {
					return successTime, errors.New("delete error")
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.QuestionIdPathIdentifier: fixtures.QuestionId,
				},
			}),
			wantCode: http.StatusInternalServerError,
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			question.DeleteQuestionById(&tc.ctx, &tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("DeleteQuestionById() wanted Code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); tc.wantErr && !ok {
				t.Errorf("DeleteQuestionById() wanted error, got %v", tc.ctx.ResponseData)
			}
			if r, ok := tc.ctx.ResponseData.(responses.SuccessResponse); !tc.wantErr && (!ok || !reflect.DeepEqual(resp, r)) {
				t.Errorf("DeleteQuestionById() wanted %v, got %v", resp, tc.ctx.ResponseData)
			}
		})
	}
}
