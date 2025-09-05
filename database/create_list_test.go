package database_test

import (
	"ama/api/application/list"
	"ama/api/constants"
	"ama/api/database"
	"ama/api/logging"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"testing"
)

func TestCreateList(t *testing.T) {
	testListName := "test list"
	logger := logging.GetLogger()
	testCases := []struct{
		name string
		l    list.List
		userId    string
		db database.Database
		wantErr bool
	}{
		{
			name: "Success - Exists",
			l: fixtures.ValidLists[0],
			userId: fixtures.UserId,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID: fixtures.UserId,
									Data: fixtures.ValidBaseUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											QueryDocuments: []test.MockDocumentConfig{
												{
													ID: fixtures.QuestionId,
													Data: fixtures.ValidDatabaseQuestion,
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
			wantErr: false,
		},
		{
			// Theoretically impossible for the uuid4 to generate a duplicate but
			// this gets 100% test coverage :flex:
			name: "Success - Exists No ID Provided",
			l: list.List{Name: testListName},
			userId: fixtures.UserId,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID: fixtures.UserId,
									Data: fixtures.ValidBaseUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.NewId: {
											QueryDocuments: []test.MockDocumentConfig{
												{
													ID: fixtures.QuestionId,
													Data: fixtures.ValidDatabaseQuestion,
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
			wantErr: false,
		},
		{
			name: "Success - Does Not Exist",
			l: list.List{ID: "fester-fiesta", Name: testListName},
			userId: fixtures.UserId,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID: fixtures.UserId,
									Data: fixtures.ValidBaseUser,
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
			name: "Empty userId",
			l: list.List{Name: testListName},
			userId: "   ",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID: fixtures.UserId,
									Data: fixtures.ValidBaseUser,
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
			name: "Read Error",
			l: list.List{Name: testListName},
			userId: fixtures.UserId,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID: fixtures.UserId,
									Data: fixtures.ValidBaseUser,
									GetErr: errors.New("cannot read user"),
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
			name: "Existence Check Error",
			l: list.List{ID: fixtures.ListId, Name: testListName},
			userId: fixtures.UserId,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID: fixtures.UserId,
									Data: fixtures.ValidBaseUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											QueryDocuments: []test.MockDocumentConfig{
												{
													ID: fixtures.QuestionId,
													Data: fixtures.ValidDatabaseQuestion,
												},
											},
											MockError: errors.New("error iterating list"),
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
			err := tc.db.CreateList(tc.userId, tc.l)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreateList() wantedErr = %v; got = %v", tc.wantErr, err)
			}
		})
	}
}
