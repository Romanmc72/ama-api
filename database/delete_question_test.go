package database_test

import (
	"ama/api/constants"
	"ama/api/database"
	"ama/api/logging"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"testing"
)

func TestDeleteQuestion(t *testing.T) {
	logger := logging.GetLogger()
	testCases := []struct {
		name    string
		db      database.Database
		wantErr bool
	}{
		{
			name: "Success",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.QuestionId: {
									ID:   fixtures.QuestionId,
									Data: fixtures.ValidDatabaseQuestion,
								},
							},
						},
					},
				}),
				logger,
			),
			wantErr: false,
		},
		{
			name: "Database Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.QuestionId: {
									ID:        fixtures.QuestionId,
									Data:      fixtures.ValidDatabaseQuestion,
									DeleteErr: errors.New("delete error"),
								},
							},
						},
					},
				}),
				logger,
			),
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.db.DeleteQuestion(fixtures.QuestionId)
			if (err != nil) != tc.wantErr {
				t.Errorf("DeleteQuestion() wantedErr = %v; got = %v", tc.wantErr, err)
			}
		})
	}
}
