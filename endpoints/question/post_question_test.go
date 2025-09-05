package question

import (
	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/test"
	"net/http"
	"testing"
)

func TestPostQuestion(t *testing.T) {
	testQuestionId := ""
	testNewQuestion := []byte(`{"prompt": "test question", "tags": ["test"]}`)
	tests := []struct {
		name              string
		inputBody         []byte
		dbShouldFail      bool
		expectedCode      int
		expectedError     bool
		expectedData      application.Question
		expectedErrorData responses.ErrorResponse
	}{
		{
			name:          "Success case",
			inputBody:     testNewQuestion,
			dbShouldFail:  false,
			expectedCode:  http.StatusCreated,
			expectedError: false,
			// TODO: Get this to work with the test
			expectedData: application.Question{
				ID:     testQuestionId,
				Prompt: "test question",
				Tags:   []string{"test"},
			},
			expectedErrorData: responses.ErrorResponse{},
		},
		{
			name:              "Invalid JSON",
			inputBody:         []byte(`{"invalid": "input"}`),
			dbShouldFail:      false,
			expectedCode:      http.StatusBadRequest,
			expectedError:     true,
			expectedData:      application.Question{},
			expectedErrorData: responses.NewError("invalid data"),
		},
		{
			name:              "Database error",
			inputBody:         testNewQuestion,
			dbShouldFail:      true,
			expectedCode:      http.StatusInternalServerError,
			expectedError:     true,
			expectedData:      application.Question{},
			expectedErrorData: responses.NewError("encountered an error writing that data"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: tt.inputBody,
			})
			mockDB := test.NewMockQuestionManager()
			mockDB.ShouldError = tt.dbShouldFail

			PostQuestion(mockCtx, mockDB)

			if mockCtx.ResponseCode != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, mockCtx.ResponseCode)
			}

			if tt.expectedError {
				if _, ok := mockCtx.ResponseData.(responses.ErrorResponse); !ok {
					t.Errorf("Expected error response, got something else %v", mockCtx.ResponseData)
				}
			} else {
				if question, ok := mockCtx.ResponseData.(application.Question); !ok {
					t.Errorf("Expected question response, got something else, %v", mockCtx.ResponseData)
				} else if question.ID != testQuestionId || question.Prompt != tt.expectedData.Prompt {
					t.Errorf("Expected question ID to be %s and prompt to be %s, got ID %s and prompt %s", testQuestionId, tt.expectedData.Prompt, question.ID, question.Prompt)
				}
			}
		})
	}
}
