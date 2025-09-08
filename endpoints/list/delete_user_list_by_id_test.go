package list_test

import (
	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/list"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"net/http"
	"testing"
)

func TestDeleteUserListByID(t *testing.T) {
	testCases := []struct {
		name     string
		db       *test.MockListManager
		ctx      *test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name: "Success",
			db: test.NewMockListManager(test.MockListManagerConfig{
				ReadUser: func(id string) (application.User, error) {
					return fixtures.ValidUser, nil
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
					constants.ListIdPathIdentifier: fixtures.ListId,
				},
			}),
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name: "Failure - Blank Path Identifiers",
			db:   test.NewMockListManager(test.MockListManagerConfig{}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: "     ",
					constants.ListIdPathIdentifier: "     ",
				},
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Read Error",
			db: test.NewMockListManager(test.MockListManagerConfig{
				ReadUser: func(id string) (application.User, error) {
					return application.User{}, errors.New("user not found")
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
					constants.ListIdPathIdentifier: fixtures.ListId,
				},
			}),
			wantCode: http.StatusNotFound,
			wantErr:  true,
		},
		{
			name: "Success - User Has No Such List",
			db: test.NewMockListManager(test.MockListManagerConfig{
				ReadUser: func(id string) (application.User, error) {
					return fixtures.ValidUser, nil
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
					constants.ListIdPathIdentifier: "list-that-does-not-exist",
				},
			}),
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name: "Failure - Cannot Delete Liked Questions",
			db: test.NewMockListManager(test.MockListManagerConfig{
				ReadUser: func(id string) (application.User, error) {
					return fixtures.ValidUser, nil
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
					constants.ListIdPathIdentifier: fixtures.NewId,
				},
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Delete Error",
			db: test.NewMockListManager(test.MockListManagerConfig{
				ReadUser: func(id string) (application.User, error) {
					return fixtures.ValidUser, nil
				},
				DeleteList: func(userId, listId string) error {
					return errors.New("unable to delete list")
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
					constants.ListIdPathIdentifier: fixtures.ListId,
				},
			}),
			wantCode: http.StatusInternalServerError,
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list.DeleteUserListByID(tc.ctx, tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("DeleteUserListByID() wanted code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); !ok && tc.wantErr && tc.ctx.ResponseData != nil {
				t.Errorf("DeleteUserListByID() wanted error Response, got = %v", tc.ctx.ResponseData)
			}
		})
	}
}
