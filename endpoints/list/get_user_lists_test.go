package list_test

import (
	"ama/api/application"
	appList "ama/api/application/list"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/list"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestGetUserLists(t *testing.T) {
	testCases := []struct {
		name     string
		db       *test.MockUserManager
		ctx      *test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name: "Success",
			db: test.NewMockUserManager(test.MockUserManagerConfig{
				ReadUser: func(id string) (application.User, error) {
					return fixtures.ValidUser, nil
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name: "Failure - Blank Path Identifiers",
			db:   test.NewMockUserManager(test.MockUserManagerConfig{}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{},
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Not Found",
			db: test.NewMockUserManager(test.MockUserManagerConfig{
				ReadUser: func(id string) (application.User, error) {
					return application.User{}, errors.New("user not found")
				},
			}),
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantCode: http.StatusNotFound,
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list.GetUserLists(tc.ctx, tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("GetUserLists() wanted code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); !ok && tc.wantErr && tc.ctx.ResponseData != nil {
				t.Errorf("GetUserLists() wanted error Response, got = %v", tc.ctx.ResponseData)
			}
			if l, ok := tc.ctx.ResponseData.([]appList.List); !tc.wantErr && (!ok || !reflect.DeepEqual(l, fixtures.ValidLists)) {
				t.Errorf("GetUserLists() wanted list = %s Response, got = %v", fixtures.ValidLists[0], tc.ctx.ResponseData)
			}
		})
	}
}
