package user_test

import (
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/endpoints/user"
	"ama/api/interfaces"
	"ama/api/test"
	"ama/api/test/fixtures"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestPutUserById(t *testing.T) {
	validUserBytes, _ := json.Marshal(fixtures.ValidBaseUser)
	testCases := []struct {
		name string
		db interfaces.UserWriter
		ctx *test.MockAPIContext
		wantCode int
		wantErr bool
	}{
		{
			name: "Success",
			db: &test.MockUserManager{},
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: validUserBytes,
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantCode: http.StatusOK,
			wantErr: false,
		},
		{
			name: "Failure - Input JSON Bind Error",
			db: &test.MockUserManager{},
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: []byte(`{"firebaseId": 69, "name": {"fruit": "banana"}}`),
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "Failure - Internal Server Error",
			db: &test.MockUserManager{
				MockUpdateUser: func(userData interfaces.UserConverter) error {
					return errors.New("cannot update user, sorry")
				},
			},
			ctx: test.NewMockAPIContext(test.MockAPIContextConfig{
				InputJSON: validUserBytes,
				Params: map[string]string{
					constants.UserIdPathIdentifier: fixtures.UserId,
				},
			}),
			wantCode: http.StatusInternalServerError,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user.PutUserByUserId(tc.ctx, tc.db)
			if tc.wantCode != tc.ctx.ResponseCode {
				t.Errorf("PutUserById() wanted code = %d; got = %d", tc.wantCode, tc.ctx.ResponseCode)
			}
			if _, ok := tc.ctx.ResponseData.(responses.ErrorResponse); !ok && !tc.wantErr && tc.ctx.ResponseData != nil {
				t.Errorf("PutUserById() wanted error Response, got = %v", tc.ctx.ResponseData)
			}
		})
	}
}
