package database_test

import (
	"ama/api/application"
	"ama/api/constants"
	"ama/api/database"
	"ama/api/logging"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"testing"
)

func TestUpdateQuestion(t *testing.T) {
	logger := logging.GetLogger()
	testCases := []struct {
		name       string
		db         database.Database
		questionId string
		question   *application.NewQuestion
		wantErr    bool
	}{
		{
			name: "Updates Question",
			db: database.ManualTestConnect(t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.QuestionId: {
									Data: fixtures.ValidNewQuestion,
								},
							},
						},
					},
				}),
				logger,
			),
			questionId: fixtures.QuestionId,
			question:   &fixtures.ValidNewQuestion,
			wantErr:    false,
		},
		{
			name: "Write Error",
			db: database.ManualTestConnect(t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.QuestionId: {
									Data: fixtures.ValidNewQuestion,
									SetErr:  errors.New("write error"),
								},
							},
						},
					},
				}),
				logger,
			),
			questionId: fixtures.QuestionId,
			question:   &fixtures.ValidNewQuestion,
			wantErr:    true,
		},
		{
			name: "Invalid Question",
			db: database.ManualTestConnect(t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {},
					},
				}),
				logger,
			),
			questionId: fixtures.QuestionId,
			question:   &application.NewQuestion{},
			wantErr:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.db.UpdateQuestion(tc.questionId, tc.question)
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateQuestion() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
