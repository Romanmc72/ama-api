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

func TestCreateQuestion(t *testing.T) {
	logger := logging.GetLogger()
	testCases := []struct {
		name    string
		db      database.Database
		q       application.NewQuestion
		wantErr bool
	}{
		{
			name: "No Prompt",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {},
					},
				}),
				logger,
			),
			q: application.NewQuestion{
				Prompt: "    ",
				Tags:   []string{"1", "2"},
			},
			wantErr: true,
		},
		{
			name: "No Tags",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {},
					},
				}),
				logger,
			),
			q: application.NewQuestion{
				Prompt: "That's a great question",
				Tags:   []string{},
			},
			wantErr: true,
		},
		{
			name: "Duplicate Tags",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {},
					},
				}),
				logger,
			),
			q: application.NewQuestion{
				Prompt: "That's a great question",
				Tags:   []string{"test", "test"},
			},
			wantErr: true,
		},
		{
			name: "Write Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {
							MockError: errors.New("cannot add document"),
						},
					},
				}),
				logger,
			),
			q:       fixtures.ValidNewQuestion,
			wantErr: true,
		},
		{
			name: "Success",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {},
					},
				}),
				logger,
			),
			q:       fixtures.ValidNewQuestion,
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.db.CreateQuestion(&tc.q)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}
