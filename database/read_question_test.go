package database_test

import (
	"ama/api/constants"
	"ama/api/database"
	"ama/api/logging"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestReadQuestion(t *testing.T) {
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
			name: "Document Not Found Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.QuestionId: {
									ID:     fixtures.QuestionId,
									Data:   fixtures.ValidDatabaseQuestion,
									GetErr: status.Error(codes.NotFound, "unable to find document"),
								},
							},
						},
					},
				}),
				logger,
			),
			wantErr: true,
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
									ID:     fixtures.QuestionId,
									Data:   fixtures.ValidDatabaseQuestion,
									GetErr: errors.New("firestore is on a lunch break and will return after 1:00PM PST"),
								},
							},
						},
					},
				}),
				logger,
			),
			wantErr: true,
		},
		{
			name: "Data Conversion Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.QuestionCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.QuestionId: {
									ID: fixtures.QuestionId,
									Data: struct {
										Prompt bool  `json:"prompt"`
										Tags   []int `json:"tags"`
									}{Prompt: false, Tags: []int{8, 19, 99}},
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
			q, err := tc.db.ReadQuestion(fixtures.QuestionId)
			if (err != nil) != tc.wantErr {
				t.Errorf("ReadQuestion() wantedErr = %v; got = %v", tc.wantErr, err)
			}
			if !tc.wantErr && !reflect.DeepEqual(q, fixtures.ValidQuestion) {
				t.Errorf("ReadQuestion() wanted = %v; got = %v", fixtures.ValidQuestion, q)
			}
		})
	}
}
