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
	"testing"
)

func TestGetQuestions(t *testing.T) {
	qSlice := []application.Question{fixtures.ValidQuestion}
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
				ReadQuestions: func(limit int, finalId string, tags []string) ([]application.Question, error) {
					return qSlice, nil
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				QueryValues: map[string][]string{
					constants.LimitParam:   {"10"},
					constants.TagParam:     {"test"},
					constants.FinalIdParam: {fixtures.QuestionId},
				},
			}),
			wantErr: false,
		},
		{
			name:     "Success - No Questions Found",
			wantCode: http.StatusOK,
			db: *test.NewMockQuestionManager(test.MockQuestionManagerConfig{
				ReadQuestions: func(limit int, finalId string, tags []string) ([]application.Question, error) {
					return []application.Question{}, nil
				},
			}),
			ctx:     *test.NewMockAPIContext(test.MockAPIContextConfig{}),
			wantErr: false,
		},
		{
			name: "Failure - Internal Server Error",
			db: *test.NewMockQuestionManager(test.MockQuestionManagerConfig{
				ReadQuestions: func(limit int, finalId string, tags []string) ([]application.Question, error) {
					return []application.Question{}, errors.New("read error")
				},
			}),
			ctx:      *test.NewMockAPIContext(test.MockAPIContextConfig{}),
			wantCode: http.StatusInternalServerError,
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			question.GetQuestions(&tc.ctx, &tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("GetQuestions() wanted Code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); tc.wantErr && !ok {
				t.Errorf("GetQuestions() wanted error, got %v", tc.ctx.ResponseData)
			}
			if _, ok := tc.ctx.ResponseData.([]application.Question); !tc.wantErr && !ok {
				t.Errorf("GetQuestions() wanted %v, got %v", qSlice, tc.ctx.ResponseData)
			}
		})
	}
}
