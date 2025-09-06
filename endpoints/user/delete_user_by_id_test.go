package user_test

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/user"
	"ama/api/test"
	"ama/api/test/fixtures"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestDeleteUserById(t *testing.T) {
	deleteTime := time.Now()
	deleteSuccessResponse := responses.NewDeleteUserResponse(true, deleteTime.Unix())
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
				DeleteUser: func(id string) (time.Time, error) {
					return deleteTime, nil
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
			name: "Failure - Internal Server Error",
			db: *test.NewMockUserManager(test.MockUserManagerConfig{
				DeleteUser: func(id string) (time.Time, error) {
					return time.Now(), errors.New("delete error")
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantCode: http.StatusInternalServerError,
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user.DeleteUserById(&tc.ctx, &tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("DeleteUserById() wanted Code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); tc.wantErr && !ok {
				t.Errorf("DeleteUserById() wanted error, got %v", tc.ctx.ResponseData)
			}
			if r, ok := tc.ctx.ResponseData.(responses.DeleteUserResponse); !tc.wantErr && (!ok || !reflect.DeepEqual(r, deleteSuccessResponse)) {
				t.Errorf("DeleteUserById() wanted %v, got %v", deleteSuccessResponse, tc.ctx.ResponseData)
			}
		})
	}
}
