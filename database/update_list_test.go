package database_test

import (
	"ama/api/application/list"
	"ama/api/constants"
	"ama/api/database"
	"ama/api/logging"
	"ama/api/test"
	"ama/api/test/fixtures"
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/firestore"
)

func TestUpdateList(t *testing.T) {
	logger := logging.GetLogger()
	txErrDb := test.NewMockDatabase(&test.MockDBConfig{
		Collections: map[string]test.MockCollectionConfig{
			constants.UserCollection: {
				Documents: map[string]test.MockDocumentConfig{
					fixtures.UserId: {
						GetErr: errors.New("get error"),
					},
				},
			},
		},
	})
	txErrDb.MockRunTransaction = func(
		ctx context.Context,
		f func(context.Context, *firestore.Transaction) error,
		t ...firestore.TransactionOption,
	) error {
		return errors.New("Something went wrong in the transaction")
	}
	testCases := []struct {
		name    string
		db      database.Database
		userId  string
		list    *list.List
		wantErr bool
	}{
		{
			name: "Updates List",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									Data: fixtures.ValidBaseUser,
									ID:   fixtures.UserId,
								},
							},
						},
					},
				}),
				logger,
			),
			userId:  fixtures.QuestionId,
			list:    &fixtures.ValidLists[0],
			wantErr: false,
		},
		{
			name: "Missing List ID",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									Data: fixtures.ValidBaseUser,
									ID:   fixtures.UserId,
								},
							},
						},
					},
				}),
				logger,
			),
			userId: fixtures.QuestionId,
			list: &list.List{
				Name: "I have a name but no ID",
			},
			wantErr: true,
		},
		{
			name: "Missing User ID",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									Data: fixtures.ValidBaseUser,
									ID:   fixtures.UserId,
								},
							},
						},
					},
				}),
				logger,
			),
			userId:  "",
			list:    &fixtures.ValidLists[0],
			wantErr: true,
		},
		{
			name: "List Not Found",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									Data:   fixtures.ValidBaseUser,
									ID:     fixtures.UserId,
									GetErr: errors.New("get error"),
								},
							},
						},
					},
				}),
				logger,
			),
			userId: "",
			list: &list.List{
				ID:   "incorrect-list-id",
				Name: "too bad this will fail",
			},
			wantErr: true,
		},
		{
			name:    "Transaction Error",
			db:      database.ManualTestConnect(t.Context(), txErrDb, logger),
			userId:  fixtures.UserId,
			list:    &fixtures.ValidLists[0],
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.db.UpdateList(tc.userId, *tc.list)
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateList() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
