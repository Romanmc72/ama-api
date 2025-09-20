package database_test

import (
	"ama/api/constants"
	"ama/api/database"
	"ama/api/logging"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"testing"

	"cloud.google.com/go/firestore"
)

func TestDeleteList(t *testing.T) {
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
									Data: fixtures.ValidBaseUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											QueryDocuments: []test.MockDocumentConfig{
												{
													ID:   fixtures.QuestionId,
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
			wantErr: false,
		},
		{
			name: "Delete Collection Error",
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
											MockError: errors.New("delete error"),
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
		{
			name: "Delete Collection Error - BulkWriter",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					BulkWriter: test.MockBulkWriter{
						MockDelete: func(doc *firestore.DocumentRef) (*firestore.BulkWriterJob, error) {
							return nil, errors.New("unable to delete")
						},
					},
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID:   fixtures.UserId,
									Data: fixtures.ValidBaseUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											QueryDocuments: []test.MockDocumentConfig{
												{
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
		{
			name: "Transaction Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					TransacitonErr: errors.New("transaction failure"),
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID:   fixtures.UserId,
									Data: fixtures.ValidBaseUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											QueryDocuments: []test.MockDocumentConfig{
												{
													ID:   fixtures.QuestionId,
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
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.db.DeleteList(fixtures.UserId, fixtures.ListId)
			if (err != nil) != tc.wantErr {
				t.Errorf("DeleteList() wantedErr = %v; got = %v", tc.wantErr, err)
			}
		})
	}
}
