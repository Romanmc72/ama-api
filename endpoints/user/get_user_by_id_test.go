package user_test

import (
	"ama/api/application"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/user"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestGetUserById(t *testing.T) {
	testCases := []struct {
		name     string
		db       test.MockUserManager
		ctx      test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name:     "Success",
			wantCode: http.StatusOK,
			db: *test.NewMockUserManager(test.MockUserManagerConfig{
				ReadUser: func(id string) (application.User, error) {
					return fixtures.ValidUser, nil
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantErr: false,
		},
		{
			name: "Failure - Blank userId",
			db:   *test.NewMockUserManager(test.MockUserManagerConfig{}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: "   ",
				},
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Not Found",
			db: *test.NewMockUserManager(test.MockUserManagerConfig{
				ReadUser: func(id string) (application.User, error) {
					return application.User{}, errors.New("user not found")
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
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
			user.GetUserByUserId(&tc.ctx, &tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("GetUserByUserId() wanted Code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); tc.wantErr && !ok {
				t.Errorf("GetUserByUserId() wanted error, got %v", tc.ctx.ResponseData)
			}
			if u, ok := tc.ctx.ResponseData.(application.User); !tc.wantErr && (!ok || !reflect.DeepEqual(u, fixtures.ValidUser)) {
				t.Errorf("GetUserByUserId() wanted %v, got %v", fixtures.ValidUser, tc.ctx.ResponseData)
			}
		})
	}
}
