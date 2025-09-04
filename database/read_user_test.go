package database_test

import (
	"ama/api/constants"
	"ama/api/database"
	"ama/api/logging"
	"ama/api/test"
	"ama/api/test/fixtures"
	"reflect"
	"testing"
)

// TODO: find out why failures are not failing
func TestReadUser(t *testing.T) {
	logger := logging.GetLogger()
	testCases := []struct {
		name    string
		id      string
		db      database.Database
		wantErr bool
	}{
		{
			name: "Read Success",
			id:   fixtures.UserValid.ID,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserValid.ID: {
									Data: fixtures.BaseUserValid,
									ID:   fixtures.UserValid.ID,
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
			id:   fixtures.UserValid.ID,
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
			name: "Conversion Error",
			id:   fixtures.UserValid.ID,
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserValid.ID: {
									Data: map[string]any{
										"name":         "mr. bad data",
										"subscription": -11,
										"settings":     []string{"lol", "these", "are", "not", "settings"},
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
			u, err := tc.db.ReadUser(tc.id)
			if (err != nil) != tc.wantErr {
				t.Errorf("ReadUser() wantedErr = %v; got = %v", tc.wantErr, err)
			}
			if !tc.wantErr && !reflect.DeepEqual(u, fixtures.UserValid) {
				t.Errorf("ReadUser() wanted = %s; got = %s", fixtures.UserValid, u)
			}
		})
	}
}
