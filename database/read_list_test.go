package database_test

import (
	"ama/api/application"
	"ama/api/application/list"
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

func TestReadList(t *testing.T) {
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
									Data: fixtures.ValidUser.BaseUser,
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
			wantErr: false,
		},
		{
			name: "User Not Found",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {},
					},
				}),
				logger,
			),
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
									ID: fixtures.UserId,
									Data: user.BaseUser{
										Name:         "ben",
										Email:        "bkenobi@jedicouncil.gov",
										Subscription: fixtures.ValidUserSubscription,
										Settings:     fixtures.ValidUserSettings,
										Lists: []list.List{
											{
												ID:   "not-the-droids-you-are-looking-for",
												Name: "Obiwan's list",
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
			name: "List Read Error",
			db: database.ManualTestConnect(
				t.Context(),
				test.NewMockDatabase(&test.MockDBConfig{
					Collections: map[string]test.MockCollectionConfig{
						constants.UserCollection: {
							Documents: map[string]test.MockDocumentConfig{
								fixtures.UserId: {
									ID:   fixtures.UserId,
									Data: fixtures.ValidUser.BaseUser,
									NestedCollections: map[string]test.MockCollectionConfig{
										fixtures.ListId: {
											MockError: errors.New("error reading questions"),
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
	}

	limit := 10
	finalId := fixtures.QuestionId
	tags := fixtures.Tags
	wantedList := fixtures.ValidLists[0]
	wantedQuestions := []application.Question{fixtures.ValidQuestion}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l, q, err := tc.db.ReadList(fixtures.UserId, fixtures.ListId, limit, finalId, tags)
			if (err != nil) != tc.wantErr {
				t.Errorf("ReadList() wantedErr = %v; got = %v", tc.wantErr, err)
			}
			if !tc.wantErr && !reflect.DeepEqual(l, wantedList) {
				t.Errorf("ReadList() wantedList = %v; got = %v", wantedList, l)
			}
			if !tc.wantErr && !reflect.DeepEqual(q, wantedQuestions) {
				t.Errorf("ReadList() wantedQuestions = %v; got = %v", wantedQuestions, q)
			}
		})
	}
}
