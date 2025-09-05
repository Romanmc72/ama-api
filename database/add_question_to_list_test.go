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

func TestAddQuestionToList(t *testing.T) {
	logger := logging.GetLogger()
	testCases := []struct {
		name     string
		userId   string
		listId   string
		question application.Question
		db       database.Database
		wantErr  bool
	}{
		{
			name:     "Success",
			userId:   fixtures.UserId,
			listId:   fixtures.ListId,
			question: fixtures.ValidQuestion,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID:   fixtures.UserId,
									Data: fixtures.ValidBaseUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											Documents: map[string]test.MockDocumentConfig{},
										},
									},
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
			name:     "Invalid Input - Blank userId",
			userId:   "",
			listId:   fixtures.ListId,
			question: fixtures.ValidQuestion,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{}),
				logger,
			),
			wantErr: true,
		},
		{
			name:     "Invalid Input - Blank listId",
			userId:   fixtures.UserId,
			listId:   "",
			question: fixtures.ValidQuestion,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{}),
				logger,
			),
			wantErr: true,
		},
		{
			name:     "Invalid Input - Blank question.ID",
			userId:   fixtures.UserId,
			listId:   fixtures.ListId,
			question: application.Question{},
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{}),
				logger,
			),
			wantErr: true,
		},
		{
			name:     "Write Error",
			userId:   fixtures.UserId,
			listId:   fixtures.ListId,
			question: fixtures.ValidQuestion,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID:   fixtures.UserId,
									Data: fixtures.ValidBaseUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											Documents: map[string]test.MockDocumentConfig{
												fixtures.QuestionId: {
													SetErr: errors.New("write error"),
												},
											},
										},
									},
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
			err := tc.db.AddQuestionToList(tc.userId, tc.listId, tc.question)
			if (err != nil) != tc.wantErr {
				t.Errorf("AddQuestionToList() wantedErr = %v; got = %v", tc.wantErr, err)
			}
		})
	}
}
