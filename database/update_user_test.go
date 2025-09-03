package database_test

import (
	"ama/api/application"
	"ama/api/application/user"
	"ama/api/constants"
	"ama/api/database"
	"ama/api/logging"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"testing"
)

func TestUpdateUser(t *testing.T) {
	userId := "test-user-id"
	testCases := []struct {
		name    string
		db      database.Database
		user    *application.User
		wantErr bool
	}{
		{
			name: "Invalid User - Empty BaseUser",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{},
				}),
				logging.GetLogger(),
			),
			user: &application.User{
				ID:       userId,
				BaseUser: user.BaseUser{},
			},
			wantErr: true,
		},
		{
			name: "Write Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								userId: {
									ID:  userId,
									Err: errors.New("write error"),
								},
							},
						},
					},
				}),
				logging.GetLogger(),
			),
			user: &application.User{
				ID:       userId,
				BaseUser: fixtures.BaseUserValid,
			},
			wantErr: true,
		},
		{
			name: "Succeeds",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								userId: {
									ID: userId,
								},
							},
						},
					},
				}),
				logging.GetLogger(),
			),
			user: &application.User{
				ID:       userId,
				BaseUser: fixtures.BaseUserValid,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.db.UpdateUser(tc.user)
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
