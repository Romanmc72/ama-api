package user_test

import (
	"ama/api/application"
	"ama/api/application/responses"
	appUser "ama/api/application/user"
	"ama/api/endpoints/user"
	"ama/api/test"
	"ama/api/test/fixtures"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestPostUser(t *testing.T) {
	validUserBytes, _ := json.Marshal(fixtures.ValidBaseUser)
	invalidUserBytes, _ := json.Marshal(fixtures.InvalidBaseUser)
	testCases := []struct {
		name     string
		db       test.MockUserManager
		ctx      test.MockAPIContext
		wantCode int
		wantErr  bool
	}{
		{
			name:     "Success",
			wantCode: http.StatusCreated,
			db: *test.NewMockUserManager(test.MockUserManagerConfig{
				CreateUser: func(userData appUser.BaseUser) (application.User, error) {
					return fixtures.ValidUser, nil
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: validUserBytes,
			}),
			wantErr: false,
		},
		{
			name: "Failure - Invalid JSON",
			db:   *test.NewMockUserManager(test.MockUserManagerConfig{}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: []byte(`{"settings": "invalid input data"`),
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name: "Failure - Invalid Input Data",
			db:   *test.NewMockUserManager(test.MockUserManagerConfig{}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: invalidUserBytes,
			}),
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
		{
			name:     "Failure - Internal Server Error",
			wantCode: http.StatusInternalServerError,
			db: *test.NewMockUserManager(test.MockUserManagerConfig{
				CreateUser: func(userData appUser.BaseUser) (application.User, error) {
					return application.User{}, errors.New("create error")
				},
			}),
			ctx: *test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: validUserBytes,
			}),
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user.PostUser(&tc.ctx, &tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("PostUser() wanted Code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); tc.wantErr && !ok {
				t.Errorf("PostUser() wanted error, got %v", tc.ctx.ResponseData)
			}
			if u, ok := tc.ctx.ResponseData.(application.User); !tc.wantErr && (!ok || !reflect.DeepEqual(u, fixtures.ValidUser)) {
				t.Errorf("PostUser() wanted %v, got %v", fixtures.ValidUser, tc.ctx.ResponseData)
			}
		})
	}
}
