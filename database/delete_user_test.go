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

func TestDeleteUser(t *testing.T) {
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
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID:   fixtures.UserId,
									Data: fixtures.ValidUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											Documents: map[string]test.MockDocumentConfig{
												fixtures.QuestionId: {
													ID:   fixtures.QuestionId,
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
			name: "Read Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{},
						},
					},
				}),
				logger,
			),
			wantErr: true,
		},
		{
			name: "Delete List Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID:   fixtures.UserId,
									Data: fixtures.ValidUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											MockError: errors.New("unable to delete collection"),
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
		// TODO: Revisit; Likely not covered, and will be very tedious to mock
		// just the delete throwing an error and not the read.
		{
			name: "Delete User Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID:        fixtures.UserId,
									Data:      fixtures.ValidUser,
									DeleteErr: errors.New("delete error"),
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											Documents: map[string]test.MockDocumentConfig{
												fixtures.QuestionId: {
													ID:   fixtures.QuestionId,
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
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.db.DeleteUser(fixtures.UserId)
			if (err != nil) != tc.wantErr {
				t.Errorf("DeleteUser() wantedErr = %v; got = %v", tc.wantErr, err)
			}
		})
	}
}
