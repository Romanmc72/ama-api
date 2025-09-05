package database_test

import (
	"ama/api/application/user"
	"ama/api/constants"
	"ama/api/database"
	"ama/api/logging"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"reflect"
	"testing"
)

func TestCreateUser(t *testing.T) {
	logger := logging.GetLogger()
	testCases := []struct{
		name string
		user user.BaseUser
		db database.Database
		wantErr bool
	}{
		{
			name: "Success",
			user: fixtures.ValidBaseUser,
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
			wantErr: false,
		},
		{
			name: "Read Error",
			user: fixtures.ValidBaseUser,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID: fixtures.UserId,
									GetErr: errors.New("cannot read the user"),
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
			user: fixtures.ValidBaseUser,
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
			name: "Write Error",
			user: fixtures.ValidBaseUser,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID: fixtures.UserId,
									Data: nil,
									SetErr: errors.New("cannot write the user for some reason"),
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
			u, err := tc.db.CreateUser(tc.user)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreateUser() wantedErr = %v; got = %v", tc.wantErr, err)
			}
			if !tc.wantErr && !reflect.DeepEqual(u, fixtures.ValidUser) {
				t.Errorf("CreateUser() wanted User = %s; got = %s", fixtures.ValidUser, u)
			}
		})
	}
}
