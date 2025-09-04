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

func TestRemoveQuestionFromList(t *testing.T) {
	logger := logging.GetLogger()
	testCases := []struct {
		name       string
		db         database.Database
		userId     string
		listId     string
		questionId string
		wantErr    bool
	}{
		{
			name: "Successfully Removes",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserValid.ID: {
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ValidLists[0].ID: {
											Documents: map[string]test.MockDocumentConfig{
												fixtures.QuestionId: {
													Data: fixtures.ValidQuestion,
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
			userId:     fixtures.UserValid.ID,
			questionId: fixtures.QuestionId,
			listId:     fixtures.ValidLists[0].ID,
			wantErr:    false,
		},
		{
			name: "Missing Required Parameters - userId",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserValid.ID: {
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ValidLists[0].ID: {
											Documents: map[string]test.MockDocumentConfig{
												fixtures.QuestionId: {
													Data: fixtures.ValidQuestion,
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
			userId:     "",
			questionId: fixtures.QuestionId,
			listId:     fixtures.ValidLists[0].ID,
			wantErr:    true,
		},
		{
			name: "Missing Required Parameters - listId",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserValid.ID: {
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ValidLists[0].ID: {
											Documents: map[string]test.MockDocumentConfig{
												fixtures.QuestionId: {
													Data: fixtures.ValidQuestion,
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
			userId:     fixtures.UserValid.ID,
			questionId: fixtures.QuestionId,
			listId:     "",
			wantErr:    true,
		},
		{
			name: "Missing Required Parameters - questionId",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserValid.ID: {
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ValidLists[0].ID: {
											Documents: map[string]test.MockDocumentConfig{
												fixtures.QuestionId: {
													Data: fixtures.ValidQuestion,
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
			userId:     fixtures.UserValid.ID,
			questionId: "",
			listId:     fixtures.ValidLists[0].ID,
			wantErr:    true,
		},
		{
			name: "Database Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserValid.ID: {
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ValidLists[0].ID: {
											Documents: map[string]test.MockDocumentConfig{
												fixtures.QuestionId: {
													Data: fixtures.ValidQuestion,
													Err:  errors.New("some database error"),
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
			userId:     fixtures.UserValid.ID,
			questionId: fixtures.QuestionId,
			listId:     fixtures.ValidLists[0].ID,
			wantErr:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.db.RemoveQuestionFromList(tc.userId, tc.listId, tc.questionId)
			if (err != nil) != tc.wantErr {
				t.Errorf("RemoveQuestionFromList() wantedErr = %v; got = %v", tc.wantErr, err)
			}
		})
	}
}
