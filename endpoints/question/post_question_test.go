package question_test

import (
	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/question"
	"ama/api/interfaces"
	"ama/api/test"
	"ama/api/test/fixtures"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestPostQuestion(t *testing.T) {
	qBytes, _ := json.Marshal(fixtures.ValidNewQuestion)
	invalidQBytes, _ := json.Marshal(fixtures.InvalidNewQuestion)
	testCases := []struct {
		name     string
		db       test.MockQuestionManager
		ctx      test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name:     "Success",
			wantCode: http.StatusCreated,
			db: *test.NewMockQuestionManager(test.MockQuestionManagerConfig{
				CreateQuestion: func(questionData interfaces.QuestionConverter) (application.Question, error) {
					return fixtures.ValidQuestion, nil
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: qBytes,
				Params: map[string]string{
					constants.QuestionIdPathIdentifier: fixtures.QuestionId,
				},
			}),
			wantErr: false,
		},
		{
			name: "Failure - Bad Request Invalid JSON",
			db:   *test.NewMockQuestionManager(test.MockQuestionManagerConfig{}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: []byte(`{"prompt": [1,2,3,4], "tags": "you are it"}`),
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Bad Request Invalid Input Data",
			db:   *test.NewMockQuestionManager(test.MockQuestionManagerConfig{}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: invalidQBytes,
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Internal Server Error",
			db: *test.NewMockQuestionManager(test.MockQuestionManagerConfig{
				CreateQuestion: func(questionData interfaces.QuestionConverter) (application.Question, error) {
					return application.Question{}, errors.New("could not create")
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: qBytes,
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
			question.PostQuestion(&tc.ctx, &tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("PostQuestion() wanted Code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); tc.wantErr && !ok {
				t.Errorf("PostQuestion() wanted error, got %v", tc.ctx.ResponseData)
			}
			if q, ok := tc.ctx.ResponseData.(application.Question); !tc.wantErr && (!ok || !reflect.DeepEqual(q, fixtures.ValidQuestion)) {
				t.Errorf("PostQuestion() wanted %v, got %v", fixtures.ValidQuestion, tc.ctx.ResponseData)
			}
		})
	}
}
